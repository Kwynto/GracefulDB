package loghelper

import (
	"io"
	"os"

	"log/slog"

	"github.com/Kwynto/GracefulDB/internal/config"
)

var iomw io.Writer = io.MultiWriter(os.Stdout, os.Stderr)

func SetupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case config.EnvDev:
		log = slog.New(
			slog.NewTextHandler(iomw, &slog.HandlerOptions{
				AddSource: true,
				Level:     slog.LevelDebug,
			}),
		)
	case config.EnvProd:
		log = slog.New(
			slog.NewJSONHandler(iomw, &slog.HandlerOptions{
				AddSource: false,
				Level:     slog.LevelInfo,
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
