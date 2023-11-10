package vqlanalyzer

import (
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/Kwynto/GracefulDB/internal/gtypes"
)

// TODO: Processing
func Processing(in *gtypes.VQuery) *gtypes.VAnswer {
	var response gtypes.VAnswer
	var msgDesc string

	switch in.Action {
	// TODO: auth
	case "auth":
		response = gtypes.VAnswer{
			Action: "response",
			Error:  0,
		}
	// TODO: read
	case "read":
		response = gtypes.VAnswer{
			Action: "response",
			Error:  0,
		}
	// TODO: store
	case "store":
		response = gtypes.VAnswer{
			Action: "response",
			Error:  0,
		}
	// TODO: delete
	case "delete":
		response = gtypes.VAnswer{
			Action: "response",
			Error:  0,
		}
	default:
		if in.Action == "" {
			msgDesc = "Empty command."
			slog.Debug(msgDesc)
			response = gtypes.VAnswer{
				Action: "response",
				// Empty command (code 30)
				Error:       30,
				Description: msgDesc,
			}
		} else {
			msgDesc = fmt.Sprintf("Unknown command: \"%s\".", in.Action)
			slog.Debug(msgDesc)
			response = gtypes.VAnswer{
				Action: "response",
				// Unknown command (code 31)
				Error:       31,
				Description: msgDesc,
			}
		}
	}

	return &response
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
