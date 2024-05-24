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
		slUResIds            = make([]uint64, 0, 4)
		slUProgressIds       = make([]uint64, 0, 4)
		slUBlacklistIds      = make([]uint64, 0, 4)
		isDelete        bool = false
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
							if len(slLineData) < 4 {
								continue
							}
							sValueId, sValueShape := slLineData[0], slLineData[3] // id, time, status, shape

							uValueShape, err := strconv.ParseUint(sValueShape, 10, 64)
							if err != nil {
								continue
							}

							uID, err := strconv.ParseUint(sValueId, 10, 64)
							if err != nil {
								continue
							}

							if uValueShape == 30 {
								slUBlacklistIds = append(slUBlacklistIds, uID)
							}

							// if uValueShape < 30 {
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
							// }
						}
					}
				}
			}

			slUResIds = slices.Compact(slUResIds)
			for _, uBlackVal := range slUBlacklistIds {
				iInd := slices.Index(slUResIds, uBlackVal)
				if iInd >= 0 {
					slUResIds = slices.Delete(slUResIds, iInd, iInd+1)
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
							if len(slLineData) < 4 {
								continue
							}
							sValueId, sValueTime, sValueShape := slLineData[0], slLineData[1], slLineData[3] // id, time, status, shape

							uValueShape, err := strconv.ParseUint(sValueShape, 10, 64)
							if err != nil {
								continue
							}

							uID, err := strconv.ParseUint(sValueId, 10, 64)
							if err != nil {
								continue
							}
							uValueTime, err := strconv.ParseUint(sValueTime, 10, 64)
							if err != nil {
								continue
							}

							if uValueShape == 30 {
								slUBlacklistIds = append(slUBlacklistIds, uID)
							}

							// if uValueShape < 30 {
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
							// }
						}
					}
				}
			}

			slUResIds = slices.Compact(slUResIds)
			for _, uBlackVal := range slUBlacklistIds {
				iInd := slices.Index(slUResIds, uBlackVal)
				if iInd >= 0 {
					slUResIds = slices.Delete(slUResIds, iInd, iInd+1)
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
							if len(slLineData) < 4 {
								continue
							}
							sValueId, sValueStatus, sValueShape := slLineData[0], slLineData[2], slLineData[3] // id, time, status, shape

							uValueShape, err := strconv.ParseUint(sValueShape, 10, 64)
							if err != nil {
								continue
							}

							uID, err := strconv.ParseUint(sValueId, 10, 64)
							if err != nil {
								continue
							}
							uValueStatus, err := strconv.ParseUint(sValueStatus, 10, 64)
							if err != nil {
								continue
							}

							if uValueShape == 30 {
								slUBlacklistIds = append(slUBlacklistIds, uID)
							}

							// if uValueShape < 30 {
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
							// }
						}
					}
				}
			}

			slUResIds = slices.Compact(slUResIds)
			for _, uBlackVal := range slUBlacklistIds {
				iInd := slices.Index(slUResIds, uBlackVal)
				if iInd >= 0 {
					slUResIds = slices.Delete(slUResIds, iInd, iInd+1)
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
							if len(slLineData) < 4 {
								continue
							}
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
		stColumnInfo, isOk := stTableInfo.Columns[stCond.Key]
		if !isOk {
			return []uint64{}
		}

		sValueCond := Encode64(stCond.Value)

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
						if len(slLineData) < 2 {
							continue
						}
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
				if len(slLineData) < 4 {
					break
				}
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

func whereSelection(slStWhere []gtypes.TConditions, stAdditionalData gtypes.TAdditionalData) []uint64 {
	// This function is complete
	var (
		slUAcc         []uint64 // = make([]uint64, 4)
		slUProgressIds []uint64 // = make([]uint64, 4)
		sSelector      string   = ""
	)

	if len(slStWhere) < 1 {
		return slUAcc
	}

	for _, stElem := range slStWhere {
		switch stElem.Type {
		case "operation":
			// progressIds = nil
			slUProgressIds = findWhereIds(stElem, stAdditionalData)
			switch sSelector {
			case "or":
				slUAcc = mergeOr(slUAcc, slUProgressIds)
			case "and":
				slUAcc = mergeAnd(slUAcc, slUProgressIds)
			default:
				slUAcc = append(slUAcc, slUProgressIds...)
				sSelector = ""
			}
		case "or":
			sSelector = "or"
		case "and":
			sSelector = "and"
		}
	}

	if len(slUAcc) >= 1 {
		if slUAcc[0] == 0 {
			slUAcc = slices.Delete(slUAcc, 0, 1)
		}
	}

	return slUAcc
}

func DeleteRows(sNameDB, sNameTable string, stDeleteIn gtypes.TDeleteStruct) ([]uint64, bool) {
	// This function is complete
	var slUWhereIds []uint64 = []uint64{}
	var slStRowsForStore []gtypes.TRowForStore
	var slSCols []string = []string{}

	if !stDeleteIn.IsWhere {
		stDeleteIn.Where[0] = gtypes.TConditions{
			Type:      "operation",
			Key:       "_id",
			Operation: ">",
			Value:     "0",
		}
		stDeleteIn.IsWhere = true
	}

	stDBInfo, isOkDB := GetDBInfo(sNameDB)
	if !isOkDB {
		return []uint64{}, false
	}
	stTableInfo, isOk := stDBInfo.Tables[sNameTable]
	if !isOk {
		return []uint64{}, false
	}

	// chacking keys
	for _, stWhereElem := range stDeleteIn.Where {
		if stWhereElem.Type == "operation" {
			if stWhereElem.Key != "_id" && stWhereElem.Key != "_time" && stWhereElem.Key != "_status" && stWhereElem.Key != "_shape" {
				_, isOk := stTableInfo.Columns[stWhereElem.Key]
				if !isOk {
					return []uint64{}, false
				}
			}
		}
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
		if uId == 0 {
			continue
		}

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

	return slUWhereIds, true
}

func SelectRows(sNameDB, sNameTable string, stSelectIn gtypes.TSelectStruct) ([]gtypes.TResponseRow, bool) {
	// - It's almost done
	var slReturnedCells = make([]string, 0, 4)
	var slStRowsForResponse = make([]gtypes.TResponseRow, 0, 4)

	if !stSelectIn.IsWhere {
		stBaseCond := gtypes.TConditions{
			Type:      "operation",
			Key:       "_id",
			Operation: ">",
			Value:     "0",
		}
		stSelectIn.Where = append(stSelectIn.Where, stBaseCond)
		stSelectIn.IsWhere = true
	}

	stDBInfo, isOkDB := GetDBInfo(sNameDB)
	if !isOkDB {
		return []gtypes.TResponseRow{}, false
	}
	stTableInfo, isOkTable := stDBInfo.Tables[sNameTable]
	if !isOkTable {
		return []gtypes.TResponseRow{}, false
	}

	// chacking keys
	for _, stWhereElem := range stSelectIn.Where {
		if stWhereElem.Type == "operation" {
			if stWhereElem.Key != "_id" && stWhereElem.Key != "_time" && stWhereElem.Key != "_status" && stWhereElem.Key != "_shape" {
				_, isOk := stTableInfo.Columns[stWhereElem.Key]
				if !isOk {
					return []gtypes.TResponseRow{}, false
				}
			}
		}
	}
	stAdditionalData := gtypes.TAdditionalData{
		Db:    sNameDB,
		Table: sNameTable,
	}
	slUWhereIds := whereSelection(stSelectIn.Where, stAdditionalData)

	for _, sColVal := range stSelectIn.Columns {
		switch sColVal {
		case "*":
			for _, stColVal := range stTableInfo.Columns {
				slReturnedCells = append(slReturnedCells, stColVal.Name)
			}
		case "_id", "_time", "_status", "_shape":
			continue
		default:
			slReturnedCells = append(slReturnedCells, sColVal)
		}

		// if sColVal == "*" {
		// 	for _, stColVal := range stTableInfo.Columns {
		// 		slReturnedCells = append(slReturnedCells, stColVal.Name)
		// 	}
		// } else {
		// 	slReturnedCells = append(slReturnedCells, sColVal)
		// }
	}
	slReturnedCells = slices.Compact(slReturnedCells)
	slReturnedCells = slices.Clip(slReturnedCells)

	// Selection by IDs
	for _, uId := range slUWhereIds {
		if uId == 0 {
			continue
		}

		var stRowForResponse = make(gtypes.TResponseRow, 0)

		time, status, shape, isOk := GetInfoById(sNameDB, sNameTable, uId)
		if !isOk {
			continue
		}

		stRowForResponse["_id"] = fmt.Sprint(uId)
		stRowForResponse["_time"] = time
		stRowForResponse["_status"] = status
		stRowForResponse["_shape"] = shape

		// TODO: Make an identifier generator for the cache.
		for _, sCol := range slReturnedCells {
			sValue, isOkVal := GetColumnById(sNameDB, sNameTable, sCol, uId)
			if !isOkVal {
				return slStRowsForResponse, false
			}
			stRowForResponse[sCol] = Decode64(sValue)
		}

		slStRowsForResponse = append(slStRowsForResponse, stRowForResponse)
	}

	return slStRowsForResponse, true
}

func UpdateRows(sNameDB, sNameTable string, stUpdateIn gtypes.TUpdaateStruct) ([]uint64, bool) {
	// This function is complete
	var slUWhereIds []uint64 = []uint64{}
	var slURowsForStore []gtypes.TRowForStore
	var slSCols []string = []string{}
	var sValue string = ""

	stDBInfo, isOkDB := GetDBInfo(sNameDB)
	if !isOkDB {
		return []uint64{}, false
	}
	stTableInfo, isOk := stDBInfo.Tables[sNameTable]
	if !isOk {
		return []uint64{}, false
	}

	// chacking keys
	for _, stWhereElem := range stUpdateIn.Where {
		if stWhereElem.Type == "operation" {
			if stWhereElem.Key != "_id" && stWhereElem.Key != "_time" && stWhereElem.Key != "_status" && stWhereElem.Key != "_shape" {
				_, isOk := stTableInfo.Columns[stWhereElem.Key]
				if !isOk {
					return []uint64{}, false
				}
			}
		}
	}

	for _, stCol := range stTableInfo.Columns {
		slSCols = append(slSCols, stCol.Name)
	}

	stAdditionalData := gtypes.TAdditionalData{
		Db:    sNameDB,
		Table: sNameTable,
	}

	slUWhereIds = whereSelection(stUpdateIn.Where, stAdditionalData)

	dtNow := time.Now().Unix()

	// Updating by changing the status of records and setting new values
	for _, uId := range slUWhereIds {
		if uId == 0 {
			continue
		}

		var stRowStore = gtypes.TRowForStore{}
		stRowStore.Id = uId
		stRowStore.Time = dtNow
		stRowStore.Status = 0
		stRowStore.Shape = 20 // this is code of update
		stRowStore.DB = sNameDB
		stRowStore.Table = sNameTable
		for _, sCol := range slSCols {
			sNewValue, isOk := stUpdateIn.Couples[sCol]
			if isOk {
				sValue = Encode64(sNewValue)
			} else {
				sValue, _ = GetColumnById(sNameDB, sNameTable, sCol, uId)
			}
			stRowStore.Row = append(stRowStore.Row, gtypes.TColumnForStore{
				Field: sCol,
				Value: sValue,
			})
		}

		slURowsForStore = append(slURowsForStore, stRowStore)
	}

	if len(slUWhereIds) > 0 {
		go InsertIntoBuffer(slURowsForStore)
	}

	return slUWhereIds, true
}

func InsertRows(sNameDB, sNameTable string, slSColumns []string, slSlSRowsIn [][]string) ([]uint64, bool) {
	// This function is complete
	var slSLSRows [][]string

	for _, sCol := range slSColumns {
		if sCol == "_id" || sCol == "_time" || sCol == "_status" || sCol == "_shape" {
			return nil, false
		}
	}

	mxStorageBlock.Lock() // заменить на другой блок
	defer mxStorageBlock.Unlock()

	stDBInfo, isOkDB := StStorageInfo.DBs[sNameDB]
	if !isOkDB {
		return nil, false
	}

	stTableInfo, isOkTable := stDBInfo.Tables[sNameTable]
	if !isOkTable {
		return nil, false
	}

	for _, sCol := range slSColumns {
		if !slices.Contains(stTableInfo.Order, sCol) {
			return nil, false
		}
	}

	dtNow := time.Now().Unix()

	var slStRowsForStore []gtypes.TRowForStore

	for _, slSRow := range slSlSRowsIn {
		var slSTRow []string

		iLCol := len(slSColumns)
		iLRow := len(slSRow)

		if iLCol != iLRow {
			if iLCol < iLRow {
				slSTRow = slSRow[:iLCol]
			}
			if iLCol > iLRow {
				slSTRow = slSRow
				for i := iLRow; i < iLCol; i++ {
					slSTRow = append(slSTRow, stTableInfo.Columns[slSColumns[i]].Specification.Default)
				}
			}
		} else {
			slSTRow = slSRow
		}
		slSLSRows = append(slSLSRows, slSTRow)
	}

	for _, slSRow := range slSLSRows {
		var stRowStore = gtypes.TRowForStore{}
		stTableInfo.Count++
		for _, stColumn := range stTableInfo.Columns {
			var stColStore = gtypes.TColumnForStore{}

			stColStore.Field = stColumn.Name

			var sVStore string
			iInd := slices.Index(slSColumns, stColumn.Name)
			if iInd != -1 {
				sVStore = Encode64(slSRow[iInd])
			} else {
				if stColumn.Specification.NotNull {
					return nil, false
				}
				sVStore = stColumn.Specification.Default
			}
			// colStore.Id = tableInfo.Count
			// colStore.Time = tNow
			stColStore.Value = sVStore
			stRowStore.Row = append(stRowStore.Row, stColStore)
		}
		stRowStore.Id = stTableInfo.Count
		stRowStore.Time = dtNow
		stRowStore.Status = 0
		stRowStore.Shape = 0
		stRowStore.DB = sNameDB
		stRowStore.Table = sNameTable

		slStRowsForStore = append(slStRowsForStore, stRowStore)
	}

	var slUResult []uint64
	for _, stRow := range slStRowsForStore {
		slUResult = append(slUResult, stRow.Id)
	}

	go InsertIntoBuffer(slStRowsForStore)

	stDBInfo.Tables[sNameTable] = stTableInfo
	StStorageInfo.DBs[sNameDB] = stDBInfo

	return slUResult, stDBInfo.Save()
}
