package rest

import (
	"context"

	"github.com/Kwynto/GracefulDB/internal/config"
	"github.com/Kwynto/GracefulDB/pkg/lib/closer"
)

func Start(cfg *config.Config) {
	// -
}

func Shutdown(ctx context.Context, c *closer.Closer) {
	c.Done()
}
