package server

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/Kwynto/GracefulDB/internal/base/basicsystem/gauth"
	"github.com/Kwynto/GracefulDB/internal/base/core"
	"github.com/Kwynto/GracefulDB/internal/config"
	"github.com/Kwynto/GracefulDB/internal/connectors/grpc"
	"github.com/Kwynto/GracefulDB/internal/connectors/rest"
	"github.com/Kwynto/GracefulDB/internal/connectors/socket"
	"github.com/Kwynto/GracefulDB/internal/connectors/websocketconn"
	"github.com/Kwynto/GracefulDB/internal/manage/webmanage"
	"github.com/Kwynto/GracefulDB/pkg/lib/closer"
)

var stopSignal = make(chan struct{}, 1)

func Run(ctx context.Context, cfg *config.Config) error {
	var closeProcs = &closer.Closer{}

	// TODO: Load the core of the system
	go core.Engine(cfg)
	closeProcs.AddHandler(core.Shutdown) // Register a shutdown handler.

	// Basic system - begin
	// Loading the authorization module
	go gauth.Start()
	closeProcs.AddHandler(gauth.Shutdown) // Register a shutdown handler.
	// Basic system - end

	// TODO: Start Socket connector
	if cfg.SocketConnector.Enable {
		go socket.Start(cfg)
		closeProcs.AddHandler(socket.Shutdown) // Register a shutdown handler.
	}

	// TODO: Start WebSocket connector
	if cfg.WebSocketConnector.Enable {
		go websocketconn.Start(cfg)
		closeProcs.AddHandler(websocketconn.Shutdown) // Register a shutdown handler.
	}

	// Start REST API connector
	if cfg.RestConnector.Enable {
		go rest.Start(cfg)
		closeProcs.AddHandler(rest.Shutdown) // Register a shutdown handler.
	}

	// Start gRPC connector
	if cfg.GrpcConnector.Enable {
		go grpc.Start(cfg)
		closeProcs.AddHandler(grpc.Shutdown) // Register a shutdown handler.
	}

	// TODO: Start web-server for manage system
	if cfg.WebServer.Enable {
		go webmanage.Start(cfg)
		closeProcs.AddHandler(webmanage.Shutdown) // Register a shutdown handler.
	}

	select {
	case <-ctx.Done(): // Waiting for a stop signal from the OS
		break
	case <-stopSignal: // Program signal from the control interface
		break
	}
	slog.Warn("The shutdown process has started.")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), cfg.ShutdownTimeOut)
	defer cancel()

	if err := closeProcs.Close(shutdownCtx); err != nil {
		return fmt.Errorf("server shutdown: %v", err)
	}

	slog.Info("All processes are stopped.")

	return nil
}

func Stop(login string) {
	msg := fmt.Sprintf("A stop signal was received from the control interface from the %s user", login)
	slog.Warn(msg, slog.String("user", login))
	stopSignal <- struct{}{}
}
