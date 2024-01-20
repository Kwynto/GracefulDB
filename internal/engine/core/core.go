package core

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/Kwynto/GracefulDB/internal/config"
	"github.com/Kwynto/GracefulDB/pkg/lib/closer"
	"github.com/Kwynto/GracefulDB/pkg/lib/ecowriter"
)

const (
	NAME_STORAGE_INFO = "storage.json"
	NAME_DB_INFO      = "db.json"
	NAME_TABLE_INFO   = "table.json"
)

type tCoreSettings struct {
	Storage    string
	BucketSize int
	FreezeMode bool
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

// Marks the database as deleted, but does not delete files.
func RemoveDB(name string) bool {
	// This function is complete
	var dbInfo tDBInfo

	dbInfoPath := fmt.Sprintf("%s%s/%s", LocalCoreSettings.Storage, name, NAME_DB_INFO)

	err := ecowriter.ReadJSON(dbInfoPath, &dbInfo)
	if err != nil {
		return false
	}

	dbInfo.LastUpdate = time.Now()
	dbInfo.Deleted = true
	err2 := ecowriter.WriteJSON(dbInfoPath, dbInfo)

	return err2 == nil
}

// Deletes the folder and database files.
func StrongRemoveDB(name string) bool {
	// This function is complete
	fullName := fmt.Sprintf("%s%s", LocalCoreSettings.Storage, name)
	err := os.Remove(fullName)

	return err == nil
}

// Creating a new database.
func CreateDB(name string) bool {
	// This function is complete
	fullName := fmt.Sprintf("%s%s", LocalCoreSettings.Storage, name)
	err := os.Mkdir(fullName, 0666)
	if err != nil {
		return false
	}

	dbInfoPath := fmt.Sprintf("%s/%s", fullName, NAME_DB_INFO)

	dbInfo := tDBInfo{
		Name:       name,
		Tables:     []string{},
		LastUpdate: time.Now(),
		Deleted:    false,
	}

	err2 := ecowriter.WriteJSON(dbInfoPath, dbInfo)
	return err2 == nil
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
	slog.Info("The core of the DBMS was started.")
}

func Shutdown(ctx context.Context, c *closer.Closer) {
	// -
	c.Done()
}
