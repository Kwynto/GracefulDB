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

// DML — язык изменения данных (Data Manipulation Language)

func (q tQuery) DMLSelect() (result string, err error) {
	// -
	op := "internal -> analyzers -> sql -> DML -> DMLSelect"
	defer func() { e.Wrapper(op, err) }()

	var (
		resultIds []uint64
		okSelect  bool
		res       gtypes.Response
		resArr    gtypes.ResponseUints
		selectIn  = gtypes.TSelectStruct{
			Where:    make([]gtypes.TConditions, 0, 4),
			Columns:  make([]string, 0, 4),
			Distinct: false,
		}
	)

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
		resArr.Ticket = newticket
		res.Ticket = newticket
	}

	state, ok := core.States[q.Ticket]
	if !ok {
		res.State = "error"
		res.Result = "unknown database"
		return ecowriter.EncodeJSON(res), errors.New("unknown database")
	}
	db := state.CurrentDB
	if db == "" {
		res.State = "error"
		res.Result = "no database selected"
		return ecowriter.EncodeJSON(res), errors.New("no database selected")
	}

	// Parsing an expression - Begin

	instruction := vqlexp.RegExpCollection["SelectWord"].ReplaceAllLiteralString(q.Instruction, "")

	orderbyStr := ""
	groupbyStr := ""
	isOrder := false
	isGroup := false

	if vqlexp.RegExpCollection["OrderbyToEnd"].MatchString(instruction) {
		orderbyStr = vqlexp.RegExpCollection["OrderbyToEnd"].FindString(instruction)
		instruction = vqlexp.RegExpCollection["OrderbyToEnd"].ReplaceAllLiteralString(instruction, "")
		orderbyStr = vqlexp.RegExpCollection["Orderby"].ReplaceAllLiteralString(orderbyStr, "")
		isOrder = true
	}

	if vqlexp.RegExpCollection["GroupbyToEnd"].MatchString(instruction) {
		groupbyStr = vqlexp.RegExpCollection["GroupbyToEnd"].FindString(instruction)
		instruction = vqlexp.RegExpCollection["GroupbyToEnd"].ReplaceAllLiteralString(instruction, "")
		groupbyStr = vqlexp.RegExpCollection["Groupby"].ReplaceAllLiteralString(groupbyStr, "")
		isGroup = true
	}

	if vqlexp.RegExpCollection["WhereToEnd"].MatchString(instruction) {
		whereStr := vqlexp.RegExpCollection["WhereToEnd"].FindString(instruction)
		instruction = vqlexp.RegExpCollection["WhereToEnd"].ReplaceAllLiteralString(instruction, "")
		whereStr = vqlexp.RegExpCollection["Where"].ReplaceAllLiteralString(whereStr, "")
		expression, err := parseWhere(whereStr)
		if err != nil {
			return `{"state":"error", "result":"condition error"}`, errors.New("condition error")
		}
		selectIn.Where = append(selectIn.Where, expression...)
		selectIn.IsWhere = true
	}

	table := vqlexp.RegExpCollection["SelectFromToEnd"].FindString(instruction)
	instruction = vqlexp.RegExpCollection["SelectFromToEnd"].ReplaceAllLiteralString(instruction, "")
	table = vqlexp.RegExpCollection["SelectFromWord"].ReplaceAllLiteralString(table, "")
	table = vqlexp.RegExpCollection["Spaces"].ReplaceAllLiteralString(table, "")
	table = vqlexp.RegExpCollection["QuotationMarks"].ReplaceAllLiteralString(table, "")
	table = vqlexp.RegExpCollection["SpecQuotationMark"].ReplaceAllLiteralString(table, "")
	if table == "" {
		return `{"state":"error", "result":"invalid table name"}`, errors.New("invalid table name")
	}

	distinctBool := vqlexp.RegExpCollection["SelectDistinctWord"].MatchString(instruction)
	if distinctBool {
		instruction = vqlexp.RegExpCollection["SelectDistinctWord"].ReplaceAllLiteralString(instruction, "")
	}
	selectIn.Distinct = distinctBool

	columnsStr := strings.TrimSpace(instruction)
	columns := vqlexp.RegExpCollection["Comma"].Split(columnsStr, -1)
	for _, col := range columns {
		col = vqlexp.RegExpCollection["Spaces"].ReplaceAllLiteralString(col, "")
		col = vqlexp.RegExpCollection["QuotationMarks"].ReplaceAllLiteralString(col, "")
		col = vqlexp.RegExpCollection["SpecQuotationMark"].ReplaceAllLiteralString(col, "")
		if col != "" {
			selectIn.Columns = append(selectIn.Columns, col)
		}
	}
	if len(selectIn.Columns) < 1 {
		return `{"state":"error", "result":"no columns"}`, errors.New("no columns")
	}

	if isOrder {
		orderbyExp, err := parseOrderBy(orderbyStr, selectIn.Columns)
		if err != nil {
			return `{"state":"error", "result":"condition error"}`, errors.New("condition error")
		}
		selectIn.Orderby = orderbyExp
		selectIn.IsOrder = isOrder
	}

	if isGroup {
		groupbyCols, err := parseGroupBy(groupbyStr, selectIn.Columns)
		if err != nil {
			return `{"state":"error", "result":"condition error"}`, errors.New("condition error")
		}
		selectIn.Groupby = append(selectIn.Groupby, groupbyCols...)
		selectIn.IsGroup = isGroup
	}

	// Parsing an expression - End

	dbInfo, okDB := core.GetDBInfo(db)
	if okDB {
		var flagsAcs gtypes.TAccessFlags
		var okFlags bool = false
		var luxUser bool = false

		_, okTable := dbInfo.Tables[table]
		if !okTable {
			return `{"state":"error", "result":"invalid table name"}`, errors.New("invalid table name")
		}

		dbAccess, okAccess := core.GetDBAccess(db)
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
			return `{"state":"error", "result":"internal error"}`, errors.New("internal error")
		}

		if !luxUser && !flagsAcs.Select {
			return `{"state":"error", "result":"not enough rights"}`, errors.New("not enough rights")
		}

		// TODO: Make an implementation in the kernel
		resultIds, okSelect = core.SelectRows(db, table, selectIn)
		if !okSelect {
			return `{"state":"error", "result":"the record(s) cannot be updated"}`, errors.New("the record cannot be updated")
		}
	} else {
		return `{"state":"error", "result":"invalid database name"}`, errors.New("invalid database name")
	}

	resArr.State = "ok"
	resArr.Result = resultIds
	return ecowriter.EncodeJSON(resArr), nil
}

func (q tQuery) DMLInsert() (result string, err error) {
	// This method is complete
	op := "internal -> analyzers -> sql -> DML -> DMLInsert"
	defer func() { e.Wrapper(op, err) }()

	var (
		resultIds []uint64
		okInsert  bool
		res       gtypes.Response
		resArr    gtypes.ResponseUints
		columnsIn = make([]string, 0)
	)

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
		resArr.Ticket = newticket
		res.Ticket = newticket
	}

	state, ok := core.States[q.Ticket]
	if !ok {
		res.State = "error"
		res.Result = "unknown database"
		return ecowriter.EncodeJSON(res), errors.New("unknown database")
	}
	db := state.CurrentDB
	if db == "" {
		res.State = "error"
		res.Result = "no database selected"
		return ecowriter.EncodeJSON(res), errors.New("no database selected")
	}

	instruction := vqlexp.RegExpCollection["InsertWord"].ReplaceAllLiteralString(q.Instruction, "")
	valuesStr := vqlexp.RegExpCollection["InsertValuesToEnd"].FindString(instruction)
	instruction = vqlexp.RegExpCollection["InsertValuesToEnd"].ReplaceAllLiteralString(instruction, "")

	columnsStr := vqlexp.RegExpCollection["InsertColParenthesis"].FindString(instruction)
	columnsStr = vqlexp.RegExpCollection["InsertParenthesis"].ReplaceAllLiteralString(columnsStr, "")
	columnsStr = vqlexp.RegExpCollection["Spaces"].ReplaceAllLiteralString(columnsStr, "")
	columnsStr = vqlexp.RegExpCollection["QuotationMarks"].ReplaceAllLiteralString(columnsStr, "")
	columnsStr = vqlexp.RegExpCollection["SpecQuotationMark"].ReplaceAllLiteralString(columnsStr, "")
	columnsIn = vqlexp.RegExpCollection["Comma"].Split(columnsStr, -1)

	table := vqlexp.RegExpCollection["InsertColParenthesis"].ReplaceAllLiteralString(instruction, "")
	table = vqlexp.RegExpCollection["Spaces"].ReplaceAllLiteralString(table, "")
	table = vqlexp.RegExpCollection["QuotationMarks"].ReplaceAllLiteralString(table, "")
	table = vqlexp.RegExpCollection["SpecQuotationMark"].ReplaceAllLiteralString(table, "")

	var rowsIn [][]string
	valuesStr = vqlexp.RegExpCollection["InsertValuesWord"].ReplaceAllLiteralString(valuesStr, "")
	valuesArr := vqlexp.RegExpCollection["InsertSplitParenthesis"].Split(valuesStr, -1)
	for _, value := range valuesArr {
		value = vqlexp.RegExpCollection["InsertParenthesis"].ReplaceAllLiteralString(value, "")
		valueIn := vqlexp.RegExpCollection["Comma"].Split(value, -1)
		var rowIn []string
		for _, val := range valueIn {
			val = strings.TrimSpace(val)
			val = strings.TrimRight(val, `"'`)
			val = strings.TrimRight(val, "`")
			val = strings.TrimLeft(val, `"'`)
			val = strings.TrimLeft(val, "`")
			rowIn = append(rowIn, val)
		}
		rowsIn = append(rowsIn, rowIn)
	}

LabelCheck:
	dbInfo, okDB := core.GetDBInfo(db)
	if okDB {
		var flagsAcs gtypes.TAccessFlags
		var okFlags bool = false
		var luxUser bool = false

		_, okTable := dbInfo.Tables[table]
		if !okTable {
			if core.LocalCoreSettings.FriendlyMode {
				if !core.CreateTable(db, table, true) {
					return `{"state":"error", "result":"invalid table name"}`, errors.New("invalid table name")
				}
				goto LabelCheck
			}
			return `{"state":"error", "result":"invalid table name"}`, errors.New("invalid table name")
		}

		if len(columnsIn) == 0 || columnsIn[0] == "" {
			// clear(columnsIn)
			columnsIn = dbInfo.Tables[table].Order
		}

		dbAccess, okAccess := core.GetDBAccess(db)
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
			return `{"state":"error", "result":"internal error"}`, errors.New("internal error")
		}

		if !luxUser && !flagsAcs.Insert {
			return `{"state":"error", "result":"not enough rights"}`, errors.New("not enough rights")
		}

		resultIds, okInsert = core.InsertRows(db, table, columnsIn, rowsIn)
		if !okInsert {
			return `{"state":"error", "result":"the record(s) cannot be inserted"}`, errors.New("the record cannot be inserted")
		}
	} else {
		if core.LocalCoreSettings.FriendlyMode {
			if !core.CreateDB(db, login, true) {
				return `{"state":"error", "result":"invalid database name"}`, errors.New("invalid database name")
			}
			goto LabelCheck
		}
		return `{"state":"error", "result":"invalid database name"}`, errors.New("invalid database name")
	}

	resArr.State = "ok"
	resArr.Result = resultIds
	return ecowriter.EncodeJSON(resArr), nil
}

func (q tQuery) DMLUpdate() (result string, err error) {
	// -
	op := "internal -> analyzers -> sql -> DML -> DMLUpdate"
	defer func() { e.Wrapper(op, err) }()

	var (
		resultIds []uint64
		okUpdate  bool
		res       gtypes.Response
		resArr    gtypes.ResponseUints
		updateIn  = gtypes.TUpdaateStruct{
			Where:   make([]gtypes.TConditions, 0, 4),
			Couples: make(map[string]string),
		}
	)

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
		resArr.Ticket = newticket
		res.Ticket = newticket
	}

	state, ok := core.States[q.Ticket]
	if !ok {
		res.State = "error"
		res.Result = "unknown database"
		return ecowriter.EncodeJSON(res), errors.New("unknown database")
	}
	db := state.CurrentDB
	if db == "" {
		res.State = "error"
		res.Result = "no database selected"
		return ecowriter.EncodeJSON(res), errors.New("no database selected")
	}

	instruction := vqlexp.RegExpCollection["UpdateWord"].ReplaceAllLiteralString(q.Instruction, "")
	whereStr := vqlexp.RegExpCollection["WhereToEnd"].FindString(instruction)
	whereStr = vqlexp.RegExpCollection["Where"].ReplaceAllLiteralString(whereStr, "")
	// columnsValuesIn.Where = whereStr

	expression, err := parseWhere(whereStr)
	if err != nil {
		return `{"state":"error", "result":"condition error"}`, errors.New("condition error")
	}
	updateIn.Where = append(updateIn.Where, expression...)

	instruction = vqlexp.RegExpCollection["WhereToEnd"].ReplaceAllLiteralString(instruction, "")

	columnsValuesStr := vqlexp.RegExpCollection["UpdateSetToEnd"].FindString(instruction)
	columnsValuesStr = vqlexp.RegExpCollection["UpdateSetWord"].ReplaceAllLiteralString(columnsValuesStr, "")
	columnsValuesArr := vqlexp.RegExpCollection["Comma"].Split(columnsValuesStr, -1)

	if len(columnsValuesArr) == 0 || columnsValuesArr[0] == "" {
		return `{"state":"error", "result":"incorrect command syntax"}`, errors.New("incorrect command syntax")
	}

	for _, colVal := range columnsValuesArr {
		colValArr := vqlexp.RegExpCollection["SignEqual"].Split(colVal, -1)
		col := colValArr[0]
		val := colValArr[1]

		col = vqlexp.RegExpCollection["Spaces"].ReplaceAllLiteralString(col, "")
		col = vqlexp.RegExpCollection["QuotationMarks"].ReplaceAllLiteralString(col, "")
		col = vqlexp.RegExpCollection["SpecQuotationMark"].ReplaceAllLiteralString(col, "")

		if len(col) == 0 {
			return `{"state":"error", "result":"incorrect syntax"}`, errors.New("incorrect syntax")
		}

		val = strings.TrimSpace(val)
		val = vqlexp.RegExpCollection["QuotationMarks"].ReplaceAllLiteralString(val, "")
		val = vqlexp.RegExpCollection["SpecQuotationMark"].ReplaceAllLiteralString(val, "")

		updateIn.Couples[col] = val
	}

	table := vqlexp.RegExpCollection["UpdateSetToEnd"].ReplaceAllLiteralString(instruction, "")
	table = vqlexp.RegExpCollection["Spaces"].ReplaceAllLiteralString(table, "")
	table = vqlexp.RegExpCollection["QuotationMarks"].ReplaceAllLiteralString(table, "")
	table = vqlexp.RegExpCollection["SpecQuotationMark"].ReplaceAllLiteralString(table, "")

LabelCheck:
	dbInfo, okDB := core.GetDBInfo(db)
	if okDB {
		var flagsAcs gtypes.TAccessFlags
		var okFlags bool = false
		var luxUser bool = false

		_, okTable := dbInfo.Tables[table]
		if !okTable {
			if core.LocalCoreSettings.FriendlyMode {
				if !core.CreateTable(db, table, true) {
					return `{"state":"error", "result":"invalid table name"}`, errors.New("invalid table name")
				}
				goto LabelCheck
			}
			return `{"state":"error", "result":"invalid table name"}`, errors.New("invalid table name")
		}

		dbAccess, okAccess := core.GetDBAccess(db)
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
			return `{"state":"error", "result":"internal error"}`, errors.New("internal error")
		}

		if !luxUser && !flagsAcs.Update {
			return `{"state":"error", "result":"not enough rights"}`, errors.New("not enough rights")
		}

		// TODO: Make an implementation in the kernel
		resultIds, okUpdate = core.UpdateRows(db, table, updateIn)
		if !okUpdate {
			return `{"state":"error", "result":"the record(s) cannot be updated"}`, errors.New("the record cannot be updated")
		}
	} else {
		if core.LocalCoreSettings.FriendlyMode {
			if !core.CreateDB(db, login, true) {
				return `{"state":"error", "result":"invalid database name"}`, errors.New("invalid database name")
			}
			goto LabelCheck
		}
		return `{"state":"error", "result":"invalid database name"}`, errors.New("invalid database name")
	}

	resArr.State = "ok"
	resArr.Result = resultIds
	return ecowriter.EncodeJSON(resArr), nil
}

func (q tQuery) DMLDelete() (result string, err error) {
	// -
	op := "internal -> analyzers -> sql -> DML -> DMLDelete"
	defer func() { e.Wrapper(op, err) }()

	return "DMLDelete", nil
}

func (q tQuery) DMLCommit() (result string, err error) {
	// -
	op := "internal -> analyzers -> sql -> DML -> DMLCommit"
	defer func() { e.Wrapper(op, err) }()

	return "DMLCommit", nil
}

func (q tQuery) DMLRollback() (result string, err error) {
	// -
	op := "internal -> analyzers -> sql -> DML -> DMLRollback"
	defer func() { e.Wrapper(op, err) }()

	return "DMLRollback", nil
}

func (q tQuery) DMLTruncateTable() (result string, err error) {
	// This method is complete
	op := "internal -> analyzers -> sql -> DML -> DMLTruncate"
	defer func() { e.Wrapper(op, err) }()

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

	state, ok := core.States[q.Ticket]
	if !ok {
		res.State = "error"
		res.Result = "unknown database"
		return ecowriter.EncodeJSON(res), errors.New("unknown database")
	}
	db := state.CurrentDB
	if db == "" {
		res.State = "error"
		res.Result = "no database selected"
		return ecowriter.EncodeJSON(res), errors.New("no database selected")
	}

	table := vqlexp.RegExpCollection["TruncateTableWord"].ReplaceAllLiteralString(q.Instruction, "")
	table = strings.TrimSpace(table)
	table = vqlexp.RegExpCollection["QuotationMarks"].ReplaceAllLiteralString(table, "")
	table = vqlexp.RegExpCollection["SpecQuotationMark"].ReplaceAllLiteralString(table, "")

LabelCheck:
	dbInfo, okDB := core.GetDBInfo(db)
	if okDB {
		var flagsAcs gtypes.TAccessFlags
		var okFlags bool = false
		var luxUser bool = false

		_, okTable := dbInfo.Tables[table]
		if !okTable {
			if core.LocalCoreSettings.FriendlyMode {
				if !core.CreateTable(db, table, true) {
					return `{"state":"error", "result":"invalid table name"}`, errors.New("invalid table name")
				}
				goto LabelCheck
			}
			return `{"state":"error", "result":"invalid table name"}`, errors.New("invalid table name")
		}

		dbAccess, okAccess := core.GetDBAccess(db)
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
			return `{"state":"error", "result":"internal error"}`, errors.New("internal error")
		}

		if !luxUser && !flagsAcs.Delete {
			return `{"state":"error", "result":"not enough rights"}`, errors.New("not enough rights")
		}

		// TODO: Make an implementation in the kernel
		if !core.TruncateTable(db, table) {
			return `{"state":"error", "result":"the table cannot be truncated"}`, errors.New("the table cannot be truncated")
		}
	} else {
		if core.LocalCoreSettings.FriendlyMode {
			if !core.CreateDB(db, login, true) {
				return `{"state":"error", "result":"invalid database name"}`, errors.New("invalid database name")
			}
			goto LabelCheck
		}
		return `{"state":"error", "result":"internal error"}`, errors.New("internal error")
	}

	res.State = "ok"
	return ecowriter.EncodeJSON(res), nil
}
