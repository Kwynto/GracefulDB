package basicsystem

import (
	"log/slog"
	"sync"

	"github.com/Kwynto/GracefulDB/internal/config"
)

func CommandSystem(cfg *config.Config, wg *sync.WaitGroup) {
	slog.Info("GracefulDB: The basic command system was started.")
	wg.Done()
}
