package core

func InsertRows(nameDB, nameTable string, columns []string, rows [][]string) bool {
	storageBlock.RLock()
	dbInfo, okDB := StorageInfo.DBs[nameDB]
	storageBlock.RUnlock()
	if !okDB {
		return false
	}

	tableInfo, okTable := dbInfo.Tables[nameTable]
	if !okTable {
		return false
	}

	// for _, row := range rowsIn {
	// 	if len(columnsIn) == len(row) {
	// 	}
	// }

	return true
}
