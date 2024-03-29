package core

import (
	"fmt"
	"os"
	"slices"
	"time"

	"github.com/Kwynto/GracefulDB/internal/engine/basicsystem/vqlexp"
)

// Marks the table as deleted, but does not delete files.
func RemoveTable(nameDB, nameTable string) bool {
	// This function is complete
	storageBlock.Lock()
	defer storageBlock.Unlock()

	dbInfo, ok := StorageInfo.DBs[nameDB]
	if !ok {
		return false
	}

	tableInfo, ok := dbInfo.Tables[nameTable]
	if !ok {
		return false
	}

	tNow := time.Now()

	tableInfo.LastUpdate = tNow
	tableInfo.Deleted = true

	dbInfo.Removed = append(dbInfo.Removed, tableInfo)
	delete(dbInfo.Tables, nameTable)
	dbInfo.LastUpdate = tNow

	StorageInfo.DBs[nameDB] = dbInfo

	return dbInfo.Save()
}

// Deletes the folder and table files, if table was mark as 'removed'
func StrongRemoveTable(nameDB, nameTable string) bool {
	// This function is complete
	storageBlock.Lock()
	defer storageBlock.Unlock()

	dbInfo, ok := StorageInfo.DBs[nameDB]
	if !ok {
		return false
	}

	for indRange, tableInfo := range dbInfo.Removed {
		if tableInfo.Name == nameTable {
			tablePath := fmt.Sprintf("%s%s/%s", LocalCoreSettings.Storage, tableInfo.Parent, tableInfo.Folder)
			err := os.RemoveAll(tablePath)
			if err != nil {
				return false
			}

			dbInfo.Removed = slices.Delete(dbInfo.Removed, indRange, indRange+1)
			dbInfo.LastUpdate = time.Now()
			StorageInfo.DBs[nameDB] = dbInfo

			return dbInfo.Save()
		}
	}

	return false
}

// Rename a table.
func RenameTable(nameDB, oldNameTable, newNameTable string, secure bool) bool {
	// This function is complete
	if secure && !vqlexp.RegExpCollection["EntityName"].MatchString(newNameTable) {
		return false
	}

	storageBlock.Lock()
	defer storageBlock.Unlock()

	dbInfo, okDB := StorageInfo.DBs[nameDB]
	if okDB {
		tableInfo, okTable := dbInfo.Tables[oldNameTable]
		if !okTable {
			return false
		}

		tNow := time.Now()

		tableInfo.Name = newNameTable
		tableInfo.LastUpdate = tNow

		delete(dbInfo.Tables, oldNameTable)
		dbInfo.Tables[newNameTable] = tableInfo
		dbInfo.LastUpdate = tNow

		StorageInfo.DBs[nameDB] = dbInfo

		return dbInfo.Save()
	}

	return false
}

func TruncateTable(nameDB, nameTable string) bool {
	storageBlock.Lock()
	defer storageBlock.Unlock()

	// TODO: написать очистку таблицы

	return true
}

// Creating a new table.
func CreateTable(nameDB, nameTable string, secure bool) bool {
	// This function is complete
	if secure && !vqlexp.RegExpCollection["EntityName"].MatchString(nameTable) {
		return false
	}

	var folderName string

	storageBlock.Lock()
	defer storageBlock.Unlock()

	dbInfo, ok := StorageInfo.DBs[nameDB]
	if !ok {
		return false
	}

	pathDB := fmt.Sprintf("%s%s/", LocalCoreSettings.Storage, dbInfo.Folder)

	for {
		folderName = GenerateName()
		if !CheckFolder(pathDB, folderName) {
			break
		}
	}

	fullTableName := fmt.Sprintf("%s%s", pathDB, folderName)
	err := os.Mkdir(fullTableName, 0666)
	if err != nil {
		return false
	}

	tableInfo := TTableInfo{
		Name:       nameTable,
		Patronymic: nameDB,
		Folder:     folderName,
		Parent:     dbInfo.Folder,
		Columns:    make(map[string]TColumnInfo),
		Removed:    make([]TColumnInfo, 0),
		Order:      make([]string, 0),
		Count:      0,
		LastUpdate: time.Now(),
		Deleted:    false,
	}

	dbInfo.Tables[nameTable] = tableInfo
	dbInfo.LastUpdate = time.Now()
	StorageInfo.DBs[nameDB] = dbInfo

	return dbInfo.Save()
}
