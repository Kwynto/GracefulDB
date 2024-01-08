package server

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/Kwynto/GracefulDB/internal/config"
)

func Test_Stop(t *testing.T) {
	t.Run("Stop() function testing", func(t *testing.T) {
		Stop("tester")
		res := <-stopSignal
		if reflect.TypeOf(res) != reflect.TypeOf(struct{}{}) {
			t.Error("Stop() error = wrong result.")
		}
	})
}

func Test_Run(t *testing.T) {
	t.Run("Run() function testing - Signal", func(t *testing.T) {
		stopSignal <- struct{}{}

		ctx, stop := context.WithTimeout(context.Background(), 5*time.Second)
		defer stop()

		time.Sleep(2 * time.Second)

		if err := Run(ctx, &config.DefaultConfig); err == nil {
			t.Error("Run() error = wrong result.")
		}
	})

	t.Run("Run() function testing - Context", func(t *testing.T) {
		ctx, stop := context.WithTimeout(context.Background(), 1*time.Second)
		defer stop()

		if err := Run(ctx, &config.DefaultConfig); err == nil {
			t.Error("Run() error = wrong result.")
		}
	})

	t.Run("Run() function testing - positive", func(t *testing.T) {
		config.MustLoad("../../config/default.yaml")
		ctx, stop := context.WithTimeout(context.Background(), 5*time.Second)
		defer stop()

		if err := Run(ctx, &config.DefaultConfig); err != nil {
			t.Errorf("Run() error = wrong result: %v", err)
		}
	})
}
