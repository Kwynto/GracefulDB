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
	// INFOFILE_STORAGE = "storage.json"
	INFOFILE_DB = "db.json"
	// INFOFILE_TABLE   = "table.json"
	// INFOFILE_COLUMN  = "column.json"
)

type tCoreSettings struct {
	Storage    string
	BucketSize int
	FreezeMode bool
}

type tStorageInfo struct {
	DBs     map[string]tDBInfo `json:"dbs"`     // [name db] tDBInfo
	Removed []tDBInfo          `json:"removed"` // Removed databases
}

func (s *tStorageInfo) Load() bool {
	// This method is complete
	var dbInfo tDBInfo

	s.DBs = make(map[string]tDBInfo)
	s.Removed = make([]tDBInfo, 0)

	files, err := os.ReadDir(LocalCoreSettings.Storage)
	if err != nil {
		return false
	}

	for _, file := range files {
		if file.IsDir() {
			nameDir := file.Name()
			dbInfoFile := fmt.Sprintf("%s%s/%s", LocalCoreSettings.Storage, nameDir, INFOFILE_DB)
			err := ecowriter.ReadJSON(dbInfoFile, &dbInfo)
			if err == nil {
				if dbInfo.Deleted {
					s.Removed = append(s.Removed, dbInfo)
				} else {
					s.DBs[dbInfo.Name] = dbInfo
				}
			}
		}
	}
	return true
}

type tDBInfo struct {
	Name       string                `json:"name"`
	Folder     string                `json:"folder"`
	Tables     map[string]tTableInfo `json:"tables"`
	Removed    []tTableInfo          `json:"removed"` // Removed tables
	LastUpdate time.Time             `json:"lastupdate"`
	Deleted    bool                  `json:"deleted"`
}

// Saving the database structure.
func (d tDBInfo) Save() bool {
	// This method is complete
	path := fmt.Sprintf("%s%s/%s", LocalCoreSettings.Storage, d.Folder, INFOFILE_DB)
	err := ecowriter.WriteJSON(path, d)
	return err == nil
}

type tTableInfo struct {
	Name       string                 `json:"name"`
	Folder     string                 `json:"folder"`
	Parent     string                 `json:"parent"`
	Columns    map[string]tColumnInfo `json:"columns"`
	Removed    []tColumnInfo          `json:"removed"` // Removed columns
	Order      []string               `json:"order"`
	LastUpdate time.Time              `json:"lastupdate"`
	Deleted    bool                   `json:"deleted"`
}

type tColumnInfo struct {
	Name       string    `json:"name"`
	Folder     string    `json:"folder"`
	Parents    string    `json:"parents"`
	BucketLog  uint8     `json:"blog"`
	BucketSize int       `json:"bsize"`
	OldRev     string    `json:"oldrev"`
	CurrentRev string    `json:"currentrev"`
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
	// DBs:     make(map[string]tDBInfo),
	// Removed: make([]tDBInfo, 0),
}

func LoadLocalCoreSettings(cfg *config.Config) tCoreSettings {
	return tCoreSettings{
		Storage:    cfg.CoreSettings.Storage,
		BucketSize: cfg.CoreSettings.BucketSize,
		FreezeMode: cfg.CoreSettings.FreezeMode,
	}
}

func Start(cfg *config.Config) {
	LocalCoreSettings = LoadLocalCoreSettings(cfg)

	if !StorageInfo.Load() {
		slog.Error("Storage activation error !!!")
	}

	slog.Info("The core of the DBMS was started.")
}

func Shutdown(ctx context.Context, c *closer.Closer) {
	// -
	c.Done()
}
