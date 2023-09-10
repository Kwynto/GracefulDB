package loghelper

import (
	"os"

	"log/slog"

	"github.com/Kwynto/GracefulDB/internal/config"
)

func SetupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case config.EnvDev:
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
				Level: slog.LevelDebug,
			}),
		)
	case config.EnvProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
				Level: slog.LevelInfo,
			}),
		)
	}
	return log
}

func Err(err error) slog.Attr {
	return slog.Attr{
		Key:   "error",
		Value: slog.StringValue(err.Error()),
	}
}
