package sqlanalyzer

import (
	"log/slog"
	"sync"

	"github.com/Kwynto/GracefulDB/internal/config"
)

func Analyzer(cfg *config.Config, log *slog.Logger, wg *sync.WaitGroup) {
	wg.Done()
}
