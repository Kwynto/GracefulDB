package core

import (
	"slices"
	"time"

	"github.com/Kwynto/GracefulDB/internal/engine/basicsystem/gtypes"
)

// func whereOper(cond gtypes.TConditions) []uint64 {
// 	// -

// 	return []uint64{}
// }

func whereSelection(acc []uint64, where []gtypes.TConditions) []uint64 {
	// -
	if len(where) < 1 {
		return acc
	}

	// head := where[0]
	// tail := where[1:]

	// switch head.Type {
	// case "operation":

	// 	resIds := whereOper(head)
	// 	acc = append(acc, resIds...)
	// case "or":
	// 	// resIds := whereOr(head)
	// 	// acc = append(acc, resIds...)
	// case "and":
	// 	// resIds := whereAnd(head)
	// 	// acc = append(acc, resIds...)
	// }

	// return whereSelection(acc, tail)
	return []uint64{}
}

func DeleteRows(nameDB, nameTable string, deleteIn gtypes.TDeleteStruct) ([]uint64, bool) {
	// - ! It's almost done
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

	whereIds = whereSelection(whereIds, deleteIn.Where) // TODO: do it

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

	go InsertIntoBuffer(rowsForStore)

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
