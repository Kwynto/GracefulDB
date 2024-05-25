package server

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/Kwynto/GracefulDB/internal/config"
	"github.com/Kwynto/GracefulDB/internal/connectors/grpc"
	"github.com/Kwynto/GracefulDB/internal/connectors/rest"
	"github.com/Kwynto/GracefulDB/internal/connectors/websocketconn"
	"github.com/Kwynto/GracefulDB/internal/engine/basicsystem/gauth"
	"github.com/Kwynto/GracefulDB/internal/engine/basicsystem/instead"
	"github.com/Kwynto/GracefulDB/internal/engine/core"
	"github.com/Kwynto/GracefulDB/internal/manage/webmanage"
	"github.com/Kwynto/GracefulDB/pkg/lib/closer"
	"github.com/Kwynto/GracefulDB/pkg/lib/e"
)

var chStopSignal = make(chan struct{}, 1)

func Run(ctx context.Context, stCfg *config.TConfig) (err error) {
	sOperation := "internal -> server -> Run"
	defer func() { e.Wrapper(sOperation, err) }()

	// TODO: Load the core of the system
	go core.Start(stCfg)
	closer.AddHandler(core.Shutdown) // Register a shutdown handler.

	// Basic system - begin
	// Loading the authorization module
	go gauth.Start()
	closer.AddHandler(gauth.Shutdown) // Register a shutdown handler.

	// Loading the caching system
	go instead.Start()
	closer.AddHandler(instead.Shutdown) // Register a shutdown handler.

	// Basic system - end

	// Start WebSocket connector
	if stCfg.WebSocketConnector.Enable {
		go websocketconn.Start(stCfg)
		closer.AddHandler(websocketconn.Shutdown) // Register a shutdown handler.
	}

	// Start REST API connector
	if stCfg.RestConnector.Enable {
		go rest.Start(stCfg)
		closer.AddHandler(rest.Shutdown) // Register a shutdown handler.
	}

	// Start gRPC connector
	if stCfg.GrpcConnector.Enable {
		go grpc.Start(stCfg)
		closer.AddHandler(grpc.Shutdown) // Register a shutdown handler.
	}

	// TODO: Start web-server for manage system
	if stCfg.WebServer.Enable {
		go webmanage.Start(stCfg)
		closer.AddHandler(webmanage.Shutdown) // Register a shutdown handler.
	}

	select {
	case <-ctx.Done(): // Waiting for a stop signal from the OS
		break
	case <-chStopSignal: // Program signal from the control interface
		break
	}
	slog.Warn("The shutdown process has started.")

	ctxShutdown, fnCancel := context.WithTimeout(context.Background(), stCfg.ShutdownTimeOut)
	defer fnCancel()

	if err := closer.Close(ctxShutdown); err != nil {
		return fmt.Errorf("server shutdown: %v", err)
	}

	slog.Info("All processes are stopped.")

	return nil
}

func Stop(sLogin string) {
	sMsg := fmt.Sprintf("A stop signal was received from the control interface from the %s user", sLogin)
	slog.Warn(sMsg, slog.String("user", sLogin))
	chStopSignal <- struct{}{}
}
