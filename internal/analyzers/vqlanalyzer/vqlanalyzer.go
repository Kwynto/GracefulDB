package vqlanalyzer

import (
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/Kwynto/GracefulDB/internal/base/basicsystem/gauth"
	"github.com/Kwynto/GracefulDB/internal/base/basicsystem/gtypes"
)

// TODO: Reading
func Reading(login *string, access *gauth.TRights, fields *gtypes.VFields) *gtypes.VData {
	// TODO: Сделать чтение из базы (возможно перенести эту функцию в ядро)
	return &gtypes.VData{}
}

// TODO: Processing
func Processing(in *gtypes.VQuery) *gtypes.VAnswer {
	var response gtypes.VAnswer = gtypes.VAnswer{
		Action: "response",
		Secret: gtypes.VSecret{
			QueryID: in.Secret.QueryID,
		},
		Error: 0,
	}

	switch in.Action {

	case "auth":
		// The branch is completed
		ticket, err := gauth.NewAuth(&in.Secret)
		if err != nil || ticket == "" {
			// Authorization error (code 432)
			response.Error = 432
			response.Description = "Authorization error"
			return &response
		}
		response.Secret.Ticket = ticket

	case "read":
		// The branch is completed
		login, access, newticket, err := gauth.CheckTicket(in.Secret.Ticket)
		if err != nil {
			slog.Debug("Authorization error", slog.String("operation", "read"))
			// Authorization error (code 440)
			response.Error = 440
			response.Description = "Authorization error"
		}
		if newticket != "" && newticket != in.Secret.Ticket {
			response.Secret.Ticket = newticket
		}
		// slog.Debug("Serving a read request", slog.String("login", login), slog.String("access", fmt.Sprint(access)))
		response.Data = *Reading(&login, &access, &in.Fields)

	// TODO: store
	case "store":
		// -

	// TODO: delete
	case "delete":
		// -

	// TODO: manage
	case "manage":
		// -

	default:
		if in.Action == "" {
			slog.Debug("Empty command.")
			// Empty command (code 430)
			response.Error = 430
			response.Description = "Empty command."
		} else {
			msgDesc := fmt.Sprintf("Unknown command: \"%s\".", in.Action)
			slog.Debug(msgDesc)
			// Unknown command (code 431)
			response.Error = 431
			response.Description = msgDesc
		}
	}

	return &response
}

func Request(instruction *[]byte) *[]byte {
	var qry *gtypes.VQuery

	if !json.Valid(*instruction) {
		slog.Debug("No valid query", slog.String("instruction", string(*instruction)))
		// ERROR 420 - Invalid request
		resB := []byte(`{"action":"response","error":420,"description":"Invalid request"}`)
		return &resB
	}

	// FIXME: Unmarshsl только для тестов, для оптимизации нужно переделать на NewDecoder.Decode
	if err := json.Unmarshal(*instruction, &qry); err != nil {
		slog.Debug("Erroneous request", slog.String("err", err.Error()))
		// ERROR 421 - Incorrect request structure
		resB := []byte(fmt.Sprintf("{\"action\":\"response\",\"error\":421,\"description\":\"%s\"}", err.Error()))
		return &resB
	}

	bAnswer, err := json.Marshal(Processing(qry))
	if err != nil {
		// ERROR 410 - Server error
		resB := []byte(fmt.Sprintf("{\"action\":\"response\",\"error\":410,\"description\":\"%s\"}", err.Error()))
		return &resB
	}

	return &bAnswer
}
