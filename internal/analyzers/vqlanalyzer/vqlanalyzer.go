package vqlanalyzer

import (
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/Kwynto/GracefulDB/internal/engine/basicsystem/gauth"
	"github.com/Kwynto/GracefulDB/internal/engine/basicsystem/gtypes"
)

type VAuth struct {
	Login  string
	Access gauth.TProfile
}

// TODO: Reading
func Reading(auth *VAuth, fields *gtypes.VFields) *gtypes.VData {
	// TODO: Сделать чтение из базы (возможно перенести эту функцию в ядро)
	return &gtypes.VData{}
}

// TODO: Storing
func Storing(auth *VAuth, in *gtypes.VQuery) *gtypes.VData {
	// TODO: Сделать запись в базу (возможно перенести эту функцию в ядро)
	return &gtypes.VData{}
}

// TODO: Deleting
func Deleting(auth *VAuth, fields *gtypes.VFields) *gtypes.VData {
	// TODO: Сделать удаление из базы (возможно перенести эту функцию в ядро)
	return &gtypes.VData{}
}

// TODO: Managing
func Managing(auth *VAuth, in *gtypes.VQuery) *gtypes.VData {
	// TODO: Сделать управление базой (возможно перенести эту функцию в ядро)
	return &gtypes.VData{}
}

// Request processing
func Processing(in *gtypes.VQuery) *gtypes.VAnswer {
	var response gtypes.VAnswer = gtypes.VAnswer{
		Action: "response",
		Secret: gtypes.VSecret{
			QueryID: in.Secret.QueryID,
		},
		Error: 0,
	}

	var auth VAuth

	if in.Action != "auth" {
		login, access, newticket, err := gauth.CheckTicket(in.Secret.Ticket)
		auth = VAuth{
			Login:  login,
			Access: access,
		}
		if newticket != "" && newticket != in.Secret.Ticket {
			response.Secret.Ticket = newticket
		}
		if err != nil {
			slog.Debug("Authorization error", slog.String("operation", "read"))
			// Authorization error (code 440)
			response.Error = 440
			response.Description = "Authorization error"
			return &response
		}
	}

	switch in.Action {
	case "auth":
		ticket, err := gauth.NewAuth(&in.Secret)
		if err != nil || ticket == "" {
			// Authorization error (code 432)
			response.Error = 432
			response.Description = "Authorization error"
			return &response
		}
		response.Secret.Ticket = ticket

	case "read":
		response.Data = *Reading(&auth, &in.Fields)

	case "store":
		response.Data = *Storing(&auth, in)

	case "delete":
		response.Data = *Deleting(&auth, &in.Fields)

	case "manage":
		response.Data = *Managing(&auth, in)

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

// Getting a clean request from the connector and returning a response
func Request(instruction *[]byte) *[]byte {
	var qry *gtypes.VQuery

	if !json.Valid(*instruction) {
		slog.Debug("No valid query", slog.String("instruction", string(*instruction)))
		// ERROR 420 - Invalid request
		resB := []byte(`{"action":"response","error":420,"description":"Invalid request"}`)
		return &resB
	}

	// FIXME: Unmarshal только для тестов, для оптимизации нужно переделать на NewDecoder.Decode
	if err := json.Unmarshal(*instruction, &qry); err != nil {
		slog.Debug("Erroneous request", slog.String("err", err.Error()))
		// ERROR 421 - Incorrect request structure
		resB := []byte(fmt.Sprintf("{\"action\":\"response\",\"error\":421,\"description\":\"%s\"}", err.Error()))
		return &resB
	}

	bAnswer, _ := json.Marshal(Processing(qry))
	// if err != nil {
	// 	// ERROR 410 - Server error
	// 	resB := []byte(fmt.Sprintf("{\"action\":\"response\",\"error\":410,\"description\":\"%s\"}", err.Error()))
	// 	return &resB
	// }

	return &bAnswer
}
