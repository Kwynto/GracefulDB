package vqlanalyzer

import (
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/Kwynto/GracefulDB/internal/base/basicsystem/gauth"
	"github.com/Kwynto/GracefulDB/internal/base/basicsystem/gtypes"
)

// TODO: Processing
func Processing(in *gtypes.VQuery) *gtypes.VAnswer {
	var response gtypes.VAnswer
	var msgDesc string

	switch in.Action {

	case "auth":
		ticket, err := gauth.NewAuth(&in.Secret)
		if err != nil || ticket == "" {
			return &gtypes.VAnswer{
				Action: "response",
				// Authorization error (code 432)
				Error:       432,
				Description: "Authorization error",
			}
		}

		response = gtypes.VAnswer{
			Action: "response",
			Secret: gtypes.VSecret{
				Ticket:  ticket,
				QueryID: in.Secret.QueryID,
			},
			Error: 0,
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

	// TODO: manage
	case "manage":
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
				// Empty command (code 430)
				Error:       430,
				Description: msgDesc,
			}
		} else {
			msgDesc = fmt.Sprintf("Unknown command: \"%s\".", in.Action)
			slog.Debug(msgDesc)
			response = gtypes.VAnswer{
				Action: "response",
				// Unknown command (code 431)
				Error:       431,
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
		// ERROR 420 - Invalid request
		return `{"action":"response","error":420,"description":"Invalid request"}`
	}

	// FIXME: Unmarshsl только для тестов, для оптимизации нужно переделать на NewDecoder.Decode
	if err := json.Unmarshal([]byte(instruction), &qry); err != nil {
		slog.Debug("Erroneous request", slog.String("err", err.Error()))
		// ERROR 421 - Incorrect request structure
		return fmt.Sprintf("{\"action\":\"response\",\"error\":421,\"description\":\"%s\"}", err.Error())
	}

	bAnswer, err := json.Marshal(Processing(qry))
	if err != nil {
		// ERROR 410 - Server error
		return fmt.Sprintf("{\"action\":\"response\",\"error\":410,\"description\":\"%s\"}", err.Error())
	}

	return string(bAnswer)
}
