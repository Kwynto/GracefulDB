package core

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"sync"
	"time"

	"github.com/Kwynto/GracefulDB/internal/config"
	"github.com/Kwynto/GracefulDB/internal/engine/basicsystem/gtypes"
	"github.com/Kwynto/GracefulDB/pkg/lib/closer"
	"github.com/Kwynto/GracefulDB/pkg/lib/ecowriter"
)

const (
	INFOFILE_DB      = "db.json"
	INFOFILE_STORAGE = "storage.json"
	POSTFIX_ID       = "_id"
)

var storageBlock sync.RWMutex

type tCoreSettings struct {
	Storage      string
	BucketSize   int64
	FriendlyMode bool
}

type tStorageInfo struct {
	DBs     map[string]TDBInfo        `json:"dbs"`     // [name db] tDBInfo
	Removed []TDBInfo                 `json:"removed"` // Removed databases
	Access  map[string]gtypes.TAccess `json:"access"`  // [name db] - TAccess
}

func GetDBInfo(nameDB string) (TDBInfo, bool) {
	// This function is complete
	storageBlock.RLock()
	defer storageBlock.RUnlock()

	info, ok := StorageInfo.DBs[nameDB]
	if !ok {
		return TDBInfo{}, false
	}

	return info, true
}

func GetDBAccess(nameDB string) (gtypes.TAccess, bool) {
	// This function is complete
	storageBlock.RLock()
	defer storageBlock.RUnlock()

	access, ok := StorageInfo.Access[nameDB]
	if !ok {
		return gtypes.TAccess{}, false
	}

	return access, true
}

func SetAccessFlags(db, user string, flags gtypes.TAccessFlags) {
	// This procedure is complete
	storageBlock.Lock()
	StorageInfo.Access[db].Flags[user] = flags
	StorageInfo.Save()
	storageBlock.Unlock()
}

func (s *tStorageInfo) Load() bool {
	// This method is complete
	storageBlock.Lock()
	defer storageBlock.Unlock()

	var dbInfo TDBInfo

	s.DBs = make(map[string]TDBInfo)
	s.Removed = make([]TDBInfo, 0)
	s.Access = make(map[string]gtypes.TAccess)

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

	infoStorageFile := fmt.Sprintf("%s%s", LocalCoreSettings.Storage, INFOFILE_STORAGE)
	errR := ecowriter.ReadJSON(infoStorageFile, &s.Access)
	if errR != nil {
		s.Access = make(map[string]gtypes.TAccess)
		err := ecowriter.WriteJSON(infoStorageFile, s.Access)
		if err != nil {
			return false
		}
	}

	return true
}

func (s *tStorageInfo) Save() bool {
	// This method is complete
	// Don't use mutex
	infoStorageFile := fmt.Sprintf("%s%s", LocalCoreSettings.Storage, INFOFILE_STORAGE)
	return ecowriter.WriteJSON(infoStorageFile, s.Access) == nil
}

type TDBInfo struct {
	Name       string                `json:"name"`
	Folder     string                `json:"folder"`
	Tables     map[string]TTableInfo `json:"tables"`
	Removed    []TTableInfo          `json:"removed"` // Removed tables
	LastUpdate time.Time             `json:"lastupdate"`
	Deleted    bool                  `json:"deleted"`
}

// Saving the database structure.
func (d TDBInfo) Save() bool {
	// This method is complete
	// Don't use mutex
	path := fmt.Sprintf("%s%s/%s", LocalCoreSettings.Storage, d.Folder, INFOFILE_DB)
	return ecowriter.WriteJSON(path, d) == nil
}

type TTableInfo struct {
	Name       string                 `json:"name"`
	Patronymic string                 `json:"patronymic"`
	Folder     string                 `json:"folder"`
	Parent     string                 `json:"parent"`
	Columns    map[string]TColumnInfo `json:"columns"`
	Removed    []TColumnInfo          `json:"removed"` // Removed columns
	Order      []string               `json:"order"`
	Count      uint64                 `json:"count"`
	LastUpdate time.Time              `json:"lastupdate"`
	Deleted    bool                   `json:"deleted"`
}

type TColumnInfo struct {
	Name          string                      `json:"name"`
	OldName       string                      `json:"oldname"` // only for core
	Folder        string                      `json:"folder"`
	Parents       string                      `json:"parents"`
	BucketLog     uint8                       `json:"blog"`
	BucketSize    int64                       `json:"bsize"`
	OldRev        string                      `json:"oldrev"`
	CurrentRev    string                      `json:"currentrev"`
	Specification gtypes.TColumnSpecification `json:"specification"`
	LastUpdate    time.Time                   `json:"lastupdate"`
	Deleted       bool                        `json:"deleted"`
}

type TState struct {
	CurrentDB string
}

var LocalCoreSettings tCoreSettings = tCoreSettings{
	Storage:      "./data/",
	BucketSize:   800,
	FriendlyMode: true,
}

var StorageInfo tStorageInfo = tStorageInfo{
	// DBs:     make(map[string]tDBInfo),
	// Removed: make([]tDBInfo, 0),
}

var States map[string]TState // ticket -> tState

func LoadLocalCoreSettings(cfg *config.Config) tCoreSettings {
	// This function is complete
	return tCoreSettings{
		Storage:      cfg.CoreSettings.Storage,
		BucketSize:   cfg.CoreSettings.BucketSize,
		FriendlyMode: cfg.CoreSettings.FriendlyMode,
	}
}

func Start(cfg *config.Config) {
	// -
	LocalCoreSettings = LoadLocalCoreSettings(cfg)
	// RegExpCollection = CompileRegExpCollection()

	if !StorageInfo.Load() {
		slog.Error("Storage activation error !!!")
	}

	States = make(map[string]TState)

	go WriteBufferService()

	slog.Info("The core of the DBMS was started.")
}

func Shutdown(ctx context.Context, c *closer.Closer) {
	// -
	if !StorageInfo.Save() {
		c.AddMsg("Failure to save access rights !!!")
	}
	signalSD <- struct{}{}
	c.Done()
}
