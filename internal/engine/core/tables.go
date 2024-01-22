package core

import (
	"fmt"
	"os"
	"slices"
	"time"

	"github.com/Kwynto/GracefulDB/pkg/lib/ecowriter"
)

// Marks the table as deleted, but does not delete files.
func RemoveTable(nameDB, nameTable string) bool {
	// This function is complete
	var dbInfo tDBInfo
	var tableInfo tTableInfo

	folderName, ok := StorageInfo.DBs[nameDB]
	if ok {
		if CheckFolderOrFile(LocalCoreSettings.Storage, folderName) {
			dbInfoPath := fmt.Sprintf("%s%s/%s", LocalCoreSettings.Storage, folderName, INFOFILE_DB)
			err := ecowriter.ReadJSON(dbInfoPath, &dbInfo)
			if err != nil {
				return false
			}

			folderTName, ok2 := dbInfo.Tables[nameTable]
			if ok2 {
				if CheckFolderOrFile(fmt.Sprintf("%s%s/", LocalCoreSettings.Storage, folderName), folderTName) {
					tableInfoPath := fmt.Sprintf("%s%s/%s/%s", LocalCoreSettings.Storage, folderName, folderTName, INFOFILE_TABLE)
					err := ecowriter.ReadJSON(tableInfoPath, &tableInfo)
					if err != nil {
						return false
					}

					tNow := time.Now()

					tableInfo.LastUpdate = tNow
					tableInfo.Deleted = true

					dbInfo.Removed = append(dbInfo.Removed, folderTName)
					delete(dbInfo.Tables, nameTable)
					dbInfo.LastUpdate = tNow

					err2 := ecowriter.WriteJSON(tableInfoPath, tableInfo)
					if err2 != nil {
						return false
					}

					err3 := ecowriter.WriteJSON(dbInfoPath, dbInfo)
					if err3 != nil {
						return false
					}
				}
			}
		} else {
			return false
		}
	}

	return true
}

// Deletes the folder and table files, if table was mark as 'removed'
func StrongRemoveTable(nameDB, nameTable string) bool {
	// This function is complete
	var dbInfo tDBInfo
	var tableInfo tTableInfo

	folderName, ok := StorageInfo.DBs[nameDB]
	if !ok {
		return false
	}

	if CheckFolderOrFile(LocalCoreSettings.Storage, folderName) {
		fullPath := fmt.Sprintf("%s%s", LocalCoreSettings.Storage, folderName)
		dbInfoPath := fmt.Sprintf("%s/%s", fullPath, INFOFILE_DB)
		err := ecowriter.ReadJSON(dbInfoPath, &dbInfo)
		if err != nil {
			return false
		}

		for indRange, folderTName := range dbInfo.Removed {
			tablePath := fmt.Sprintf("%s/%s", fullPath, folderTName)
			tableInfoPath := fmt.Sprintf("%s/%s", tablePath, INFOFILE_TABLE)
			err := ecowriter.ReadJSON(tableInfoPath, &tableInfo)
			if err != nil {
				return false
			}
			if tableInfo.Name == nameTable {
				err := os.Remove(tablePath)
				if err != nil {
					return false
				}

				slices.Delete(dbInfo.Removed, indRange, indRange+1)
				ecowriter.WriteJSON(dbInfoPath, dbInfo)
				return true
			}
		}
	}

	return false
}

// Creating a new table.
func CreateTable(nameDB, nameTable string) bool {
	// This function is complete
	var folderName string
	var dbInfo tDBInfo = tDBInfo{}
	var tableInfo tTableInfo = tTableInfo{}

	folderDB, ok := StorageInfo.DBs[nameDB]
	if !ok {
		return false
	}

	if !CheckFolderOrFile(LocalCoreSettings.Storage, folderDB) {
		return false
	}

	pathDB := fmt.Sprintf("%s%s/", LocalCoreSettings.Storage, folderDB)

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

	tableInfo = tTableInfo{
		Name:       nameTable,
		Columns:    make(map[string]string),
		Removed:    []string{},
		LastUpdate: time.Now(),
		Deleted:    false,
	}
	tableInfoPath := fmt.Sprintf("%s/%s", fullTableName, INFOFILE_TABLE)
	if ecowriter.WriteJSON(tableInfoPath, &tableInfo) != nil {
		return false
	}

	dbInfoPath := fmt.Sprintf("%s%s", pathDB, INFOFILE_DB)
	if ecowriter.ReadJSON(dbInfoPath, &dbInfo) != nil {
		return false
	}
	dbInfo.Tables[nameTable] = folderName
	dbInfo.LastUpdate = time.Now()
	err2 := ecowriter.WriteJSON(dbInfoPath, dbInfo)

	return err2 == nil
}
