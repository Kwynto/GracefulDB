package server

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/Kwynto/GracefulDB/internal/base/core"
	"github.com/Kwynto/GracefulDB/internal/config"
	"github.com/Kwynto/GracefulDB/internal/connectors/grpc"
	"github.com/Kwynto/GracefulDB/internal/connectors/rest"
	"github.com/Kwynto/GracefulDB/internal/connectors/socketconnector"
	"github.com/Kwynto/GracefulDB/internal/manage/webmanage"
	"github.com/Kwynto/GracefulDB/pkg/lib/closer"
)

func Run(ctx context.Context, cfg *config.Config) error {
	var closeProcs = &closer.Closer{}

	// TODO: Load the core of the system
	go core.Engine(cfg)
	closeProcs.AddHandler(core.Shutdown) // Register a shutdown handler.

	// TODO: Start Socket connector
	if cfg.SocketConnector.Enable {
		go socketconnector.Start(cfg)
		closeProcs.AddHandler(socketconnector.Shutdown) // Register a shutdown handler.
	}

	// TODO: Start REST API connector
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

	// Waiting for a stop signal from the OS
	<-ctx.Done()
	slog.Warn("The shutdown process has started.")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), cfg.ShutdownTimeOut)
	defer cancel()

	if err := closeProcs.Close(shutdownCtx); err != nil {
		return fmt.Errorf("server shutdown: %v", err)
	}

	slog.Info("All processes are stopped.")

	return nil
}
