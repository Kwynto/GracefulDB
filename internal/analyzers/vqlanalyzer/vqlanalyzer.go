package vqlanalyzer

import (
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/Kwynto/GracefulDB/internal/gtypes"
)

// TODO: Processing
func Processing(in *gtypes.VQuery) *gtypes.VAnswer {
	return &gtypes.VAnswer{
		Action: "response",
		Secret: gtypes.VSecret{},
		Data:   gtypes.VData{},
		Error:  0,
	}
}

func Request(instruction string) string {
	var qry *gtypes.VQuery

	if !json.Valid([]byte(instruction)) {
		slog.Debug("No valid query", slog.String("instruction", instruction))
		// ERROR 10 - Invalid request
		return `{"action":"response","error":10,"description":"Invalid request"}`
	}

	// FIXME: Unmarshsl только для тестов, для оптимизации нужно переделать на NewDecoder.Decode
	if err := json.Unmarshal([]byte(instruction), &qry); err != nil {
		slog.Debug("Erroneous request", slog.String("err", err.Error()))
		// ERROR 11 - Incorrect request structure
		return fmt.Sprintf("{\"action\":\"response\",\"error\":11,\"description\":\"%s\"}", err.Error())
	}

	bAnswer, err := json.Marshal(Processing(qry))
	if err != nil {
		// ERROR 20 - Server error
		return fmt.Sprintf("{\"action\":\"response\",\"error\":20,\"description\":\"%s\"}", err.Error())
	}

	return string(bAnswer)
}
