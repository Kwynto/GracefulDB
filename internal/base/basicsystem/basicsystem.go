package basicsystem

import (
	"context"
	"log/slog"

	"github.com/Kwynto/GracefulDB/internal/config"
	"github.com/Kwynto/GracefulDB/internal/gtypes"
	"github.com/Kwynto/GracefulDB/pkg/lib/closer"
)

func Processing(in *gtypes.VQuery) *gtypes.VAnswer {
	return &gtypes.VAnswer{
		Action: "response",
		Secret: gtypes.VSecret{},
		Data:   gtypes.VData{},
		Error:  0,
	}
}

func CommandSystem(cfg *config.Config) {
	slog.Info("GracefulDB: The basic command system was started.")
}

func Shutdown(ctx context.Context, c *closer.Closer) {
	c.Done()
}
