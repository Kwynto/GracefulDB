package basicsystem

import (
	"log/slog"
	"sync"

	"github.com/Kwynto/GracefulDB/internal/config"
)

func CommandSystem(cfg *config.Config, log *slog.Logger, wg *sync.WaitGroup) {
	wg.Done()
}
