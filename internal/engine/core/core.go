package core

import (
	"context"
	"crypto/rand"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/Kwynto/GracefulDB/internal/config"
	"github.com/Kwynto/GracefulDB/pkg/lib/closer"
	"github.com/Kwynto/GracefulDB/pkg/lib/ecowriter"
)

const (
	INFOFILE_STORAGE = "storage.json"
	INFOFILE_DB      = "db.json"
	INFOFILE_TABLE   = "table.json"
	INFOFILE_COLUMN  = "column.json"
)

type tCoreSettings struct {
	Storage    string
	BucketSize int
	FreezeMode bool
}

type tStorageInfo struct {
	DBs map[string]string `json:"dbs"` // [name db] name folder
}

type tDBInfo struct {
	Name       string    `json:"name"`
	Tables     []string  `json:"tables"`
	LastUpdate time.Time `json:"lastupdate"`
	Deleted    bool      `json:"deleted"`
}

type tCoreFile struct {
	Descriptor *os.File
	Expire     time.Duration
}

type tCoreProcessing struct {
	FileDescriptors map[string]tCoreFile
}

var LocalCoreSettings tCoreSettings = tCoreSettings{
	Storage:    "./data/",
	BucketSize: 800,
	FreezeMode: false,
}

var CoreProcessing tCoreProcessing

var StorageInfo tStorageInfo = tStorageInfo{
	DBs: make(map[string]string, 0),
}

// Name generation
func GenerateName() string {
	// This function is complete
	b := make([]byte, 16)
	rand.Read(b)

	return fmt.Sprintf("%x", b)
}

// Checking the folder name
func CheckFolderOrFile(patch, name string) bool {
	// This function is complete
	fullPath := fmt.Sprintf("%s%s", patch, name)
	_, err := os.Stat(fullPath)

	return os.IsExist(err)
}

// Marks the database as deleted, but does not delete files.
func RemoveDB(name string) bool {
	// This function is complete
	var dbInfo tDBInfo

	folderName, ok := StorageInfo.DBs[name]
	if ok {
		if CheckFolderOrFile(fmt.Sprintf("%s%s", LocalCoreSettings.Storage, folderName), folderName) {
			dbInfoPath := fmt.Sprintf("%s%s/%s", LocalCoreSettings.Storage, folderName, INFOFILE_DB)
			err := ecowriter.ReadJSON(dbInfoPath, &dbInfo)
			if err != nil {
				return false
			}
			dbInfo.LastUpdate = time.Now()
			dbInfo.Deleted = true
			err2 := ecowriter.WriteJSON(dbInfoPath, &dbInfo)
			if err2 != nil {
				return false
			}
		} else {
			return false
		}
	}

	return true
}

// Deletes the folder and database files.
func StrongRemoveDB(name string) bool {
	// This function is complete
	folderName, ok := StorageInfo.DBs[name]
	if ok {
		if CheckFolderOrFile(LocalCoreSettings.Storage, folderName) {
			fullPath := fmt.Sprintf("%s%s", LocalCoreSettings.Storage, folderName)
			err := os.Remove(fullPath)
			if err != nil {
				return false
			}

			delete(StorageInfo.DBs, name)
			storagePath := fmt.Sprintf("%s%s", LocalCoreSettings.Storage, INFOFILE_STORAGE)
			ecowriter.WriteJSON(storagePath, StorageInfo)
			return true
		}
	}

	return false
}

// Creating a new database.
func CreateDB(name string) bool {
	// This function is complete

	_, ok := StorageInfo.DBs[name]
	if ok {
		return false
	}

	var folderName string

	for {
		folderName = GenerateName()
		if !CheckFolderOrFile(LocalCoreSettings.Storage, folderName) {
			break
		}
	}

	fullName := fmt.Sprintf("%s%s", LocalCoreSettings.Storage, folderName)
	err := os.Mkdir(fullName, 0666)
	if err != nil {
		return false
	}

	dbInfoPath := fmt.Sprintf("%s/%s", fullName, INFOFILE_DB)

	dbInfo := tDBInfo{
		Name:       name,
		Tables:     []string{},
		LastUpdate: time.Now(),
		Deleted:    false,
	}

	err2 := ecowriter.WriteJSON(dbInfoPath, dbInfo)
	if err2 != nil {
		return false
	}
	StorageInfo.DBs[name] = folderName
	storagePath := fmt.Sprintf("%s%s", LocalCoreSettings.Storage, INFOFILE_STORAGE)
	ecowriter.WriteJSON(storagePath, StorageInfo)

	return true
}

func LoadLocalCoreSettings(cfg *config.Config) tCoreSettings {
	return tCoreSettings{
		Storage:    cfg.CoreSettings.Storage,
		BucketSize: cfg.CoreSettings.BucketSize,
		FreezeMode: cfg.CoreSettings.FreezeMode,
	}
}

func Engine(cfg *config.Config) {
	LocalCoreSettings = LoadLocalCoreSettings(cfg)

	storagePath := fmt.Sprintf("%s%s", LocalCoreSettings.Storage, INFOFILE_STORAGE)
	err := ecowriter.ReadJSON(storagePath, &StorageInfo)
	if err != nil {
		StorageInfo.DBs = make(map[string]string, 0)
		ecowriter.WriteJSON(storagePath, StorageInfo)
	}

	slog.Info("The core of the DBMS was started.")
}

func Shutdown(ctx context.Context, c *closer.Closer) {
	storagePath := fmt.Sprintf("%s%s", LocalCoreSettings.Storage, INFOFILE_STORAGE)
	ecowriter.WriteJSON(storagePath, StorageInfo)

	c.Done()
}
