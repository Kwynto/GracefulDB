package loghelper

import (
	"fmt"
	"io"
	"log"
	"os"

	"log/slog"

	"github.com/Kwynto/GracefulDB/internal/config"
)

var IoFile io.Writer
var IoMultiWriter io.Writer
var LogHandler slog.Handler
var LogServerError *log.Logger

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
		LogHandler = slog.NewTextHandler(
			IoMultiWriter,
			&slog.HandlerOptions{
				AddSource: true,
				Level:     slog.LevelDebug,
			})
		nlog = slog.New(LogHandler)
	case config.EnvProd:
		LogHandler = slog.NewJSONHandler(
			IoMultiWriter,
			&slog.HandlerOptions{
				AddSource: false,
				Level:     slog.LevelInfo,
			})
		nlog = slog.New(LogHandler)
	default:
		LogHandler = slog.NewJSONHandler(
			IoMultiWriter,
			&slog.HandlerOptions{
				AddSource: false,
				Level:     slog.LevelInfo,
			})
		nlog = slog.New(LogHandler)
	}

	return nlog
}

func Err(err error) slog.Attr {
	return slog.Attr{
		Key:   "error",
		Value: slog.StringValue(err.Error()),
	}
}

func Init(cfg *config.Config) {
	inlog := SetupLogger(cfg)
	slog.SetDefault(inlog)
	LogServerError = slog.NewLogLogger(LogHandler, slog.LevelError)
}
