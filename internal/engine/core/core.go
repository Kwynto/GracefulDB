package core

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/Kwynto/GracefulDB/internal/config"
	"github.com/Kwynto/GracefulDB/pkg/lib/closer"
)

type tCoreSettings struct {
	Storage    string
	BucketSize int
	FreezeMode bool
}

type tDBInfo struct {
	Name   string   `json:"name"`
	Tables []string `json:"tables"`
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

func RemoveDB(name string) bool {
	fullName := fmt.Sprintf("%s%s", LocalCoreSettings.Storage, name)
	err := os.Remove(fullName)

	return err == nil
}

func CreateDB(name string) bool {
	fullName := fmt.Sprintf("%s%s", LocalCoreSettings.Storage, name)
	err := os.Mkdir(fullName, 0666)
	if err != nil {
		return false
	}

	dbInfoPath := fmt.Sprintf("%s/%s", fullName, "info.json")
	fo, err := os.OpenFile(dbInfoPath, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	if err != nil {
		return false
	}
	defer fo.Close()

	dbInfo := tDBInfo{
		Name:   name,
		Tables: []string{},
	}
	bytesDBInfo, err := json.Marshal(dbInfo)
	if err != nil {
		return false
	}
	fo.Write(bytesDBInfo)

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
	slog.Info("The core of the DBMS was started.")
}

func Shutdown(ctx context.Context, c *closer.Closer) {
	// -
	c.Done()
}
