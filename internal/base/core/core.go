package core

import (
	"context"
	"log/slog"

	"github.com/Kwynto/GracefulDB/internal/config"
)

func Engine(cfg *config.Config) {
	slog.Info("GracefulDB: The core of the system was started.")
}

func Shutdown(ctx context.Context) error {
	// -
	// for i := 0; i < 10; i++ {
	// 	fmt.Print(".")
	// 	time.Sleep(1 * time.Second)
	// }
	return nil
}
