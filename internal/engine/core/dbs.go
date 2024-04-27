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
	storageBlock.Lock()
	defer storageBlock.Unlock()

	dbInfo, ok := StorageInfo.DBs[nameDB]
	if ok {
		dbInfo.LastUpdate = time.Now()
		dbInfo.Deleted = true

		StorageInfo.Removed = append(StorageInfo.Removed, dbInfo)
		delete(StorageInfo.DBs, nameDB)
		delete(StorageInfo.Access, nameDB)
		StorageInfo.Save()

		return dbInfo.Save()
	}

	return false
}

// Deletes the folder and database files, if DB was mark as 'removed'
func StrongRemoveDB(nameDB string) bool {
	// This function is complete
	storageBlock.Lock()
	defer storageBlock.Unlock()

	for indRange, dbInfo := range StorageInfo.Removed {
		if dbInfo.Name == nameDB {
			// dbPath := fmt.Sprintf("%s%s", LocalCoreSettings.Storage, dbInfo.Folder)
			dbPath := filepath.Join(LocalCoreSettings.Storage, dbInfo.Folder)
			err := os.RemoveAll(dbPath)
			if err != nil {
				return false
			}

			StorageInfo.Removed = slices.Delete(StorageInfo.Removed, indRange, indRange+1)

			return true
		}
	}

	return false
}

// Rename a database.
func RenameDB(oldName, newName string, secure bool) bool {
	// This function is complete
	if secure && !vqlexp.RegExpCollection["EntityName"].MatchString(newName) {
		return false
	}

	storageBlock.Lock()
	defer storageBlock.Unlock()

	dbInfo, okDB := StorageInfo.DBs[oldName]
	dbAccess, okAccess := StorageInfo.Access[oldName]

	if okDB && okAccess {
		dbInfo.Name = newName
		dbInfo.LastUpdate = time.Now()

		delete(StorageInfo.DBs, oldName)
		delete(StorageInfo.Access, oldName)

		StorageInfo.DBs[newName] = dbInfo
		StorageInfo.Access[newName] = dbAccess
		StorageInfo.Save()

		return dbInfo.Save()
	}

	return false
}

// Creating a new database.
func CreateDB(nameDB string, owner string, secure bool) bool {
	// This function is complete
	if secure && !vqlexp.RegExpCollection["EntityName"].MatchString(nameDB) {
		return false
	}

	var folderDB string

	storageBlock.Lock()
	defer storageBlock.Unlock()

	_, ok := StorageInfo.DBs[nameDB]
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

	StorageInfo.DBs[nameDB] = dbInfo
	StorageInfo.Access[nameDB] = dbAccess
	StorageInfo.Save()

	return dbInfo.Save()
}
