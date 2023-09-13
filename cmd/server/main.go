package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"log/slog"

	"github.com/Kwynto/GracefulDB/internal/analyzers/sqlanalyzer"
	"github.com/Kwynto/GracefulDB/internal/base/basicsystem"
	"github.com/Kwynto/GracefulDB/internal/base/core"
	"github.com/Kwynto/GracefulDB/internal/config"
	"github.com/Kwynto/GracefulDB/internal/connectors/grpc"
	"github.com/Kwynto/GracefulDB/internal/connectors/rest"
	"github.com/Kwynto/GracefulDB/internal/connectors/socketconnector"
	"github.com/Kwynto/GracefulDB/internal/lib/helpers/loghelper"
	"github.com/Kwynto/GracefulDB/internal/manage/webmanage"
)

func runServer(ctx context.Context, cfg *config.Config) error {
	// TODO: Load the core of the system
	go core.Engine(cfg)

	// TODO: Run the basic command system
	go basicsystem.CommandSystem(cfg)

	// TODO: Start the language analyzer (SQL)
	go sqlanalyzer.Analyzer(cfg)

	// TODO: Start Socket connector
	go socketconnector.Start(cfg)

	// TODO: Start REST API connector
	go rest.Start(cfg)

	// TODO: Start gRPC connector
	go grpc.Start(cfg)

	// TODO: Start web-server for manage system
	go webmanage.Start(cfg)

	// We are waiting for a stop signal from the OS
	<-ctx.Done()
	slog.Warn("The shutdown process has started.")

	processShutdown := make(chan struct{}, 1)

	shutdownCtx, cancel := context.WithTimeout(context.Background(), cfg.ShutdownTimeOut)
	defer cancel()

	// Stopping all processes.
	go func() {
		// TODO: Alternate stopping of all processes
		time.Sleep(13 * time.Second)

		processShutdown <- struct{}{}
	}()

	select {
	case <-shutdownCtx.Done():
		return fmt.Errorf("server shutdown: %w", ctx.Err())
	case <-processShutdown:
		slog.Info("All processes are stopped.")
	}

	return nil
}

func main() {
	// Init config
	configPath := os.Getenv("CONFIG_PATH")
	cfg := config.MustLoad(configPath)

	// Init logger: slog
	loghelper.Init(cfg)
	slog.Info("starting GracefulDB", slog.String("env", cfg.Env))
	slog.Debug("debug messages are enabled")

	// Signal tracking
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	if err := runServer(ctx, cfg); err != nil {
		slog.Error("An unexpected error occurred while the server was running.", slog.String("err", err.Error()))
	}

	slog.Info("GracefulDB has finished its work and will miss you.")
}
