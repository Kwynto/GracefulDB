package vqlanalyzer

import (
	"errors"
	"strings"

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
			Orderby: gtypes.TOrderBy{
				Cols: make([]string, 0, 4),
				Sort: make([]uint8, 0, 4),
			},
			Groupby:  make([]string, 0, 4),
			Where:    make([]gtypes.TConditions, 0, 4),
			Columns:  make([]string, 0, 4),
			IsOrder:  false,
			IsGroup:  false,
			IsWhere:  false,
			Distinct: false,
		}
	)

	// Pre checking

	login, db, access, newticket, err := preChecker(q.Ticket)
	if err != nil {
		res.State = "error"
		res.Result = err.Error()
		return ecowriter.EncodeJSON(res), err
	}

	if newticket != "" {
		resArr.Ticket = newticket
		res.Ticket = newticket
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
	table = trimQuotationMarks(table)
	if table == "" {
		return `{"state":"error", "result":"invalid table name"}`, errors.New("invalid table name")
	}

	if !core.IfExistTable(db, table) {
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
		col = trimQuotationMarks(col)
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

	// Post checking

	luxUser, flagsAcs, err := dourPostChecker(db, table, login, access)
	if err != nil {
		res.State = "error"
		res.Result = err.Error()
		return ecowriter.EncodeJSON(res), err
	}

	if !luxUser && !flagsAcs.Select {
		return `{"state":"error", "result":"not enough rights"}`, errors.New("not enough rights")
	}

	// Request execution

	// TODO: Make an implementation in the kernel
	resultIds, okSelect = core.SelectRows(db, table, selectIn)
	if !okSelect {
		return `{"state":"error", "result":"the record(s) cannot be updated"}`, errors.New("the record cannot be updated")
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

	// Pre checking

	login, db, access, newticket, err := preChecker(q.Ticket)
	if err != nil {
		res.State = "error"
		res.Result = err.Error()
		return ecowriter.EncodeJSON(res), err
	}

	if newticket != "" {
		resArr.Ticket = newticket
		res.Ticket = newticket
	}

	// Parsing an expression - Begin

	instruction := vqlexp.RegExpCollection["InsertWord"].ReplaceAllLiteralString(q.Instruction, "")
	valuesStr := vqlexp.RegExpCollection["InsertValuesToEnd"].FindString(instruction)
	instruction = vqlexp.RegExpCollection["InsertValuesToEnd"].ReplaceAllLiteralString(instruction, "")

	columnsStr := vqlexp.RegExpCollection["InsertColParenthesis"].FindString(instruction)
	columnsStr = vqlexp.RegExpCollection["InsertParenthesis"].ReplaceAllLiteralString(columnsStr, "")
	columnsStr = vqlexp.RegExpCollection["Spaces"].ReplaceAllLiteralString(columnsStr, "")
	columnsStr = trimQuotationMarks(columnsStr)
	columnsIn = vqlexp.RegExpCollection["Comma"].Split(columnsStr, -1)

	table := vqlexp.RegExpCollection["InsertColParenthesis"].ReplaceAllLiteralString(instruction, "")
	table = vqlexp.RegExpCollection["Spaces"].ReplaceAllLiteralString(table, "")
	table = trimQuotationMarks(table)

	if !core.IfExistTable(db, table) {
		return `{"state":"error", "result":"invalid table name"}`, errors.New("invalid table name")
	}

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

	if len(columnsIn) == 0 || columnsIn[0] == "" {
		dbInfo, okDB := core.GetDBInfo(db)
		if okDB {
			columnsIn = dbInfo.Tables[table].Order
		} else {
			return `{"state":"error", "result":"invalid database name"}`, errors.New("invalid database name")
		}
	}

	// Parsing an expression - End

	// Post checking

	luxUser, flagsAcs, err := friendlyPostChecker(db, table, login, access)
	if err != nil {
		res.State = "error"
		res.Result = err.Error()
		return ecowriter.EncodeJSON(res), err
	}

	if !luxUser && !flagsAcs.Insert {
		return `{"state":"error", "result":"not enough rights"}`, errors.New("not enough rights")
	}

	// Request execution

	resultIds, okInsert = core.InsertRows(db, table, columnsIn, rowsIn)
	if !okInsert {
		return `{"state":"error", "result":"the record(s) cannot be inserted"}`, errors.New("the record cannot be inserted")
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

	// Pre checking

	login, db, access, newticket, err := preChecker(q.Ticket)
	if err != nil {
		res.State = "error"
		res.Result = err.Error()
		return ecowriter.EncodeJSON(res), err
	}

	if newticket != "" {
		resArr.Ticket = newticket
		res.Ticket = newticket
	}

	// Parsing an expression - Begin

	instruction := vqlexp.RegExpCollection["UpdateWord"].ReplaceAllLiteralString(q.Instruction, "")
	whereStr := vqlexp.RegExpCollection["WhereToEnd"].FindString(instruction)
	whereStr = vqlexp.RegExpCollection["Where"].ReplaceAllLiteralString(whereStr, "")
	// columnsValuesIn.Where = whereStr

	expression, err := parseWhere(whereStr)
	if err != nil {
		return `{"state":"error", "result":"condition error"}`, err
	}
	if len(expression) == 0 {
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
		col = trimQuotationMarks(col)

		if len(col) == 0 {
			return `{"state":"error", "result":"incorrect syntax"}`, errors.New("incorrect syntax")
		}

		val = strings.TrimSpace(val)
		val = trimQuotationMarks(val)

		updateIn.Couples[col] = val
	}

	table := vqlexp.RegExpCollection["UpdateSetToEnd"].ReplaceAllLiteralString(instruction, "")
	table = vqlexp.RegExpCollection["Spaces"].ReplaceAllLiteralString(table, "")
	table = trimQuotationMarks(table)

	if !core.IfExistTable(db, table) {
		return `{"state":"error", "result":"invalid table name"}`, errors.New("invalid table name")
	}

	// Parsing an expression - End

	// Post checking

	luxUser, flagsAcs, err := friendlyPostChecker(db, table, login, access)
	if err != nil {
		res.State = "error"
		res.Result = err.Error()
		return ecowriter.EncodeJSON(res), err
	}

	if !luxUser && !flagsAcs.Update {
		return `{"state":"error", "result":"not enough rights"}`, errors.New("not enough rights")
	}

	// Request execution

	// TODO: Make an implementation in the kernel
	resultIds, okUpdate = core.UpdateRows(db, table, updateIn)
	if !okUpdate {
		return `{"state":"error", "result":"the record(s) cannot be updated"}`, errors.New("the record cannot be updated")
	}

	resArr.State = "ok"
	resArr.Result = resultIds
	return ecowriter.EncodeJSON(resArr), nil
}

func (q tQuery) DMLDelete() (result string, err error) {
	// This function is complete
	op := "internal -> analyzers -> sql -> DML -> DMLDelete"
	defer func() { e.Wrapper(op, err) }()

	var (
		resultIds []uint64
		okDel     bool
		res       gtypes.Response
		resArr    gtypes.ResponseUints
		deleteIn  = gtypes.TDeleteStruct{
			Where:   make([]gtypes.TConditions, 0, 4),
			IsWhere: false,
		}
	)

	// Pre checking

	login, db, access, newticket, err := preChecker(q.Ticket)
	if err != nil {
		res.State = "error"
		res.Result = err.Error()
		return ecowriter.EncodeJSON(res), err
	}

	if newticket != "" {
		resArr.Ticket = newticket
		res.Ticket = newticket
	}

	// Parsing an expression - Begin

	instruction := vqlexp.RegExpCollection["DeleteWord"].ReplaceAllLiteralString(q.Instruction, "")

	if vqlexp.RegExpCollection["WhereToEnd"].MatchString(instruction) {
		whereStr := vqlexp.RegExpCollection["WhereToEnd"].FindString(instruction)
		instruction = vqlexp.RegExpCollection["WhereToEnd"].ReplaceAllLiteralString(instruction, "")
		whereStr = vqlexp.RegExpCollection["Where"].ReplaceAllLiteralString(whereStr, "")
		expression, err := parseWhere(whereStr)
		if err != nil {
			return `{"state":"error", "result":"condition error"}`, errors.New("condition error")
		}
		deleteIn.Where = append(deleteIn.Where, expression...)
		deleteIn.IsWhere = true
	}

	table := vqlexp.RegExpCollection["Spaces"].ReplaceAllLiteralString(instruction, "")
	table = trimQuotationMarks(table)
	if table == "" {
		return `{"state":"error", "result":"invalid table name"}`, errors.New("invalid table name")
	}

	if !core.IfExistTable(db, table) {
		return `{"state":"error", "result":"invalid table name"}`, errors.New("invalid table name")
	}

	// Parsing an expression - End

	// Post checking

	luxUser, flagsAcs, err := dourPostChecker(db, table, login, access)
	if err != nil {
		res.State = "error"
		res.Result = err.Error()
		return ecowriter.EncodeJSON(res), err
	}

	if !luxUser && !flagsAcs.Delete {
		return `{"state":"error", "result":"not enough rights"}`, errors.New("not enough rights")
	}

	// Request execution

	// TODO: Make an implementation in the kernel
	resultIds, okDel = core.DeleteRows(db, table, deleteIn)
	if !okDel {
		return `{"state":"error", "result":"the record(s) cannot be updated"}`, errors.New("the record cannot be updated")
	}

	resArr.State = "ok"
	resArr.Result = resultIds
	return ecowriter.EncodeJSON(resArr), nil
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

	table := vqlexp.RegExpCollection["TruncateTableWord"].ReplaceAllLiteralString(q.Instruction, "")
	table = strings.TrimSpace(table)
	table = trimQuotationMarks(table)

	if !core.IfExistTable(db, table) {
		return `{"state":"error", "result":"invalid table name"}`, errors.New("invalid table name")
	}

	// Parsing an expression - End

	// Post checking

	luxUser, flagsAcs, err := friendlyPostChecker(db, table, login, access)
	if err != nil {
		res.State = "error"
		res.Result = err.Error()
		return ecowriter.EncodeJSON(res), err
	}

	if !luxUser && !flagsAcs.Delete {
		return `{"state":"error", "result":"not enough rights"}`, errors.New("not enough rights")
	}

	// Request execution

	// TODO: Make an implementation in the kernel
	if !core.TruncateTable(db, table) {
		return `{"state":"error", "result":"the table cannot be truncated"}`, errors.New("the table cannot be truncated")
	}

	res.State = "ok"
	return ecowriter.EncodeJSON(res), nil
}
