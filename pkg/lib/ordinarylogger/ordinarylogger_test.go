package ordinarylogger

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"os"
	"path/filepath"
	"reflect"
	"testing"
	"time"
)

const (
	TEST_LOG_PATH = ""
	TEST_LOG_FILE = "test.log"
)

func Test_Handle(t *testing.T) {
	sNameFile := filepath.Join(TEST_LOG_PATH, fmt.Sprintf("%s%s", "develop", ".log"))
	ioFile, _ := os.OpenFile(sNameFile, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)

	LogHandler = newOrdinaryHandler(os.Stdout, ioFile, "develop")
	inlog := slog.New(LogHandler)

	t.Run("Handle() function testing #1", func(t *testing.T) {
		rec := slog.Record{
			Time:    time.Now(),
			Message: "test msg",
			Level:   slog.LevelDebug,
		}
		err := inlog.Handler().Handle(context.Background(), rec)
		if err != nil {
			t.Errorf("Handle() error: type mismatch: %v", err)
		}
	})

	t.Run("Handle() function testing #2", func(t *testing.T) {
		rec := slog.Record{
			Time:    time.Now(),
			Message: "test msg",
			Level:   slog.LevelInfo,
		}
		err := inlog.Handler().Handle(context.Background(), rec)
		if err != nil {
			t.Errorf("Handle() error: type mismatch: %v", err)
		}
	})

	t.Run("Handle() function testing #3", func(t *testing.T) {
		rec := slog.Record{
			Time:    time.Now(),
			Message: "test msg",
			Level:   slog.LevelWarn,
		}
		err := inlog.Handler().Handle(context.Background(), rec)
		if err != nil {
			t.Errorf("Handle() error: type mismatch: %v", err)
		}
	})

	t.Run("Handle() function testing #4", func(t *testing.T) {
		rec := slog.Record{
			Time:    time.Now(),
			Message: "test msg",
			Level:   slog.LevelError,
		}
		err := inlog.Handler().Handle(context.Background(), rec)
		if err != nil {
			t.Errorf("Handle() error: type mismatch: %v", err)
		}
	})

	t.Run("Handle() function testing #5", func(t *testing.T) {
		rec := slog.Record{
			Time:    time.Now(),
			Message: "test msg",
			Level:   slog.LevelError,
		}
		rec.AddAttrs(slog.Attr{
			Key:   "test",
			Value: slog.StringValue("test error"),
		})
		err := inlog.Handler().Handle(context.Background(), rec)
		if err != nil {
			t.Errorf("Handle() error: type mismatch: %v", err)
		}
	})

	sNameFile2 := filepath.Join(TEST_LOG_PATH, fmt.Sprintf("%s%s", "working", ".log"))
	ioFile2, _ := os.OpenFile(sNameFile2, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)

	LogHandler = newOrdinaryHandler(os.Stdout, ioFile2, "working")
	inlog = slog.New(LogHandler)

	t.Run("Handle() function testing #6", func(t *testing.T) {
		rec := slog.Record{
			Time:    time.Now(),
			Message: "test msg",
			Level:   slog.LevelDebug,
		}
		err := inlog.Handler().Handle(context.Background(), rec)
		if err != nil {
			t.Errorf("Handle() error: type mismatch: %v", err)
		}
	})

	t.Run("Handle() function testing #7", func(t *testing.T) {
		rec := slog.Record{
			Time:    time.Now(),
			Message: "test msg",
			Level:   slog.LevelInfo,
		}
		err := inlog.Handler().Handle(context.Background(), rec)
		if err != nil {
			t.Errorf("Handle() error: type mismatch: %v", err)
		}
	})

	t.Run("Handle() function testing #8", func(t *testing.T) {
		rec := slog.Record{
			Time:    time.Now(),
			Message: "test msg",
			Level:   slog.LevelWarn,
		}
		err := inlog.Handler().Handle(context.Background(), rec)
		if err != nil {
			t.Errorf("Handle() error: type mismatch: %v", err)
		}
	})

	t.Run("Handle() function testing #9", func(t *testing.T) {
		rec := slog.Record{
			Time:    time.Now(),
			Message: "test msg",
			Level:   slog.LevelError,
		}
		err := inlog.Handler().Handle(context.Background(), rec)
		if err != nil {
			t.Errorf("Handle() error: type mismatch: %v", err)
		}
	})

	t.Run("Handle() function testing #10", func(t *testing.T) {
		rec := slog.Record{
			Time:    time.Now(),
			Message: "test msg",
			Level:   slog.LevelError,
		}
		rec.AddAttrs(slog.Attr{
			Key:   "test",
			Value: slog.StringValue("test error"),
		})

		err := inlog.Handler().Handle(context.Background(), rec)
		if err != nil {
			t.Errorf("Handle() error: type mismatch: %v", err)
		}
	})
}

func Test_newOrdinaryHandler(t *testing.T) {
	sNameFile := fmt.Sprintf("%s%s%s", TEST_LOG_PATH, "develop", ".log")
	iof, _ := os.OpenFile(sNameFile, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)

	logHandler1 := newOrdinaryHandler(os.Stdout, iof, "develop")

	t.Run("newOrdinaryHandler() function testing", func(t *testing.T) {
		if r1 := reflect.TypeOf(logHandler1); r1 != reflect.TypeOf(&TStOrdinaryHandler{}) {
			t.Errorf("newOrdinaryHandler() error: type mismatch: %v", r1)
		}
	})

	logHandler2 := newOrdinaryHandler(os.Stdout, iof, "working")

	t.Run("newOrdinaryHandler() function testing", func(t *testing.T) {
		if r1 := reflect.TypeOf(logHandler2); r1 != reflect.TypeOf(&TStOrdinaryHandler{}) {
			t.Errorf("newOrdinaryHandler() error: type mismatch: %v", r1)
		}
	})
}

func Test_Init(t *testing.T) {
	t.Run("Init() function testing", func(t *testing.T) {
		Init(TEST_LOG_PATH, "develop")

		if r1 := reflect.TypeOf(LogServerError); r1 != reflect.TypeOf(&log.Logger{}) {
			t.Errorf("Init() error: type mismatch: %v", r1)
		}

		if r1 := reflect.TypeOf(slog.Default().Handler()); r1 != reflect.TypeOf(&TStOrdinaryHandler{}) {
			t.Errorf("Init() error: type of handler mismatch in slog.Default: %v", r1)
		}
	})
}
