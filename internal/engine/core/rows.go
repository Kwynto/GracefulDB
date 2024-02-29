package core

import (
	"slices"
	"time"
)

func InsertRows(nameDB, nameTable string, columns []string, rows [][]string) ([]uint64, bool) {
	for _, row := range rows {
		if len(columns) != len(row) {
			return nil, false
		}
	}

	for _, col := range columns {
		if col == "_id" || col == "_time" {
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

	var rowsForStore []tRowForStore

	for _, row := range rows {
		var rowStore = tRowForStore{}
		tableInfo.Count++
		rowCount := tableInfo.Count
		for _, column := range tableInfo.Columns {
			var colStore = tColumnForStore{}

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
			colStore.Id = rowCount
			colStore.Time = tNow
			colStore.Value = vStore
			rowStore.Row = append(rowStore.Row, colStore)
		}
		rowStore.Id = rowCount
		rowStore.Time = tNow

		rowsForStore = append(rowsForStore, rowStore)
	}

	var result []uint64
	for _, row := range rowsForStore {
		result = append(result, row.Id)
	}

	go InsertForBuffer(nameDB, nameTable, rowsForStore)

	dbInfo.Tables[nameTable] = tableInfo
	StorageInfo.DBs[nameDB] = dbInfo

	return result, dbInfo.Save()
}
