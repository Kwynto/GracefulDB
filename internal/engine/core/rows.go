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
	"github.com/Kwynto/GracefulDB/internal/engine/basicsystem/instead"
)

// type tMValues map[uint64]string
// type tMReverseValues map[string]uint64

func findWhereIds(stCond gtypes.TConditions, stAdditionalData gtypes.TAdditionalData) []uint64 {
	// This function is complete
	var (
		slUResIds            = make([]uint64, 0, 4)
		slUProgressIds       = make([]uint64, 0, 4)
		slUBlacklistIds      = make([]uint64, 0, 4)
		isDelete        bool = false
	)

	if stAdditionalData.Stamp <= 0 {
		stAdditionalData.Stamp = time.Now().Unix()
	}

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
						sFullNameFile := filepath.Join(sFolderPath, sFileName)

						slCache, isOkFile := instead.LoadFile(sFullNameFile, stAdditionalData.Stamp)
						if !isOkFile {
							continue
						}

						for _, slLine := range slCache {
							if len(slLine) < 4 {
								continue
							}
							sValueId, sValueShape := slLine[0], slLine[3] // id, time, status, shape

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

			slUResIds = slices.Compact(slUResIds)
			for _, uBlackVal := range slUBlacklistIds {
				iInd := slices.Index(slUResIds, uBlackVal)
				if iInd >= 0 {
					slUResIds = slices.Delete(slUResIds, iInd, iInd+1)
				}
			}
		}

		if stCond.Key == "_time" {
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
						sFullNameFile := filepath.Join(sFolderPath, sFileName)

						slCache, isOkFile := instead.LoadFile(sFullNameFile, stAdditionalData.Stamp)
						if !isOkFile {
							continue
						}

						for _, slLine := range slCache {
							if len(slLine) < 4 {
								continue
							}
							sValueId, sValueTime, sValueShape := slLine[0], slLine[1], slLine[3] // id, time, status, shape

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

			slUResIds = slices.Compact(slUResIds)
			for _, uBlackVal := range slUBlacklistIds {
				iInd := slices.Index(slUResIds, uBlackVal)
				if iInd >= 0 {
					slUResIds = slices.Delete(slUResIds, iInd, iInd+1)
				}
			}
		}

		if stCond.Key == "_status" {
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
						sFullNameFile := filepath.Join(sFolderPath, sFileName)

						slCache, isOkFile := instead.LoadFile(sFullNameFile, stAdditionalData.Stamp)
						if !isOkFile {
							continue
						}

						for _, slLine := range slCache {
							if len(slLine) < 4 {
								continue
							}
							sValueId, sValueStatus, sValueShape := slLine[0], slLine[2], slLine[3] // id, time, status, shape

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

			slUResIds = slices.Compact(slUResIds)
			for _, uBlackVal := range slUBlacklistIds {
				iInd := slices.Index(slUResIds, uBlackVal)
				if iInd >= 0 {
					slUResIds = slices.Delete(slUResIds, iInd, iInd+1)
				}
			}
		}

		if stCond.Key == "_shape" {
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
						sFullNameFile := filepath.Join(sFolderPath, sFileName)

						slCache, isOkFile := instead.LoadFile(sFullNameFile, stAdditionalData.Stamp)
						if !isOkFile {
							continue
						}

						for _, slLine := range slCache {
							if len(slLine) < 4 {
								continue
							}
							sValueId, sValueShape := slLine[0], slLine[3] // id, time, status, shape

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

		sFolderPath := filepath.Join(StLocalCoreSettings.Storage, stColumnInfo.Parents, stColumnInfo.Folder)
		slFiles, err := os.ReadDir(sFolderPath)
		if err != nil {
			return []uint64{}
		}

		for _, fVal := range slFiles {
			if !fVal.IsDir() {
				sFileName := fVal.Name()
				if strings.Contains(sFileName, stTableInfo.CurrentRev) {
					sFullNameFile := filepath.Join(sFolderPath, sFileName)

					slCache, isOkFile := instead.LoadFile(sFullNameFile, stAdditionalData.Stamp)
					if !isOkFile {
						continue
					}

					for _, slLine := range slCache {
						if len(slLine) < 2 {
							continue
						}
						sValueId, sValueData := slLine[0], slLine[1] // id, [data]

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

		sFolderSysPath := filepath.Join(StLocalCoreSettings.Storage, stTableInfo.Parent, stTableInfo.Folder, "service")

		for _, uIDVal := range slUProgressIds {
			uMaxBucket := Pow(2, stTableInfo.BucketLog)
			uHashID := uIDVal % uMaxBucket
			if uHashID == 0 {
				uHashID = uMaxBucket
			}

			sFullNameFile := filepath.Join(sFolderSysPath, fmt.Sprintf("%s_%d", stTableInfo.CurrentRev, uHashID))

			slCache, isOkFile := instead.LoadFile(sFullNameFile, stAdditionalData.Stamp)
			if !isOkFile {
				continue
			}

			isDelete = false

			for _, slLine := range slCache {
				if len(slLine) < 4 {
					continue
				}
				sValueId, sValueShape := slLine[0], slLine[3] // id, time, status, shape

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

func WhereSelection(slStWhere []gtypes.TConditions, stAdditionalData gtypes.TAdditionalData) []uint64 {
	// This function is complete
	var (
		slUAcc                = make([]uint64, 0, 4)
		slUProgressIds        = make([]uint64, 0, 4)
		sSelector      string = ""
	)

	if len(slStWhere) < 1 {
		return slUAcc
	}

	for _, stElem := range slStWhere {
		switch stElem.Type {
		case "operation":
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

func OrderByVQL(uIds []uint64, stOrderByExp gtypes.TOrderBy, stAdditionalData gtypes.TAdditionalData) []uint64 {
	// This function is complete
	var (
		mValues         = make(map[uint64]string)
		mReversValues   = make(map[string]uint64)
		slSortingString = make([]string, 0, 4)
		slUResIds       = make([]uint64, 0, 4)
	)

	if !stOrderByExp.Is {
		return uIds
	}

	sKey := stOrderByExp.Cols[0]
	uSort := stOrderByExp.Sort[0]

	if uSort == 0 {
		return uIds
	}

	if sKey == "_id" {
		switch uSort {
		case 1:
			slices.Sort(uIds)
		case 2:
			slices.Sort(uIds)
			slices.Reverse(uIds)
		}
		return uIds
	}

	if stAdditionalData.Stamp <= 0 {
		stAdditionalData.Stamp = time.Now().Unix()
	}

	stDBInfo, isOk := GetDBInfo(stAdditionalData.Db)
	if !isOk {
		return uIds
	}
	stTableInfo := stDBInfo.Tables[stAdditionalData.Table]

	if sKey == "_time" || sKey == "_status" || sKey == "_shape" {
		sFolderPath := filepath.Join(StLocalCoreSettings.Storage, stTableInfo.Parent, stTableInfo.Folder, "service")

		slFiles, err := os.ReadDir(sFolderPath)
		if err != nil {
			return uIds
		}

		for _, fVal := range slFiles {
			if !fVal.IsDir() {
				sFileName := fVal.Name()
				if strings.Contains(sFileName, stTableInfo.CurrentRev) {
					sFullNameFile := filepath.Join(sFolderPath, sFileName)

					slCache, isOkFile := instead.LoadFile(sFullNameFile, stAdditionalData.Stamp)
					if !isOkFile {
						continue
					}

					for _, slLine := range slCache {
						if len(slLine) < 4 {
							continue
						}
						sValueId, sValueTime, sValueStatus, sValueShape := slLine[0], slLine[1], slLine[2], slLine[3] // id, time, status, shape

						uValueShape, err := strconv.ParseUint(sValueShape, 10, 64)
						if err != nil {
							continue
						}

						uID, err := strconv.ParseUint(sValueId, 10, 64)
						if err != nil {
							continue
						}

						for _, idVal := range uIds {
							if idVal == uID {
								if uValueShape != 30 {
									switch sKey {
									case "_time":
										mValues[idVal] = sValueTime
									case "_status":
										mValues[idVal] = sValueStatus
									case "_shape":
										mValues[idVal] = sValueShape
									}
								}
								break
							}
						}
					}
				}
			}
		}
	} else {
		stColumnInfo, isOk := stTableInfo.Columns[sKey]
		if !isOk {
			return uIds
		}

		sFolderPath := filepath.Join(StLocalCoreSettings.Storage, stColumnInfo.Parents, stColumnInfo.Folder)

		slFiles, err := os.ReadDir(sFolderPath)
		if err != nil {
			return uIds
		}

		for _, fVal := range slFiles {
			if !fVal.IsDir() {
				sFileName := fVal.Name()
				if strings.Contains(sFileName, stTableInfo.CurrentRev) {
					sFullNameFile := filepath.Join(sFolderPath, sFileName)

					slCache, isOkFile := instead.LoadFile(sFullNameFile, stAdditionalData.Stamp)
					if !isOkFile {
						continue
					}

					for _, slLine := range slCache {
						if len(slLine) < 2 {
							continue
						}
						sValueId, sValueData := slLine[0], slLine[1] // id, [data]

						uID, err := strconv.ParseUint(sValueId, 10, 64)
						if err != nil {
							continue
						}

						for _, idVal := range uIds {
							if idVal == uID {
								mValues[idVal] = sValueData
								break
							}
						}
					}
				}
			}
		}
	}

	// Ordering
	for uId, sVal := range mValues {
		mReversValues[sVal] = uId
		slSortingString = append(slSortingString, sVal)
	}

	switch uSort {
	case 1:
		slices.Sort(slSortingString)
	case 2:
		slices.Sort(slSortingString)
		slices.Reverse(slSortingString)
	}

	for _, v := range slSortingString {
		id := mReversValues[v]
		slUResIds = append(slUResIds, id)
	}

	return slUResIds
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

	dtNow := time.Now().Unix()

	stAdditionalData := gtypes.TAdditionalData{
		Db:    sNameDB,
		Table: sNameTable,
		Stamp: dtNow,
	}

	slUWhereIds = WhereSelection(stDeleteIn.Where, stAdditionalData)

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
	// NOTE: for SQL
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
		Stamp: time.Now().Unix(),
	}
	slUWhereIds := WhereSelection(stSelectIn.Where, stAdditionalData)

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

	}
	slReturnedCells = slices.Compact(slReturnedCells)
	slReturnedCells = slices.Clip(slReturnedCells)

	// Selection by IDs
	for _, uId := range slUWhereIds {
		if uId == 0 {
			continue
		}

		var stRowForResponse = make(gtypes.TResponseRow, 0)

		time, status, shape, isOk := GetInfoById(uId, stAdditionalData)
		if !isOk {
			continue
		}

		stRowForResponse["_id"] = fmt.Sprint(uId)
		stRowForResponse["_time"] = time
		stRowForResponse["_status"] = status
		stRowForResponse["_shape"] = shape

		for _, sCol := range slReturnedCells {
			sValue, isOkVal := GetColumnById(sCol, uId, stAdditionalData)
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
	var (
		slUWhereIds     []uint64 = []uint64{}
		slURowsForStore []gtypes.TRowForStore
		slSCols         []string = []string{}
		sValue          string   = ""
	)

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

	dtNow := time.Now().Unix()

	stAdditionalData := gtypes.TAdditionalData{
		Db:    sNameDB,
		Table: sNameTable,
		Stamp: dtNow,
	}

	slUWhereIds = WhereSelection(stUpdateIn.Where, stAdditionalData)

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
				sValue, _ = GetColumnById(sCol, uId, stAdditionalData)
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
