package core

import (
	"context"
	"log/slog"

	"github.com/Kwynto/GracefulDB/internal/config"
	"github.com/Kwynto/GracefulDB/pkg/lib/closer"
)

func Engine(cfg *config.Config) {
	slog.Info("GracefulDB: The core of the system was started.")
}

func Shutdown(ctx context.Context, c *closer.Closer) {
	// -
	// for i := 0; i < 10; i++ {
	// 	fmt.Print(".")
	// 	time.Sleep(1 * time.Second)
	// }
	// c.AddMsg("Imitation of an error")
	c.Done()
}
