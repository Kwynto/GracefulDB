package main

import (
	"os"
	"sync"

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
	go core.Engine(cfg, log, &wg)

	// TODO: Run the basic command system
	wg.Add(1)
	go basicsystem.CommandSystem(cfg, log, &wg)

	// TODO: Start the language analyzer (SQL)
	wg.Add(1)
	go sqlanalyzer.Analyzer(cfg, log, &wg)

	// TODO: Start Socket connector
	wg.Add(1)
	go socketconnector.Start(cfg, log, &wg)

	// TODO: Start REST API connector
	wg.Add(1)
	go rest.Start(cfg, log, &wg)

	// TODO: Start gRPC connector
	wg.Add(1)
	go grpc.Start(cfg, log, &wg)

	// TODO: Start web-server for manage system
	wg.Add(1)
	go webmanage.Start(cfg, log, &wg)

	// TODO:: Signal tracking

	// Wait for all processes to complete
	wg.Wait()
	log.Info("GracefulDB has finished its work and will miss you.")
}
