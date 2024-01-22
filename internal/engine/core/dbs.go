package core

import (
	"fmt"
	"os"
	"slices"
	"time"

	"github.com/Kwynto/GracefulDB/pkg/lib/ecowriter"
)

// Marks the database as deleted, but does not delete files.
func RemoveDB(name string) bool {
	// This function is complete
	var dbInfo tDBInfo

	folderName, ok := StorageInfo.DBs[name]
	if ok {
		if CheckFolderOrFile(LocalCoreSettings.Storage, folderName) {
			dbInfoPath := fmt.Sprintf("%s%s/%s", LocalCoreSettings.Storage, folderName, INFOFILE_DB)
			err := ecowriter.ReadJSON(dbInfoPath, &dbInfo)
			if err != nil {
				return false
			}
			dbInfo.LastUpdate = time.Now()
			dbInfo.Deleted = true
			err2 := ecowriter.WriteJSON(dbInfoPath, dbInfo)
			if err2 != nil {
				return false
			}

			StorageInfo.Removed = append(StorageInfo.Removed, folderName)
			delete(StorageInfo.DBs, name)
			storagePath := fmt.Sprintf("%s%s", LocalCoreSettings.Storage, INFOFILE_STORAGE)
			err3 := ecowriter.WriteJSON(storagePath, StorageInfo)
			if err3 != nil {
				return false
			}
		} else {
			return false
		}
	}

	return true
}

// Deletes the folder and database files, if DB was mark as 'removed'
func StrongRemoveDB(name string) bool {
	// This function is complete
	var dbInfo tDBInfo

	for indRange, folderName := range StorageInfo.Removed {
		if CheckFolderOrFile(LocalCoreSettings.Storage, folderName) {
			fullPath := fmt.Sprintf("%s%s", LocalCoreSettings.Storage, folderName)
			dbInfoPath := fmt.Sprintf("%s/%s", fullPath, INFOFILE_DB)
			err := ecowriter.ReadJSON(dbInfoPath, &dbInfo)
			if err != nil {
				return false
			}

			if dbInfo.Name == name {
				err := os.Remove(fullPath)
				if err != nil {
					return false
				}
				slices.Delete(StorageInfo.Removed, indRange, indRange+1)
				storagePath := fmt.Sprintf("%s%s", LocalCoreSettings.Storage, INFOFILE_STORAGE)
				ecowriter.WriteJSON(storagePath, StorageInfo)
				return true
			}
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
		Tables:     make(map[string]string),
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
