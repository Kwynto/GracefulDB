package core

import (
	"log/slog"
	"sync"

	"github.com/Kwynto/GracefulDB/internal/config"
)

func Engine(cfg *config.Config, wg *sync.WaitGroup) {
	slog.Info("GracefulDB: The core of the system was started.")
	wg.Done()
}
