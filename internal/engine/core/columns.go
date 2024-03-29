package core

import (
	"fmt"
	"os"
	"slices"
	"time"

	"github.com/Kwynto/GracefulDB/internal/engine/basicsystem/gtypes"
	"github.com/Kwynto/GracefulDB/internal/engine/basicsystem/vqlexp"
)

// Marks the column as deleted, but does not delete files.
func RemoveColumn(nameDB, nameTable, nameColumn string) bool {
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

	columnInfo, ok := tableInfo.Columns[nameColumn]
	if !ok {
		return false
	}

	tNow := time.Now()

	columnInfo.LastUpdate = tNow
	columnInfo.Deleted = true

	tableInfo.Removed = append(tableInfo.Removed, columnInfo)
	delete(tableInfo.Columns, nameColumn)
	ind := slices.Index(tableInfo.Order, nameColumn)
	if ind != -1 {
		tableInfo.Order = slices.Delete(tableInfo.Order, ind, ind+1)
	}
	tableInfo.LastUpdate = tNow

	dbInfo.Tables[nameTable] = tableInfo
	dbInfo.LastUpdate = tNow

	StorageInfo.DBs[nameDB] = dbInfo

	return dbInfo.Save()
}

// Deletes the folder and column files, if column was mark as 'removed'
func StrongRemoveColumn(nameDB, nameTable, nameColumn string) bool {
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

	for indRange, columnInfo := range tableInfo.Removed {
		if columnInfo.Name == nameColumn {
			columnPath := fmt.Sprintf("%s%s/%s", LocalCoreSettings.Storage, columnInfo.Parents, columnInfo.Folder)
			err := os.RemoveAll(columnPath)
			if err != nil {
				return false
			}

			tNow := time.Now()

			tableInfo.Removed = slices.Delete(tableInfo.Removed, indRange, indRange+1)
			tableInfo.LastUpdate = tNow

			dbInfo.Tables[nameTable] = tableInfo
			dbInfo.LastUpdate = tNow

			StorageInfo.DBs[nameDB] = dbInfo

			return dbInfo.Save()
		}
	}

	return false
}

// Changing a column
func ChangeColumn(nameDB, nameTable string, newDataColumn gtypes.TColumnForWrite, secure bool) bool {
	// This function is complete
	if secure && !vqlexp.RegExpCollection["EntityName"].MatchString(newDataColumn.Name) {
		return false
	}

	if newDataColumn.Spec.Default != "" {
		newDataColumn.Spec.Default = Encode64(newDataColumn.Spec.Default)
	}

	storageBlock.Lock()
	defer storageBlock.Unlock()

	dbInfo, okDB := StorageInfo.DBs[nameDB]
	if !okDB {
		return false
	}

	tableInfo, okTable := dbInfo.Tables[nameTable]
	if !okTable {
		return false
	}

	var name string
	if newDataColumn.IsChName {
		name = newDataColumn.OldName
	} else {
		name = newDataColumn.Name
	}
	columnInfo, okCol := tableInfo.Columns[name]
	if !okCol {
		return false
	}

	tNow := time.Now()

	if newDataColumn.IsChName {
		columnInfo.OldName = columnInfo.Name
		columnInfo.Name = newDataColumn.Name
	}

	columnInfo.Specification.Default = newDataColumn.Spec.Default
	columnInfo.Specification.NotNull = newDataColumn.Spec.NotNull
	columnInfo.Specification.Unique = newDataColumn.Spec.Unique
	columnInfo.LastUpdate = tNow

	if newDataColumn.IsChName {
		delete(tableInfo.Columns, newDataColumn.OldName)
		tableInfo.Columns[newDataColumn.Name] = columnInfo
		i := slices.Index(tableInfo.Order, newDataColumn.OldName)
		if i > -1 {
			tableInfo.Order[i] = newDataColumn.Name
		} else {
			tableInfo.Order = append(tableInfo.Order, newDataColumn.Name)
		}
	} else {
		tableInfo.Columns[newDataColumn.Name] = columnInfo
	}
	tableInfo.LastUpdate = tNow

	dbInfo.Tables[nameTable] = tableInfo
	dbInfo.LastUpdate = tNow

	StorageInfo.DBs[nameDB] = dbInfo

	return dbInfo.Save()
}

func GetDescriptionColumn(db, table, column string) gtypes.TDescColumn {
	// This function is complete
	dbInfo, _ := GetDBInfo(db)
	col := dbInfo.Tables[table].Columns[column]

	return gtypes.TDescColumn{
		DB:         db,
		Table:      table,
		Column:     column,
		Path:       fmt.Sprintf("%s%s/%s/", LocalCoreSettings.Storage, col.Parents, col.Folder),
		Spec:       col.Specification,
		CurrentRev: col.CurrentRev,
		BucketSize: col.BucketSize,
		BucketLog:  col.BucketLog,
	}
}

// Creating a new column
func CreateColumn(nameDB, nameTable, nameColumn string, secure bool, specification gtypes.TColumnSpecification) bool {
	// This function is complete
	if secure && !vqlexp.RegExpCollection["EntityName"].MatchString(nameColumn) {
		return false
	}

	var folderName string

	if specification.Default != "" {
		specification.Default = Encode64(specification.Default)
	}

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

	if _, ok := tableInfo.Columns[nameColumn]; ok {
		return false
	}

	pathTable := fmt.Sprintf("%s%s/%s/", LocalCoreSettings.Storage, tableInfo.Parent, tableInfo.Folder)

	for {
		folderName = GenerateName()
		if !CheckFolder(pathTable, folderName) {
			break
		}
	}

	fullColumnName := fmt.Sprintf("%s%s", pathTable, folderName)
	err := os.Mkdir(fullColumnName, 0666)
	if err != nil {
		return false
	}

	tNow := time.Now()

	columnInfo := TColumnInfo{
		Name:          nameColumn,
		OldName:       "",
		Folder:        folderName,
		Parents:       fmt.Sprintf("%s/%s", tableInfo.Parent, tableInfo.Folder),
		BucketLog:     2,
		BucketSize:    LocalCoreSettings.BucketSize,
		OldRev:        "",
		CurrentRev:    GenerateRev(),
		Specification: specification,
		LastUpdate:    tNow,
		Deleted:       false,
	}

	tableInfo.Columns[nameColumn] = columnInfo
	tableInfo.Order = append(tableInfo.Order, nameColumn)
	tableInfo.LastUpdate = tNow

	dbInfo.Tables[nameTable] = tableInfo
	dbInfo.LastUpdate = tNow

	StorageInfo.DBs[nameDB] = dbInfo

	return dbInfo.Save()
}
