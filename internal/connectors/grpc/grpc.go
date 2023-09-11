package grpc

import (
	"sync"

	"github.com/Kwynto/GracefulDB/internal/config"
)

func Start(cfg *config.Config, wg *sync.WaitGroup) {
	wg.Done()
}
