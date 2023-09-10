package main

import (
	"os"
	"sync"

	"log/slog"

	"github.com/Kwynto/GracefulDB/internal/config"
	"github.com/Kwynto/GracefulDB/internal/lib/helpers/loghelper"
)

var wg sync.WaitGroup

func main() {
	// Init variables
	wg = sync.WaitGroup{}

	// Init config: cleanenv
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = "./config/default.yaml"
	}
	cfg := config.MustLoad(configPath)

	// Init logger: slog
	log := loghelper.SetupLogger(cfg)
	log.Info("starting GracefulDB", slog.String("env", cfg.Env))
	log.Debug("debug messages are enabled")

	// TODO: Load the core of the system
	wg.Add(1)

	// TODO: Run the basic command system
	wg.Add(1)

	// TODO: Start the language analyzer (SQL)
	wg.Add(1)

	// TODO: Start Socket connector
	wg.Add(1)

	// TODO: Start REST API connector
	wg.Add(1)

	// TODO: Start gRPC connector
	wg.Add(1)

	// TODO: Start web-server for manage system
	wg.Add(1)

	// TODO:: Signal tracking

	// Wait for all processes to complete
	wg.Wait()
	log.Info("GracefulDB has finished its work and will miss you.")
}
