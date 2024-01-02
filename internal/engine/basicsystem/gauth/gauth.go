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

	"github.com/Kwynto/GracefulDB/internal/engine/basicsystem/gtypes"
	"github.com/Kwynto/GracefulDB/pkg/lib/closer"
	"github.com/Kwynto/GracefulDB/pkg/lib/e"
)

const (
	AUTH_FILE   = "./config/auth.json"
	ACCESS_FILE = "./config/access.json"

	DEFAULT_USER     = "root"
	DEFAULT_PASSWORD = "toor"
)

type TRole int

const (
	SYSTEM   TRole = iota
	ADMIN          // All administrator rights
	MANAGER        // User management rights only
	ENGINEER       // Only the rights to control the engine and to force the launch of diagnostic processes.
	USER           // Limited rights of a regular user.
)

func (t TRole) String() string {
	return [...]string{"SYSTEM", "ADMIN", "MANAGER", "ENGINEER", "USER"}[t]
}

func (t TRole) IsSystem() bool {
	return t == SYSTEM
}

func (t TRole) IsAdmin() bool {
	return t == ADMIN
}

func (t TRole) IsManager() bool {
	return t == MANAGER
}

func (t TRole) IsEngineer() bool {
	return t == ENGINEER
}

func (t TRole) IsUser() bool {
	return t == USER
}

func (t TRole) IsNotUser() bool {
	return t != USER
}

type TStatus int

const (
	UNDEFINED TStatus = iota
	NEW
	ACTIVE
	BANED
)

func (t TStatus) String() string {
	return [...]string{"UNDEFINED", "NEW", "ACTIVE", "BANED"}[t]
}

func (t TStatus) IsBad() bool {
	return t < 1 || t > 2
}

func (t TStatus) IsGood() bool {
	return t > 0 && t < 3
}

type TProfile struct {
	Description string
	Status      TStatus
	Roles       []TRole
}

// Chacking of authorization.
func (t TProfile) IsAllowed(rules []TRole) bool {
	if t.Status.IsBad() {
		return false
	}

	for _, role := range t.Roles {
		if role == ADMIN {
			return true
		}
		for _, rule := range rules {
			if role == rule && role != 0 {
				return true
			}
		}
	}

	return false
}

type tAuth map[string]string // map[tLogin]tHach

type tTicket map[string]string       // login - ticket
type tReversTicket map[string]string // ticket - login

type tAccess map[string]TProfile // map[tLogin]TProfile

var (
	HashMap   tAuth = make(tAuth, 0)
	AccessMap       = make(tAccess, 0)

	ticketMap          tTicket       = make(tTicket, 0)
	oldTicketMap       tTicket       = make(tTicket, 0)
	reversTicketMap    tReversTicket = make(tReversTicket, 0)
	reversOldTicketMap tReversTicket = make(tReversTicket, 0)
)

var block sync.RWMutex

// Internal functions

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
func addUser(login string, password string, access TProfile) error {
	// This function is complete
	block.RLock()
	_, ok := HashMap[login]
	block.RUnlock()

	if ok {
		return errors.New("unable to create a user")
	}

	block.Lock()
	h := sha256.Sum256([]byte(password))
	HashMap[login] = fmt.Sprintf("%x", h)

	AccessMap[login] = access

	hashSave()
	accessSave()
	block.Unlock()

	return nil
}

// Updating a user - internal
func updateUser(login string, password string, access TProfile) error {
	// This function is complete
	block.RLock()
	_, ok := HashMap[login]
	block.RUnlock()

	if !ok {
		return errors.New("unable to update user")
	}

	block.Lock()
	h := sha256.Sum256([]byte(password))
	HashMap[login] = fmt.Sprintf("%x", h)

	if login != "root" {
		AccessMap[login] = access
	}

	hashSave()
	accessSave()
	block.Unlock()

	return nil
}

// Deleting a user - internal
func deleteUser(login string) error {
	// This function is complete
	block.RLock()
	_, ok := HashMap[login]
	block.RUnlock()

	if !ok {
		return errors.New("it is not possible to delete a user")
	}

	block.Lock()
	defer block.Unlock()

	if ticket, ok := oldTicketMap[login]; ok {
		delete(reversOldTicketMap, ticket)
		delete(oldTicketMap, login)
	}

	if ticket, ok := ticketMap[login]; ok {
		delete(reversTicketMap, ticket)
		delete(ticketMap, login)
	}

	delete(AccessMap, login)
	delete(HashMap, login)

	hashSave()
	accessSave()

	return nil
}

// Blocking the user - internal
func blockUser(login string) error {
	// This function is complete
	block.RLock()
	_, ok := HashMap[login]
	block.RUnlock()

	if !ok {
		return errors.New("it is not possible to block a user")
	}

	block.Lock()
	defer block.Unlock()

	access, ok := AccessMap[login]
	if !ok {
		return errors.New("it is not possible to block a user")
	}

	access.Status = BANED
	AccessMap[login] = access

	hashSave()
	accessSave()

	return nil
}

// UnBlocking the user - internal
func unblockUser(login string) error {
	// This function is complete
	block.RLock()
	_, ok := HashMap[login]
	block.RUnlock()

	if !ok {
		return errors.New("it is not possible to unblock a user")
	}

	block.Lock()
	defer block.Unlock()

	access, ok := AccessMap[login]
	if !ok {
		return errors.New("it is not possible to unblock a user")
	}

	access.Status = ACTIVE
	AccessMap[login] = access

	hashSave()
	accessSave()

	return nil
}

// Updating a profile of a user - internal
func updateProfile(login string, access TProfile) error {
	// This function is complete
	block.RLock()
	_, ok := HashMap[login]
	block.RUnlock()

	if !ok {
		return errors.New("unable to update user")
	}

	block.Lock()

	if login != "root" {
		AccessMap[login] = access
	}

	accessSave()
	block.Unlock()

	return nil
}

// Public functions

// Adding a user
func AddUser(login string, password string, access TProfile) (err error) {
	// This function is complete
	op := "internal -> engine -> gAuth -> AddUser"
	defer func() { e.Wrapper(op, err) }()

	if login != "root" {
		return addUser(login, password, access)
	}
	return errors.New("unable to create a user")
}

// Updating a user
func UpdateUser(login string, password string, access TProfile) (err error) {
	// This function is complete
	op := "internal -> engine -> gAuth -> UpdateUser"
	defer func() { e.Wrapper(op, err) }()

	return updateUser(login, password, access)
}

// Deleting a user
func DeleteUser(login string) (err error) {
	// This function is complete
	op := "internal -> engine -> gAuth -> DeleteUser"
	defer func() { e.Wrapper(op, err) }()

	if login != "root" {
		return deleteUser(login)
	}
	return errors.New("it is not possible to delete a user")
}

// Blocking the user
func BlockUser(login string) (err error) {
	// This function is complete
	op := "internal -> engine -> gAuth -> BlockUser"
	defer func() { e.Wrapper(op, err) }()

	if login != "root" {
		return blockUser(login)
	}
	return errors.New("it is not possible to delete a user")
}

// Unblocking the user
func UnblockUser(login string) (err error) {
	// This function is complete
	op := "internal -> engine -> gAuth -> UnblockUser"
	defer func() { e.Wrapper(op, err) }()

	if login != "root" {
		return unblockUser(login)
	}
	return errors.New("it is not possible to delete a user")
}

// Updating a profile of a user
func UpdateProfile(login string, access TProfile) (err error) {
	// This function is complete
	op := "internal -> engine -> gAuth -> UpdateProfile"
	defer func() { e.Wrapper(op, err) }()

	return updateProfile(login, access)
}

// User verification
func CheckUser(user string, password string) bool {
	// This function is complete
	dbPass, ok := HashMap[user]
	if !ok {
		return false
	}

	h := sha256.Sum256([]byte(password))
	pass := fmt.Sprintf("%x", h)

	return dbPass == pass
}

// Get user's profile and access rights
func GetProfile(user string) (prof TProfile, err error) {
	// This function is complete
	op := "internal -> engine -> gAuth -> GetProfile"
	defer func() { e.Wrapper(op, err) }()

	access, ok := AccessMap[user]
	if ok {
		return access, nil
	}
	return TProfile{}, errors.New("profile error")
}

// Verifying the authenticity of the ticket and obtaining access rights.
func CheckTicket(ticket string) (login string, access TProfile, newticket string, err error) {
	// This function is complete
	op := "internal -> engine -> gAuth -> CheckTicket"
	defer func() { e.Wrapper(op, err) }()

	block.RLock()
	defer block.RUnlock()

	login, ok1 := reversTicketMap[ticket]
	if ok1 {
		access, ok2 := AccessMap[login]
		if ok2 {
			return login, access, newticket, nil
		}
		return "", TProfile{}, "", errors.New("authorization failed")
	}

	login, ok4 := reversOldTicketMap[ticket]
	if ok4 {
		access, ok5 := AccessMap[login]
		if ok5 {
			newticket, ok6 := ticketMap[login]
			if ok6 {
				return login, access, newticket, nil
			}
			return "", TProfile{}, "", errors.New("authorization failed")
		}
		return "", TProfile{}, "", errors.New("authorization failed")
	}

	return "", TProfile{}, "", errors.New("authorization failed")
}

// Authorization verification and ticket issuance
func NewAuth(secret *gtypes.VSecret) (ticket string, err error) {
	// This function is complete
	op := "internal -> engine -> gAuth -> NewAuth"
	defer func() { e.Wrapper(op, err) }()

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

	dbPass, ok := HashMap[secret.Login]
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

		HashMap[DEFAULT_USER] = pass

		encoder := json.NewEncoder(tempFile)
		if err := encoder.Encode(HashMap); err != nil {
			slog.Debug("Error writing authentication data", slog.String("file", AUTH_FILE))
		}
	} else {
		decoder := json.NewDecoder(tempFile)
		if err := decoder.Decode(&HashMap); err != nil {
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
	if err := encoder.Encode(HashMap); err != nil {
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
		AccessMap[DEFAULT_USER] = TProfile{
			Description: "This is the main user.",
			Status:      ACTIVE,
			Roles:       []TRole{ADMIN},
		}

		encoder := json.NewEncoder(tempFile)
		if err := encoder.Encode(AccessMap); err != nil {
			slog.Debug("Error writing authentication data", slog.String("file", ACCESS_FILE))
		}
	} else {
		decoder := json.NewDecoder(tempFile)
		if err := decoder.Decode(&AccessMap); err != nil {
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
	if err := encoder.Encode(AccessMap); err != nil {
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
