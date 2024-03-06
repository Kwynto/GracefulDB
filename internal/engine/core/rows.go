package core

import (
	"slices"
	"time"
)

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

	var rowsForStore []tRowForStore

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
		var rowStore = tRowForStore{}
		tableInfo.Count++
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
			colStore.Id = tableInfo.Count
			colStore.Time = tNow
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
