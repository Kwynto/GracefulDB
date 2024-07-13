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
	"slices"
	"sync/atomic"
	"time"

	"github.com/Kwynto/GracefulDB/pkg/lib/incolor"
)

const (
	S_ENV_DEV  = "develop"
	S_ENV_WORK = "working"
)

type TStOrdinaryHandler struct {
	slog.Handler
	lScreen *log.Logger
	lFile   *log.Logger
	env     string
	uCount  atomic.Uint64
}

type TRecQueue struct {
	handler    *TStOrdinaryHandler
	sMsgScreen string
	sMsgFile   string
	uLine      uint64
}

var IoFile *os.File // io.Writer
var LogHandler slog.Handler
var LogServerError *log.Logger

var chRecQueue = make(chan TRecQueue, 1024)

func (h *TStOrdinaryHandler) Handle(ctx context.Context, r slog.Record) error {
	uLine := h.uCount.Add(1)
	go prepareLog(h, r, uLine)

	return nil
}

func prepareLog(h *TStOrdinaryHandler, r slog.Record, uLine uint64) {
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
	case S_ENV_DEV:
		for k, v := range mFields {
			sAttrsFileOut = fmt.Sprintf("%s %s=\"%s\"", sAttrsFileOut, k, v.(string))
		}
	case S_ENV_WORK:
		for k, v := range mFields {
			sAttrsFileOut = fmt.Sprintf("%s, \"%s\":\"%s\"", sAttrsFileOut, k, v.(string))
		}
	}

	sTimeScreen := r.Time.Format("[2006-01-02 15:04:05.000 -0700]")

	switch h.env {
	case S_ENV_DEV:
		sTimeStrFile := r.Time.Format("2006-01-02 15:04:05.000 -0700")
		sFileOut = fmt.Sprintf("time=%s level=%v msg=\"%s\"%s", sTimeStrFile, r.Level, r.Message, sAttrsFileOut)
	case S_ENV_WORK:
		sTimeStrFile := r.Time.Format("2006-01-02T15:04:05.000000000-0700")
		sFileOut = fmt.Sprintf("{\"time\":\"%s\", \"level\":\"%v\", \"msg\":\"%s\"%s}", sTimeStrFile, r.Level, r.Message, sAttrsFileOut)
	}

	sMsg := incolor.StringCyan(r.Message)

	stRec := TRecQueue{
		handler:    h,
		sMsgScreen: fmt.Sprint(sTimeScreen, " ", sLevel, " ", sMsg, " ", incolor.StringWhite(sAttrsScreenOut)),
		sMsgFile:   fmt.Sprint(sFileOut),
		uLine:      uLine,
	}

	chRecQueue <- stRec
}

func recordingQueue() {
	var arRecInd = make([]uint64, 1024)
	var mRecs = make(map[uint64]TRecQueue, 1024)
	var uInd uint64

	for {
		arRecInd = arRecInd[:0]
		clear(mRecs)
		uInd = 0

		iCount := len(chRecQueue)
		for i := 0; i < iCount; i++ {
			stRec := <-chRecQueue
			uInd = stRec.uLine
			arRecInd = append(arRecInd, uInd)
			mRecs[uInd] = stRec
		}

		slices.Sort(arRecInd)

		for _, uLine := range arRecInd {
			stRec := mRecs[uLine]
			stRec.handler.lScreen.Println(stRec.sMsgScreen)
			stRec.handler.lFile.Println(stRec.sMsgFile)
		}

		time.Sleep(100 * time.Millisecond)
	}
}

func newOrdinaryHandler(outScreen io.Writer, outFile io.Writer, env string) *TStOrdinaryHandler {
	var level slog.Level
	switch env {
	case S_ENV_DEV:
		level = slog.LevelDebug
	case S_ENV_WORK:
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
	h.uCount.Store(0)

	return h
}

func Init(logPath, logEnv string) {
	sNameFile := filepath.Join(logPath, fmt.Sprintf("%s%s", logEnv, ".log"))
	ioFile, _ := os.OpenFile(sNameFile, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)

	LogHandler = newOrdinaryHandler(os.Stdout, ioFile, logEnv)
	newLog := slog.New(LogHandler)

	slog.SetDefault(newLog)
	LogServerError = slog.NewLogLogger(LogHandler, slog.LevelError)
	go recordingQueue()
}
