package vqlanalyzer

import (
	"encoding/json"
	"log/slog"

	"github.com/Kwynto/GracefulDB/internal/base/basicsystem"
	"github.com/Kwynto/GracefulDB/internal/gtypes"
)

// TODO: Request
func Request(instruction string) string {

	var qry *gtypes.VQuery

	if !json.Valid([]byte(instruction)) {
		slog.Debug("No valid query", slog.String("instruction", instruction))
		// ERROR 10 - No valid query
		return `{"action":"response","error":10}`
	}

	// FIXME: Unmarshsl только для тестов, для оптимизации нужно переделать на NewDecoder.Decode
	if err := json.Unmarshal([]byte(instruction), &qry); err != nil {
		slog.Debug("Erroneous request", slog.String("err", err.Error()))
		// ERROR 11 - Incorrect request structure
		return `{"action":"response","error":11}`
	}

	bAnswer, err := json.Marshal(basicsystem.Processing(qry))
	if err != nil {
		// ERROR 20 - Server error
		return `{"action":"response","error":20}`
	}

	return string(bAnswer)
}

// func Analyzer(cfg *config.Config) {
// 	// -
// }

// func Shutdown(ctx context.Context, c *closer.Closer) {
// 	c.Done()
// }
