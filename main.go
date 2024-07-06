package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	_ "embed"

	"github.com/joho/godotenv"

	"github.com/Kwynto/GracefulDB/internal/config"
	"github.com/Kwynto/GracefulDB/internal/server"

	_ "github.com/Kwynto/GracefulDB/assets"

	"github.com/Kwynto/GracefulDB/pkg/lib/incolor"
	"github.com/Kwynto/GracefulDB/pkg/lib/ordinarylogger"
)

var (
	//go:embed LICENSE
	sLicense string
)

func main() {
	// Greeting
	fmt.Println(incolor.StringYellowH(sLicense))

	// Init config
	errDotEnv := godotenv.Load()
	sConfigPath := os.Getenv("CONFIG_PATH")
	config.MustLoad(sConfigPath)

	// if config.StDefaultConfig.Env == "test" {
	// 	fmt.Println("You should set up the configuration file correctly.")
	// 	os.Exit(0)
	// }

	// Init logger: slog
	ordinarylogger.Init(config.StDefaultConfig.LogPath, config.StDefaultConfig.Env)
	slog.Info("Starting GracefulDB", slog.String("env", config.StDefaultConfig.Env))
	slog.Info("Configuration loaded", slog.String("file", config.SDisplayConfigPath))
	if errDotEnv == nil {
		slog.Info("The environment variables were read from the env-file. Don't forget, you can use OS environment variables, they take precedence over env-files.")
	}

	if config.StDefaultConfig.Env == config.ENV_DEV {
		slog.Info("Developer mode is active.")
		slog.Warn("Perhaps you should set up the configuration file correctly.")
	}
	slog.Debug("debug messages are enabled")

	// Signal tracking
	ctxSignal, fnStopSignal := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer fnStopSignal()

	if err := server.Run(ctxSignal, &config.StDefaultConfig); err != nil {
		slog.Error("An unexpected error occurred while the server was running.", slog.String("err", err.Error()))
	}

	slog.Info("GracefulDB has finished its work and will miss you.")
}
