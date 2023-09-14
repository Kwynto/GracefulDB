package basicsystem

import (
	"context"
	"log/slog"

	"github.com/Kwynto/GracefulDB/internal/config"
	"github.com/Kwynto/GracefulDB/pkg/lib/closer"
)

func CommandSystem(cfg *config.Config) {
	slog.Info("GracefulDB: The basic command system was started.")
}

func Shutdown(ctx context.Context, c *closer.Closer) {
	c.Done()
}
