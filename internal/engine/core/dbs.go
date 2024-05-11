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
func RemoveDB(sNameDB string) bool {
	// This function is complete
	mxStorageBlock.Lock()
	defer mxStorageBlock.Unlock()

	stDBInfo, isOk := StStorageInfo.DBs[sNameDB]
	if isOk {
		stDBInfo.LastUpdate = time.Now()
		stDBInfo.Deleted = true

		StStorageInfo.Removed = append(StStorageInfo.Removed, stDBInfo)
		delete(StStorageInfo.DBs, sNameDB)
		delete(StStorageInfo.Access, sNameDB)
		StStorageInfo.Save()

		return stDBInfo.Save()
	}

	return false
}

// Deletes the folder and database files, if DB was mark as 'removed'
func StrongRemoveDB(sNameDB string) bool {
	// This function is complete
	mxStorageBlock.Lock()
	defer mxStorageBlock.Unlock()

	for iInd, stDBInfo := range StStorageInfo.Removed {
		if stDBInfo.Name == sNameDB {
			// dbPath := fmt.Sprintf("%s%s", LocalCoreSettings.Storage, dbInfo.Folder)
			sDBPath := filepath.Join(StLocalCoreSettings.Storage, stDBInfo.Folder)
			err := os.RemoveAll(sDBPath)
			if err != nil {
				return false
			}

			StStorageInfo.Removed = slices.Delete(StStorageInfo.Removed, iInd, iInd+1)

			return true
		}
	}

	return false
}

// Rename a database.
func RenameDB(sOldName, sNewName string, isSecure bool) bool {
	// This function is complete
	if isSecure && !vqlexp.MRegExpCollection["EntityName"].MatchString(sNewName) {
		return false
	}

	mxStorageBlock.Lock()
	defer mxStorageBlock.Unlock()

	stDBInfo, isOkDB := StStorageInfo.DBs[sOldName]
	stDBAccess, isOkAccess := StStorageInfo.Access[sOldName]

	if isOkDB && isOkAccess {
		stDBInfo.Name = sNewName
		stDBInfo.LastUpdate = time.Now()

		delete(StStorageInfo.DBs, sOldName)
		delete(StStorageInfo.Access, sOldName)

		StStorageInfo.DBs[sNewName] = stDBInfo
		StStorageInfo.Access[sNewName] = stDBAccess
		StStorageInfo.Save()

		return stDBInfo.Save()
	}

	return false
}

// Creating a new database.
func CreateDB(sNameDB string, sOwner string, isSecure bool) bool {
	// This function is complete
	if isSecure && !vqlexp.MRegExpCollection["EntityName"].MatchString(sNameDB) {
		return false
	}

	var sFolderDB string

	mxStorageBlock.Lock()
	defer mxStorageBlock.Unlock()

	_, isOk := StStorageInfo.DBs[sNameDB]
	if isOk {
		return false
	}

	for {
		sFolderDB = GenerateName()
		if !CheckFolder(StLocalCoreSettings.Storage, sFolderDB) {
			break
		}
	}

	// fullNameFolderDB := fmt.Sprintf("%s%s", LocalCoreSettings.Storage, folderDB)
	sFullNameFolderDB := filepath.Join(StLocalCoreSettings.Storage, sFolderDB)
	err := os.Mkdir(sFullNameFolderDB, 0666)
	if err != nil {
		return false
	}

	stDBInfo := TDBInfo{
		Name:       sNameDB,
		Folder:     sFolderDB,
		Tables:     make(map[string]TTableInfo),
		Removed:    make([]TTableInfo, 0),
		LastUpdate: time.Now(),
		Deleted:    false,
	}

	stDBAccess := gtypes.TAccess{
		Owner: sOwner,
		Flags: make(map[string]gtypes.TAccessFlags),
	}

	StStorageInfo.DBs[sNameDB] = stDBInfo
	StStorageInfo.Access[sNameDB] = stDBAccess
	StStorageInfo.Save()

	return stDBInfo.Save()
}
