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
		return ecowriter.EncodeString(gtypes.Response{
			State:  "error",
			Result: "an empty ticket",
		}), errors.New("an empty ticket")
	}

	login, access, newticket, err := gauth.CheckTicket(q.Ticket)
	if err != nil {
		return ecowriter.EncodeString(gtypes.Response{
			State:  "error",
			Result: err.Error(),
		}), err
	}

	if access.Status.IsBad() {
		return ecowriter.EncodeString(gtypes.Response{
			State:  "error",
			Result: "auth error",
		}), errors.New("auth error")
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
					return ecowriter.EncodeString(gtypes.Response{
						State:  "error",
						Result: "not enough rights",
					}), errors.New("not enough rights")
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
		return ecowriter.EncodeString(gtypes.Response{
			State:  "error",
			Result: "an empty ticket",
		}), errors.New("an empty ticket")
	}

	login, access, newticket, err := gauth.CheckTicket(q.Ticket)
	if err != nil {
		return ecowriter.EncodeString(gtypes.Response{
			State:  "error",
			Result: err.Error(),
		}), err
	}

	if access.Status.IsBad() {
		return ecowriter.EncodeString(gtypes.Response{
			State:  "error",
			Result: "auth error",
		}), errors.New("auth error")
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
						return ecowriter.EncodeString(gtypes.Response{
							State:  "error",
							Result: "not enough rights",
						}), errors.New("not enough rights")
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
				return ecowriter.EncodeString(gtypes.Response{
					State:  "error",
					Result: "not enough rights",
				}), errors.New("not enough rights")
			}

			if !core.RemoveTable(db, table) {
				return ecowriter.EncodeString(gtypes.Response{
					State:  "error",
					Result: "the table cannot be deleted",
				}), errors.New("the table cannot be deleted")
			}
		}

		if !luxUser && !flagsAcs.Create {
			return ecowriter.EncodeString(gtypes.Response{
				State:  "error",
				Result: "not enough rights",
			}), errors.New("not enough rights")
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
	}

	res.State = "ok"
	return ecowriter.EncodeString(res), nil
}

func (q tQuery) DDLCreate() (result string, err error) {
	// This method is complete
	op := "internal -> analyzers -> sql -> DDL -> DDLCreate"
	defer func() { e.Wrapper(op, err) }()

	var res gtypes.Response

	isDB := core.RegExpCollection["CreateDatabaseWord"].MatchString(q.Instruction)
	isTable := core.RegExpCollection["CreateTableWord"].MatchString(q.Instruction)

	if isDB {
		return q.DDLCreateDB()
	} else if isTable {
		return q.DDLCreateTable()
	}

	res.State = "error"
	res.Result = "unknown command"
	return ecowriter.EncodeString(res), errors.New("unknown command")
}

func (q tQuery) DDLAlter() (result string, err error) {
	// -
	op := "internal -> analyzers -> sql -> DDL -> DDLAlter"
	defer func() { e.Wrapper(op, err) }()

	return "DDLAlter", nil
}

func (q tQuery) DDLDrop() (result string, err error) {
	// -
	op := "internal -> analyzers -> sql -> DDL -> DDLDrop"
	defer func() { e.Wrapper(op, err) }()

	return "DDLDrop", nil
}
