package ordinarylogger

import (
	"context"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"os"
	"reflect"
	"testing"
	"time"
)

const (
	TEST_LOG_PATH = ""
	TEST_LOG_FILE = "test.log"
)

func Test_Handle(t *testing.T) {
	inlog := setupLogger(TEST_LOG_PATH, "dev")

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

	inlog = setupLogger(TEST_LOG_PATH, "prod")

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
	iof := openLogFile(fmt.Sprintf("%s%s%s", TEST_LOG_PATH, "dev", ".log"))

	logHandler1 := newOrdinaryHandler(os.Stdout, iof, "dev")

	t.Run("newOrdinaryHandler() function testing", func(t *testing.T) {
		if r1 := reflect.TypeOf(logHandler1); r1 != reflect.TypeOf(&TStOrdinaryHandler{}) {
			t.Errorf("newOrdinaryHandler() error: type mismatch: %v", r1)
		}
	})

	logHandler2 := newOrdinaryHandler(os.Stdout, iof, "prod")

	t.Run("newOrdinaryHandler() function testing", func(t *testing.T) {
		if r1 := reflect.TypeOf(logHandler2); r1 != reflect.TypeOf(&TStOrdinaryHandler{}) {
			t.Errorf("newOrdinaryHandler() error: type mismatch: %v", r1)
		}
	})
}

func Test_openLogFile(t *testing.T) {
	t.Run("openLogFile() function testing", func(t *testing.T) {
		iof := openLogFile(fmt.Sprintf("%s%s%s", TEST_LOG_PATH, "dev", ".log"))

		if r1 := reflect.TypeOf(iof); r1 != reflect.TypeOf(&os.File{}) {
			t.Errorf("openLogFile() error: type mismatch: %v", r1)
		}
	})
}

func Test_setupLogger(t *testing.T) {
	t.Run("setupLogger() function testing", func(t *testing.T) {
		inlog := setupLogger(TEST_LOG_PATH, "dev")

		if r1 := reflect.TypeOf(inlog); r1 != reflect.TypeOf(&slog.Logger{}) {
			t.Errorf("setupLogger() error: type mismatch: %v", r1)
		}

		if r1 := reflect.TypeOf(inlog.Handler()); r1 != reflect.TypeOf(&TStOrdinaryHandler{}) {
			t.Errorf("setupLogger() error: type of handler mismatch: %v", r1)
		}
	})
}

func Test_Err(t *testing.T) {
	t.Run("Err() function testing", func(t *testing.T) {
		fakeErr := errors.New("fake error")
		res := Err(fakeErr)

		if r1 := reflect.TypeOf(res); r1 != reflect.TypeOf(slog.Attr{}) {
			t.Errorf("Err() error: type mismatch: %v", r1)
		}

		if (res.Key != "error") || !res.Value.Equal(slog.StringValue(fakeErr.Error())) {
			t.Errorf("Err() error: broken attribute")
		}
	})
}

func Test_Init(t *testing.T) {
	t.Run("Init() function testing", func(t *testing.T) {
		Init(TEST_LOG_PATH, "dev")

		if r1 := reflect.TypeOf(LogServerError); r1 != reflect.TypeOf(&log.Logger{}) {
			t.Errorf("Init() error: type mismatch: %v", r1)
		}

		if r1 := reflect.TypeOf(slog.Default().Handler()); r1 != reflect.TypeOf(&TStOrdinaryHandler{}) {
			t.Errorf("Init() error: type of handler mismatch in slog.Default: %v", r1)
		}
	})
}
