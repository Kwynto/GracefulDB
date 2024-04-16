package core

import (
	"fmt"
	"os"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/Kwynto/GracefulDB/internal/engine/basicsystem/gtypes"
)

func findWhereIds(cond gtypes.TConditions, additionalData gtypes.TAdditionalData) []uint64 {
	// -

	var (
		resIds               = make([]uint64, 4)
		progressIds []uint64 = make([]uint64, 4)
	)

	if cond.Type != "operation" {
		return []uint64{}
	}

	if cond.Key == "_id" {
		switch cond.Operation {
		case "<=":
			value, err := strconv.ParseUint(cond.Value, 10, 64)
			if err != nil {
				return []uint64{}
			}
			countVal := StorageInfo.DBs[additionalData.Db].Tables[additionalData.Table].Count
			if countVal < value {
				value = countVal
			}
			for i := uint64(1); i <= value; i++ {
				progressIds = append(progressIds, i)
			}
		case ">=":
			value, err := strconv.ParseUint(cond.Value, 10, 64)
			if err != nil {
				return []uint64{}
			}
			countVal := StorageInfo.DBs[additionalData.Db].Tables[additionalData.Table].Count
			for i := value; i <= countVal; i++ {
				progressIds = append(progressIds, i)
			}
		case "<":
			value, err := strconv.ParseUint(cond.Value, 10, 64)
			if err != nil {
				return []uint64{}
			}
			countVal := StorageInfo.DBs[additionalData.Db].Tables[additionalData.Table].Count
			if countVal < value {
				value = countVal
			}
			for i := uint64(1); i < value; i++ {
				progressIds = append(progressIds, i)
			}
		case ">":
			value, err := strconv.ParseUint(cond.Value, 10, 64)
			if err != nil {
				return []uint64{}
			}
			value++
			countVal := StorageInfo.DBs[additionalData.Db].Tables[additionalData.Table].Count
			for i := value; i <= countVal; i++ {
				progressIds = append(progressIds, i)
			}
		case "=":
			if value, err := strconv.ParseUint(cond.Value, 10, 64); err == nil {
				progressIds = append(progressIds, value)
			}
			// case "like":
			// 	return []uint64{}
			// case "regexp":
			// 	return []uint64{}
		}
	}

	if cond.Key == "_time" {
		tableInfo := StorageInfo.DBs[additionalData.Db].Tables[additionalData.Table]
		folderPath := fmt.Sprintf("%s%s/%s/service", LocalCoreSettings.Storage, tableInfo.Parent, tableInfo.Folder)

		valueCond, err := strconv.ParseUint(cond.Value, 10, 64)
		if err != nil {
			return []uint64{}
		}

		files, err := os.ReadDir(folderPath)
		if err != nil {
			return []uint64{}
		}

		for _, file := range files {
			if !file.IsDir() {
				fileName := file.Name()
				if strings.Contains(fileName, tableInfo.CurrentRev) {
					fullNameFile := fmt.Sprintf("%s/%s", folderPath, fileName)
					fileText, err := FileRead(fullNameFile)
					if err != nil {
						continue
					}
					fileData := strings.Split(fileText, "\n")
					for _, line := range fileData {
						lineData := strings.Split(line, "|")
						valueId, valueTime := lineData[0], lineData[1] // id, time, status, shape
						value, err := strconv.ParseUint(valueId, 10, 64)
						if err != nil {
							continue
						}
						valueTimeBase, err := strconv.ParseUint(valueTime, 10, 64)
						if err != nil {
							continue
						}

						switch cond.Operation {
						case "<=":
							if valueTimeBase <= valueCond {
								progressIds = append(progressIds, value)
							}
						case ">=":
							if valueTimeBase >= valueCond {
								progressIds = append(progressIds, value)
							}
						case "<":
							if valueTimeBase < valueCond {
								progressIds = append(progressIds, value)
							}
						case ">":
							if valueTimeBase > valueCond {
								progressIds = append(progressIds, value)
							}
						case "=":
							if valueTimeBase == valueCond {
								progressIds = append(progressIds, value)
							}

							// case "like":
							// 	return []uint64{}
							// case "regexp":
							// 	return []uint64{}
						}
					}
				}
			}
		}
	}

	// TODO: do it
	// if cond.Key == "_status" {
	// }

	// if cond.Key == "_shape" {
	// }

	slices.Sort(progressIds)
	progressIds = slices.Clip(progressIds)

	// TODO: make a check of all IDs before returning values
	resIds = append(resIds, progressIds...)

	// resIds = slices.Clip(resIds)
	return resIds
}

func mergeOr(first, second []uint64) []uint64 {
	// This function is complete
	resIds := append(first, second...)
	slices.Sort(resIds)
	resIds = slices.Compact(resIds)

	resIds = slices.Clip(resIds)
	return resIds
}

func mergeAnd(first, second []uint64) []uint64 {
	// This function is complete
	var resIds = make([]uint64, 4)

	for _, sElem := range second {
		if slices.Contains(first, sElem) {
			resIds = append(resIds, sElem)
		}
	}

	resIds = slices.Clip(resIds)
	return resIds
}

func whereSelection(where []gtypes.TConditions, additionalData gtypes.TAdditionalData) []uint64 {
	// - It's almost done
	var (
		acc         []uint64 // = make([]uint64, 4)
		progressIds []uint64 // = make([]uint64, 4)
		selector    string   = ""
	)

	if len(where) < 1 {
		return acc
	}

	for _, elem := range where {
		switch elem.Type {
		case "operation":
			progressIds = nil
			progressIds = findWhereIds(elem, additionalData) // TODO: do it
			switch selector {
			case "or":
				acc = mergeOr(acc, progressIds)
			case "and":
				acc = mergeAnd(acc, progressIds)
			default:
				acc = append(acc, progressIds...)
				selector = ""
			}
		case "or":
			selector = "or"
		case "and":
			selector = "and"
		}
	}

	return acc
}

func DeleteRows(nameDB, nameTable string, deleteIn gtypes.TDeleteStruct) ([]uint64, bool) {
	// This function is complete
	var whereIds []uint64 = []uint64{}
	var rowsForStore []gtypes.TRowForStore
	var cols []string = []string{}

	if !deleteIn.IsWhere {
		deleteIn.Where[0] = gtypes.TConditions{
			Type:      "operation",
			Key:       "_id",
			Operation: ">",
			Value:     "0",
		}
		deleteIn.IsWhere = true
	}

	// chacking keys
	for _, whereElem := range deleteIn.Where {
		if whereElem.Type == "operation" {
			if whereElem.Key != "_id" && whereElem.Key != "_time" && whereElem.Key != "_status" && whereElem.Key != "_shape" {
				_, ok := StorageInfo.DBs[nameDB].Tables[nameTable].Columns[whereElem.Key]
				if !ok {
					return []uint64{}, false
				}
			}
		}
	}

	dbInfo, okDB := GetDBInfo(nameDB)
	if !okDB {
		return []uint64{}, false
	}
	tableInfo, ok := dbInfo.Tables[nameTable]
	if !ok {
		return []uint64{}, false
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
		rowStore.Shape = 3 // this is code of delete
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

	return whereIds, true
}

func SelectRows(nameDB, nameTable string, updateIn gtypes.TSelectStruct) ([]uint64, bool) {
	// -
	return []uint64{}, true
}

func UpdateRows(nameDB, nameTable string, updateIn gtypes.TUpdaateStruct) ([]uint64, bool) {
	// -
	return []uint64{}, true
}

func InsertRows(nameDB, nameTable string, columns []string, rowsin [][]string) ([]uint64, bool) {
	// This function is complete
	var rows [][]string

	for _, col := range columns {
		if col == "_id" || col == "_time" || col == "_status" || col == "_shape" {
			return nil, false
		}
	}

	storageBlock.Lock() // заменить на другой блок
	defer storageBlock.Unlock()

	dbInfo, okDB := StorageInfo.DBs[nameDB]
	if !okDB {
		return nil, false
	}

	tableInfo, okTable := dbInfo.Tables[nameTable]
	if !okTable {
		return nil, false
	}

	for _, col := range columns {
		if !slices.Contains(tableInfo.Order, col) {
			return nil, false
		}
	}

	tNow := time.Now().Unix()

	var rowsForStore []gtypes.TRowForStore

	for _, row := range rowsin {
		var trow []string

		lCol := len(columns)
		lRow := len(row)

		if lCol != lRow {
			if lCol < lRow {
				trow = row[:lCol]
			}
			if lCol > lRow {
				trow = row
				for i := lRow; i < lCol; i++ {
					trow = append(trow, tableInfo.Columns[columns[i]].Specification.Default)
				}
			}
		} else {
			trow = row
		}
		rows = append(rows, trow)
	}

	for _, row := range rows {
		var rowStore = gtypes.TRowForStore{}
		tableInfo.Count++
		for _, column := range tableInfo.Columns {
			var colStore = gtypes.TColumnForStore{}

			colStore.Field = column.Name

			var vStore string
			ind := slices.Index(columns, column.Name)
			if ind != -1 {
				vStore = Encode64(row[ind])
			} else {
				if column.Specification.NotNull {
					fmt.Println("Point 4")
					return nil, false
				}
				vStore = column.Specification.Default
			}
			// colStore.Id = tableInfo.Count
			// colStore.Time = tNow
			colStore.Value = vStore
			rowStore.Row = append(rowStore.Row, colStore)
		}
		rowStore.Id = tableInfo.Count
		rowStore.Time = tNow
		rowStore.Status = 0
		rowStore.Shape = 0
		rowStore.DB = nameDB
		rowStore.Table = nameTable

		rowsForStore = append(rowsForStore, rowStore)
	}

	var result []uint64
	for _, row := range rowsForStore {
		result = append(result, row.Id)
	}

	go InsertIntoBuffer(rowsForStore)

	dbInfo.Tables[nameTable] = tableInfo
	StorageInfo.DBs[nameDB] = dbInfo

	return result, dbInfo.Save()
}
