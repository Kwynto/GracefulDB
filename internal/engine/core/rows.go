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

func findWhereIds(stCond gtypes.TConditions, stAdditionalData gtypes.TAdditionalData) []uint64 {
	// This function is complete
	var (
		slUResIds           = make([]uint64, 4)
		slUProgressIds      = make([]uint64, 4)
		isDelete       bool = false
	)

	if stCond.Type != "operation" {
		return []uint64{}
	}

	stDBInfo, isOk := GetDBInfo(stAdditionalData.Db)
	if !isOk {
		return []uint64{}
	}
	stTableInfo := stDBInfo.Tables[stAdditionalData.Table]

	if stCond.Key == "_id" || stCond.Key == "_time" || stCond.Key == "_status" || stCond.Key == "_shape" {
		if stCond.Key == "_id" {
			// folderPath := fmt.Sprintf("%s%s/%s/service", LocalCoreSettings.Storage, tableInfo.Parent, tableInfo.Folder)
			sFolderPath := filepath.Join(StLocalCoreSettings.Storage, stTableInfo.Parent, stTableInfo.Folder, "service")

			uValueCond, err := strconv.ParseUint(stCond.Value, 10, 64)
			if err != nil {
				return []uint64{}
			}

			slFiles, err := os.ReadDir(sFolderPath)
			if err != nil {
				return []uint64{}
			}

			for _, fVal := range slFiles {
				if !fVal.IsDir() {
					sFileName := fVal.Name()
					if strings.Contains(sFileName, stTableInfo.CurrentRev) {
						// fullNameFile := fmt.Sprintf("%s/%s", folderPath, fileName)
						sFullNameFile := filepath.Join(sFolderPath, sFileName)
						sFileText, err := ecowriter.FileRead(sFullNameFile)
						if err != nil {
							continue
						}
						slSFileData := strings.Split(sFileText, "\n")
						for _, sLine := range slSFileData {
							slLineData := strings.Split(sLine, "|")
							sValueId, sValueShape := slLineData[0], slLineData[3] // id, time, status, shape

							uValueShape, err := strconv.ParseUint(sValueShape, 10, 64)
							if err != nil {
								continue
							}

							if uValueShape < 30 {
								uID, err := strconv.ParseUint(sValueId, 10, 64)
								if err != nil {
									continue
								}

								switch {
								case stCond.Operation == "<=" && uID <= uValueCond:
									slUResIds = append(slUResIds, uID)
								case stCond.Operation == ">=" && uID >= uValueCond:
									slUResIds = append(slUResIds, uID)
								case stCond.Operation == "<" && uID < uValueCond:
									slUResIds = append(slUResIds, uID)
								case stCond.Operation == ">" && uID > uValueCond:
									slUResIds = append(slUResIds, uID)
								case stCond.Operation == "=" && uID == uValueCond:
									slUResIds = append(slUResIds, uID)
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

		if stCond.Key == "_time" {
			// folderPath := fmt.Sprintf("%s%s/%s/service", LocalCoreSettings.Storage, tableInfo.Parent, tableInfo.Folder)
			sFolderPath := filepath.Join(StLocalCoreSettings.Storage, stTableInfo.Parent, stTableInfo.Folder, "service")

			uValueCond, err := strconv.ParseUint(stCond.Value, 10, 64)
			if err != nil {
				return []uint64{}
			}

			slFiles, err := os.ReadDir(sFolderPath)
			if err != nil {
				return []uint64{}
			}

			for _, fVal := range slFiles {
				if !fVal.IsDir() {
					sFileName := fVal.Name()
					if strings.Contains(sFileName, stTableInfo.CurrentRev) {
						// fullNameFile := fmt.Sprintf("%s/%s", folderPath, fileName)
						sFullNameFile := filepath.Join(sFolderPath, sFileName)
						sFileText, err := ecowriter.FileRead(sFullNameFile)
						if err != nil {
							continue
						}
						slFileData := strings.Split(sFileText, "\n")
						for _, sLine := range slFileData {
							slLineData := strings.Split(sLine, "|")
							sValueId, sValueTime, sValueShape := slLineData[0], slLineData[1], slLineData[3] // id, time, status, shape

							uValueShape, err := strconv.ParseUint(sValueShape, 10, 64)
							if err != nil {
								continue
							}

							if uValueShape < 30 {
								uID, err := strconv.ParseUint(sValueId, 10, 64)
								if err != nil {
									continue
								}
								uValueTime, err := strconv.ParseUint(sValueTime, 10, 64)
								if err != nil {
									continue
								}

								switch {
								case stCond.Operation == "<=" && uValueTime <= uValueCond:
									slUResIds = append(slUResIds, uID)
								case stCond.Operation == ">=" && uValueTime >= uValueCond:
									slUResIds = append(slUResIds, uID)
								case stCond.Operation == "<" && uValueTime < uValueCond:
									slUResIds = append(slUResIds, uID)
								case stCond.Operation == ">" && uValueTime > uValueCond:
									slUResIds = append(slUResIds, uID)
								case stCond.Operation == "=" && uValueTime == uValueCond:
									slUResIds = append(slUResIds, uID)
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

		if stCond.Key == "_status" {
			// folderPath := fmt.Sprintf("%s%s/%s/service", LocalCoreSettings.Storage, tableInfo.Parent, tableInfo.Folder)
			sFolderPath := filepath.Join(StLocalCoreSettings.Storage, stTableInfo.Parent, stTableInfo.Folder, "service")

			uValueCond, err := strconv.ParseUint(stCond.Value, 10, 64)
			if err != nil {
				return []uint64{}
			}

			slFiles, err := os.ReadDir(sFolderPath)
			if err != nil {
				return []uint64{}
			}

			for _, fVal := range slFiles {
				if !fVal.IsDir() {
					sFileName := fVal.Name()
					if strings.Contains(sFileName, stTableInfo.CurrentRev) {
						// fullNameFile := fmt.Sprintf("%s/%s", folderPath, fileName)
						sFullNameFile := filepath.Join(sFolderPath, sFileName)
						sFileText, err := ecowriter.FileRead(sFullNameFile)
						if err != nil {
							continue
						}
						slFileData := strings.Split(sFileText, "\n")
						for _, sVal := range slFileData {
							slLineData := strings.Split(sVal, "|")
							sValueId, sValueStatus, sValueShape := slLineData[0], slLineData[2], slLineData[3] // id, time, status, shape

							uValueShape, err := strconv.ParseUint(sValueShape, 10, 64)
							if err != nil {
								continue
							}

							if uValueShape < 30 {
								uID, err := strconv.ParseUint(sValueId, 10, 64)
								if err != nil {
									continue
								}
								uValueStatus, err := strconv.ParseUint(sValueStatus, 10, 64)
								if err != nil {
									continue
								}

								switch {
								case stCond.Operation == "<=" && uValueStatus <= uValueCond:
									slUResIds = append(slUResIds, uID)
								case stCond.Operation == ">=" && uValueStatus >= uValueCond:
									slUResIds = append(slUResIds, uID)
								case stCond.Operation == "<" && uValueStatus < uValueCond:
									slUResIds = append(slUResIds, uID)
								case stCond.Operation == ">" && uValueStatus > uValueCond:
									slUResIds = append(slUResIds, uID)
								case stCond.Operation == "=" && uValueStatus == uValueCond:
									slUResIds = append(slUResIds, uID)
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

		if stCond.Key == "_shape" {
			// folderPath := fmt.Sprintf("%s%s/%s/service", LocalCoreSettings.Storage, tableInfo.Parent, tableInfo.Folder)
			sFolderPath := filepath.Join(StLocalCoreSettings.Storage, stTableInfo.Parent, stTableInfo.Folder, "service")

			uValueCond, err := strconv.ParseUint(stCond.Value, 10, 64)
			if err != nil {
				return []uint64{}
			}

			slFiles, err := os.ReadDir(sFolderPath)
			if err != nil {
				return []uint64{}
			}

			for _, fVal := range slFiles {
				if !fVal.IsDir() {
					sFileName := fVal.Name()
					if strings.Contains(sFileName, stTableInfo.CurrentRev) {
						// fullNameFile := fmt.Sprintf("%s/%s", folderPath, fileName)
						sFullNameFile := filepath.Join(sFolderPath, sFileName)
						sFileText, err := ecowriter.FileRead(sFullNameFile)
						if err != nil {
							continue
						}
						slFileData := strings.Split(sFileText, "\n")
						for _, sLine := range slFileData {
							slLineData := strings.Split(sLine, "|")
							sValueId, sValueShape := slLineData[0], slLineData[3] // id, time, status, shape
							uID, err := strconv.ParseUint(sValueId, 10, 64)
							if err != nil {
								continue
							}
							uValueShape, err := strconv.ParseUint(sValueShape, 10, 64)
							if err != nil {
								continue
							}

							switch {
							case stCond.Operation == "<=" && uValueShape <= uValueCond:
								slUResIds = append(slUResIds, uID)
							case stCond.Operation == ">=" && uValueShape >= uValueCond:
								slUResIds = append(slUResIds, uID)
							case stCond.Operation == "<" && uValueShape < uValueCond:
								slUResIds = append(slUResIds, uID)
							case stCond.Operation == ">" && uValueShape > uValueCond:
								slUResIds = append(slUResIds, uID)
							case stCond.Operation == "=" && uValueShape == uValueCond:
								slUResIds = append(slUResIds, uID)
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
		sValueCond := Encode64(stCond.Value)

		stColumnInfo, isOk := stTableInfo.Columns[stCond.Key]
		if !isOk {
			return []uint64{}
		}

		// folderPath := fmt.Sprintf("%s%s/%s", LocalCoreSettings.Storage, columnInfo.Parents, columnInfo.Folder)
		sFolderPath := filepath.Join(StLocalCoreSettings.Storage, stColumnInfo.Parents, stColumnInfo.Folder)
		slFiles, err := os.ReadDir(sFolderPath)
		if err != nil {
			return []uint64{}
		}

		for _, fVal := range slFiles {
			if !fVal.IsDir() {
				sFileName := fVal.Name()
				if strings.Contains(sFileName, stTableInfo.CurrentRev) {
					// fullNameFile := fmt.Sprintf("%s/%s", folderPath, fileName)
					sFullNameFile := filepath.Join(sFolderPath, sFileName)
					sFileText, err := ecowriter.FileRead(sFullNameFile)
					if err != nil {
						continue
					}
					slFileData := strings.Split(sFileText, "\n")
					for _, sLine := range slFileData {
						slLineData := strings.Split(sLine, "|")
						sValueId, sValueData := slLineData[0], slLineData[1] // id, [data]
						uID, err := strconv.ParseUint(sValueId, 10, 64)
						if err != nil {
							continue
						}

						switch {
						case stCond.Operation == "<=" && sValueData <= sValueCond:
							slUProgressIds = append(slUProgressIds, uID)
						case stCond.Operation == ">=" && sValueData >= sValueCond:
							slUProgressIds = append(slUProgressIds, uID)
						case stCond.Operation == "<" && sValueData < sValueCond:
							slUProgressIds = append(slUProgressIds, uID)
						case stCond.Operation == ">" && sValueData > sValueCond:
							slUProgressIds = append(slUProgressIds, uID)
						case stCond.Operation == "=" && sValueData == sValueCond:
							slUProgressIds = append(slUProgressIds, uID)
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
		slices.Sort(slUProgressIds)
		slUProgressIds = slices.Compact(slUProgressIds)
		slUProgressIds = slices.Clip(slUProgressIds)

		// folderSysPath := fmt.Sprintf("%s%s/%s/service", LocalCoreSettings.Storage, tableInfo.Parent, tableInfo.Folder)
		sFolderSysPath := filepath.Join(StLocalCoreSettings.Storage, stTableInfo.Parent, stTableInfo.Folder, "service")

		for _, uIDVal := range slUProgressIds {
			uMaxBucket := Pow(2, stTableInfo.BucketLog)
			uHashID := uIDVal % uMaxBucket
			if uHashID == 0 {
				uHashID = uMaxBucket
			}

			// fullNameFile := fmt.Sprintf("%s/%s_%d", folderSysPath, tableInfo.CurrentRev, hashid)
			sFullNameFile := filepath.Join(sFolderSysPath, fmt.Sprintf("%s_%d", stTableInfo.CurrentRev, uHashID))
			sFileText, err := ecowriter.FileRead(sFullNameFile)
			if err != nil {
				continue
			}

			slFileData := strings.Split(sFileText, "\n")

			isDelete = false

			for _, sLine := range slFileData {
				slLineData := strings.Split(sLine, "|")
				sValueId, sValueShape := slLineData[0], slLineData[3] // id, time, status, shape

				uValue, err := strconv.ParseUint(sValueId, 10, 64)
				if err != nil {
					continue
				}
				uValueShape, err := strconv.ParseUint(sValueShape, 10, 64)
				if err != nil {
					continue
				}

				if uValue == uIDVal && uValueShape == 30 {
					isDelete = true
				}
			}

			if !isDelete {
				slUResIds = append(slUResIds, uIDVal)
			}
		}
	}

	slices.Sort(slUResIds)
	slUResIds = slices.Compact(slUResIds)
	slUResIds = slices.Clip(slUResIds)

	return slUResIds
}

func mergeOr(slUFirst, slUSecond []uint64) []uint64 {
	// This function is complete
	slUResIds := append(slUFirst, slUSecond...)
	slices.Sort(slUResIds)
	slUResIds = slices.Compact(slUResIds)

	slUResIds = slices.Clip(slUResIds)
	return slUResIds
}

func mergeAnd(slUFirst, slUSecond []uint64) []uint64 {
	// This function is complete
	var slUResIds = make([]uint64, 4)

	for _, uElem := range slUSecond {
		if slices.Contains(slUFirst, uElem) {
			slUResIds = append(slUResIds, uElem)
		}
	}

	slUResIds = slices.Clip(slUResIds)
	return slUResIds
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

	mxStorageBlock.Lock() // заменить на другой блок
	defer mxStorageBlock.Unlock()

	dbInfo, okDB := StStorageInfo.DBs[nameDB]
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
	StStorageInfo.DBs[nameDB] = dbInfo

	return result, dbInfo.Save()
}
