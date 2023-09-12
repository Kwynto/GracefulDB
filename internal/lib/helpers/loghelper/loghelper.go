package loghelper

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"log/slog"
	"os"

	"github.com/fatih/color"

	"github.com/Kwynto/GracefulDB/internal/config"
)

type GDBPrettyHandler struct {
	slog.Handler
	lScreen *log.Logger
	lFile   *log.Logger
	env     string
}

var IoFile io.Writer
var LogHandler slog.Handler
var LogServerError *log.Logger

func (h *GDBPrettyHandler) Handle(ctx context.Context, r slog.Record) error {
	var strFileOut string

	level := r.Level.String() + ":"

	switch r.Level {
	case slog.LevelDebug:
		level = color.MagentaString(level)
	case slog.LevelInfo:
		level = color.GreenString(level)
	case slog.LevelWarn:
		level = color.YellowString(level)
	case slog.LevelError:
		level = color.RedString(level)
	}

	fields := make(map[string]interface{}, r.NumAttrs())
	r.Attrs(func(a slog.Attr) bool {
		fields[a.Key] = a.Value.Any()

		return true
	})

	b, err := json.MarshalIndent(fields, "", "  ")
	if err != nil {
		return err
	}

	timeStrScreen := r.Time.Format("[2006-01-02 15:04:05.000 -0700]")

	switch h.env {
	case config.EnvDev:
		timeStrFile := r.Time.Format("2006-01-02 15:04:05.000 -0700")
		strFileOut = fmt.Sprintf("time=%s level=%v msg=\"%s\" %s", timeStrFile, r.Level, r.Message, string(b))
	case config.EnvProd:
		timeStrFile := r.Time.Format("2006-01-02T15:04:05.000000000-0700")
		strFileOut = fmt.Sprintf("{\"time\":\"%s\",\"level\":\"%v\",\"msg\":\"%s\", \"attributes\":%v}", timeStrFile, r.Level, r.Message, string(b))
	}

	msg := color.CyanString(r.Message)

	h.lScreen.Println(timeStrScreen, level, msg, color.WhiteString(string(b)))
	h.lFile.Println(strFileOut)

	return nil
}

func NewGDBPrettyHandler(outScreen io.Writer, outFile io.Writer, env string) *GDBPrettyHandler {
	var gbdlevel slog.Level
	switch env {
	case config.EnvDev:
		gbdlevel = slog.LevelDebug
	case config.EnvProd:
		gbdlevel = slog.LevelInfo
	}

	h := &GDBPrettyHandler{
		Handler: slog.NewJSONHandler(outScreen, &slog.HandlerOptions{
			Level: gbdlevel,
		}),
		lScreen: log.New(outScreen, "", 0),
		lFile:   log.New(outFile, "", 0),
		env:     env,
	}

	return h
}

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

	LogHandler = NewGDBPrettyHandler(os.Stdout, IoFile, cfg.Env)
	nlog = slog.New(LogHandler)

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
