package vqlanalyzer

import (
	"errors"
	"strings"

	"github.com/Kwynto/GracefulDB/internal/engine/basicsystem/gauth"
	"github.com/Kwynto/GracefulDB/internal/engine/basicsystem/gtypes"
	"github.com/Kwynto/GracefulDB/internal/engine/core"
	"github.com/Kwynto/GracefulDB/pkg/lib/e"
	"github.com/Kwynto/GracefulDB/pkg/lib/ecowriter"
)

// DML — язык изменения данных (Data Manipulation Language)

func (q tQuery) DMLSelect() (result string, err error) {
	// -
	op := "internal -> analyzers -> sql -> DML -> DMLSelect"
	defer func() { e.Wrapper(op, err) }()

	return "DMLSelect", nil
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

	instruction := core.RegExpCollection["InsertWord"].ReplaceAllLiteralString(q.Instruction, "")
	valuesStr := core.RegExpCollection["InsertValuesToEnd"].FindString(instruction)
	instruction = core.RegExpCollection["InsertValuesToEnd"].ReplaceAllLiteralString(instruction, "")

	columnsStr := core.RegExpCollection["InsertColParenthesis"].FindString(instruction)
	columnsStr = core.RegExpCollection["InsertParenthesis"].ReplaceAllLiteralString(columnsStr, "")
	columnsStr = core.RegExpCollection["Spaces"].ReplaceAllLiteralString(columnsStr, "")
	columnsStr = core.RegExpCollection["QuotationMarks"].ReplaceAllLiteralString(columnsStr, "")
	columnsStr = core.RegExpCollection["SpecQuotationMark"].ReplaceAllLiteralString(columnsStr, "")
	columnsIn = core.RegExpCollection["Comma"].Split(columnsStr, -1)

	table := core.RegExpCollection["InsertColParenthesis"].ReplaceAllLiteralString(instruction, "")
	table = core.RegExpCollection["Spaces"].ReplaceAllLiteralString(table, "")
	table = core.RegExpCollection["QuotationMarks"].ReplaceAllLiteralString(table, "")
	table = core.RegExpCollection["SpecQuotationMark"].ReplaceAllLiteralString(table, "")

	var rowsIn [][]string
	valuesStr = core.RegExpCollection["InsertValuesWord"].ReplaceAllLiteralString(valuesStr, "")
	valuesArr := core.RegExpCollection["InsertSplitParenthesis"].Split(valuesStr, -1)
	for _, value := range valuesArr {
		value = core.RegExpCollection["InsertParenthesis"].ReplaceAllLiteralString(value, "")
		valueIn := core.RegExpCollection["Comma"].Split(value, -1)
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
		expression = make([]gtypes.TConditions, 0, 4)
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

	instruction := core.RegExpCollection["UpdateWord"].ReplaceAllLiteralString(q.Instruction, "")
	whereStr := core.RegExpCollection["WhereToEnd"].FindString(instruction)
	whereStr = core.RegExpCollection["Where"].ReplaceAllLiteralString(whereStr, "")
	// columnsValuesIn.Where = whereStr

	for {
		headCond := core.RegExpCollection["WhereExpression"].ReplaceAllLiteralString(whereStr, "")
		condition := core.RegExpCollection["WhereOperationConditions"].Split(headCond, -1)
		keyIn := condition[0]
		valueIn := condition[1]

		keyIn = core.RegExpCollection["Spaces"].ReplaceAllLiteralString(keyIn, "")
		keyIn = core.RegExpCollection["QuotationMarks"].ReplaceAllLiteralString(keyIn, "")
		keyIn = core.RegExpCollection["SpecQuotationMark"].ReplaceAllLiteralString(keyIn, "")

		valueIn = strings.TrimSpace(valueIn)
		valueIn = core.RegExpCollection["QuotationMarks"].ReplaceAllLiteralString(valueIn, "")
		valueIn = core.RegExpCollection["SpecQuotationMark"].ReplaceAllLiteralString(valueIn, "")

		if keyIn == "" {
			return `{"state":"error", "result":"condition error"}`, errors.New("condition error")
		}
		if valueIn == "" {
			return `{"state":"error", "result":"condition error"}`, errors.New("condition error")
		}

		exp := gtypes.TConditions{
			Type:  "operation",
			Key:   keyIn,
			Value: valueIn,
		}

		if core.RegExpCollection["WhereOperation_<="].MatchString(headCond) {
			exp.Operation = "<="
		} else if core.RegExpCollection["WhereOperation_>="].MatchString(headCond) {
			exp.Operation = ">="
		} else if core.RegExpCollection["WhereOperation_<"].MatchString(headCond) {
			exp.Operation = "<"
		} else if core.RegExpCollection["WhereOperation_>"].MatchString(headCond) {
			exp.Operation = ">"
		} else if core.RegExpCollection["WhereOperation_="].MatchString(headCond) {
			exp.Operation = "="
		} else if core.RegExpCollection["WhereOperation_LIKE"].MatchString(headCond) {
			exp.Operation = "like"
		} else {
			return `{"state":"error", "result":"condition error"}`, errors.New("condition error")
		}
		expression = append(expression, exp)

		whereStr = core.RegExpCollection["WhereExpression"].FindString(whereStr)
		logicOper := core.RegExpCollection["WhereExpression_And_Or_Word"].FindString(whereStr)
		// logicOper = strings.TrimSpace(logicOper)

		if core.RegExpCollection["OR"].MatchString(logicOper) {
			expression = append(expression, gtypes.TConditions{
				Type: "or",
			})
		} else if core.RegExpCollection["AND"].MatchString(logicOper) {
			expression = append(expression, gtypes.TConditions{
				Type: "and",
			})
		} else {
			break
		}

		whereStr = core.RegExpCollection["WhereExpression_And_Or_Word"].ReplaceAllLiteralString(whereStr, "")
	}
	updateIn.Where = append(updateIn.Where, expression...)

	instruction = core.RegExpCollection["WhereToEnd"].ReplaceAllLiteralString(instruction, "")

	columnsValuesStr := core.RegExpCollection["UpdateSetToEnd"].FindString(instruction)
	columnsValuesStr = core.RegExpCollection["UpdateSetWord"].ReplaceAllLiteralString(columnsValuesStr, "")
	columnsValuesArr := core.RegExpCollection["Comma"].Split(columnsValuesStr, -1)

	if len(columnsValuesArr) == 0 || columnsValuesArr[0] == "" {
		return `{"state":"error", "result":"incorrect command syntax"}`, errors.New("incorrect command syntax")
	}

	for _, colVal := range columnsValuesArr {
		colValArr := core.RegExpCollection["SignEqual"].Split(colVal, -1)
		col := colValArr[0]
		val := colValArr[1]

		col = core.RegExpCollection["Spaces"].ReplaceAllLiteralString(col, "")
		col = core.RegExpCollection["QuotationMarks"].ReplaceAllLiteralString(col, "")
		col = core.RegExpCollection["SpecQuotationMark"].ReplaceAllLiteralString(col, "")

		if len(col) == 0 {
			return `{"state":"error", "result":"incorrect syntax"}`, errors.New("incorrect syntax")
		}

		val = strings.TrimSpace(val)
		val = core.RegExpCollection["QuotationMarks"].ReplaceAllLiteralString(val, "")
		val = core.RegExpCollection["SpecQuotationMark"].ReplaceAllLiteralString(val, "")

		updateIn.Couples[col] = val
	}

	table := core.RegExpCollection["UpdateSetToEnd"].ReplaceAllLiteralString(instruction, "")
	table = core.RegExpCollection["Spaces"].ReplaceAllLiteralString(table, "")
	table = core.RegExpCollection["QuotationMarks"].ReplaceAllLiteralString(table, "")
	table = core.RegExpCollection["SpecQuotationMark"].ReplaceAllLiteralString(table, "")

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

	table := core.RegExpCollection["TruncateTableWord"].ReplaceAllLiteralString(q.Instruction, "")
	table = strings.TrimSpace(table)
	table = core.RegExpCollection["QuotationMarks"].ReplaceAllLiteralString(table, "")
	table = core.RegExpCollection["SpecQuotationMark"].ReplaceAllLiteralString(table, "")

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
