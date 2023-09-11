package sqlanalyzer

import (
	"sync"

	"github.com/Kwynto/GracefulDB/internal/config"
)

func Analyzer(cfg *config.Config, wg *sync.WaitGroup) {
	wg.Done()
}
