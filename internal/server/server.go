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
	"github.com/Kwynto/GracefulDB/internal/engine/core"
	"github.com/Kwynto/GracefulDB/internal/manage/webmanage"
	"github.com/Kwynto/GracefulDB/pkg/lib/closer"
)

var stopSignal = make(chan struct{}, 1)

func Run(ctx context.Context, cfg *config.Config) error {
	// TODO: Load the core of the system
	go core.Engine(cfg)
	closer.AddHandler(core.Shutdown) // Register a shutdown handler.

	// Basic system - begin
	// Loading the authorization module
	go gauth.Start()
	closer.AddHandler(gauth.Shutdown) // Register a shutdown handler.
	// Basic system - end

	// Start WebSocket connector
	if cfg.WebSocketConnector.Enable {
		go websocketconn.Start(cfg)
		closer.AddHandler(websocketconn.Shutdown) // Register a shutdown handler.
	}

	// Start REST API connector
	if cfg.RestConnector.Enable {
		go rest.Start(cfg)
		closer.AddHandler(rest.Shutdown) // Register a shutdown handler.
	}

	// Start gRPC connector
	if cfg.GrpcConnector.Enable {
		go grpc.Start(cfg)
		closer.AddHandler(grpc.Shutdown) // Register a shutdown handler.
	}

	// TODO: Start web-server for manage system
	if cfg.WebServer.Enable {
		go webmanage.Start(cfg)
		closer.AddHandler(webmanage.Shutdown) // Register a shutdown handler.
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

	if err := closer.Close(shutdownCtx); err != nil {
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
