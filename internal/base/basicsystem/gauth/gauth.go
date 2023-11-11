package gauth

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/gob"
	"errors"
	"fmt"
	"log/slog"
	"os"

	"github.com/Kwynto/GracefulDB/internal/base/basicsystem/gtypes"
	"github.com/Kwynto/GracefulDB/pkg/lib/closer"
)

const (
	AUTH_FILE   = "./config/auth.gob"
	ACCESS_FILE = "./config/access.gob"
)

type tAuth map[string]string

var hashMap tAuth = make(tAuth, 0)
var accessMap = make(tAuth, 0)

var ticketMap tAuth = make(tAuth, 0)
var oldTicketMap tAuth = make(tAuth, 0)

// var block sync.RWMutex

func generateTicket() string {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return ""
	}
	return fmt.Sprintf("%x", b)
}

func NewAuth(secret *gtypes.VSecret) (string, error) {
	var pass string

	if secret.Login == "" {
		slog.Debug("Authorization error", slog.String("login", secret.Login))
		return "", errors.New("authorization error")
	}

	if len(secret.Hash) == 32 {
		pass = secret.Hash
	} else {
		h := sha256.Sum256([]byte(secret.Password))
		pass = fmt.Sprintf("%x", h)
	}

	dbPass, ok := hashMap[secret.Login]
	if !ok {
		slog.Debug("Authorization error", slog.String("login", secret.Login))
		return "", errors.New("authorization error")
	}

	if pass == dbPass {
		newTicket := generateTicket()

		oldTicket, ok := ticketMap[secret.Login]
		if ok {
			oldTicketMap[secret.Login] = oldTicket
		}
		ticketMap[secret.Login] = newTicket

		return newTicket, nil
	}

	slog.Debug("Authorization error", slog.String("login", secret.Login))
	return "", errors.New("authorization error")
}

func hashLoad() {
	var isFNotEx bool = false

	// check if file exists
	if _, err := os.Stat(AUTH_FILE); os.IsNotExist(err) {
		isFNotEx = true
	}

	tempFile, err := os.OpenFile(AUTH_FILE, os.O_CREATE|os.O_APPEND, os.ModePerm)
	if err != nil {
		slog.Error("The authorization file cannot be opened", slog.String("err", err.Error()))
	}
	defer tempFile.Close()

	if isFNotEx {
		h := sha256.Sum256([]byte("toor"))
		pass := fmt.Sprintf("%x", h)

		hashMap["root"] = pass

		encoder := gob.NewEncoder(tempFile)
		if err := encoder.Encode(hashMap); err != nil {
			slog.Debug("Error writing authorization data")
		}
	}

	decoder := gob.NewDecoder(tempFile)
	decoder.Decode(&hashMap)
}

func hashSave() {
	tempFile, err := os.OpenFile(AUTH_FILE, os.O_CREATE|os.O_APPEND, os.ModePerm)
	if err != nil {
		slog.Error("The authorization file cannot be opened", slog.String("err", err.Error()))
	}
	defer tempFile.Close()

	encoder := gob.NewEncoder(tempFile)
	if err := encoder.Encode(hashMap); err != nil {
		slog.Debug("Error writing authorization data")
	}
}

func accessLoad() {
	var isFNotEx bool = false

	// check if file exists
	if _, err := os.Stat(ACCESS_FILE); os.IsNotExist(err) {
		isFNotEx = true
	}

	tempFile, err := os.OpenFile(ACCESS_FILE, os.O_CREATE|os.O_APPEND, os.ModePerm)
	if err != nil {
		slog.Error("The authorization file cannot be opened", slog.String("err", err.Error()))
	}
	defer tempFile.Close()

	if isFNotEx {
		accessMap["root"] = "admin"

		encoder := gob.NewEncoder(tempFile)
		if err := encoder.Encode(accessMap); err != nil {
			slog.Debug("Error writing authorization data")
		}
	}

	decoder := gob.NewDecoder(tempFile)
	decoder.Decode(&accessMap)
}

func accessSave() {
	tempFile, err := os.OpenFile(ACCESS_FILE, os.O_CREATE|os.O_APPEND, os.ModePerm)
	if err != nil {
		slog.Error("The authorization file cannot be opened", slog.String("err", err.Error()))
	}
	defer tempFile.Close()

	encoder := gob.NewEncoder(tempFile)
	if err := encoder.Encode(accessMap); err != nil {
		slog.Debug("Error writing authorization data")
	}
}

// Package initialization
func Start() {
	hashLoad()
	accessLoad()
}

func Shutdown(ctx context.Context, c *closer.Closer) {
	hashSave()
	accessSave()

	c.Done()
}
