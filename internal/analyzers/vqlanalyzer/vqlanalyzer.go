package vqlanalyzer

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"

	"github.com/Kwynto/GracefulDB/internal/gtypes"
)

func generateTicket() string {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return ""
	}
	return fmt.Sprintf("%x", b)
}

func NewAuth(secret *gtypes.VSecret) (string, error) {
	// FIXME: Сделать настоящую авторизацию
	var pass string

	if len(secret.Hash) != 32 {
		pass = secret.Hash
	} else {
		h := sha256.Sum256([]byte(secret.Password))
		pass = fmt.Sprintf("%x", h)
	}

	// FIXME: delete it
	slog.Debug(pass)

	// dbPass := secret.Login

	// if pass == dbPass {
	// 	return generateTicket(), nil
	// }

	if secret.Login == "root" && secret.Password == "toor" {
		return generateTicket(), nil
	}
	slog.Debug("Authorization error", slog.String("login", secret.Login), slog.String("pass", secret.Password))
	return "", errors.New("authorization error")
}

// TODO: Processing
func Processing(in *gtypes.VQuery) *gtypes.VAnswer {
	var response gtypes.VAnswer
	var msgDesc string

	switch in.Action {

	// TODO: auth
	case "auth":
		ticket, err := NewAuth(&in.Secret)
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
