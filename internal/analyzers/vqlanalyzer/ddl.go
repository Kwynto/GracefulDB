package vqlanalyzer

import (
	"errors"
	"strings"

	"github.com/Kwynto/GracefulDB/internal/engine/basicsystem/gauth"
	"github.com/Kwynto/GracefulDB/internal/engine/basicsystem/gtypes"
	"github.com/Kwynto/GracefulDB/internal/engine/basicsystem/vqlexp"
	"github.com/Kwynto/GracefulDB/internal/engine/core"
	"github.com/Kwynto/GracefulDB/pkg/lib/e"
	"github.com/Kwynto/GracefulDB/pkg/lib/ecowriter"
)

// DDL — Data Definition Language (язык определения данных)

func (q tQuery) DDLCreateDB() (result string, err error) {
	// This method is complete
	var stRes gtypes.TResponse

	// Pre checking

	if q.Ticket == "" {
		return `{"state":"error", "result":"an empty ticket"}`, errors.New("an empty ticket")
	}

	sLogin, stAccess, sNewTicket, err := gauth.CheckTicket(q.Ticket)
	if err != nil {
		return `{"state":"error", "result":"authorization failed"}`, err
	}

	if stAccess.Status.IsBad() {
		return `{"state":"error", "result":"auth error"}`, errors.New("auth error")
	}

	if sNewTicket != "" {
		stRes.Ticket = sNewTicket
	}

	// Parsing an expression - Begin

	isINE := vqlexp.RegExpCollection["IfNotExistsWord"].MatchString(q.Instruction)

	sDB := vqlexp.RegExpCollection["CreateDatabaseWord"].ReplaceAllLiteralString(q.Instruction, "")
	if isINE {
		sDB = vqlexp.RegExpCollection["IfNotExistsWord"].ReplaceAllLiteralString(sDB, "")
	}
	sDB = strings.TrimSpace(sDB)
	sDB = trimQuotationMarks(sDB)

	// Parsing an expression - End

	// Post checking

	_, isOk := core.GetDBInfo(sDB)
	if isOk {
		if isINE {
			stRes.State = "error"
			stRes.Result = "the database exists"
			return ecowriter.EncodeJSON(stRes), errors.New("the database exists")
		}

		if !core.LocalCoreSettings.FriendlyMode {
			stRes.State = "error"
			stRes.Result = "the database exists"
			return ecowriter.EncodeJSON(stRes), errors.New("the database exists")
		}

		stDBAccess, isOk := core.GetDBAccess(sDB)
		if isOk {
			if stDBAccess.Owner != sLogin {
				var isLuxUser bool = false
				for iRole := range stAccess.Roles {
					if iRole == int(gauth.ADMIN) || iRole == int(gauth.ENGINEER) {
						isLuxUser = true
						break
					}
				}
				if !isLuxUser {
					return `{"state":"error", "result":"not enough rights"}`, errors.New("not enough rights")
				}
			}
		}

		if !core.RemoveDB(sDB) {
			stRes.State = "error"
			stRes.Result = "the database cannot be deleted"
			return ecowriter.EncodeJSON(stRes), errors.New("the database cannot be deleted")
		}
	}

	// Execution

	if !core.CreateDB(sDB, sLogin, true) {
		stRes.State = "error"
		stRes.Result = "invalid database name"
		return ecowriter.EncodeJSON(stRes), errors.New("invalid database name")
	}

	stRes.State = "ok"
	return ecowriter.EncodeJSON(stRes), nil
}

func (q tQuery) DDLCreateTable() (result string, err error) {
	// This method is complete
	var stRes gtypes.TResponse

	// Pre checking

	sLogin, sDB, stAccess, sNewTicket, err := preChecker(q.Ticket)
	if err != nil {
		stRes.State = "error"
		stRes.Result = err.Error()
		return ecowriter.EncodeJSON(stRes), err
	}

	if sNewTicket != "" {
		stRes.Ticket = sNewTicket
	}

	// Parsing an expression - Begin

	isINE := vqlexp.RegExpCollection["IfNotExistsWord"].MatchString(q.Instruction)

	sTable := vqlexp.RegExpCollection["CreateTableWord"].ReplaceAllLiteralString(q.Instruction, "")
	if isINE {
		sTable = vqlexp.RegExpCollection["IfNotExistsWord"].ReplaceAllLiteralString(sTable, "")
	}

	sColumns := vqlexp.RegExpCollection["TableColumns"].FindString(sTable)
	sColumns = vqlexp.RegExpCollection["TableParenthesis"].ReplaceAllLiteralString(sColumns, "")
	slColumnsIn := vqlexp.RegExpCollection["Comma"].Split(sColumns, -1)

	sTable = vqlexp.RegExpCollection["TableColumns"].ReplaceAllLiteralString(sTable, "")
	sTable = strings.TrimSpace(sTable)
	sTable = trimQuotationMarks(sTable)

	// Parsing an expression - End

	// Post checking, post parsing and execution

	stDBInfo, isOkDB := core.GetDBInfo(sDB)
	if isOkDB {
		var stFlagsAcs gtypes.TAccessFlags
		var isOkFlags bool = false
		var isLuxUser bool = false

		stDBAccess, isOkAccess := core.GetDBAccess(sDB)
		if isOkAccess {
			stFlagsAcs, isOkFlags = stDBAccess.Flags[sLogin]
			if stDBAccess.Owner != sLogin {
				for iRole := range stAccess.Roles {
					if iRole == int(gauth.ADMIN) || iRole == int(gauth.ENGINEER) {
						isLuxUser = true
						break
					}
				}
				if !isLuxUser {
					if !isOkFlags {
						return `{"state":"error", "result":"not enough rights"}`, errors.New("not enough rights")
					}
				}
			} else {
				isLuxUser = true
			}
		} else {
			stRes.State = "error"
			stRes.Result = "internal error"
			return ecowriter.EncodeJSON(stRes), errors.New("internal error")
		}

		_, isOkTable := stDBInfo.Tables[sTable]
		if isOkTable {
			if isINE {
				stRes.State = "error"
				stRes.Result = "the table exists"
				return ecowriter.EncodeJSON(stRes), errors.New("the table exists")
			}

			if !core.LocalCoreSettings.FriendlyMode {
				stRes.State = "error"
				stRes.Result = "the table exists"
				return ecowriter.EncodeJSON(stRes), errors.New("the table exists")
			}

			if !isLuxUser && !(stFlagsAcs.Delete && stFlagsAcs.Create) {
				return `{"state":"error", "result":"not enough rights"}`, errors.New("not enough rights")
			}

			if !core.RemoveTable(sDB, sTable) {
				return `{"state":"error", "result":"not enough rights"}`, errors.New("the table cannot be deleted")
			}
		}

		if !isLuxUser && !stFlagsAcs.Create {
			return `{"state":"error", "result":"not enough rights"}`, errors.New("not enough rights")
		}

		if !core.CreateTable(sDB, sTable, true) {
			stRes.State = "error"
			stRes.Result = "invalid database name or table name"
			return ecowriter.EncodeJSON(stRes), errors.New("invalid database name or table name")
		}

		stDBInfo, _ = core.GetDBInfo(sDB)
		stTableInfo := stDBInfo.Tables[sTable]

		var slColumns = []gtypes.TColumnForWrite{}

		for _, sColumn := range slColumnsIn {
			stCol := gtypes.TColumnForWrite{
				Name: "",
				Spec: gtypes.TColumnSpecification{
					Default: "",
					NotNull: false,
					Unique:  false,
				},
			}
			if vqlexp.RegExpCollection["ColumnUnique"].MatchString(sColumn) {
				sColumn = vqlexp.RegExpCollection["ColumnUnique"].ReplaceAllLiteralString(sColumn, "")
				stCol.Spec.Unique = true
			}
			if vqlexp.RegExpCollection["ColumnNotNull"].MatchString(sColumn) {
				sColumn = vqlexp.RegExpCollection["ColumnNotNull"].ReplaceAllLiteralString(sColumn, "")
				stCol.Spec.NotNull = true
			}
			if vqlexp.RegExpCollection["ColumnDefault"].MatchString(sColumn) {
				sColDef := vqlexp.RegExpCollection["ColumnDefault"].FindString(sColumn)
				sColumn = vqlexp.RegExpCollection["ColumnDefault"].ReplaceAllLiteralString(sColumn, "")

				sColDef = vqlexp.RegExpCollection["ColumnDefaultWord"].ReplaceAllLiteralString(sColDef, "")
				sColDef = strings.TrimSpace(sColDef)
				sColDef = trimQuotationMarks(sColDef)

				if stCol.Spec.Unique {
					stCol.Spec.Default = ""
				} else {
					stCol.Spec.Default = sColDef
				}
			}

			sColumn = vqlexp.RegExpCollection["Spaces"].ReplaceAllLiteralString(sColumn, "")
			sColumn = trimQuotationMarks(sColumn)
			stCol.Name = sColumn

			slColumns = append(slColumns, stCol)
		}

		for _, stColumn := range slColumns {
			if _, isOkCol := stTableInfo.Columns[stColumn.Name]; isOkCol {
				if !core.LocalCoreSettings.FriendlyMode {
					stRes.State = "error"
					stRes.Result = "the column exists"
					return ecowriter.EncodeJSON(stRes), errors.New("the column exists")
				}
				core.RemoveColumn(sDB, sTable, stColumn.Name)
			}

			core.CreateColumn(sDB, sTable, stColumn.Name, true, stColumn.Spec)
		}
	} else {
		stRes.State = "error"
		stRes.Result = "internal error"
		return ecowriter.EncodeJSON(stRes), errors.New("internal error")
	}

	stRes.State = "ok"
	return ecowriter.EncodeJSON(stRes), nil
}

func (q tQuery) DDLCreate() (result string, err error) {
	// This method is complete
	op := "internal -> analyzers -> sql -> DDL -> DDLCreate"
	defer func() { e.Wrapper(op, err) }()

	isDB := vqlexp.RegExpCollection["CreateDatabaseWord"].MatchString(q.Instruction)
	isTable := vqlexp.RegExpCollection["CreateTableWord"].MatchString(q.Instruction)

	if isDB {
		return q.DDLCreateDB()
	} else if isTable {
		return q.DDLCreateTable()
	}

	return `{"state":"error", "result":"unknown command"}`, errors.New("unknown command")
}

func (q tQuery) DDLAlterDB() (result string, err error) {
	// This method is complete
	var res gtypes.TResponse

	// Pre checking

	if q.Ticket == "" {
		return `{"state":"error", "result":"an empty ticket"}`, errors.New("an empty ticket")
	}

	login, access, newticket, err := gauth.CheckTicket(q.Ticket)
	if err != nil {
		return `{"state":"error", "result":"authorization failed"}`, err
	}

	if access.Status.IsBad() {
		return `{"state":"error", "result":"auth error"}`, errors.New("auth error")
	}

	if newticket != "" {
		res.Ticket = newticket
	}

	// Parsing an expression - Begin

	isRT := vqlexp.RegExpCollection["AlterDatabaseRenameTo"].MatchString(q.Instruction)
	if !isRT {
		return `{"state":"error", "result":"invalid command format"}`, errors.New("invalid command format")
	}

	oldDBName := vqlexp.RegExpCollection["AlterDatabaseRenameTo"].FindString(q.Instruction)
	oldDBName = vqlexp.RegExpCollection["AlterDatabaseWord"].ReplaceAllLiteralString(oldDBName, "")
	oldDBName = vqlexp.RegExpCollection["RenameTo"].ReplaceAllLiteralString(oldDBName, "")
	oldDBName = trimQuotationMarks(oldDBName)
	oldDBName = strings.TrimSpace(oldDBName)

	newDBName := vqlexp.RegExpCollection["AlterDatabaseRenameTo"].ReplaceAllLiteralString(q.Instruction, "")
	newDBName = trimQuotationMarks(newDBName)
	newDBName = strings.TrimSpace(newDBName)

	if oldDBName == "" || newDBName == "" {
		return `{"state":"error", "result":"invalid command format"}`, errors.New("invalid command format")
	}

	// Parsing an expression - End

	// Post checking and execution

	_, ok := core.GetDBInfo(oldDBName)
	if ok {
		dbAccess, ok := core.GetDBAccess(oldDBName)
		if ok {
			flagsAcs, okFlags := dbAccess.Flags[login]
			if dbAccess.Owner != login {
				var luxUser bool = false
				for role := range access.Roles {
					if role == int(gauth.ADMIN) || role == int(gauth.ENGINEER) {
						luxUser = true
						break
					}
				}
				if !luxUser {
					if okFlags {
						if !flagsAcs.Update {
							return `{"state":"error", "result":"not enough rights"}`, errors.New("not enough rights")
						}
					} else {
						return `{"state":"error", "result":"not enough rights"}`, errors.New("not enough rights")
					}
				}
			}
			if !core.RenameDB(oldDBName, newDBName, true) {
				res.State = "error"
				res.Result = "the database cannot be renamed"
				return ecowriter.EncodeJSON(res), errors.New("the database cannot be renamed")
			}
		} else {
			res.State = "error"
			res.Result = "internal error"
			return ecowriter.EncodeJSON(res), errors.New("internal error")
		}
	} else {
		res.State = "error"
		res.Result = "invalid database name"
		return ecowriter.EncodeJSON(res), errors.New("invalid database name")
	}

	res.State = "ok"
	return ecowriter.EncodeJSON(res), nil
}

func (q tQuery) DDLAlterTableAdd() (result string, err error) {
	// This method is complete
	var res gtypes.TResponse

	// Pre checking

	login, db, access, newticket, err := preChecker(q.Ticket)
	if err != nil {
		res.State = "error"
		res.Result = err.Error()
		return ecowriter.EncodeJSON(res), err
	}

	if newticket != "" {
		res.Ticket = newticket
	}

	// Parsing an expression - Begin

	tableName := vqlexp.RegExpCollection["AlterTableAdd"].FindString(q.Instruction)
	tableName = vqlexp.RegExpCollection["AlterTableWord"].ReplaceAllLiteralString(tableName, "")
	tableName = vqlexp.RegExpCollection["ADD"].ReplaceAllLiteralString(tableName, "")
	tableName = trimQuotationMarks(tableName)
	tableName = strings.TrimSpace(tableName)

	columnsStr := vqlexp.RegExpCollection["AlterTableAdd"].ReplaceAllLiteralString(q.Instruction, "")
	columnsStr = vqlexp.RegExpCollection["TableParenthesis"].ReplaceAllLiteralString(columnsStr, "")
	columnsIn := vqlexp.RegExpCollection["Comma"].Split(columnsStr, -1)

	var columns = []gtypes.TColumnForWrite{}

	for _, column := range columnsIn {
		col := gtypes.TColumnForWrite{
			Name: "",
			Spec: gtypes.TColumnSpecification{
				Default: "",
				NotNull: false,
				Unique:  false,
			},
		}
		if vqlexp.RegExpCollection["ColumnUnique"].MatchString(column) {
			column = vqlexp.RegExpCollection["ColumnUnique"].ReplaceAllLiteralString(column, "")
			col.Spec.Unique = true
		}
		if vqlexp.RegExpCollection["ColumnNotNull"].MatchString(column) {
			column = vqlexp.RegExpCollection["ColumnNotNull"].ReplaceAllLiteralString(column, "")
			col.Spec.NotNull = true
		}
		if vqlexp.RegExpCollection["ColumnDefault"].MatchString(column) {
			ColDef := vqlexp.RegExpCollection["ColumnDefault"].FindString(column)
			column = vqlexp.RegExpCollection["ColumnDefault"].ReplaceAllLiteralString(column, "")

			ColDef = vqlexp.RegExpCollection["ColumnDefaultWord"].ReplaceAllLiteralString(ColDef, "")
			ColDef = strings.TrimSpace(ColDef)
			ColDef = trimQuotationMarks(ColDef)

			if col.Spec.Unique {
				col.Spec.Default = ""
			} else {
				col.Spec.Default = ColDef
			}
		}

		column = vqlexp.RegExpCollection["Spaces"].ReplaceAllLiteralString(column, "")
		column = trimQuotationMarks(column)
		col.Name = column

		columns = append(columns, col)
	}

	if len(columns) < 1 {
		return `{"state":"error", "result":"invalid command format"}`, errors.New("invalid command format")
	}

	// Parsing an expression - End

	// Post checking

	luxUser, flagsAcs, err := dourPostChecker(db, tableName, login, access)
	if err != nil {
		res.State = "error"
		res.Result = err.Error()
		return ecowriter.EncodeJSON(res), err
	}

	if !luxUser {
		if !(flagsAcs.Alter && flagsAcs.Create) {
			return `{"state":"error", "result":"not enough rights"}`, errors.New("not enough rights")
		}
	}

	// Request execution

	for _, colName := range columns {
		if !core.CreateColumn(db, tableName, colName.Name, true, colName.Spec) {
			res.State = "error"
			res.Result = "the column cannot be added"
			return ecowriter.EncodeJSON(res), errors.New("the column cannot be added")
		}
	}

	res.State = "ok"
	return ecowriter.EncodeJSON(res), nil
}

func (q tQuery) DDLAlterTableDrop() (result string, err error) {
	// This method is complete
	var res gtypes.TResponse

	// Pre checking

	login, db, access, newticket, err := preChecker(q.Ticket)
	if err != nil {
		res.State = "error"
		res.Result = err.Error()
		return ecowriter.EncodeJSON(res), err
	}

	if newticket != "" {
		res.Ticket = newticket
	}

	// Parsing an expression - Begin

	tableName := vqlexp.RegExpCollection["AlterTableDrop"].FindString(q.Instruction)
	tableName = vqlexp.RegExpCollection["AlterTableWord"].ReplaceAllLiteralString(tableName, "")
	tableName = vqlexp.RegExpCollection["DROP"].ReplaceAllLiteralString(tableName, "")
	tableName = trimQuotationMarks(tableName)
	tableName = strings.TrimSpace(tableName)

	if tableName == "" {
		return `{"state":"error", "result":"invalid command format"}`, errors.New("invalid command format")
	}

	columnsStr := vqlexp.RegExpCollection["AlterTableDrop"].ReplaceAllLiteralString(q.Instruction, "")
	columnsStr = vqlexp.RegExpCollection["TableParenthesis"].ReplaceAllLiteralString(columnsStr, "")
	columnsIn := vqlexp.RegExpCollection["Comma"].Split(columnsStr, -1)

	var columns = []string{}

	for _, column := range columnsIn {
		column = vqlexp.RegExpCollection["Spaces"].ReplaceAllLiteralString(column, "")
		column = trimQuotationMarks(column)

		columns = append(columns, column)
	}

	if len(columns) < 1 {
		return `{"state":"error", "result":"invalid command format"}`, errors.New("invalid command format")
	}

	// Parsing an expression - End

	// Post checking

	luxUser, flagsAcs, err := dourPostChecker(db, tableName, login, access)
	if err != nil {
		res.State = "error"
		res.Result = err.Error()
		return ecowriter.EncodeJSON(res), err
	}

	if !luxUser {
		if !(flagsAcs.Alter && flagsAcs.Drop) {
			return `{"state":"error", "result":"not enough rights"}`, errors.New("not enough rights")
		}
	}

	// Request execution
	for _, colName := range columns {
		if !core.RemoveColumn(db, tableName, colName) {
			res.State = "error"
			res.Result = "the column cannot be deleted"
			return ecowriter.EncodeJSON(res), errors.New("the column cannot be deleted")
		}
	}

	res.State = "ok"
	return ecowriter.EncodeJSON(res), nil
}

func (q tQuery) DDLAlterTableModify() (result string, err error) {
	// This method is complete
	var res gtypes.TResponse

	// Pre checking

	login, db, access, newticket, err := preChecker(q.Ticket)
	if err != nil {
		res.State = "error"
		res.Result = err.Error()
		return ecowriter.EncodeJSON(res), err
	}

	if newticket != "" {
		res.Ticket = newticket
	}

	// Parsing an expression - Begin

	tableName := vqlexp.RegExpCollection["AlterTableModify"].FindString(q.Instruction)
	tableName = vqlexp.RegExpCollection["AlterTableWord"].ReplaceAllLiteralString(tableName, "")
	tableName = vqlexp.RegExpCollection["MODIFY"].ReplaceAllLiteralString(tableName, "")
	tableName = trimQuotationMarks(tableName)
	tableName = strings.TrimSpace(tableName)

	columnsStr := vqlexp.RegExpCollection["AlterTableModify"].ReplaceAllLiteralString(q.Instruction, "")
	columnsStr = vqlexp.RegExpCollection["TableParenthesis"].ReplaceAllLiteralString(columnsStr, "")
	columnsIn := vqlexp.RegExpCollection["Comma"].Split(columnsStr, -1)

	var columns = []gtypes.TColumnForWrite{}

	for _, column := range columnsIn {
		col := gtypes.TColumnForWrite{
			Name:    "",
			OldName: "",
			Spec: gtypes.TColumnSpecification{
				Default: "",
				NotNull: false,
				Unique:  false,
			},
			IsChName: false,
			// IsChDefault: false,
			// IsChNotNull: false,
			// IsChUniqut:  false,
		}

		if vqlexp.RegExpCollection["ColumnUnique"].MatchString(column) {
			column = vqlexp.RegExpCollection["ColumnUnique"].ReplaceAllLiteralString(column, "")
			col.Spec.Unique = true
			// col.IsChUniqut = true
		}
		if vqlexp.RegExpCollection["ColumnNotNull"].MatchString(column) {
			column = vqlexp.RegExpCollection["ColumnNotNull"].ReplaceAllLiteralString(column, "")
			col.Spec.NotNull = true
			// col.IsChNotNull = true
		}
		if vqlexp.RegExpCollection["ColumnDefault"].MatchString(column) {
			ColDef := vqlexp.RegExpCollection["ColumnDefault"].FindString(column)
			column = vqlexp.RegExpCollection["ColumnDefault"].ReplaceAllLiteralString(column, "")

			ColDef = vqlexp.RegExpCollection["ColumnDefaultWord"].ReplaceAllLiteralString(ColDef, "")
			ColDef = strings.TrimSpace(ColDef)
			ColDef = trimQuotationMarks(ColDef)

			if col.Spec.Unique {
				col.Spec.Default = ""
			} else {
				col.Spec.Default = ColDef
			}
			// col.IsChDefault = true
		}

		if vqlexp.RegExpCollection["RenameTo"].MatchString(column) {
			names := vqlexp.RegExpCollection["RenameTo"].Split(column, -1)
			oldName := names[0]
			newName := names[1]

			oldName = vqlexp.RegExpCollection["Spaces"].ReplaceAllLiteralString(oldName, "")
			oldName = trimQuotationMarks(oldName)

			newName = vqlexp.RegExpCollection["Spaces"].ReplaceAllLiteralString(newName, "")
			newName = trimQuotationMarks(newName)

			if newName != oldName {
				col.Name = newName
				col.OldName = oldName
				col.IsChName = true
			} else {
				col.Name = oldName
			}
		} else {
			column = vqlexp.RegExpCollection["Spaces"].ReplaceAllLiteralString(column, "")
			column = trimQuotationMarks(column)

			col.Name = column
		}

		if col.Name != "" {
			columns = append(columns, col)
		}
	}

	if len(columns) < 1 {
		return `{"state":"error", "result":"invalid command format"}`, errors.New("invalid command format")
	}

	// Parsing an expression - End

	// Post checking

	luxUser, flagsAcs, err := dourPostChecker(db, tableName, login, access)
	if err != nil {
		res.State = "error"
		res.Result = err.Error()
		return ecowriter.EncodeJSON(res), err
	}

	if !luxUser {
		if !flagsAcs.Alter {
			return `{"state":"error", "result":"not enough rights"}`, errors.New("not enough rights")
		}
	}

	// Request execution

	for _, column := range columns {
		if !core.ChangeColumn(db, tableName, column, true) {
			res.State = "error"
			res.Result = "the column cannot be changed"
			return ecowriter.EncodeJSON(res), errors.New("the column cannot be changed")
		}
	}

	res.State = "ok"
	return ecowriter.EncodeJSON(res), nil
}

func (q tQuery) DDLAlterTableRenameTo() (result string, err error) {
	// This method is complete
	var res gtypes.TResponse

	// Pre checking

	login, db, access, newticket, err := preChecker(q.Ticket)
	if err != nil {
		res.State = "error"
		res.Result = err.Error()
		return ecowriter.EncodeJSON(res), err
	}

	if newticket != "" {
		res.Ticket = newticket
	}

	// Parsing an expression - Begin

	isRT := vqlexp.RegExpCollection["AlterTableRenameTo"].MatchString(q.Instruction)
	if !isRT {
		return `{"state":"error", "result":"invalid command format"}`, errors.New("invalid command format")
	}

	oldTableName := vqlexp.RegExpCollection["AlterTableRenameTo"].FindString(q.Instruction)
	oldTableName = vqlexp.RegExpCollection["AlterTableWord"].ReplaceAllLiteralString(oldTableName, "")
	oldTableName = vqlexp.RegExpCollection["RenameTo"].ReplaceAllLiteralString(oldTableName, "")
	oldTableName = trimQuotationMarks(oldTableName)
	oldTableName = strings.TrimSpace(oldTableName)

	newTableName := vqlexp.RegExpCollection["AlterTableRenameTo"].ReplaceAllLiteralString(q.Instruction, "")
	newTableName = trimQuotationMarks(newTableName)
	newTableName = strings.TrimSpace(newTableName)

	if oldTableName == "" || newTableName == "" {
		return `{"state":"error", "result":"invalid command format"}`, errors.New("invalid command format")
	}

	// Parsing an expression - End

	// Post checking

	luxUser, flagsAcs, err := dourPostChecker(db, oldTableName, login, access)
	if err != nil {
		res.State = "error"
		res.Result = err.Error()
		return ecowriter.EncodeJSON(res), err
	}

	if !luxUser {
		if !flagsAcs.Alter {
			return `{"state":"error", "result":"not enough rights"}`, errors.New("not enough rights")
		}
	}

	// Request execution

	if !core.RenameTable(db, oldTableName, newTableName, true) {
		res.State = "error"
		res.Result = "the database cannot be renamed"
		return ecowriter.EncodeJSON(res), errors.New("the database cannot be renamed")
	}

	res.State = "ok"
	return ecowriter.EncodeJSON(res), nil
}

func (q tQuery) DDLAlterTable() (result string, err error) {
	// This method is complete
	isAdd := vqlexp.RegExpCollection["AlterTableAdd"].MatchString(q.Instruction)
	isDrop := vqlexp.RegExpCollection["AlterTableDrop"].MatchString(q.Instruction)
	isModify := vqlexp.RegExpCollection["AlterTableModify"].MatchString(q.Instruction)
	isRT := vqlexp.RegExpCollection["AlterTableRenameTo"].MatchString(q.Instruction)

	if isAdd {
		return q.DDLAlterTableAdd()
	} else if isDrop {
		return q.DDLAlterTableDrop()
	} else if isModify {
		return q.DDLAlterTableModify()
	} else if isRT {
		return q.DDLAlterTableRenameTo()
	}

	return `{"state":"error", "result":"unknown command"}`, errors.New("unknown command")
}

func (q tQuery) DDLAlter() (result string, err error) {
	// -
	op := "internal -> analyzers -> sql -> DDL -> DDLAlter"
	defer func() { e.Wrapper(op, err) }()

	isDB := vqlexp.RegExpCollection["AlterDatabaseWord"].MatchString(q.Instruction)
	isTable := vqlexp.RegExpCollection["AlterTableWord"].MatchString(q.Instruction)

	if isDB {
		return q.DDLAlterDB()
	} else if isTable {
		return q.DDLAlterTable()
	}

	return `{"state":"error", "result":"unknown command"}`, errors.New("unknown command")
}

func (q tQuery) DDLDropDB() (result string, err error) {
	// This method is complete
	var res gtypes.TResponse

	// Pre checking

	if q.Ticket == "" {
		return `{"state":"error", "result":"an empty ticket"}`, errors.New("an empty ticket")
	}

	login, access, newticket, err := gauth.CheckTicket(q.Ticket)
	if err != nil {
		return `{"state":"error", "result":"authorization failed"}`, err
	}

	if access.Status.IsBad() {
		return `{"state":"error", "result":"auth error"}`, errors.New("auth error")
	}

	if newticket != "" {
		res.Ticket = newticket
	}

	// Parsing an expression - Begin

	isIE := vqlexp.RegExpCollection["IfExistsWord"].MatchString(q.Instruction)

	db := vqlexp.RegExpCollection["DropDatabaseWord"].ReplaceAllLiteralString(q.Instruction, "")
	if isIE {
		db = vqlexp.RegExpCollection["IfExistsWord"].ReplaceAllLiteralString(db, "")
	}
	db = strings.TrimSpace(db)
	db = trimQuotationMarks(db)

	// Parsing an expression - End

	// Post checking

	_, ok := core.GetDBInfo(db)
	if !ok {
		if isIE {
			res.State = "error"
			res.Result = "the database not exists"
			return ecowriter.EncodeJSON(res), errors.New("the database not exists")
		}

		res.State = "ok"
		return ecowriter.EncodeJSON(res), nil
	}

	dbAccess, ok := core.GetDBAccess(db)
	if ok {
		if dbAccess.Owner != login {
			var luxUser bool = false
			for role := range access.Roles {
				if role == int(gauth.ADMIN) || role == int(gauth.ENGINEER) {
					luxUser = true
					break
				}
			}
			if !luxUser {
				return `{"state":"error", "result":"not enough rights"}`, errors.New("not enough rights")
			}
		}
	}

	// Request execution

	if !core.RemoveDB(db) {
		res.State = "error"
		res.Result = "the database cannot be deleted"
		return ecowriter.EncodeJSON(res), errors.New("the database cannot be deleted")
	}

	res.State = "ok"
	return ecowriter.EncodeJSON(res), nil
}

func (q tQuery) DDLDropTable() (result string, err error) {
	// This method is complete
	var res gtypes.TResponse

	// Pre checking

	login, db, access, newticket, err := preChecker(q.Ticket)
	if err != nil {
		res.State = "error"
		res.Result = err.Error()
		return ecowriter.EncodeJSON(res), err
	}

	if newticket != "" {
		res.Ticket = newticket
	}

	// Parsing an expression - Begin

	isIE := vqlexp.RegExpCollection["IfExistsWord"].MatchString(q.Instruction)

	table := vqlexp.RegExpCollection["DropTableWord"].ReplaceAllLiteralString(q.Instruction, "")
	if isIE {
		table = vqlexp.RegExpCollection["IfExistsWord"].ReplaceAllLiteralString(table, "")
	}
	table = strings.TrimSpace(table)
	table = trimQuotationMarks(table)

	// Parsing an expression - End

	// Post checking

	luxUser, flagsAcs, err := dourPostChecker(db, table, login, access)
	if err != nil {
		res.State = "error"
		res.Result = err.Error()
		return ecowriter.EncodeJSON(res), err
	}

	if !luxUser && !flagsAcs.Drop {
		return `{"state":"error", "result":"not enough rights"}`, errors.New("not enough rights")
	}

	// Request execution

	if !core.RemoveTable(db, table) {
		return `{"state":"error", "result":"the table cannot be deleted"}`, errors.New("the table cannot be deleted")
	}

	res.State = "ok"
	return ecowriter.EncodeJSON(res), nil
}

func (q tQuery) DDLDrop() (result string, err error) {
	// This method is complete
	op := "internal -> analyzers -> sql -> DDL -> DDLDrop"
	defer func() { e.Wrapper(op, err) }()

	isDB := vqlexp.RegExpCollection["DropDatabaseWord"].MatchString(q.Instruction)
	isTable := vqlexp.RegExpCollection["DropTableWord"].MatchString(q.Instruction)

	if isDB {
		return q.DDLDropDB()
	} else if isTable {
		return q.DDLDropTable()
	}

	return `{"state":"error", "result":"unknown command"}`, errors.New("unknown command")
}
