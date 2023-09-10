package socketconnector

import (
	"log/slog"
	"sync"

	"github.com/Kwynto/GracefulDB/internal/config"
)

func Start(cfg *config.Config, log *slog.Logger, wg *sync.WaitGroup) {
	wg.Done()
}
