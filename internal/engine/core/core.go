package core

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"regexp"
	"time"

	"github.com/Kwynto/GracefulDB/internal/config"
	"github.com/Kwynto/GracefulDB/pkg/lib/closer"
	"github.com/Kwynto/GracefulDB/pkg/lib/ecowriter"
)

const (
	INFOFILE_DB = "db.json"
	POSTFIX_ID  = "_id"
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
	Patronymic string                 `json:"patronymic"`
	Folder     string                 `json:"folder"`
	Parent     string                 `json:"parent"`
	Columns    map[string]tColumnInfo `json:"columns"`
	Removed    []tColumnInfo          `json:"removed"` // Removed columns
	Order      []string               `json:"order"`
	Count      uint64                 `json:"count"`
	LastUpdate time.Time              `json:"lastupdate"`
	Deleted    bool                   `json:"deleted"`
}

type tColumnInfo struct {
	Name          string               `json:"name"`
	Folder        string               `json:"folder"`
	Parents       string               `json:"parents"`
	BucketLog     uint8                `json:"blog"`
	BucketSize    int                  `json:"bsize"`
	OldRev        string               `json:"oldrev"`
	CurrentRev    string               `json:"currentrev"`
	Specification tColumnSpecification `json:"specification"`
	LastUpdate    time.Time            `json:"lastupdate"`
	Deleted       bool                 `json:"deleted"`
}

type tColumnSpecification struct {
	Default string `json:"default"`
	NotNull bool   `json:"notnull"`
	Unique  bool   `json:"unique"`
}

type tRegExpCollection map[string]*regexp.Regexp

func (r tRegExpCollection) CompileExp(name string, expr string) tRegExpCollection {
	// This method is completes
	re, err := regexp.Compile(expr)
	if err != nil {
		return r
	}
	r[name] = re

	return r
}

type tCoreFile struct {
	Descriptor *os.File
	Expire     time.Duration
}

type tCoreProcessing struct {
	FileDescriptors map[string]tCoreFile
}

type TState struct {
	CurrentDB string
}

var LocalCoreSettings tCoreSettings = tCoreSettings{
	Storage:    "./data/",
	BucketSize: 800,
	FreezeMode: false,
}

var RegExpCollection tRegExpCollection

var ParsingOrder = [...]string{
	"SearchSelect",
	"SearchInsert",
	"SearchUpdate",

	"SearchUse",
	"SearchAuth",

	"SearchDelete",
	"SearchTruncate",
	"SearchCommit",
	"SearchRollback",

	"SearchCreate",
	"SearchAlter",
	"SearchDrop",

	"SearchGrant",
	"SearchRevoke",
}

var StorageInfo tStorageInfo = tStorageInfo{
	// DBs:     make(map[string]tDBInfo),
	// Removed: make([]tDBInfo, 0),
}

var States map[string]TState // ticket -> tState

var CoreProcessing tCoreProcessing

func LoadLocalCoreSettings(cfg *config.Config) tCoreSettings {
	// This function is complete
	return tCoreSettings{
		Storage:    cfg.CoreSettings.Storage,
		BucketSize: cfg.CoreSettings.BucketSize,
		FreezeMode: cfg.CoreSettings.FreezeMode,
	}
}

func CompileRegExpCollection() tRegExpCollection {
	// -
	var recol tRegExpCollection = make(tRegExpCollection)
	// recol = recol.CompileExp("LineBreak", `(?m)\n`)
	// recol = recol.CompileExp("HeadCleaner", `(?m)^\s*\n*\s*`)
	// recol = recol.CompileExp("AnyCommand", `(?m)^[a-zA-Z].*;\s*`)
	recol = recol.CompileExp("EntityName", `(?m)^[a-zA-Z][a-zA-Z0-9_-]*$`)
	recol = recol.CompileExp("QuotationMarks", `(?m)[\'\"]`)
	recol = recol.CompileExp("SpecQuotationMark", "(?m)[`]")

	// DDL TODO: Разработать шаблоны
	recol = recol.CompileExp("SearchCreate", `(?m)^;`)
	recol = recol.CompileExp("SearchAlter", `(?m)^;`)
	recol = recol.CompileExp("SearchDrop", `(?m)^;`)
	// DML TODO: Разработать шаблоны
	recol = recol.CompileExp("SearchSelect", `(?m)^;`)
	recol = recol.CompileExp("SearchInsert", `(?m)^;`)
	recol = recol.CompileExp("SearchUpdate", `(?m)^;`)
	recol = recol.CompileExp("SearchDelete", `(?m)^;`)
	recol = recol.CompileExp("SearchTruncate", `(?m)^;`)
	recol = recol.CompileExp("SearchCommit", `(?m)^;`)
	recol = recol.CompileExp("SearchRollback", `(?m)^;`)
	// DCL
	recol = recol.CompileExp("SearchUse", `(?m)^[uU][sS][eE]\s*[a-zA-Z][a-zA-Z0-1]+\s*`)
	recol = recol.CompileExp("SearchGrant", `(?m)^[gG][rR][aA][nN][tT].*`)
	recol = recol.CompileExp("SearchRevoke", `(?m)^[rR][eE][vV][oO][kK][eE].*`)

	recol = recol.CompileExp("SearchAuth", `(?m)^[aA][uU][tT][hH].+`)
	// recol = recol.CompileExp("Auth", `(?m)^[aA][uU][tT][hH]`)
	recol = recol.CompileExp("Login", `(?m)[lL][oO][gG][iI][nN]\s+\S+(\s+|$)`)
	recol = recol.CompileExp("LoginWord", `(?m)[lL][oO][gG][iI][nN]`)
	recol = recol.CompileExp("Password", `(?m)[pP][aA][sS][sS][wW][oO][rR][dD]\s+\S+(\s+|$)`)
	recol = recol.CompileExp("PasswordWord", `(?m)[pP][aA][sS][sS][wW][oO][rR][dD]`)

	return recol
}

func Start(cfg *config.Config) {
	// -
	LocalCoreSettings = LoadLocalCoreSettings(cfg)
	RegExpCollection = CompileRegExpCollection()

	if !StorageInfo.Load() {
		slog.Error("Storage activation error !!!")
	}

	States = make(map[string]TState)

	slog.Info("The core of the DBMS was started.")

	// This block needs to delete
	// fmt.Println(CreateDB("ExampleDB"))
	// fmt.Println(CreateTable("ExampleDB", "ExampleTable"))
	// fmt.Println(CreateColumn("ExampleDB", "ExampleTable", "example"))
	// fmt.Println(StorageInfo)
}

func Shutdown(ctx context.Context, c *closer.Closer) {
	// -
	c.Done()
}
