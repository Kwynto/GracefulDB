package core

import (
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"time"

	"github.com/Kwynto/GracefulDB/internal/engine/basicsystem/gtypes"
	"github.com/Kwynto/GracefulDB/internal/engine/basicsystem/instead"
	"github.com/Kwynto/GracefulDB/internal/engine/basicsystem/vqlexp"
)

// Marks the column as deleted, but does not delete files.
func RemoveColumn(sNameDB, sNameTable, sNameColumn string) bool {
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

	stColumnInfo, isOk := stTableInfo.Columns[sNameColumn]
	if !isOk {
		return false
	}

	dtNow := time.Now()

	stColumnInfo.LastUpdate = dtNow
	stColumnInfo.Deleted = true

	stTableInfo.Removed = append(stTableInfo.Removed, stColumnInfo)
	delete(stTableInfo.Columns, sNameColumn)
	iInd := slices.Index(stTableInfo.Order, sNameColumn)
	if iInd != -1 {
		stTableInfo.Order = slices.Delete(stTableInfo.Order, iInd, iInd+1)
	}
	stTableInfo.LastUpdate = dtNow

	stDBInfo.Tables[sNameTable] = stTableInfo
	stDBInfo.LastUpdate = dtNow

	StStorageInfo.DBs[sNameDB] = stDBInfo

	return stDBInfo.Save()
}

// Deletes the folder and column files, if column was mark as 'removed'
func StrongRemoveColumn(sNameDB, sNameTable, sNameColumn string) bool {
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

	for iColumnInd, stColumnVal := range stTableInfo.Removed {
		if stColumnVal.Name == sNameColumn {
			sColumnPath := filepath.Join(StLocalCoreSettings.Storage, stColumnVal.Parents, stColumnVal.Folder)
			err := os.RemoveAll(sColumnPath)
			if err != nil {
				return false
			}

			dtNow := time.Now()

			stTableInfo.Removed = slices.Delete(stTableInfo.Removed, iColumnInd, iColumnInd+1)
			stTableInfo.LastUpdate = dtNow

			stDBInfo.Tables[sNameTable] = stTableInfo
			stDBInfo.LastUpdate = dtNow

			StStorageInfo.DBs[sNameDB] = stDBInfo

			return stDBInfo.Save()
		}
	}

	return false
}

// Changing a column
func ChangeColumn(sNameDB, sNameTable string, stNewDataColumn gtypes.TColumnForWrite, isSecure bool) bool {
	// This function is complete
	if isSecure && !vqlexp.MRegExpCollection["EntityName"].MatchString(stNewDataColumn.Name) {
		return false
	}

	if stNewDataColumn.Spec.Default != "" {
		stNewDataColumn.Spec.Default = Encode64(stNewDataColumn.Spec.Default)
	}

	mxStorageBlock.Lock()
	defer mxStorageBlock.Unlock()

	stDBInfo, isOkDB := StStorageInfo.DBs[sNameDB]
	if !isOkDB {
		return false
	}

	stTableInfo, isOkTable := stDBInfo.Tables[sNameTable]
	if !isOkTable {
		return false
	}

	var sName string
	if stNewDataColumn.IsChName {
		sName = stNewDataColumn.OldName
	} else {
		sName = stNewDataColumn.Name
	}
	stColumnInfo, isOkCol := stTableInfo.Columns[sName]
	if !isOkCol {
		return false
	}

	dtNow := time.Now()

	if stNewDataColumn.IsChName {
		stColumnInfo.OldName = stColumnInfo.Name
		stColumnInfo.Name = stNewDataColumn.Name
	}

	stColumnInfo.Specification.Default = stNewDataColumn.Spec.Default
	stColumnInfo.Specification.NotNull = stNewDataColumn.Spec.NotNull
	stColumnInfo.Specification.Unique = stNewDataColumn.Spec.Unique
	stColumnInfo.LastUpdate = dtNow

	if stNewDataColumn.IsChName {
		delete(stTableInfo.Columns, stNewDataColumn.OldName)
		stTableInfo.Columns[stNewDataColumn.Name] = stColumnInfo
		i := slices.Index(stTableInfo.Order, stNewDataColumn.OldName)
		if i > -1 {
			stTableInfo.Order[i] = stNewDataColumn.Name
		} else {
			stTableInfo.Order = append(stTableInfo.Order, stNewDataColumn.Name)
		}
	} else {
		stTableInfo.Columns[stNewDataColumn.Name] = stColumnInfo
	}
	stTableInfo.LastUpdate = dtNow

	stDBInfo.Tables[sNameTable] = stTableInfo
	stDBInfo.LastUpdate = dtNow

	StStorageInfo.DBs[sNameDB] = stDBInfo

	return stDBInfo.Save()
}

// func GetDescriptionColumn(db, table, column string) gtypes.TDescColumn {
// 	// This function is complete
// 	dbInfo, _ := GetDBInfo(db)
// 	tableInfo := dbInfo.Tables[table]
// 	col := tableInfo.Columns[column]

// 	return gtypes.TDescColumn{
// 		DB:         db,
// 		Table:      table,
// 		Column:     column,
// 		Path:       fmt.Sprintf("%s%s/%s/", LocalCoreSettings.Storage, col.Parents, col.Folder),
// 		Spec:       col.Specification,
// 		CurrentRev: tableInfo.CurrentRev,
// 		BucketSize: tableInfo.BucketSize,
// 		BucketLog:  tableInfo.BucketLog,
// 	}
// }

// Get up-to-date cell data
func GetColumnById(sNameColumn string, uIdRow uint64, stAddData gtypes.TAdditionalData) (string, bool) {
	// This function is complete
	var sResValue string

	stDBInfo, isOkDB := GetDBInfo(stAddData.Db)
	if !isOkDB {
		return "", false
	}
	stTableInfo, isOk := stDBInfo.Tables[stAddData.Table]
	if !isOk {
		return "", false
	}
	stColumnInfo, isOk := stTableInfo.Columns[sNameColumn]
	if !isOk {
		return "", false
	}

	sFolderPath := filepath.Join(StLocalCoreSettings.Storage, stColumnInfo.Parents, stColumnInfo.Folder)

	uMaxBucket := Pow(2, stTableInfo.BucketLog)
	uHashId := uIdRow % uMaxBucket
	if uHashId == 0 {
		uHashId = uMaxBucket
	}

	sFullNameFile := filepath.Join(sFolderPath, fmt.Sprintf("%s_%d", stTableInfo.CurrentRev, uHashId))

	slCache, isOkFile := instead.LoadFile(sFullNameFile, stAddData.Stamp)
	if !isOkFile {
		return "", false
	}

	for _, slLine := range slCache {
		if len(slLine) < 2 {
			continue
		}
		sValueId, sValueData := slLine[0], slLine[1] // id, [data]

		uId, err := strconv.ParseUint(sValueId, 10, 64)
		if err != nil {
			continue
		}
		if uId == uIdRow {
			sResValue = sValueData
		}
	}

	return sResValue, true
}

func GetInfoById(uIdRow uint64, stAddData gtypes.TAdditionalData) (sTime string, sStatus string, sShape string, isOk bool) {
	sIdRow := strconv.FormatUint(uIdRow, 10)

	stDBInfo, isOkDB := GetDBInfo(stAddData.Db)
	if !isOkDB {
		return sTime, sStatus, sShape, false
	}
	stTableInfo, isOkTable := stDBInfo.Tables[stAddData.Table]
	if !isOkTable {
		return sTime, sStatus, sShape, false
	}

	uMaxBucket := Pow(2, stTableInfo.BucketLog)
	uHashId := uIdRow % uMaxBucket
	if uHashId == 0 {
		uHashId = uMaxBucket
	}

	sFullNameFile := filepath.Join(StLocalCoreSettings.Storage, stTableInfo.Parent, stTableInfo.Folder, "service", fmt.Sprintf("%s_%d", stTableInfo.CurrentRev, uHashId))

	slCache, isOkFile := instead.LoadFile(sFullNameFile, stAddData.Stamp)
	if !isOkFile {
		return sTime, sStatus, sShape, false
	}

	for _, sLine := range slCache {
		if len(sLine) < 4 {
			return sTime, sStatus, sShape, false
		}
		if sIdRow != sLine[0] {
			continue
		}
		sTime, sStatus, sShape = sLine[1], sLine[2], sLine[3] // time, status, shape
	}

	return sTime, sStatus, sShape, true
}

// Creating a new column
func CreateColumn(sNameDB, sNameTable, sNameColumn string, isSecure bool, stSpecification gtypes.TColumnSpecification) bool {
	// This function is complete
	if isSecure && !vqlexp.MRegExpCollection["EntityName"].MatchString(sNameColumn) {
		return false
	}

	var sFolderName string

	if stSpecification.Default != "" {
		stSpecification.Default = Encode64(stSpecification.Default)
	}

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

	if _, isOk := stTableInfo.Columns[sNameColumn]; isOk {
		return false
	}

	sPathTable := filepath.Join(StLocalCoreSettings.Storage, stTableInfo.Parent, stTableInfo.Folder)

	for {
		sFolderName = GenerateName()
		if !CheckFolder(sPathTable, sFolderName) {
			break
		}
	}

	sFullColumnName := filepath.Join(sPathTable, sFolderName)
	err := os.Mkdir(sFullColumnName, 0666)
	if err != nil {
		return false
	}

	dtNow := time.Now()

	stColumnInfo := TColumnInfo{
		Name:          sNameColumn,
		OldName:       "",
		Folder:        sFolderName,
		Parents:       filepath.Join(stTableInfo.Parent, stTableInfo.Folder),
		Specification: stSpecification,
		LastUpdate:    dtNow,
		Deleted:       false,
	}

	stTableInfo.Columns[sNameColumn] = stColumnInfo
	stTableInfo.Order = append(stTableInfo.Order, sNameColumn)
	stTableInfo.LastUpdate = dtNow

	stDBInfo.Tables[sNameTable] = stTableInfo
	stDBInfo.LastUpdate = dtNow

	StStorageInfo.DBs[sNameDB] = stDBInfo

	return stDBInfo.Save()
}
