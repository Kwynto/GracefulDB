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

var (
	SAuthFile   = AUTH_FILE
	SAccessFile = ACCESS_FILE
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

	for _, iRole := range t.Roles {
		if iRole == ADMIN {
			return true
		}
		for _, iRule := range rules {
			if iRole == iRule && iRole != 0 {
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
	MHash   tAuth   = make(tAuth, 0)
	MAccess tAccess = make(tAccess, 0)

	mTicket    tTicket = make(tTicket, 0)
	mOldTicket tTicket = make(tTicket, 0)

	mReversTicket    tReversTicket = make(tReversTicket, 0)
	mReversOldTicket tReversTicket = make(tReversTicket, 0)
)

var mxAuth sync.RWMutex

// Internal functions

// Ticket generation
func GenerateTicket() string {
	// This function is complete
	slB := make([]byte, 32)
	rand.Read(slB)
	// if err != nil {
	// 	return ""
	// }
	return fmt.Sprintf("%x", slB)
}

// Getting a ticket - internal
func getTicket(sLogin string) (string, error) {
	sTicket, isOk := mTicket[sLogin]
	if isOk {
		return sTicket, nil
	}
	return "", errors.New("invalid ticket")
}

// Adding a user - internal
func addUser(sLogin string, sPassword string, stAccess TProfile) error {
	// This function is complete
	mxAuth.RLock()
	_, isOk := MHash[sLogin]
	mxAuth.RUnlock()

	if isOk {
		return errors.New("unable to create a user")
	}

	mxAuth.Lock()
	arBH := sha256.Sum256([]byte(sPassword))
	MHash[sLogin] = fmt.Sprintf("%x", arBH)

	MAccess[sLogin] = stAccess

	hashSave()
	accessSave()
	mxAuth.Unlock()

	return nil
}

// Updating a user - internal
func updateUser(sLogin string, sPassword string, stAccess TProfile) error {
	// This function is complete
	mxAuth.RLock()
	_, isOk := MHash[sLogin]
	mxAuth.RUnlock()

	if !isOk {
		return errors.New("unable to update user")
	}

	mxAuth.Lock()
	arBH := sha256.Sum256([]byte(sPassword))
	MHash[sLogin] = fmt.Sprintf("%x", arBH)

	if sLogin != "root" {
		MAccess[sLogin] = stAccess
	}

	hashSave()
	accessSave()
	mxAuth.Unlock()

	return nil
}

// Deleting a user - internal
func deleteUser(sLogin string) error {
	// This function is complete
	mxAuth.RLock()
	_, isOk := MHash[sLogin]
	mxAuth.RUnlock()

	if !isOk {
		return errors.New("it is not possible to delete a user")
	}

	mxAuth.Lock()
	defer mxAuth.Unlock()

	if sTicket, isOk := mOldTicket[sLogin]; isOk {
		delete(mReversOldTicket, sTicket)
		delete(mOldTicket, sLogin)
	}

	if sTicket, isOk := mTicket[sLogin]; isOk {
		delete(mReversTicket, sTicket)
		delete(mTicket, sLogin)
	}

	delete(MAccess, sLogin)
	delete(MHash, sLogin)

	hashSave()
	accessSave()

	return nil
}

// Blocking the user - internal
func blockUser(sLogin string) error {
	// This function is complete
	mxAuth.RLock()
	_, isOk := MHash[sLogin]
	mxAuth.RUnlock()

	if !isOk {
		return errors.New("it is not possible to block a user")
	}

	mxAuth.Lock()
	defer mxAuth.Unlock()

	stAccess, isOk := MAccess[sLogin]
	if !isOk {
		return errors.New("it is not possible to block a user")
	}

	stAccess.Status = BANED
	MAccess[sLogin] = stAccess

	hashSave()
	accessSave()

	return nil
}

// UnBlocking the user - internal
func unblockUser(sLogin string) error {
	// This function is complete
	mxAuth.RLock()
	_, isOk := MHash[sLogin]
	mxAuth.RUnlock()

	if !isOk {
		return errors.New("it is not possible to unblock a user")
	}

	mxAuth.Lock()
	defer mxAuth.Unlock()

	stAccess, isOk := MAccess[sLogin]
	if !isOk {
		return errors.New("it is not possible to unblock a user")
	}

	stAccess.Status = ACTIVE
	MAccess[sLogin] = stAccess

	hashSave()
	accessSave()

	return nil
}

// Updating a profile of a user - internal
func updateProfile(sLogin string, stAccess TProfile) error {
	// This function is complete
	mxAuth.RLock()
	_, isOk := MHash[sLogin]
	mxAuth.RUnlock()

	if !isOk {
		return errors.New("unable to update user")
	}

	mxAuth.Lock()

	if sLogin != "root" {
		MAccess[sLogin] = stAccess
	}

	accessSave()
	mxAuth.Unlock()

	return nil
}

// Public functions

// Getting a ticket
func GetTicket(login string) (ticket string, err error) {
	// This function is complete
	sOperation := "internal -> engine -> gAuth -> GetTicket"
	defer func() { e.Wrapper(sOperation, err) }()

	return getTicket(login)
}

// Adding a user
func AddUser(sLogin string, sPassword string, stAccess TProfile) (err error) {
	// This function is complete
	sOperation := "internal -> engine -> gAuth -> AddUser"
	defer func() { e.Wrapper(sOperation, err) }()

	if sLogin != "root" {
		return addUser(sLogin, sPassword, stAccess)
	}
	return errors.New("unable to create a user")
}

// Updating a user
func UpdateUser(sLogin string, sPassword string, stAccess TProfile) (err error) {
	// This function is complete
	sOperation := "internal -> engine -> gAuth -> UpdateUser"
	defer func() { e.Wrapper(sOperation, err) }()

	return updateUser(sLogin, sPassword, stAccess)
}

// Deleting a user
func DeleteUser(sLogin string) (err error) {
	// This function is complete
	sOperation := "internal -> engine -> gAuth -> DeleteUser"
	defer func() { e.Wrapper(sOperation, err) }()

	if sLogin != "root" {
		return deleteUser(sLogin)
	}
	return errors.New("it is not possible to delete a user")
}

// Blocking the user
func BlockUser(sLogin string) (err error) {
	// This function is complete
	sOperation := "internal -> engine -> gAuth -> BlockUser"
	defer func() { e.Wrapper(sOperation, err) }()

	if sLogin != "root" {
		return blockUser(sLogin)
	}
	return errors.New("it is not possible to delete a user")
}

// Unblocking the user
func UnblockUser(sLogin string) (err error) {
	// This function is complete
	sOperation := "internal -> engine -> gAuth -> UnblockUser"
	defer func() { e.Wrapper(sOperation, err) }()

	if sLogin != "root" {
		return unblockUser(sLogin)
	}
	return errors.New("it is not possible to delete a user")
}

// Updating a profile of a user
func UpdateProfile(sLogin string, stAccess TProfile) (err error) {
	// This function is complete
	sOperation := "internal -> engine -> gAuth -> UpdateProfile"
	defer func() { e.Wrapper(sOperation, err) }()

	return updateProfile(sLogin, stAccess)
}

// User verification
func CheckUser(sUser string, sPassword string) bool {
	// This function is complete
	sDBPass, isOk := MHash[sUser]
	if !isOk {
		return false
	}

	arBH := sha256.Sum256([]byte(sPassword))
	sPass := fmt.Sprintf("%x", arBH)

	return sDBPass == sPass
}

// Get user's profile and access rights
func GetProfile(sUser string) (stProf TProfile, err error) {
	// This function is complete
	sOperation := "internal -> engine -> gAuth -> GetProfile"
	defer func() { e.Wrapper(sOperation, err) }()

	stAccess, isOk := MAccess[sUser]
	if isOk {
		return stAccess, nil
	}
	return stProf, errors.New("profile error")
}

// Verifying the authenticity of the ticket and obtaining access rights.
func CheckTicket(sTicket string) (sLogin string, stAccess TProfile, sNewTicket string, err error) {
	// This function is complete
	sOperation := "internal -> engine -> gAuth -> CheckTicket"
	defer func() { e.Wrapper(sOperation, err) }()

	mxAuth.RLock()
	defer mxAuth.RUnlock()

	sLogin, isOk1 := mReversTicket[sTicket]
	if isOk1 {
		stAccess, isOk2 := MAccess[sLogin]
		if isOk2 {
			return sLogin, stAccess, sNewTicket, nil
		}
		return "", TProfile{}, "", errors.New("authorization failed")
	}

	sLogin, isOk4 := mReversOldTicket[sTicket]
	if isOk4 {
		stAccess, isOk5 := MAccess[sLogin]
		if isOk5 {
			sNewTicket, isOk6 := mTicket[sLogin]
			if isOk6 {
				return sLogin, stAccess, sNewTicket, nil
			}
			return "", TProfile{}, "", errors.New("authorization failed")
		}
		return "", TProfile{}, "", errors.New("authorization failed")
	}

	return "", TProfile{}, "", errors.New("authorization failed")
}

// Authorization verification and ticket issuance
func NewAuth(stSecret *gtypes.TSecret) (sTicket string, err error) {
	// This function is complete
	sOperation := "internal -> engine -> gAuth -> NewAuth"
	defer func() { e.Wrapper(sOperation, err) }()

	var sPass string

	if stSecret.Login == "" {
		slog.Debug("Authentication error", slog.String("login", stSecret.Login))
		return "", errors.New("authentication error")
	}

	if len(stSecret.Hash) == 64 {
		sPass = stSecret.Hash
	} else {
		arBH := sha256.Sum256([]byte(stSecret.Password))
		sPass = fmt.Sprintf("%x", arBH)
	}

	mxAuth.Lock()
	defer mxAuth.Unlock()

	sDBPass, isOk := MHash[stSecret.Login]
	if !isOk {
		slog.Debug("Authentication error", slog.String("login", stSecret.Login))
		return "", errors.New("authentication error")
	}

	if sPass == sDBPass {
		sNewTicket := GenerateTicket()

		// don't change this construction
		sOldTicket, isOk := mTicket[stSecret.Login]
		if isOk {
			sSecondOldTicket, isOk := mOldTicket[stSecret.Login]
			if isOk {
				delete(mReversOldTicket, sSecondOldTicket)
			}
			mReversOldTicket[sOldTicket] = stSecret.Login
			mOldTicket[stSecret.Login] = sOldTicket
			delete(mReversTicket, sOldTicket)
		}
		mTicket[stSecret.Login] = sNewTicket
		mReversTicket[sNewTicket] = stSecret.Login
		// end construction

		return sNewTicket, nil
	}

	slog.Debug("Authentication error", slog.String("login", stSecret.Login))
	return "", errors.New("authentication error")
}

// Loading hashs of users from a file
func hashLoad() {
	// This function is complete
	_, err := os.Stat(SAuthFile)
	isFNotEx := os.IsNotExist(err)

	fAuthFile, err := os.OpenFile(SAuthFile, os.O_CREATE|os.O_APPEND, os.ModePerm)
	if err != nil {
		slog.Error("The authentication file cannot be opened", slog.String("file", SAuthFile), slog.String("err", err.Error()))
	}
	defer fAuthFile.Close()

	if isFNotEx {
		arBH := sha256.Sum256([]byte(DEFAULT_PASSWORD))
		sPass := fmt.Sprintf("%x", arBH)

		MHash[DEFAULT_USER] = sPass

		jeAuthFile := json.NewEncoder(fAuthFile)
		if err := jeAuthFile.Encode(MHash); err != nil {
			slog.Debug("Error writing authentication data", slog.String("file", SAuthFile))
		}
	} else {
		jdAuthFile := json.NewDecoder(fAuthFile)
		if err := jdAuthFile.Decode(&MHash); err != nil {
			slog.Debug("Error loading the configuration file", slog.String("file", SAuthFile))
		}
	}
}

// Saving hashs of users in a file
func hashSave() {
	// This function is complete
	if _, err := os.Stat(SAuthFile); !os.IsNotExist(err) {
		if err2 := os.Remove(SAuthFile); err2 != nil {
			slog.Warn("Unable to delete file", slog.String("file", SAuthFile), slog.String("err", err2.Error()))
		}
	}

	fAuthFile, err := os.OpenFile(SAuthFile, os.O_CREATE|os.O_APPEND, os.ModePerm)
	if err != nil {
		slog.Error("The authentication file cannot be opened", slog.String("file", SAuthFile), slog.String("err", err.Error()))
	}
	defer fAuthFile.Close()

	jeAuthFile := json.NewEncoder(fAuthFile)
	if err := jeAuthFile.Encode(MHash); err != nil {
		slog.Warn("Error writing authentication data", slog.String("file", SAuthFile), slog.String("err", err.Error()))
	}
}

// Loading access of users from a file
func accessLoad() {
	// This function is complete
	_, err := os.Stat(SAccessFile)
	isFNotEx := os.IsNotExist(err)

	fAccessFile, err := os.OpenFile(SAccessFile, os.O_CREATE|os.O_APPEND, os.ModePerm)
	if err != nil {
		slog.Error("The authentication file cannot be opened", slog.String("file", SAccessFile), slog.String("err", err.Error()))
	}
	defer fAccessFile.Close()

	if isFNotEx {
		MAccess[DEFAULT_USER] = TProfile{
			Description: "This is the main user.",
			Status:      ACTIVE,
			Roles:       []TRole{ADMIN},
		}

		jeAccessFile := json.NewEncoder(fAccessFile)
		if err := jeAccessFile.Encode(MAccess); err != nil {
			slog.Debug("Error writing authentication data", slog.String("file", SAccessFile))
		}
	} else {
		jdAccessFile := json.NewDecoder(fAccessFile)
		if err := jdAccessFile.Decode(&MAccess); err != nil {
			slog.Debug("Error loading the configuration file", slog.String("file", SAccessFile))
		}
	}
}

// Saving access of users in a file
func accessSave() {
	// This function is complete
	if _, err := os.Stat(SAccessFile); !os.IsNotExist(err) {
		if err2 := os.Remove(SAccessFile); err2 != nil {
			slog.Warn("Unable to delete file", slog.String("file", SAccessFile), slog.String("err", err2.Error()))
		}
	}

	fAccessFile, err := os.OpenFile(SAccessFile, os.O_CREATE|os.O_APPEND, os.ModePerm)
	if err != nil {
		slog.Error("The authentication file cannot be opened", slog.String("file", SAccessFile), slog.String("err", err.Error()))
	}
	defer fAccessFile.Close()

	jeAccessFile := json.NewEncoder(fAccessFile)
	if err := jeAccessFile.Encode(MAccess); err != nil {
		slog.Warn("Error writing authentication data", slog.String("file", SAccessFile), slog.String("err", err.Error()))
	}
}

// Checking the root user's password for the default value.
func checkingTheDefaultPassword() bool {
	arBDefH := sha256.Sum256([]byte(DEFAULT_PASSWORD))
	sDefPass := fmt.Sprintf("%x", arBDefH)

	sPass := MHash["root"]

	return sDefPass == sPass
}

// Package initialization
func Start() {
	// This function is complete
	mxAuth.Lock()
	hashLoad()
	accessLoad()
	mxAuth.Unlock()
	slog.Info("The authentication system is running.")
	if checkingTheDefaultPassword() {
		sWarnMsg := fmt.Sprintf("The 'root' user has a default password of '%s'. Please change your password!", DEFAULT_PASSWORD)
		slog.Warn(sWarnMsg)
	}
}

// Shutting down the service
func Shutdown(ctx context.Context, c *closer.TCloser) {
	// This function is complete
	mxAuth.Lock()
	hashSave()
	accessSave()
	mxAuth.Unlock()
	slog.Info("The authentication system is stopped.")
	c.Done()
}
