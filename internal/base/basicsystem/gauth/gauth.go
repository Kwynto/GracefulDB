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
	"sync"

	"github.com/Kwynto/GracefulDB/internal/base/basicsystem/gtypes"
	"github.com/Kwynto/GracefulDB/pkg/lib/closer"
)

const (
	AUTH_FILE   = "./config/auth.gob"
	ACCESS_FILE = "./config/access.gob"
)

type tAuth map[string]string
type tReversAuth map[string]string

var (
	hashMap   tAuth = make(tAuth, 0)
	accessMap       = make(tAuth, 0)

	ticketMap          tAuth       = make(tAuth, 0)
	oldTicketMap       tAuth       = make(tAuth, 0)
	reversTicketMap    tReversAuth = make(tReversAuth, 0)
	reversOldTicketMap tReversAuth = make(tReversAuth, 0)
)

var block sync.RWMutex

func generateTicket() string {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return ""
	}
	return fmt.Sprintf("%x", b)
}

func addUser(login string, password string, access string) error {
	block.RLock()
	_, ok := hashMap[login]
	block.RUnlock()

	if ok {
		return errors.New("unable to create a user")
	}

	block.Lock()
	h := sha256.Sum256([]byte(password))
	hashMap[login] = fmt.Sprintf("%x", h)

	accessMap[login] = access

	hashSave()
	accessSave()
	block.Unlock()

	return nil
}

func AddUser(login string, password string, access string) error {
	if login != "root" {
		return addUser(login, password, access)
	}
	return errors.New("unable to create a user")
}

func updateUser(login string, password string, access string) error {
	block.RLock()
	_, ok := hashMap[login]
	block.RUnlock()

	if !ok {
		return errors.New("unable to update user")
	}

	block.Lock()
	h := sha256.Sum256([]byte(password))
	hashMap[login] = fmt.Sprintf("%x", h)

	if login != "root" {
		accessMap[login] = access
	}

	hashSave()
	accessSave()
	block.Unlock()

	return nil
}

func UpdateUser(login string, password string, access string) error {
	return updateUser(login, password, access)
}

func deleteUser(login string) error {
	block.RLock()
	_, ok := hashMap[login]
	block.RUnlock()

	if ok {
		block.Lock()
		defer block.Unlock()

		delete(oldTicketMap, login)
		delete(ticketMap, login)
		delete(accessMap, login)
		delete(hashMap, login)

		hashSave()
		accessSave()

		return nil
	}

	return errors.New("it is not possible to delete a user")
}

func DeleteUser(login string) error {
	if login != "root" {
		return deleteUser(login)
	}
	return errors.New("it is not possible to delete a user")
}

func CheckTicket(ticket string) (login string, access string, newticket string, err error) {
	block.RLock()
	defer block.RUnlock()

	login, ok1 := reversTicketMap[ticket]
	if ok1 {
		access, ok2 := accessMap[login]
		if ok2 {
			// oldTicket, ok3 := oldTicketMap[login]
			// if ok3 {
			// 	delete(oldTicketMap, login)
			// 	delete(reversOldTicketMap, oldTicket)
			// }
			return login, access, newticket, nil
		}
		return "", "", "", errors.New("authorization failed")
	}

	login, ok4 := reversOldTicketMap[ticket]
	if ok4 {
		access, ok5 := accessMap[login]
		if ok5 {
			newticket, ok6 := ticketMap[login]
			if ok6 {
				return login, access, newticket, nil
			}
			return "", "", "", errors.New("authorization failed")
		}
		return "", "", "", errors.New("authorization failed")
	}

	return "", "", "", errors.New("authorization failed")
}

func NewAuth(secret *gtypes.VSecret) (string, error) {
	var pass string

	if secret.Login == "" {
		slog.Debug("Authentication error", slog.String("login", secret.Login))
		return "", errors.New("authentication error")
	}

	if len(secret.Hash) == 32 {
		pass = secret.Hash
	} else {
		h := sha256.Sum256([]byte(secret.Password))
		pass = fmt.Sprintf("%x", h)
	}

	block.Lock()
	defer block.Unlock()

	dbPass, ok := hashMap[secret.Login]
	if !ok {
		slog.Debug("Authentication error", slog.String("login", secret.Login))
		return "", errors.New("authentication error")
	}

	if pass == dbPass {
		newTicket := generateTicket()

		// don't change this construction
		oldTicket, ok := ticketMap[secret.Login]
		if ok {
			secondOT, ok := oldTicketMap[secret.Login]
			if ok {
				delete(reversOldTicketMap, secondOT)
			}
			reversOldTicketMap[oldTicket] = secret.Login
			oldTicketMap[secret.Login] = oldTicket
		}
		ticketMap[secret.Login] = newTicket
		reversTicketMap[newTicket] = secret.Login
		// end construction

		return newTicket, nil
	}

	slog.Debug("Authentication error", slog.String("login", secret.Login))
	return "", errors.New("authentication error")
}

func hashLoad() {
	var isFNotEx bool = false

	// check if file exists
	if _, err := os.Stat(AUTH_FILE); os.IsNotExist(err) {
		isFNotEx = true
	}

	tempFile, err := os.OpenFile(AUTH_FILE, os.O_CREATE|os.O_APPEND, os.ModePerm)
	if err != nil {
		slog.Error("The authentication file cannot be opened", slog.String("err", err.Error()))
	}
	defer tempFile.Close()

	if isFNotEx {
		h := sha256.Sum256([]byte("toor"))
		pass := fmt.Sprintf("%x", h)

		hashMap["root"] = pass

		encoder := gob.NewEncoder(tempFile)
		if err := encoder.Encode(hashMap); err != nil {
			slog.Debug("Error writing authentication data")
		}
	}

	decoder := gob.NewDecoder(tempFile)
	decoder.Decode(&hashMap)
}

func hashSave() {
	tempFile, err := os.OpenFile(AUTH_FILE, os.O_CREATE|os.O_APPEND, os.ModePerm)
	if err != nil {
		slog.Error("The authentication file cannot be opened", slog.String("err", err.Error()))
	}
	defer tempFile.Close()

	encoder := gob.NewEncoder(tempFile)
	if err := encoder.Encode(hashMap); err != nil {
		slog.Debug("Error writing authentication data")
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
		slog.Error("The authentication file cannot be opened", slog.String("err", err.Error()))
	}
	defer tempFile.Close()

	if isFNotEx {
		accessMap["root"] = "admin"

		encoder := gob.NewEncoder(tempFile)
		if err := encoder.Encode(accessMap); err != nil {
			slog.Debug("Error writing authentication data")
		}
	}

	decoder := gob.NewDecoder(tempFile)
	decoder.Decode(&accessMap)
}

func accessSave() {
	tempFile, err := os.OpenFile(ACCESS_FILE, os.O_CREATE|os.O_APPEND, os.ModePerm)
	if err != nil {
		slog.Error("The authentication file cannot be opened", slog.String("err", err.Error()))
	}
	defer tempFile.Close()

	encoder := gob.NewEncoder(tempFile)
	if err := encoder.Encode(accessMap); err != nil {
		slog.Debug("Error writing authentication data")
	}
}

// Package initialization
func Start() {
	block.Lock()
	hashLoad()
	accessLoad()
	block.Unlock()
}

func Shutdown(ctx context.Context, c *closer.Closer) {
	block.Lock()
	hashSave()
	accessSave()
	block.Unlock()

	c.Done()
}
