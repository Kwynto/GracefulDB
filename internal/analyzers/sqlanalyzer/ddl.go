package sqlanalyzer

import (
	"errors"
	"strings"

	"github.com/Kwynto/GracefulDB/internal/engine/basicsystem/gauth"
	"github.com/Kwynto/GracefulDB/internal/engine/basicsystem/gtypes"
	"github.com/Kwynto/GracefulDB/internal/engine/core"
	"github.com/Kwynto/GracefulDB/pkg/lib/e"
	"github.com/Kwynto/GracefulDB/pkg/lib/ecowriter"
)

// DDL — язык определения данных (Data Definition Language)

func (q tQuery) DDLCreateDB() (result string, err error) {
	// This method is complete
	var res gtypes.Response

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

	isINE := core.RegExpCollection["IfNotExistsWord"].MatchString(q.Instruction)

	db := core.RegExpCollection["CreateDatabaseWord"].ReplaceAllLiteralString(q.Instruction, "")
	if isINE {
		db = core.RegExpCollection["IfNotExistsWord"].ReplaceAllLiteralString(db, "")
	}
	db = strings.TrimSpace(db)
	db = core.RegExpCollection["QuotationMarks"].ReplaceAllLiteralString(db, "")
	db = core.RegExpCollection["SpecQuotationMark"].ReplaceAllLiteralString(db, "")

	_, ok := core.StorageInfo.DBs[db]
	if ok {
		if isINE {
			res.State = "error"
			res.Result = "the database exists"
			return ecowriter.EncodeString(res), errors.New("the database exists")
		}

		dbAccess, ok := core.StorageInfo.Access[db]
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

		if !core.RemoveDB(db) {
			res.State = "error"
			res.Result = "the database cannot be deleted"
			return ecowriter.EncodeString(res), errors.New("the database cannot be deleted")
		}
	}

	if !core.CreateDB(db, login, true) {
		res.State = "error"
		res.Result = "invalid database name"
		return ecowriter.EncodeString(res), errors.New("invalid database name")
	}

	res.State = "ok"
	return ecowriter.EncodeString(res), nil
}

func (q tQuery) DDLCreateTable() (result string, err error) {
	// This method is complete
	var res gtypes.Response

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

	isINE := core.RegExpCollection["IfNotExistsWord"].MatchString(q.Instruction)

	table := core.RegExpCollection["CreateTableWord"].ReplaceAllLiteralString(q.Instruction, "")
	if isINE {
		table = core.RegExpCollection["IfNotExistsWord"].ReplaceAllLiteralString(table, "")
	}

	columnsStr := core.RegExpCollection["TableColumns"].FindString(table)
	columnsStr = core.RegExpCollection["TableParenthesis"].ReplaceAllLiteralString(columnsStr, "")
	columnsIn := core.RegExpCollection["Comma"].Split(columnsStr, -1)

	table = core.RegExpCollection["TableColumns"].ReplaceAllLiteralString(table, "")
	table = strings.TrimSpace(table)
	table = core.RegExpCollection["QuotationMarks"].ReplaceAllLiteralString(table, "")
	table = core.RegExpCollection["SpecQuotationMark"].ReplaceAllLiteralString(table, "")

	state, ok := core.States[q.Ticket]
	if !ok {
		res.State = "error"
		res.Result = "unknown database"
		return ecowriter.EncodeString(res), errors.New("unknown database")
	}
	db := state.CurrentDB
	if db == "" {
		res.State = "error"
		res.Result = "no database selected"
		return ecowriter.EncodeString(res), errors.New("no database selected")
	}

	dbInfo, okDB := core.StorageInfo.DBs[db]
	if okDB {
		var flagsAcs gtypes.TAccessFlags
		var okFlags bool = false
		var luxUser bool = false

		dbAccess, okAccess := core.StorageInfo.Access[db]
		if okAccess {
			flagsAcs, okFlags = dbAccess.Flags[login]
			if dbAccess.Owner != login {
				for role := range access.Roles {
					if role == int(gauth.ADMIN) || role == int(gauth.ENGINEER) {
						luxUser = true
						break
					}
				}
				if !luxUser {
					if !okFlags {
						return `{"state":"error", "result":"not enough rights"}`, errors.New("not enough rights")
					}
				}
			} else {
				luxUser = true
			}
		} else {
			res.State = "error"
			res.Result = "internal error"
			return ecowriter.EncodeString(res), errors.New("internal error")
		}

		_, okTable := dbInfo.Tables[table]
		if okTable {
			if isINE {
				res.State = "error"
				res.Result = "the table exists"
				return ecowriter.EncodeString(res), errors.New("the table exists")
			}

			if !luxUser && !(flagsAcs.Delete && flagsAcs.Create) {
				return `{"state":"error", "result":"not enough rights"}`, errors.New("not enough rights")
			}

			if !core.RemoveTable(db, table) {
				return `{"state":"error", "result":"not enough rights"}`, errors.New("the table cannot be deleted")
			}
		}

		if !luxUser && !flagsAcs.Create {
			return `{"state":"error", "result":"not enough rights"}`, errors.New("not enough rights")
		}

		if !core.CreateTable(db, table, true) {
			res.State = "error"
			res.Result = "invalid database name or table name"
			return ecowriter.EncodeString(res), errors.New("invalid database name or table name")
		}

		var columns = []core.TColumnForWrite{}

		for _, column := range columnsIn {
			col := core.TColumnForWrite{
				Name: "",
				Spec: core.TColumnSpecification{
					Default: "",
					NotNull: false,
					Unique:  false,
				},
			}
			if core.RegExpCollection["ColumnUnique"].MatchString(column) {
				column = core.RegExpCollection["ColumnUnique"].ReplaceAllLiteralString(column, "")
				col.Spec.Unique = true
			}
			if core.RegExpCollection["ColumnNotNull"].MatchString(column) {
				column = core.RegExpCollection["ColumnNotNull"].ReplaceAllLiteralString(column, "")
				col.Spec.NotNull = true
			}
			if core.RegExpCollection["ColumnDefault"].MatchString(column) {
				ColDef := core.RegExpCollection["ColumnDefault"].FindString(column)
				column = core.RegExpCollection["ColumnDefault"].ReplaceAllLiteralString(column, "")

				ColDef = core.RegExpCollection["ColumnDefaultWord"].ReplaceAllLiteralString(ColDef, "")
				ColDef = strings.TrimSpace(ColDef)
				ColDef = core.RegExpCollection["QuotationMarks"].ReplaceAllLiteralString(ColDef, "")
				ColDef = core.RegExpCollection["SpecQuotationMark"].ReplaceAllLiteralString(ColDef, "")

				if col.Spec.Unique {
					col.Spec.Default = ""
				} else {
					col.Spec.Default = ColDef
				}
			}

			column = core.RegExpCollection["Spaces"].ReplaceAllLiteralString(column, "")
			column = core.RegExpCollection["QuotationMarks"].ReplaceAllLiteralString(column, "")
			column = core.RegExpCollection["SpecQuotationMark"].ReplaceAllLiteralString(column, "")
			col.Name = column

			columns = append(columns, col)
		}

		for _, column := range columns {
			core.CreateColumn(db, table, column.Name, true, column.Spec)
		}
	} else {
		res.State = "error"
		res.Result = "internal error"
		return ecowriter.EncodeString(res), errors.New("internal error")
	}

	res.State = "ok"
	return ecowriter.EncodeString(res), nil
}

func (q tQuery) DDLCreate() (result string, err error) {
	// This method is complete
	op := "internal -> analyzers -> sql -> DDL -> DDLCreate"
	defer func() { e.Wrapper(op, err) }()

	isDB := core.RegExpCollection["CreateDatabaseWord"].MatchString(q.Instruction)
	isTable := core.RegExpCollection["CreateTableWord"].MatchString(q.Instruction)

	if isDB {
		return q.DDLCreateDB()
	} else if isTable {
		return q.DDLCreateTable()
	}

	return `{"state":"error", "result":"unknown command"}`, errors.New("unknown command")
}

func (q tQuery) DDLAlterDB() (result string, err error) {
	// This method is complete
	var res gtypes.Response

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

	isRT := core.RegExpCollection["AlterDatabaseRenameTo"].MatchString(q.Instruction)

	oldDBName := core.RegExpCollection["AlterDatabaseRenameTo"].FindString(q.Instruction)
	oldDBName = core.RegExpCollection["AlterDatabaseWord"].ReplaceAllLiteralString(oldDBName, "")
	oldDBName = core.RegExpCollection["RenameTo"].ReplaceAllLiteralString(oldDBName, "")
	oldDBName = core.RegExpCollection["QuotationMarks"].ReplaceAllLiteralString(oldDBName, "")
	oldDBName = core.RegExpCollection["SpecQuotationMark"].ReplaceAllLiteralString(oldDBName, "")
	oldDBName = strings.TrimSpace(oldDBName)

	newDBName := core.RegExpCollection["AlterDatabaseRenameTo"].ReplaceAllLiteralString(q.Instruction, "")
	newDBName = core.RegExpCollection["QuotationMarks"].ReplaceAllLiteralString(newDBName, "")
	newDBName = core.RegExpCollection["SpecQuotationMark"].ReplaceAllLiteralString(newDBName, "")
	newDBName = strings.TrimSpace(newDBName)

	if !isRT {
		return `{"state":"error", "result":"invalid command format"}`, errors.New("invalid command format")
	}

	if oldDBName == "" || newDBName == "" {
		return `{"state":"error", "result":"invalid command format"}`, errors.New("invalid command format")
	}

	_, ok := core.StorageInfo.DBs[oldDBName]
	if ok {
		dbAccess, ok := core.StorageInfo.Access[oldDBName]
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
				return ecowriter.EncodeString(res), errors.New("the database cannot be renamed")
			}
		} else {
			res.State = "error"
			res.Result = "internal error"
			return ecowriter.EncodeString(res), errors.New("internal error")
		}
	} else {
		res.State = "error"
		res.Result = "invalid database name"
		return ecowriter.EncodeString(res), errors.New("invalid database name")
	}

	res.State = "ok"
	return ecowriter.EncodeString(res), nil
}

func (q tQuery) DDLAlterTableAdd() (result string, err error) {
	// This method is complete
	var res gtypes.Response

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

	tableName := core.RegExpCollection["AlterTableAdd"].FindString(q.Instruction)
	tableName = core.RegExpCollection["AlterTableWord"].ReplaceAllLiteralString(tableName, "")
	tableName = core.RegExpCollection["ADD"].ReplaceAllLiteralString(tableName, "")
	tableName = core.RegExpCollection["QuotationMarks"].ReplaceAllLiteralString(tableName, "")
	tableName = core.RegExpCollection["SpecQuotationMark"].ReplaceAllLiteralString(tableName, "")
	tableName = strings.TrimSpace(tableName)

	columnsStr := core.RegExpCollection["AlterTableAdd"].ReplaceAllLiteralString(q.Instruction, "")
	columnsStr = core.RegExpCollection["TableParenthesis"].ReplaceAllLiteralString(columnsStr, "")
	columnsIn := core.RegExpCollection["Comma"].Split(columnsStr, -1)

	var columns = []core.TColumnForWrite{}

	for _, column := range columnsIn {
		col := core.TColumnForWrite{
			Name: "",
			Spec: core.TColumnSpecification{
				Default: "",
				NotNull: false,
				Unique:  false,
			},
		}
		if core.RegExpCollection["ColumnUnique"].MatchString(column) {
			column = core.RegExpCollection["ColumnUnique"].ReplaceAllLiteralString(column, "")
			col.Spec.Unique = true
		}
		if core.RegExpCollection["ColumnNotNull"].MatchString(column) {
			column = core.RegExpCollection["ColumnNotNull"].ReplaceAllLiteralString(column, "")
			col.Spec.NotNull = true
		}
		if core.RegExpCollection["ColumnDefault"].MatchString(column) {
			ColDef := core.RegExpCollection["ColumnDefault"].FindString(column)
			column = core.RegExpCollection["ColumnDefault"].ReplaceAllLiteralString(column, "")

			ColDef = core.RegExpCollection["ColumnDefaultWord"].ReplaceAllLiteralString(ColDef, "")
			ColDef = strings.TrimSpace(ColDef)
			ColDef = core.RegExpCollection["QuotationMarks"].ReplaceAllLiteralString(ColDef, "")
			ColDef = core.RegExpCollection["SpecQuotationMark"].ReplaceAllLiteralString(ColDef, "")

			if col.Spec.Unique {
				col.Spec.Default = ""
			} else {
				col.Spec.Default = ColDef
			}
		}

		column = core.RegExpCollection["Spaces"].ReplaceAllLiteralString(column, "")
		column = core.RegExpCollection["QuotationMarks"].ReplaceAllLiteralString(column, "")
		column = core.RegExpCollection["SpecQuotationMark"].ReplaceAllLiteralString(column, "")
		col.Name = column

		columns = append(columns, col)
	}

	if len(columns) < 1 {
		return `{"state":"error", "result":"invalid command format"}`, errors.New("invalid command format")
	}

	state, ok := core.States[q.Ticket]
	if !ok {
		res.State = "error"
		res.Result = "unknown database"
		return ecowriter.EncodeString(res), errors.New("unknown database")
	}
	db := state.CurrentDB
	if db == "" {
		res.State = "error"
		res.Result = "no database selected"
		return ecowriter.EncodeString(res), errors.New("no database selected")
	}

	_, okDB := core.StorageInfo.DBs[db]
	if okDB {
		dbAccess, okAccess := core.StorageInfo.Access[db]
		if okAccess {
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
						if !(flagsAcs.Alter && flagsAcs.Create) {
							return `{"state":"error", "result":"not enough rights"}`, errors.New("not enough rights")
						}
					} else {
						return `{"state":"error", "result":"not enough rights"}`, errors.New("not enough rights")
					}
				}
			}
			for _, colName := range columns {
				if !core.CreateColumn(db, tableName, colName.Name, true, colName.Spec) {
					res.State = "error"
					res.Result = "the column cannot be added"
					return ecowriter.EncodeString(res), errors.New("the column cannot be added")
				}
			}
		} else {
			res.State = "error"
			res.Result = "internal error"
			return ecowriter.EncodeString(res), errors.New("internal error")
		}
	} else {
		res.State = "error"
		res.Result = "invalid database name"
		return ecowriter.EncodeString(res), errors.New("invalid database name")
	}
	// TODO: тута проверка прав и выполнение

	res.State = "ok"
	return ecowriter.EncodeString(res), nil
}

func (q tQuery) DDLAlterTableDrop() (result string, err error) {
	// This method is complete
	var res gtypes.Response

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

	tableName := core.RegExpCollection["AlterTableDrop"].FindString(q.Instruction)
	tableName = core.RegExpCollection["AlterTableWord"].ReplaceAllLiteralString(tableName, "")
	tableName = core.RegExpCollection["DROP"].ReplaceAllLiteralString(tableName, "")
	tableName = core.RegExpCollection["QuotationMarks"].ReplaceAllLiteralString(tableName, "")
	tableName = core.RegExpCollection["SpecQuotationMark"].ReplaceAllLiteralString(tableName, "")
	tableName = strings.TrimSpace(tableName)

	if tableName == "" {
		return `{"state":"error", "result":"invalid command format"}`, errors.New("invalid command format")
	}

	columnsStr := core.RegExpCollection["AlterTableDrop"].ReplaceAllLiteralString(q.Instruction, "")
	columnsStr = core.RegExpCollection["TableParenthesis"].ReplaceAllLiteralString(columnsStr, "")
	columnsIn := core.RegExpCollection["Comma"].Split(columnsStr, -1)

	var columns = []string{}

	for _, column := range columnsIn {
		column = core.RegExpCollection["Spaces"].ReplaceAllLiteralString(column, "")
		column = core.RegExpCollection["QuotationMarks"].ReplaceAllLiteralString(column, "")
		column = core.RegExpCollection["SpecQuotationMark"].ReplaceAllLiteralString(column, "")

		columns = append(columns, column)
	}

	if len(columns) < 1 {
		return `{"state":"error", "result":"invalid command format"}`, errors.New("invalid command format")
	}

	state, ok := core.States[q.Ticket]
	if !ok {
		res.State = "error"
		res.Result = "unknown database"
		return ecowriter.EncodeString(res), errors.New("unknown database")
	}
	db := state.CurrentDB
	if db == "" {
		res.State = "error"
		res.Result = "no database selected"
		return ecowriter.EncodeString(res), errors.New("no database selected")
	}

	_, okDB := core.StorageInfo.DBs[db]
	if okDB {
		dbAccess, okAccess := core.StorageInfo.Access[db]
		if okAccess {
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
						if !(flagsAcs.Alter && flagsAcs.Drop) {
							return `{"state":"error", "result":"not enough rights"}`, errors.New("not enough rights")
						}
					} else {
						return `{"state":"error", "result":"not enough rights"}`, errors.New("not enough rights")
					}
				}
			}
			for _, colName := range columns {
				if !core.RemoveColumn(db, tableName, colName) {
					res.State = "error"
					res.Result = "the column cannot be deleted"
					return ecowriter.EncodeString(res), errors.New("the column cannot be deleted")
				}
			}
		} else {
			res.State = "error"
			res.Result = "internal error"
			return ecowriter.EncodeString(res), errors.New("internal error")
		}
	} else {
		res.State = "error"
		res.Result = "invalid database name"
		return ecowriter.EncodeString(res), errors.New("invalid database name")
	}

	res.State = "ok"
	return ecowriter.EncodeString(res), nil
}

func (q tQuery) DDLAlterTableModify() (result string, err error) {
	// -
	var res gtypes.Response

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

	res.State = "ok"
	return ecowriter.EncodeString(res), nil
}

func (q tQuery) DDLAlterTableRenameTo() (result string, err error) {
	// -
	var res gtypes.Response

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

	res.State = "ok"
	return ecowriter.EncodeString(res), nil
}

func (q tQuery) DDLAlterTable() (result string, err error) {
	// This method is complete

	isAdd := core.RegExpCollection["AlterTableAdd"].MatchString(q.Instruction)
	isDrop := core.RegExpCollection["AlterTableDrop"].MatchString(q.Instruction)
	isModify := core.RegExpCollection["AlterTableModify"].MatchString(q.Instruction)
	isRT := core.RegExpCollection["AlterTableRenameTo"].MatchString(q.Instruction)

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

	isDB := core.RegExpCollection["AlterDatabaseWord"].MatchString(q.Instruction)
	isTable := core.RegExpCollection["AlterTableWord"].MatchString(q.Instruction)

	if isDB {
		return q.DDLAlterDB()
	} else if isTable {
		return q.DDLAlterTable()
	}

	return `{"state":"error", "result":"unknown command"}`, errors.New("unknown command")
}

func (q tQuery) DDLDropDB() (result string, err error) {
	// This method is complete
	var res gtypes.Response

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

	isIE := core.RegExpCollection["IfExistsWord"].MatchString(q.Instruction)

	db := core.RegExpCollection["DropDatabaseWord"].ReplaceAllLiteralString(q.Instruction, "")
	if isIE {
		db = core.RegExpCollection["IfExistsWord"].ReplaceAllLiteralString(db, "")
	}
	db = strings.TrimSpace(db)
	db = core.RegExpCollection["QuotationMarks"].ReplaceAllLiteralString(db, "")
	db = core.RegExpCollection["SpecQuotationMark"].ReplaceAllLiteralString(db, "")

	_, ok := core.StorageInfo.DBs[db]
	if !ok {
		if isIE {
			res.State = "error"
			res.Result = "the database not exists"
			return ecowriter.EncodeString(res), errors.New("the database not exists")
		}

		res.State = "ok"
		return ecowriter.EncodeString(res), nil
	}

	dbAccess, ok := core.StorageInfo.Access[db]
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

	if !core.RemoveDB(db) {
		res.State = "error"
		res.Result = "the database cannot be deleted"
		return ecowriter.EncodeString(res), errors.New("the database cannot be deleted")
	}

	res.State = "ok"
	return ecowriter.EncodeString(res), nil
}

func (q tQuery) DDLDropTable() (result string, err error) {
	// This method is complete
	var res gtypes.Response

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

	isIE := core.RegExpCollection["IfExistsWord"].MatchString(q.Instruction)

	table := core.RegExpCollection["DropTableWord"].ReplaceAllLiteralString(q.Instruction, "")
	if isIE {
		table = core.RegExpCollection["IfExistsWord"].ReplaceAllLiteralString(table, "")
	}
	table = strings.TrimSpace(table)
	table = core.RegExpCollection["QuotationMarks"].ReplaceAllLiteralString(table, "")
	table = core.RegExpCollection["SpecQuotationMark"].ReplaceAllLiteralString(table, "")

	state, ok := core.States[q.Ticket]
	if !ok {
		res.State = "error"
		res.Result = "unknown database"
		return ecowriter.EncodeString(res), errors.New("unknown database")
	}
	db := state.CurrentDB
	if db == "" {
		res.State = "error"
		res.Result = "no database selected"
		return ecowriter.EncodeString(res), errors.New("no database selected")
	}

	dbInfo, okDB := core.StorageInfo.DBs[db]
	if okDB {
		var flagsAcs gtypes.TAccessFlags
		var okFlags bool = false
		var luxUser bool = false

		dbAccess, okAccess := core.StorageInfo.Access[db]
		if okAccess {
			flagsAcs, okFlags = dbAccess.Flags[login]
			if dbAccess.Owner != login {
				for role := range access.Roles {
					if role == int(gauth.ADMIN) || role == int(gauth.ENGINEER) {
						luxUser = true
						break
					}
				}
				if !luxUser {
					if !okFlags {
						return `{"state":"error", "result":"not enough rights"}`, errors.New("not enough rights")
					}
				}
			} else {
				luxUser = true
			}
		} else {
			res.State = "error"
			res.Result = "internal error"
			return ecowriter.EncodeString(res), errors.New("internal error")
		}

		_, okTable := dbInfo.Tables[table]
		if !okTable {
			if isIE {
				res.State = "error"
				res.Result = "the table not exists"
				return ecowriter.EncodeString(res), errors.New("the table not exists")
			}

			res.State = "ok"
			return ecowriter.EncodeString(res), nil
		}

		if !luxUser && !flagsAcs.Drop {
			return `{"state":"error", "result":"not enough rights"}`, errors.New("not enough rights")
		}

		if !core.RemoveTable(db, table) {
			return `{"state":"error", "result":"the table cannot be deleted"}`, errors.New("the table cannot be deleted")
		}
	} else {
		res.State = "error"
		res.Result = "internal error"
		return ecowriter.EncodeString(res), errors.New("internal error")
	}

	res.State = "ok"
	return ecowriter.EncodeString(res), nil
}

func (q tQuery) DDLDrop() (result string, err error) {
	// This method is complete
	op := "internal -> analyzers -> sql -> DDL -> DDLDrop"
	defer func() { e.Wrapper(op, err) }()

	isDB := core.RegExpCollection["DropDatabaseWord"].MatchString(q.Instruction)
	isTable := core.RegExpCollection["DropTableWord"].MatchString(q.Instruction)

	if isDB {
		return q.DDLDropDB()
	} else if isTable {
		return q.DDLDropTable()
	}

	return `{"state":"error", "result":"unknown command"}`, errors.New("unknown command")
}
