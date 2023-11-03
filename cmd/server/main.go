package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"log/slog"

	"github.com/Kwynto/GracefulDB/internal/config"
	"github.com/Kwynto/GracefulDB/internal/server"

	"github.com/Kwynto/GracefulDB/pkg/lib/prettylogger"
)

func main() {
	// Init config
	configPath := os.Getenv("CONFIG_PATH")
	cfg := config.MustLoad(configPath)

	// Init logger: slog
	prettylogger.Init(cfg.LogPath, cfg.Env)
	slog.Info("starting GracefulDB", slog.String("env", cfg.Env))
	slog.Debug("debug messages are enabled")

	// Signal tracking
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	if err := server.Run(ctx, cfg); err != nil {
		slog.Error("An unexpected error occurred while the server was running.", slog.String("err", err.Error()))
	}

	slog.Info("GracefulDB has finished its work and will miss you.")
}
