package core

import (
	"os"
	"path/filepath"
	"slices"
	"time"

	"github.com/Kwynto/GracefulDB/internal/engine/basicsystem/gtypes"
	"github.com/Kwynto/GracefulDB/internal/engine/basicsystem/vqlexp"
)

// Marks the table as deleted, but does not delete files.
func RemoveTable(sNameDB, sNameTable string) bool {
	// This function is complete
	mxStorageBlock.Lock()
	defer mxStorageBlock.Unlock()

	stDBInfo, isOk := StStorageInfo.DBs[sNameDB]
	if !isOk {
		return false
	}

	stTableInfo, isOk := stDBInfo.Tables[sNameTable]
	if !isOk {
		return false
	}

	dtNow := time.Now()

	stTableInfo.LastUpdate = dtNow
	stTableInfo.Deleted = true

	stDBInfo.Removed = append(stDBInfo.Removed, stTableInfo)
	delete(stDBInfo.Tables, sNameTable)
	stDBInfo.LastUpdate = dtNow

	StStorageInfo.DBs[sNameDB] = stDBInfo

	return stDBInfo.Save()
}

// Deletes the folder and table files, if table was mark as 'removed'
func StrongRemoveTable(sNameDB, sNameTable string) bool {
	// This function is complete
	mxStorageBlock.Lock()
	defer mxStorageBlock.Unlock()

	stDBInfo, isOk := StStorageInfo.DBs[sNameDB]
	if !isOk {
		return false
	}

	for iRangeInd, stTableInfo := range stDBInfo.Removed {
		if stTableInfo.Name == sNameTable {
			// tablePath := fmt.Sprintf("%s%s/%s", LocalCoreSettings.Storage, tableInfo.Parent, tableInfo.Folder)
			sTablePath := filepath.Join(StLocalCoreSettings.Storage, stTableInfo.Parent, stTableInfo.Folder)
			err := os.RemoveAll(sTablePath)
			if err != nil {
				return false
			}

			stDBInfo.Removed = slices.Delete(stDBInfo.Removed, iRangeInd, iRangeInd+1)
			stDBInfo.LastUpdate = time.Now()
			StStorageInfo.DBs[sNameDB] = stDBInfo

			return stDBInfo.Save()
		}
	}

	return false
}

// Rename a table.
func RenameTable(sNameDB, sOldNameTable, sNewNameTable string, isSecure bool) bool {
	// This function is complete
	if isSecure && !vqlexp.MRegExpCollection["EntityName"].MatchString(sNewNameTable) {
		return false
	}

	mxStorageBlock.Lock()
	defer mxStorageBlock.Unlock()

	stDBInfo, isOkDB := StStorageInfo.DBs[sNameDB]
	if isOkDB {
		stTableInfo, isOkTable := stDBInfo.Tables[sOldNameTable]
		if !isOkTable {
			return false
		}

		dtNow := time.Now()

		stTableInfo.Name = sNewNameTable
		stTableInfo.LastUpdate = dtNow

		delete(stDBInfo.Tables, sOldNameTable)
		stDBInfo.Tables[sNewNameTable] = stTableInfo
		stDBInfo.LastUpdate = dtNow

		StStorageInfo.DBs[sNameDB] = stDBInfo

		return stDBInfo.Save()
	}

	return false
}

func TruncateTable(sNameDB, sNameTable string) bool {
	// This function is complete
	var slUWhereIds []uint64 = []uint64{}
	var slStRowsForStore []gtypes.TRowForStore
	var slSCols []string = []string{}
	var stDeleteIn gtypes.TDeleteStruct = gtypes.TDeleteStruct{
		Where:   make([]gtypes.TConditions, 2),
		IsWhere: false,
	}

	stDeleteIn.Where[0] = gtypes.TConditions{
		Type:      "operation",
		Key:       "_id",
		Operation: ">",
		Value:     "0",
	}
	stDeleteIn.IsWhere = true

	stDBInfo, isOkDB := GetDBInfo(sNameDB)
	if !isOkDB {
		return false
	}
	stTableInfo, isOk := stDBInfo.Tables[sNameTable]
	if !isOk {
		return false
	}
	for _, stCol := range stTableInfo.Columns {
		slSCols = append(slSCols, stCol.Name)
	}

	stAdditionalData := gtypes.TAdditionalData{
		Db:    sNameDB,
		Table: sNameTable,
	}

	slUWhereIds = whereSelection(stDeleteIn.Where, stAdditionalData)

	dtNow := time.Now().Unix()

	// Deleting by changing the status of records and setting zero values
	for _, uId := range slUWhereIds {
		var stRowStore = gtypes.TRowForStore{}
		stRowStore.Id = uId
		stRowStore.Time = dtNow
		stRowStore.Status = 0
		stRowStore.Shape = 30 // this is code of delete
		stRowStore.DB = sNameDB
		stRowStore.Table = sNameTable
		for _, sCol := range slSCols {
			stRowStore.Row = append(stRowStore.Row, gtypes.TColumnForStore{
				Field: sCol,
				Value: "",
			})
		}

		slStRowsForStore = append(slStRowsForStore, stRowStore)
	}

	if len(slUWhereIds) > 0 {
		go InsertIntoBuffer(slStRowsForStore)
	}

	return true
}

// Creating a new table.
func CreateTable(sNameDB, sNameTable string, isSecure bool) bool {
	// This function is complete
	if isSecure && !vqlexp.MRegExpCollection["EntityName"].MatchString(sNameTable) {
		return false
	}

	var sFolderName string

	mxStorageBlock.Lock()
	defer mxStorageBlock.Unlock()

	stDBInfo, isOk := StStorageInfo.DBs[sNameDB]
	if !isOk {
		return false
	}

	// pathDB := fmt.Sprintf("%s%s/", LocalCoreSettings.Storage, dbInfo.Folder)
	sPathDB := filepath.Join(StLocalCoreSettings.Storage, stDBInfo.Folder)

	for {
		sFolderName = GenerateName()
		if !CheckFolder(sPathDB, sFolderName) {
			break
		}
	}

	// fullTableName := fmt.Sprintf("%s%s", pathDB, folderName)
	sFullTableName := filepath.Join(sPathDB, sFolderName)
	err1 := os.Mkdir(sFullTableName, 0666)
	if err1 != nil {
		return false
	}

	// serviceName := fmt.Sprintf("%s/service", fullTableName)
	sServiceName := filepath.Join(sFullTableName, "service")
	err2 := os.Mkdir(sServiceName, 0666)
	if err2 != nil {
		return false
	}

	stTableInfo := TTableInfo{
		Name:       sNameTable,
		Patronymic: sNameDB,
		Folder:     sFolderName,
		Parent:     stDBInfo.Folder,
		Columns:    make(map[string]TColumnInfo),
		Removed:    make([]TColumnInfo, 0),
		Order:      make([]string, 0),
		BucketLog:  2,
		BucketSize: StLocalCoreSettings.BucketSize,
		OldRev:     "",
		CurrentRev: GenerateRev(),
		Count:      0,
		LastUpdate: time.Now(),
		Deleted:    false,
	}

	stDBInfo.Tables[sNameTable] = stTableInfo
	stDBInfo.LastUpdate = time.Now()
	StStorageInfo.DBs[sNameDB] = stDBInfo

	return stDBInfo.Save()
}

func IfExistTable(sDB, sTable string) bool {
	stDBInfo, isOk := GetDBInfo(sDB)
	if !isOk {
		return false
	}
	_, isOkTab := stDBInfo.Tables[sTable]

	return isOkTab
}
