package core

import (
	"fmt"
	"os"
	"time"

	"github.com/Kwynto/GracefulDB/pkg/lib/ecowriter"
)

// Creating a new column.
func CreateColumn(nameDB, nameTable, nameColumn string) bool {
	// This function is complete
	var folderName string
	var dbInfo tDBInfo = tDBInfo{}
	var tableInfo tTableInfo = tTableInfo{}
	var columnInfo tColumnInfo = tColumnInfo{}

	folderDB, ok := StorageInfo.DBs[nameDB]
	if !ok {
		return false
	}

	if !CheckFolderOrFile(LocalCoreSettings.Storage, folderDB) {
		return false
	}

	pathDB := fmt.Sprintf("%s%s/", LocalCoreSettings.Storage, folderDB)
	dbInfoPath := fmt.Sprintf("%s%s", pathDB, INFOFILE_DB)
	if ecowriter.ReadJSON(dbInfoPath, &dbInfo) != nil {
		return false
	}

	folderTable, ok := dbInfo.Tables[nameTable]
	if !ok {
		return false
	}

	if !CheckFolderOrFile(pathDB, folderTable) {
		return false
	}

	pathTable := fmt.Sprintf("%s%s/", pathDB, folderTable)
	tableInfoPath := fmt.Sprintf("%s%s", pathTable, INFOFILE_TABLE)
	if ecowriter.ReadJSON(tableInfoPath, &tableInfo) != nil {
		return false
	}

	for {
		folderName = GenerateName()
		if !CheckFolderOrFile(pathTable, folderName) {
			break
		}
	}

	fullColumnName := fmt.Sprintf("%s%s", pathTable, folderName)
	err := os.Mkdir(fullColumnName, 0666)
	if err != nil {
		return false
	}

	tNow := time.Now()

	columnInfo = tColumnInfo{
		Name:       nameColumn,
		BucketLog:  2,
		BucketSize: LocalCoreSettings.BucketSize,
		OldRev:     "",
		CurrentRev: GenerateRev(),
		LastUpdate: tNow,
		Deleted:    false,
	}

	columnInfoPath := fmt.Sprintf("%s/%s", fullColumnName, INFOFILE_COLUMN)
	if ecowriter.WriteJSON(columnInfoPath, &columnInfo) != nil {
		return false
	}

	tableInfo.Columns[nameColumn] = folderName
	tableInfo.LastUpdate = tNow
	err2 := ecowriter.WriteJSON(tableInfoPath, tableInfo)

	return err2 == nil
}
