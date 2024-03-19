package main

import (
	"context"
	"embed"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"log/slog"

	"github.com/Kwynto/GracefulDB/internal/config"
	"github.com/Kwynto/GracefulDB/internal/manage/webmanage"
	"github.com/Kwynto/GracefulDB/internal/server"

	"github.com/Kwynto/GracefulDB/pkg/lib/colorterm"
	"github.com/Kwynto/GracefulDB/pkg/lib/ordinarylogger"
)

var (
	//go:embed LICENSE
	license string

	//go:embed ui/html
	uiHtmlDir embed.FS

	//go:embed ui/static
	uiStaticDir embed.FS
)

func main() {
	// Greeting
	fmt.Println(colorterm.StringYellowH(license))

	// Set UI
	webmanage.SetUiDirs(&uiHtmlDir, &uiStaticDir)

	// Init config
	configPath := os.Getenv("CONFIG_PATH")
	config.MustLoad(configPath)

	if config.DefaultConfig.Env == "test" {
		fmt.Println("You should set up the configuration file correctly.")
		os.Exit(0)
	}

	// Init logger: slog
	ordinarylogger.Init(config.DefaultConfig.LogPath, config.DefaultConfig.Env)
	slog.Info("Starting GracefulDB", slog.String("env", config.DefaultConfig.Env))
	slog.Info("Configuration loaded", slog.String("file", config.DisplayConfigPath))
	slog.Debug("debug messages are enabled")

	// Signal tracking
	startCtx := context.Background()
	ctx, stop := signal.NotifyContext(startCtx, syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	if err := server.Run(ctx, &config.DefaultConfig); err != nil {
		slog.Error("An unexpected error occurred while the server was running.", slog.String("err", err.Error()))
	}

	slog.Info("GracefulDB has finished its work and will miss you.")
}
