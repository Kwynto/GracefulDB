package loghelper

import (
	"fmt"
	"io"
	"log"
	"os"

	"log/slog"

	"github.com/Kwynto/GracefulDB/internal/config"
)

var IoMultiWriter io.Writer

func OpenLogFile(name string) (io.Writer, error) {
	fo, err := os.OpenFile(name, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	if err != nil {
		return nil, err
	}

	return fo, nil
}

func SetupLogger(cfg *config.Config) *slog.Logger {

	var nlog *slog.Logger

	IoFile, err := OpenLogFile(fmt.Sprintf("%s%s%s", cfg.LogPath, cfg.Env, ".log"))
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	IoMultiWriter = io.MultiWriter(os.Stdout, IoFile)

	switch cfg.Env {
	case config.EnvDev:
		nlog = slog.New(
			slog.NewTextHandler(IoMultiWriter, &slog.HandlerOptions{
				AddSource: true,
				Level:     slog.LevelDebug,
			}),
		)
	case config.EnvProd:
		nlog = slog.New(
			slog.NewJSONHandler(IoMultiWriter, &slog.HandlerOptions{
				AddSource: false,
				Level:     slog.LevelInfo,
			}),
		)
	default:
		nlog = slog.New(
			slog.NewJSONHandler(IoMultiWriter, &slog.HandlerOptions{
				AddSource: false,
				Level:     slog.LevelInfo,
			}),
		)
	}

	return nlog
}

func Err(err error) slog.Attr {
	return slog.Attr{
		Key:   "error",
		Value: slog.StringValue(err.Error()),
	}
}
