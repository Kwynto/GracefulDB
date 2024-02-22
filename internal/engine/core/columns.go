package core

import (
	"fmt"
	"os"
	"slices"
	"time"
)

// Marks the column as deleted, but does not delete files.
func RemoveColumn(nameDB, nameTable, nameColumn string) bool {
	// This function is complete
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

// Changing a column.
func ChangeColumn(nameDB, nameTable string, newDataColumn TColumnForWrite, secure bool) bool {
	// This function is complete
	if secure && !RegExpCollection["EntityName"].MatchString(newDataColumn.Name) {
		return false
	}

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
	if newDataColumn.IsChDefault {
		columnInfo.Specification.Default = newDataColumn.Spec.Default
	}
	if newDataColumn.IsChNotNull {
		columnInfo.Specification.NotNull = newDataColumn.Spec.NotNull
	}
	if newDataColumn.IsChUniqut {
		columnInfo.Specification.Unique = newDataColumn.Spec.Unique
	}
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
	// StorageInfo.Save() // rewriting only Access

	return dbInfo.Save()
}

// Creating a new column.
func CreateColumn(nameDB, nameTable, nameColumn string, secure bool, specification TColumnSpecification) bool {
	// This function is complete
	if secure && !RegExpCollection["EntityName"].MatchString(nameDB) &&
		!RegExpCollection["EntityName"].MatchString(nameTable) &&
		!RegExpCollection["EntityName"].MatchString(nameColumn) {
		return false
	}

	var folderName string

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

	columnInfo := tColumnInfo{
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
