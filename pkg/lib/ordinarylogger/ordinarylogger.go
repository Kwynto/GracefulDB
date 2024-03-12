package ordinarylogger

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"log/slog"
	"os"

	"github.com/fatih/color"
)

const (
	EnvDev  = "dev"
	EnvProd = "prod"
)

type OrdinaryHandler struct {
	slog.Handler
	lScreen *log.Logger
	lFile   *log.Logger
	env     string
}

var IoFile *os.File // io.Writer
var LogHandler slog.Handler
var LogServerError *log.Logger

func (h *OrdinaryHandler) Handle(ctx context.Context, r slog.Record) error {
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

	byteAttrs, _ := json.MarshalIndent(fields, "", "  ")
	// if err != nil {
	// 	return err
	// }

	strAttrsScreenOut := string(byteAttrs)
	if strAttrsScreenOut == "{}" {
		strAttrsScreenOut = ""
	}

	strAttrsFileOut := ""
	switch h.env {
	case EnvDev:
		for k, v := range fields {
			strAttrsFileOut = fmt.Sprintf("%s %s=\"%s\"", strAttrsFileOut, k, v.(string))
		}
	case EnvProd:
		for k, v := range fields {
			strAttrsFileOut = fmt.Sprintf("%s, \"%s\":\"%s\"", strAttrsFileOut, k, v.(string))
		}
	}

	timeStrScreen := r.Time.Format("[2006-01-02 15:04:05.000 -0700]")

	switch h.env {
	case EnvDev:
		timeStrFile := r.Time.Format("2006-01-02 15:04:05.000 -0700")
		strFileOut = fmt.Sprintf("time=%s level=%v msg=\"%s\"%s", timeStrFile, r.Level, r.Message, strAttrsFileOut)
	case EnvProd:
		timeStrFile := r.Time.Format("2006-01-02T15:04:05.000000000-0700")
		strFileOut = fmt.Sprintf("{\"time\":\"%s\",\"level\":\"%v\",\"msg\":\"%s\"%s}", timeStrFile, r.Level, r.Message, strAttrsFileOut)
	}

	msg := color.CyanString(r.Message)

	h.lScreen.Println(timeStrScreen, level, msg, color.WhiteString(strAttrsScreenOut))
	h.lFile.Println(strFileOut)

	return nil
}

func newOrdinaryHandler(outScreen io.Writer, outFile io.Writer, env string) *OrdinaryHandler {
	var gbdlevel slog.Level
	switch env {
	case EnvDev:
		gbdlevel = slog.LevelDebug
	case EnvProd:
		gbdlevel = slog.LevelInfo
	}

	h := &OrdinaryHandler{
		Handler: slog.NewJSONHandler(outScreen, &slog.HandlerOptions{
			Level: gbdlevel,
		}),
		lScreen: log.New(outScreen, "", 0),
		lFile:   log.New(outFile, "", 0),
		env:     env,
	}

	return h
}

func openLogFile(name string) *os.File {
	fo, _ := os.OpenFile(name, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)

	return fo
}

func setupLogger(logPath, logEnv string) *slog.Logger {
	var nlog *slog.Logger

	IoFile := openLogFile(fmt.Sprintf("%s%s%s", logPath, logEnv, ".log"))

	LogHandler = newOrdinaryHandler(os.Stdout, IoFile, logEnv)
	nlog = slog.New(LogHandler)

	return nlog
}

func Err(err error) slog.Attr {
	return slog.Attr{
		Key:   "error",
		Value: slog.StringValue(err.Error()),
	}
}

func Init(logPath, logEnv string) {
	inlog := setupLogger(logPath, logEnv)
	slog.SetDefault(inlog)
	LogServerError = slog.NewLogLogger(LogHandler, slog.LevelError)
}
