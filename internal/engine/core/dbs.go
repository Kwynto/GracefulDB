package core

import (
	"os"
	"path/filepath"
	"slices"
	"time"

	"github.com/Kwynto/GracefulDB/internal/engine/basicsystem/gtypes"
	"github.com/Kwynto/GracefulDB/internal/engine/basicsystem/vqlexp"
)

// Marks the database as deleted, but does not delete files.
func RemoveDB(nameDB string) bool {
	// This function is complete
	mxStorageBlock.Lock()
	defer mxStorageBlock.Unlock()

	dbInfo, ok := StStorageInfo.DBs[nameDB]
	if ok {
		dbInfo.LastUpdate = time.Now()
		dbInfo.Deleted = true

		StStorageInfo.Removed = append(StStorageInfo.Removed, dbInfo)
		delete(StStorageInfo.DBs, nameDB)
		delete(StStorageInfo.Access, nameDB)
		StStorageInfo.Save()

		return dbInfo.Save()
	}

	return false
}

// Deletes the folder and database files, if DB was mark as 'removed'
func StrongRemoveDB(nameDB string) bool {
	// This function is complete
	mxStorageBlock.Lock()
	defer mxStorageBlock.Unlock()

	for indRange, dbInfo := range StStorageInfo.Removed {
		if dbInfo.Name == nameDB {
			// dbPath := fmt.Sprintf("%s%s", LocalCoreSettings.Storage, dbInfo.Folder)
			dbPath := filepath.Join(LocalCoreSettings.Storage, dbInfo.Folder)
			err := os.RemoveAll(dbPath)
			if err != nil {
				return false
			}

			StStorageInfo.Removed = slices.Delete(StStorageInfo.Removed, indRange, indRange+1)

			return true
		}
	}

	return false
}

// Rename a database.
func RenameDB(oldName, newName string, secure bool) bool {
	// This function is complete
	if secure && !vqlexp.MRegExpCollection["EntityName"].MatchString(newName) {
		return false
	}

	mxStorageBlock.Lock()
	defer mxStorageBlock.Unlock()

	dbInfo, okDB := StStorageInfo.DBs[oldName]
	dbAccess, okAccess := StStorageInfo.Access[oldName]

	if okDB && okAccess {
		dbInfo.Name = newName
		dbInfo.LastUpdate = time.Now()

		delete(StStorageInfo.DBs, oldName)
		delete(StStorageInfo.Access, oldName)

		StStorageInfo.DBs[newName] = dbInfo
		StStorageInfo.Access[newName] = dbAccess
		StStorageInfo.Save()

		return dbInfo.Save()
	}

	return false
}

// Creating a new database.
func CreateDB(nameDB string, owner string, secure bool) bool {
	// This function is complete
	if secure && !vqlexp.MRegExpCollection["EntityName"].MatchString(nameDB) {
		return false
	}

	var folderDB string

	mxStorageBlock.Lock()
	defer mxStorageBlock.Unlock()

	_, ok := StStorageInfo.DBs[nameDB]
	if ok {
		return false
	}

	for {
		folderDB = GenerateName()
		if !CheckFolder(LocalCoreSettings.Storage, folderDB) {
			break
		}
	}

	// fullNameFolderDB := fmt.Sprintf("%s%s", LocalCoreSettings.Storage, folderDB)
	fullNameFolderDB := filepath.Join(LocalCoreSettings.Storage, folderDB)
	err := os.Mkdir(fullNameFolderDB, 0666)
	if err != nil {
		return false
	}

	dbInfo := TDBInfo{
		Name:       nameDB,
		Folder:     folderDB,
		Tables:     make(map[string]TTableInfo),
		Removed:    make([]TTableInfo, 0),
		LastUpdate: time.Now(),
		Deleted:    false,
	}

	dbAccess := gtypes.TAccess{
		Owner: owner,
		Flags: make(map[string]gtypes.TAccessFlags),
	}

	StStorageInfo.DBs[nameDB] = dbInfo
	StStorageInfo.Access[nameDB] = dbAccess
	StStorageInfo.Save()

	return dbInfo.Save()
}
