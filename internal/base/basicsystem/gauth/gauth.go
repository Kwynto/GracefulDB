package gauth

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"sync"

	"github.com/Kwynto/GracefulDB/internal/base/basicsystem/gtypes"
	"github.com/Kwynto/GracefulDB/pkg/lib/closer"
)

const (
	AUTH_FILE   = "./config/auth.json"
	ACCESS_FILE = "./config/access.json"

	DEFAULT_USER     = "root"
	DEFAULT_PASSWORD = "toor"
)

type tRole int

const (
	SYSTEM tRole = iota
	ADMIN
	ENGINEER
	MANAGER
	USER
	WEBUSER
)

func (t tRole) String() string {
	return [...]string{"SYSTEM", "ADMIN", "ENGINEER", "MANAGER", "USER", "WEBUSER"}[t]
}

type tRights struct {
	Role  tRole
	Rules []string // []tRule
}

type tAuth map[string]string // map[tLogin]tHach

type tTicket map[string]string       // login - ticket
type tReversTicket map[string]string // ticket - login

type tAccess map[string]tRights // map[tLogin]tRights

var (
	hashMap   tAuth = make(tAuth, 0)
	accessMap       = make(tAccess, 0)

	ticketMap          tTicket       = make(tTicket, 0)
	oldTicketMap       tTicket       = make(tTicket, 0)
	reversTicketMap    tReversTicket = make(tReversTicket, 0)
	reversOldTicketMap tReversTicket = make(tReversTicket, 0)
)

var block sync.RWMutex

// Ticket generation
func generateTicket() string {
	// This function is complete
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return ""
	}
	return fmt.Sprintf("%x", b)
}

// Adding a user - internal
func addUser(login string, password string, access tRights) error {
	// This function is complete
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

// Updating a user - internal
func updateUser(login string, password string, access tRights) error {
	// This function is complete
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

// Deleting a user - internal
func deleteUser(login string) error {
	block.RLock()
	_, ok := hashMap[login]
	block.RUnlock()

	if !ok {
		return errors.New("it is not possible to delete a user")
	}

	block.Lock()
	defer block.Unlock()

	delete(hashMap, login)
	delete(accessMap, login)

	// TODO: Need to write deleting reverses
	delete(ticketMap, login)
	delete(oldTicketMap, login)

	hashSave()
	accessSave()

	return nil
}

// Adding a user
func AddUser(login string, password string, access tRights) error {
	// This function is complete
	if login != "root" {
		return addUser(login, password, access)
	}
	return errors.New("unable to create a user")
}

// Updating a user
func UpdateUser(login string, password string, access tRights) error {
	// This function is complete
	return updateUser(login, password, access)
}

// Deleting a user
func DeleteUser(login string) error {
	// This function is complete
	if login != "root" {
		return deleteUser(login)
	}
	return errors.New("it is not possible to delete a user")
}

func CheckTicket(ticket string) (login string, access tRights, newticket string, err error) {
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
		return "", tRights{}, "", errors.New("authorization failed")
	}

	login, ok4 := reversOldTicketMap[ticket]
	if ok4 {
		access, ok5 := accessMap[login]
		if ok5 {
			newticket, ok6 := ticketMap[login]
			if ok6 {
				return login, access, newticket, nil
			}
			return "", tRights{}, "", errors.New("authorization failed")
		}
		return "", tRights{}, "", errors.New("authorization failed")
	}

	return "", tRights{}, "", errors.New("authorization failed")
}

// Authorization verification and ticket issuance
func NewAuth(secret *gtypes.VSecret) (string, error) {
	// This function is complete
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
			delete(reversTicketMap, oldTicket)
		}
		ticketMap[secret.Login] = newTicket
		reversTicketMap[newTicket] = secret.Login
		// end construction

		return newTicket, nil
	}

	slog.Debug("Authentication error", slog.String("login", secret.Login))
	return "", errors.New("authentication error")
}

// Loading hashs of users from a file
func hashLoad() {
	// This function is complete
	_, err := os.Stat(AUTH_FILE)
	isFNotEx := os.IsNotExist(err)

	tempFile, err := os.OpenFile(AUTH_FILE, os.O_CREATE|os.O_APPEND, os.ModePerm)
	if err != nil {
		slog.Error("The authentication file cannot be opened", slog.String("file", AUTH_FILE), slog.String("err", err.Error()))
	}
	defer tempFile.Close()

	if isFNotEx {
		h := sha256.Sum256([]byte(DEFAULT_PASSWORD))
		pass := fmt.Sprintf("%x", h)

		hashMap[DEFAULT_USER] = pass

		encoder := json.NewEncoder(tempFile)
		if err := encoder.Encode(hashMap); err != nil {
			slog.Debug("Error writing authentication data", slog.String("file", AUTH_FILE))
		}
	} else {
		decoder := json.NewDecoder(tempFile)
		if err := decoder.Decode(&hashMap); err != nil {
			slog.Debug("Error loading the configuration file", slog.String("file", AUTH_FILE))
		}
	}
}

// Saving hashs of users in a file
func hashSave() {
	// This function is complete
	if _, err := os.Stat(AUTH_FILE); !os.IsNotExist(err) {
		if err2 := os.Remove(AUTH_FILE); err2 != nil {
			slog.Warn("Unable to delete file", slog.String("file", AUTH_FILE), slog.String("err", err2.Error()))
		}
	}

	tempFile, err := os.OpenFile(AUTH_FILE, os.O_CREATE|os.O_APPEND, os.ModePerm)
	if err != nil {
		slog.Error("The authentication file cannot be opened", slog.String("file", AUTH_FILE), slog.String("err", err.Error()))
	}
	defer tempFile.Close()

	encoder := json.NewEncoder(tempFile)
	if err := encoder.Encode(hashMap); err != nil {
		slog.Warn("Error writing authentication data", slog.String("file", AUTH_FILE), slog.String("err", err.Error()))
	}
}

// Loading access of users from a file
func accessLoad() {
	// This function is complete
	_, err := os.Stat(ACCESS_FILE)
	isFNotEx := os.IsNotExist(err)

	tempFile, err := os.OpenFile(ACCESS_FILE, os.O_CREATE|os.O_APPEND, os.ModePerm)
	if err != nil {
		slog.Error("The authentication file cannot be opened", slog.String("file", ACCESS_FILE), slog.String("err", err.Error()))
	}
	defer tempFile.Close()

	if isFNotEx {
		accessMap[DEFAULT_USER] = tRights{
			Role:  ADMIN,
			Rules: []string{},
		}

		encoder := json.NewEncoder(tempFile)
		if err := encoder.Encode(accessMap); err != nil {
			slog.Debug("Error writing authentication data", slog.String("file", ACCESS_FILE))
		}
	} else {
		decoder := json.NewDecoder(tempFile)
		if err := decoder.Decode(&accessMap); err != nil {
			slog.Debug("Error loading the configuration file", slog.String("file", ACCESS_FILE))
		}
	}
}

// Saving access of users in a file
func accessSave() {
	// This function is complete
	if _, err := os.Stat(ACCESS_FILE); !os.IsNotExist(err) {
		if err2 := os.Remove(ACCESS_FILE); err2 != nil {
			slog.Warn("Unable to delete file", slog.String("file", ACCESS_FILE), slog.String("err", err2.Error()))
		}
	}

	tempFile, err := os.OpenFile(ACCESS_FILE, os.O_CREATE|os.O_APPEND, os.ModePerm)
	if err != nil {
		slog.Error("The authentication file cannot be opened", slog.String("file", ACCESS_FILE), slog.String("err", err.Error()))
	}
	defer tempFile.Close()

	encoder := json.NewEncoder(tempFile)
	if err := encoder.Encode(accessMap); err != nil {
		slog.Warn("Error writing authentication data", slog.String("file", ACCESS_FILE), slog.String("err", err.Error()))
	}
}

// Package initialization
func Start() {
	// This function is complete
	block.Lock()
	hashLoad()
	accessLoad()
	block.Unlock()
	slog.Info("The authentication system is running.")
}

// Shutting down the service
func Shutdown(ctx context.Context, c *closer.Closer) {
	// This function is complete
	block.Lock()
	hashSave()
	accessSave()
	block.Unlock()
	slog.Info("The authentication system is stopped.")
	c.Done()
}
