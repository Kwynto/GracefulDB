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
			// tablePath := fmt.Sprintf("%s%s/%s", LocalCoreSettings.Storage, tableInfo.Parent, tableInfo.Folder)
			tablePath := filepath.Join(LocalCoreSettings.Storage, tableInfo.Parent, tableInfo.Folder)
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
	// This function is complete
	var whereIds []uint64 = []uint64{}
	var rowsForStore []gtypes.TRowForStore
	var cols []string = []string{}
	var deleteIn gtypes.TDeleteStruct = gtypes.TDeleteStruct{
		Where:   make([]gtypes.TConditions, 2),
		IsWhere: false,
	}

	deleteIn.Where[0] = gtypes.TConditions{
		Type:      "operation",
		Key:       "_id",
		Operation: ">",
		Value:     "0",
	}
	deleteIn.IsWhere = true

	dbInfo, okDB := GetDBInfo(nameDB)
	if !okDB {
		return false
	}
	tableInfo, ok := dbInfo.Tables[nameTable]
	if !ok {
		return false
	}
	for _, col := range tableInfo.Columns {
		cols = append(cols, col.Name)
	}

	additionalData := gtypes.TAdditionalData{
		Db:    nameDB,
		Table: nameTable,
	}

	whereIds = whereSelection(deleteIn.Where, additionalData)

	tNow := time.Now().Unix()

	// Deleting by changing the status of records and setting zero values
	for _, id := range whereIds {
		var rowStore = gtypes.TRowForStore{}
		rowStore.Id = id
		rowStore.Time = tNow
		rowStore.Status = 0
		rowStore.Shape = 30 // this is code of delete
		rowStore.DB = nameDB
		rowStore.Table = nameTable
		for _, col := range cols {
			rowStore.Row = append(rowStore.Row, gtypes.TColumnForStore{
				Field: col,
				Value: "",
			})
		}

		rowsForStore = append(rowsForStore, rowStore)
	}

	if len(whereIds) > 0 {
		go InsertIntoBuffer(rowsForStore)
	}

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

	// pathDB := fmt.Sprintf("%s%s/", LocalCoreSettings.Storage, dbInfo.Folder)
	pathDB := filepath.Join(LocalCoreSettings.Storage, dbInfo.Folder)

	for {
		folderName = GenerateName()
		if !CheckFolder(pathDB, folderName) {
			break
		}
	}

	// fullTableName := fmt.Sprintf("%s%s", pathDB, folderName)
	fullTableName := filepath.Join(pathDB, folderName)
	err1 := os.Mkdir(fullTableName, 0666)
	if err1 != nil {
		return false
	}

	// serviceName := fmt.Sprintf("%s/service", fullTableName)
	serviceName := filepath.Join(fullTableName, "service")
	err2 := os.Mkdir(serviceName, 0666)
	if err2 != nil {
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
		BucketLog:  2,
		BucketSize: LocalCoreSettings.BucketSize,
		OldRev:     "",
		CurrentRev: GenerateRev(),
		Count:      0,
		LastUpdate: time.Now(),
		Deleted:    false,
	}

	dbInfo.Tables[nameTable] = tableInfo
	dbInfo.LastUpdate = time.Now()
	StorageInfo.DBs[nameDB] = dbInfo

	return dbInfo.Save()
}

func IfExistTable(db, table string) bool {
	dbInfo, okDb := GetDBInfo(db)
	if !okDb {
		return false
	}
	_, okTab := dbInfo.Tables[table]

	return okTab
}
