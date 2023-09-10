package core

import (
	"log/slog"
	"sync"

	"github.com/Kwynto/GracefulDB/internal/config"
)

func Engine(cfg *config.Config, log *slog.Logger, wg *sync.WaitGroup) {
	wg.Done()
}
