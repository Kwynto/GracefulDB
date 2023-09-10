package loghelper

import (
	"fmt"
	"io"
	"log"
	"os"

	"log/slog"

	"github.com/Kwynto/GracefulDB/internal/config"
	"github.com/Kwynto/GracefulDB/internal/lib/helpers/fileshelper"
)

func OpenLogFile(name string) (io.Writer, error) {
	if !fileshelper.FileExists(name) {
		if err := fileshelper.CreateFile(name); err != nil {
			return nil, err
		}
	}

	fo, err := os.Open(name)
	if err != nil {
		return nil, err
	}

	return fo, nil
}

func SetupLogger(cfg *config.Config) *slog.Logger {
	var (
		iof  io.Writer
		iomw io.Writer
		nlog *slog.Logger
	)

	iof, err := OpenLogFile(fmt.Sprintf("%s%s%s", cfg.LogPath, cfg.Env, ".log"))
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	iomw = io.MultiWriter(os.Stdout, iof)

	switch cfg.Env {
	case config.EnvDev:
		nlog = slog.New(
			slog.NewTextHandler(iomw, &slog.HandlerOptions{
				AddSource: true,
				Level:     slog.LevelDebug,
			}),
		)
	case config.EnvProd:
		nlog = slog.New(
			slog.NewJSONHandler(iomw, &slog.HandlerOptions{
				AddSource: false,
				Level:     slog.LevelInfo,
			}),
		)
	default:
		nlog = slog.New(
			slog.NewJSONHandler(iomw, &slog.HandlerOptions{
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
