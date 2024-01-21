package core

import (
	"fmt"
	"os"
	"time"

	"github.com/Kwynto/GracefulDB/pkg/lib/ecowriter"
)

// Creating a new table.
func CreateTable(nameDB, nameTable string) bool {
	// This function is complete
	var folderName string
	var dbInfo tDBInfo = tDBInfo{}

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

	tableInfo := tTableInfo{
		Name:       nameTable,
		Columns:    make(map[string]string),
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
