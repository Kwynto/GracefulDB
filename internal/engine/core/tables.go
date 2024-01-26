package core

import (
	"fmt"
	"os"
	"slices"
	"time"
)

// Marks the table as deleted, but does not delete files.
func RemoveTable(nameDB, nameTable string) bool {
	// This function is complete
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

// Creating a new table.
func CreateTable(nameDB, nameTable string) bool {
	// This function is complete
	var folderName string

	dbInfo, ok := StorageInfo.DBs[nameDB]
	if !ok {
		return false
	}

	if !CheckFolderOrFile(LocalCoreSettings.Storage, dbInfo.Folder) {
		return false
	}

	pathDB := fmt.Sprintf("%s%s/", LocalCoreSettings.Storage, dbInfo.Folder)

	for {
		folderName = GenerateName()
		if !CheckFolderOrFile(pathDB, folderName) {
			break
		}
	}

	fullTableName := fmt.Sprintf("%s%s", pathDB, folderName)
	err := os.Mkdir(fullTableName, 0666)
	if err != nil {
		return false
	}

	tableInfo := tTableInfo{
		Name:       nameTable,
		Folder:     folderName,
		Parent:     dbInfo.Folder,
		Columns:    make(map[string]tColumnInfo),
		Removed:    make([]tColumnInfo, 0),
		Order:      make([]string, 0),
		LastUpdate: time.Now(),
		Deleted:    false,
	}

	dbInfo.Tables[nameTable] = tableInfo
	dbInfo.LastUpdate = time.Now()
	StorageInfo.DBs[nameDB] = dbInfo

	return dbInfo.Save()
}