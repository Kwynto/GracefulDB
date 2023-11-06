package vqlanalyzer

import (
	"context"

	"github.com/Kwynto/GracefulDB/internal/config"
	"github.com/Kwynto/GracefulDB/pkg/lib/closer"
)

// TODO: Request
func Request(instruction string) string {
	return "There should be a response from the processed request."
}

func Analyzer(cfg *config.Config) {
	// -
}

func Shutdown(ctx context.Context, c *closer.Closer) {
	c.Done()
}
