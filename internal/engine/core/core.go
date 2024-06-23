package core

import (
	"context"
	"log/slog"
	"os"
	"path/filepath"
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

var mxStorageBlock sync.RWMutex

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

func GetDBInfo(sNameDB string) (TDBInfo, bool) {
	// This function is complete
	mxStorageBlock.RLock()
	defer mxStorageBlock.RUnlock()

	stDBInfo, isOk := StStorageInfo.DBs[sNameDB]
	if !isOk {
		return TDBInfo{}, false
	}

	return stDBInfo, true
}

func GetDBAccess(sNameDB string) (gtypes.TAccess, bool) {
	// This function is complete
	mxStorageBlock.RLock()
	defer mxStorageBlock.RUnlock()

	stAccess, isOk := StStorageInfo.Access[sNameDB]
	if !isOk {
		return gtypes.TAccess{}, false
	}

	return stAccess, true
}

func SetAccessFlags(sDB, sUser string, stFlags gtypes.TAccessFlags) {
	// This procedure is complete
	mxStorageBlock.Lock()
	StStorageInfo.Access[sDB].Flags[sUser] = stFlags
	StStorageInfo.Save()
	mxStorageBlock.Unlock()
}

func (s *tStorageInfo) Load() bool {
	// This method is complete
	mxStorageBlock.Lock()
	defer mxStorageBlock.Unlock()

	var stDBInfo TDBInfo

	s.DBs = make(map[string]TDBInfo)
	s.Removed = make([]TDBInfo, 0)
	s.Access = make(map[string]gtypes.TAccess)

	slFiles, err := os.ReadDir(StLocalCoreSettings.Storage)
	if err != nil {
		slog.Error(err.Error())
		return false
	}

	for _, file := range slFiles {
		if file.IsDir() {
			sNameDir := file.Name()
			sDBInfoFile := filepath.Join(StLocalCoreSettings.Storage, sNameDir, INFOFILE_DB)
			err := ecowriter.ReadJSON(sDBInfoFile, &stDBInfo)
			if err == nil {
				if stDBInfo.Deleted {
					s.Removed = append(s.Removed, stDBInfo)
				} else {
					s.DBs[stDBInfo.Name] = stDBInfo
				}
			}
		}
	}

	sInfoStorageFile := filepath.Join(StLocalCoreSettings.Storage, INFOFILE_STORAGE)
	errR := ecowriter.ReadJSON(sInfoStorageFile, &s.Access)
	if errR != nil {
		s.Access = make(map[string]gtypes.TAccess)
		err := ecowriter.WriteJSON(sInfoStorageFile, s.Access)
		if err != nil {
			return false
		}
	}

	return true
}

func (s *tStorageInfo) Save() bool {
	// This method is complete
	// Don't use mutex
	sInfoStorageFile := filepath.Join(StLocalCoreSettings.Storage, INFOFILE_STORAGE)
	return ecowriter.WriteJSON(sInfoStorageFile, s.Access) == nil
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
	sPath := filepath.Join(StLocalCoreSettings.Storage, d.Folder, INFOFILE_DB)
	return ecowriter.WriteJSON(sPath, d) == nil
}

type TTableInfo struct {
	Name       string                 `json:"name"`
	Patronymic string                 `json:"patronymic"`
	Folder     string                 `json:"folder"`
	Parent     string                 `json:"parent"`
	Columns    map[string]TColumnInfo `json:"columns"`
	Removed    []TColumnInfo          `json:"removed"` // Removed columns
	Order      []string               `json:"order"`
	BucketLog  uint8                  `json:"blog"`
	BucketSize int64                  `json:"bsize"`
	OldRev     string                 `json:"oldrev"`
	CurrentRev string                 `json:"currentrev"`
	Count      uint64                 `json:"count"`
	LastUpdate time.Time              `json:"lastupdate"`
	Deleted    bool                   `json:"deleted"`
}

type TColumnInfo struct {
	Name          string                      `json:"name"`
	OldName       string                      `json:"oldname"` // only for core
	Folder        string                      `json:"folder"`
	Parents       string                      `json:"parents"`
	Specification gtypes.TColumnSpecification `json:"specification"`
	LastUpdate    time.Time                   `json:"lastupdate"`
	Deleted       bool                        `json:"deleted"`
}

type TState struct {
	CurrentDB string
}

var StLocalCoreSettings tCoreSettings = tCoreSettings{
	Storage:      "./data",
	BucketSize:   800,
	FriendlyMode: true,
}

var StStorageInfo tStorageInfo = tStorageInfo{
	// DBs:     make(map[string]tDBInfo),
	// Removed: make([]tDBInfo, 0),
}

var MStates map[string]TState // ticket -> tState

func LoadLocalCoreSettings(cfg *config.TConfig) tCoreSettings {
	// This function is complete
	return tCoreSettings{
		Storage:      cfg.CoreSettings.Storage,
		BucketSize:   cfg.CoreSettings.BucketSize,
		FriendlyMode: cfg.CoreSettings.FriendlyMode,
	}
}

func Start(cfg *config.TConfig) {
	// -
	StLocalCoreSettings = LoadLocalCoreSettings(cfg)

	if !StStorageInfo.Load() {
		slog.Error("Storage activation error !!!")
	}

	MStates = make(map[string]TState)

	go WriteBufferService()

	slog.Info("The core of the DBMS was started.")
}

func Shutdown(ctx context.Context, c *closer.TCloser) {
	// -
	if !StStorageInfo.Save() {
		c.AddMsg("Failure to save access rights !!!")
	}
	chSignalShutdown <- struct{}{}
	c.Done()
}
