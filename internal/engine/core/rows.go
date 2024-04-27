package core

import (
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/Kwynto/GracefulDB/internal/engine/basicsystem/gtypes"
	"github.com/Kwynto/GracefulDB/pkg/lib/ecowriter"
)

func findWhereIds(cond gtypes.TConditions, additionalData gtypes.TAdditionalData) []uint64 {
	// This function is complete
	var (
		resIds               = make([]uint64, 4)
		progressIds []uint64 = make([]uint64, 4)
		isDelete    bool     = false
	)

	if cond.Type != "operation" {
		return []uint64{}
	}

	dbInfo, ok := GetDBInfo(additionalData.Db)
	if !ok {
		return []uint64{}
	}
	tableInfo := dbInfo.Tables[additionalData.Table]

	if cond.Key == "_id" || cond.Key == "_time" || cond.Key == "_status" || cond.Key == "_shape" {
		if cond.Key == "_id" {
			// folderPath := fmt.Sprintf("%s%s/%s/service", LocalCoreSettings.Storage, tableInfo.Parent, tableInfo.Folder)
			folderPath := filepath.Join(LocalCoreSettings.Storage, tableInfo.Parent, tableInfo.Folder, "service")

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
						// fullNameFile := fmt.Sprintf("%s/%s", folderPath, fileName)
						fullNameFile := filepath.Join(folderPath, fileName)
						fileText, err := ecowriter.FileRead(fullNameFile)
						if err != nil {
							continue
						}
						fileData := strings.Split(fileText, "\n")
						for _, line := range fileData {
							lineData := strings.Split(line, "|")
							valueId, valueShape := lineData[0], lineData[3] // id, time, status, shape

							valueShapeBase, err := strconv.ParseUint(valueShape, 10, 64)
							if err != nil {
								continue
							}

							if valueShapeBase < 30 {
								value, err := strconv.ParseUint(valueId, 10, 64)
								if err != nil {
									continue
								}

								switch {
								case cond.Operation == "<=" && value <= valueCond:
									resIds = append(resIds, value)
								case cond.Operation == ">=" && value >= valueCond:
									resIds = append(resIds, value)
								case cond.Operation == "<" && value < valueCond:
									resIds = append(resIds, value)
								case cond.Operation == ">" && value > valueCond:
									resIds = append(resIds, value)
								case cond.Operation == "=" && value == valueCond:
									resIds = append(resIds, value)
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
		}

		if cond.Key == "_time" {
			// folderPath := fmt.Sprintf("%s%s/%s/service", LocalCoreSettings.Storage, tableInfo.Parent, tableInfo.Folder)
			folderPath := filepath.Join(LocalCoreSettings.Storage, tableInfo.Parent, tableInfo.Folder, "service")

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
						// fullNameFile := fmt.Sprintf("%s/%s", folderPath, fileName)
						fullNameFile := filepath.Join(folderPath, fileName)
						fileText, err := ecowriter.FileRead(fullNameFile)
						if err != nil {
							continue
						}
						fileData := strings.Split(fileText, "\n")
						for _, line := range fileData {
							lineData := strings.Split(line, "|")
							valueId, valueTime, valueShape := lineData[0], lineData[1], lineData[3] // id, time, status, shape

							valueShapeBase, err := strconv.ParseUint(valueShape, 10, 64)
							if err != nil {
								continue
							}

							if valueShapeBase < 30 {
								value, err := strconv.ParseUint(valueId, 10, 64)
								if err != nil {
									continue
								}
								valueTimeBase, err := strconv.ParseUint(valueTime, 10, 64)
								if err != nil {
									continue
								}

								switch {
								case cond.Operation == "<=" && valueTimeBase <= valueCond:
									resIds = append(resIds, value)
								case cond.Operation == ">=" && valueTimeBase >= valueCond:
									resIds = append(resIds, value)
								case cond.Operation == "<" && valueTimeBase < valueCond:
									resIds = append(resIds, value)
								case cond.Operation == ">" && valueTimeBase > valueCond:
									resIds = append(resIds, value)
								case cond.Operation == "=" && valueTimeBase == valueCond:
									resIds = append(resIds, value)
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
		}

		if cond.Key == "_status" {
			// folderPath := fmt.Sprintf("%s%s/%s/service", LocalCoreSettings.Storage, tableInfo.Parent, tableInfo.Folder)
			folderPath := filepath.Join(LocalCoreSettings.Storage, tableInfo.Parent, tableInfo.Folder, "service")

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
						// fullNameFile := fmt.Sprintf("%s/%s", folderPath, fileName)
						fullNameFile := filepath.Join(folderPath, fileName)
						fileText, err := ecowriter.FileRead(fullNameFile)
						if err != nil {
							continue
						}
						fileData := strings.Split(fileText, "\n")
						for _, line := range fileData {
							lineData := strings.Split(line, "|")
							valueId, valueStatus, valueShape := lineData[0], lineData[2], lineData[3] // id, time, status, shape

							valueShapeBase, err := strconv.ParseUint(valueShape, 10, 64)
							if err != nil {
								continue
							}

							if valueShapeBase < 30 {
								value, err := strconv.ParseUint(valueId, 10, 64)
								if err != nil {
									continue
								}
								valueStatusBase, err := strconv.ParseUint(valueStatus, 10, 64)
								if err != nil {
									continue
								}

								switch {
								case cond.Operation == "<=" && valueStatusBase <= valueCond:
									resIds = append(resIds, value)
								case cond.Operation == ">=" && valueStatusBase >= valueCond:
									resIds = append(resIds, value)
								case cond.Operation == "<" && valueStatusBase < valueCond:
									resIds = append(resIds, value)
								case cond.Operation == ">" && valueStatusBase > valueCond:
									resIds = append(resIds, value)
								case cond.Operation == "=" && valueStatusBase == valueCond:
									resIds = append(resIds, value)
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
		}

		if cond.Key == "_shape" {
			// folderPath := fmt.Sprintf("%s%s/%s/service", LocalCoreSettings.Storage, tableInfo.Parent, tableInfo.Folder)
			folderPath := filepath.Join(LocalCoreSettings.Storage, tableInfo.Parent, tableInfo.Folder, "service")

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
						// fullNameFile := fmt.Sprintf("%s/%s", folderPath, fileName)
						fullNameFile := filepath.Join(folderPath, fileName)
						fileText, err := ecowriter.FileRead(fullNameFile)
						if err != nil {
							continue
						}
						fileData := strings.Split(fileText, "\n")
						for _, line := range fileData {
							lineData := strings.Split(line, "|")
							valueId, valueShape := lineData[0], lineData[3] // id, time, status, shape
							value, err := strconv.ParseUint(valueId, 10, 64)
							if err != nil {
								continue
							}
							valueShapeBase, err := strconv.ParseUint(valueShape, 10, 64)
							if err != nil {
								continue
							}

							switch {
							case cond.Operation == "<=" && valueShapeBase <= valueCond:
								resIds = append(resIds, value)
							case cond.Operation == ">=" && valueShapeBase >= valueCond:
								resIds = append(resIds, value)
							case cond.Operation == "<" && valueShapeBase < valueCond:
								resIds = append(resIds, value)
							case cond.Operation == ">" && valueShapeBase > valueCond:
								resIds = append(resIds, value)
							case cond.Operation == "=" && valueShapeBase == valueCond:
								resIds = append(resIds, value)
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
	} else {
		valueCond := Encode64(cond.Value)

		columnInfo, ok := tableInfo.Columns[cond.Key]
		if !ok {
			return []uint64{}
		}

		// folderPath := fmt.Sprintf("%s%s/%s", LocalCoreSettings.Storage, columnInfo.Parents, columnInfo.Folder)
		folderPath := filepath.Join(LocalCoreSettings.Storage, columnInfo.Parents, columnInfo.Folder)
		files, err := os.ReadDir(folderPath)
		if err != nil {
			return []uint64{}
		}

		for _, file := range files {
			if !file.IsDir() {
				fileName := file.Name()
				if strings.Contains(fileName, tableInfo.CurrentRev) {
					// fullNameFile := fmt.Sprintf("%s/%s", folderPath, fileName)
					fullNameFile := filepath.Join(folderPath, fileName)
					fileText, err := ecowriter.FileRead(fullNameFile)
					if err != nil {
						continue
					}
					fileData := strings.Split(fileText, "\n")
					for _, line := range fileData {
						lineData := strings.Split(line, "|")
						valueId, valueData := lineData[0], lineData[1] // id, [data]
						value, err := strconv.ParseUint(valueId, 10, 64)
						if err != nil {
							continue
						}

						switch {
						case cond.Operation == "<=" && valueData <= valueCond:
							progressIds = append(progressIds, value)
						case cond.Operation == ">=" && valueData >= valueCond:
							progressIds = append(progressIds, value)
						case cond.Operation == "<" && valueData < valueCond:
							progressIds = append(progressIds, value)
						case cond.Operation == ">" && valueData > valueCond:
							progressIds = append(progressIds, value)
						case cond.Operation == "=" && valueData == valueCond:
							progressIds = append(progressIds, value)
							// case "like":
							// 	return []uint64{}
							// case "regexp":
							// 	return []uint64{}
						}
					}
				}
			}
		}

		// checking on system records
		slices.Sort(progressIds)
		progressIds = slices.Compact(progressIds)
		progressIds = slices.Clip(progressIds)

		// folderSysPath := fmt.Sprintf("%s%s/%s/service", LocalCoreSettings.Storage, tableInfo.Parent, tableInfo.Folder)
		folderSysPath := filepath.Join(LocalCoreSettings.Storage, tableInfo.Parent, tableInfo.Folder, "service")

		for _, id := range progressIds {
			maxBucket := Pow(2, tableInfo.BucketLog)
			hashid := id % maxBucket
			if hashid == 0 {
				hashid = maxBucket
			}

			// fullNameFile := fmt.Sprintf("%s/%s_%d", folderSysPath, tableInfo.CurrentRev, hashid)
			fullNameFile := filepath.Join(folderSysPath, fmt.Sprintf("%s_%d", tableInfo.CurrentRev, hashid))
			fileText, err := ecowriter.FileRead(fullNameFile)
			if err != nil {
				continue
			}

			fileData := strings.Split(fileText, "\n")

			isDelete = false

			for _, line := range fileData {
				lineData := strings.Split(line, "|")
				valueId, valueShape := lineData[0], lineData[3] // id, time, status, shape

				value, err := strconv.ParseUint(valueId, 10, 64)
				if err != nil {
					continue
				}
				valueShapeBase, err := strconv.ParseUint(valueShape, 10, 64)
				if err != nil {
					continue
				}

				if value == id && valueShapeBase == 30 {
					isDelete = true
				}
			}

			if !isDelete {
				resIds = append(resIds, id)
			}
		}
	}

	slices.Sort(resIds)
	resIds = slices.Compact(resIds)
	resIds = slices.Clip(resIds)

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
	// This function is complete
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
			// progressIds = nil
			progressIds = findWhereIds(elem, additionalData)
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

	dbInfo, okDB := GetDBInfo(nameDB)
	if !okDB {
		return []uint64{}, false
	}
	tableInfo, ok := dbInfo.Tables[nameTable]
	if !ok {
		return []uint64{}, false
	}

	// chacking keys
	for _, whereElem := range deleteIn.Where {
		if whereElem.Type == "operation" {
			if whereElem.Key != "_id" && whereElem.Key != "_time" && whereElem.Key != "_status" && whereElem.Key != "_shape" {
				_, ok := tableInfo.Columns[whereElem.Key]
				if !ok {
					return []uint64{}, false
				}
			}
		}
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

	return whereIds, true
}

func SelectRows(nameDB, nameTable string, updateIn gtypes.TSelectStruct) ([]gtypes.TResponseRow, bool) {
	// - It's almost done

	// TODO: do it

	return []gtypes.TResponseRow{}, true
}

func UpdateRows(nameDB, nameTable string, updateIn gtypes.TUpdaateStruct) ([]uint64, bool) {
	// This function is complete
	var whereIds []uint64 = []uint64{}
	var rowsForStore []gtypes.TRowForStore
	var cols []string = []string{}
	var value string = ""

	dbInfo, okDB := GetDBInfo(nameDB)
	if !okDB {
		return []uint64{}, false
	}
	tableInfo, ok := dbInfo.Tables[nameTable]
	if !ok {
		return []uint64{}, false
	}

	// chacking keys
	for _, whereElem := range updateIn.Where {
		if whereElem.Type == "operation" {
			if whereElem.Key != "_id" && whereElem.Key != "_time" && whereElem.Key != "_status" && whereElem.Key != "_shape" {
				_, ok := tableInfo.Columns[whereElem.Key]
				if !ok {
					return []uint64{}, false
				}
			}
		}
	}

	for _, col := range tableInfo.Columns {
		cols = append(cols, col.Name)
	}

	additionalData := gtypes.TAdditionalData{
		Db:    nameDB,
		Table: nameTable,
	}

	whereIds = whereSelection(updateIn.Where, additionalData)

	tNow := time.Now().Unix()

	// Updating by changing the status of records and setting new values
	for _, id := range whereIds {
		var rowStore = gtypes.TRowForStore{}
		rowStore.Id = id
		rowStore.Time = tNow
		rowStore.Status = 0
		rowStore.Shape = 20 // this is code of update
		rowStore.DB = nameDB
		rowStore.Table = nameTable
		for _, col := range cols {
			newValue, ok := updateIn.Couples[col]
			if ok {
				value = newValue
			} else {
				value, _ = GetColumnById(nameDB, nameTable, col, id)
				// okCol := false
				// value, okCol = GetColumnById(nameDB, nameTable, col, id)
				// if !okCol {
				// 	value = ""
				// }
			}
			rowStore.Row = append(rowStore.Row, gtypes.TColumnForStore{
				Field: col,
				Value: value,
			})
		}

		rowsForStore = append(rowsForStore, rowStore)
	}

	if len(whereIds) > 0 {
		go InsertIntoBuffer(rowsForStore)
	}

	return whereIds, true
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
