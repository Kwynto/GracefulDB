package ordinarylogger

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/Kwynto/GracefulDB/pkg/lib/incolor"
)

const (
	SEnvDev  = "dev"
	SEnvProd = "prod"
)

type TStOrdinaryHandler struct {
	slog.Handler
	lScreen *log.Logger
	lFile   *log.Logger
	env     string
}

type TRecQueue struct {
	handler    *TStOrdinaryHandler
	sMsgScreen string
	sMsgFile   string
}

var IoFile *os.File // io.Writer
var LogHandler slog.Handler
var LogServerError *log.Logger

var chRecQueue = make(chan TRecQueue, 1024)

func (h *TStOrdinaryHandler) Handle(ctx context.Context, r slog.Record) error {
	var sFileOut string

	sLevel := r.Level.String() + ":"

	switch r.Level {
	case slog.LevelDebug:
		sLevel = incolor.StringMagenta(sLevel)
	case slog.LevelInfo:
		sLevel = incolor.StringGreen(sLevel)
	case slog.LevelWarn:
		sLevel = incolor.StringYellow(sLevel)
	case slog.LevelError:
		sLevel = incolor.StringRed(sLevel)
	}

	mFields := make(map[string]interface{}, r.NumAttrs())
	r.Attrs(func(a slog.Attr) bool {
		mFields[a.Key] = a.Value.Any()

		return true
	})

	bSlAttrs, _ := json.MarshalIndent(mFields, "", "  ")

	sAttrsScreenOut := string(bSlAttrs)
	if sAttrsScreenOut == "{}" {
		sAttrsScreenOut = ""
	}

	sAttrsFileOut := ""
	switch h.env {
	case SEnvDev:
		for k, v := range mFields {
			sAttrsFileOut = fmt.Sprintf("%s %s=\"%s\"", sAttrsFileOut, k, v.(string))
		}
	case SEnvProd:
		for k, v := range mFields {
			sAttrsFileOut = fmt.Sprintf("%s, \"%s\":\"%s\"", sAttrsFileOut, k, v.(string))
		}
	}

	sTimeScreen := r.Time.Format("[2006-01-02 15:04:05.000 -0700]")

	switch h.env {
	case SEnvDev:
		sTimeStrFile := r.Time.Format("2006-01-02 15:04:05.000 -0700")
		sFileOut = fmt.Sprintf("time=%s level=%v msg=\"%s\"%s", sTimeStrFile, r.Level, r.Message, sAttrsFileOut)
	case SEnvProd:
		sTimeStrFile := r.Time.Format("2006-01-02T15:04:05.000000000-0700")
		sFileOut = fmt.Sprintf("{\"time\":\"%s\",\"level\":\"%v\",\"msg\":\"%s\"%s}", sTimeStrFile, r.Level, r.Message, sAttrsFileOut)
	}

	sMsg := incolor.StringCyan(r.Message)

	stRec := TRecQueue{
		handler:    h,
		sMsgScreen: fmt.Sprint(sTimeScreen, sLevel, sMsg, incolor.StringWhite(sAttrsScreenOut)),
		sMsgFile:   fmt.Sprint(sFileOut),
	}

	chRecQueue <- stRec

	return nil
}

func newOrdinaryHandler(outScreen io.Writer, outFile io.Writer, env string) *TStOrdinaryHandler {
	var level slog.Level
	switch env {
	case SEnvDev:
		level = slog.LevelDebug
	case SEnvProd:
		level = slog.LevelInfo
	}

	h := &TStOrdinaryHandler{
		Handler: slog.NewJSONHandler(outScreen, &slog.HandlerOptions{
			Level: level,
		}),
		lScreen: log.New(outScreen, "", 0),
		lFile:   log.New(outFile, "", 0),
		env:     env,
	}

	return h
}

func openLogFile(name string) *os.File {
	f, _ := os.OpenFile(name, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	return f
}

func setupLogger(logPath, logEnv string) *slog.Logger {
	var newLog *slog.Logger

	ioFile := openLogFile(filepath.Join(logPath, fmt.Sprintf("%s%s", logEnv, ".log")))

	LogHandler = newOrdinaryHandler(os.Stdout, ioFile, logEnv)
	newLog = slog.New(LogHandler)

	return newLog
}

func recordingQueue() {
	for {
		stRec := <-chRecQueue
		stRec.handler.lScreen.Println(stRec.sMsgScreen)
		stRec.handler.lFile.Println(stRec.sMsgFile)
	}
}

func Init(logPath, logEnv string) {
	newLog := setupLogger(logPath, logEnv)
	slog.SetDefault(newLog)
	LogServerError = slog.NewLogLogger(LogHandler, slog.LevelError)
	go recordingQueue()
}
