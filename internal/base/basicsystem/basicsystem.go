package basicsystem

import (
	"context"
	"log/slog"

	"github.com/Kwynto/GracefulDB/internal/config"
)

func CommandSystem(cfg *config.Config) {
	slog.Info("GracefulDB: The basic command system was started.")
}

func Shutdown(ctx context.Context) error {
	return nil
}
