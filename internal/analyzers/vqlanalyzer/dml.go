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

// DML — Data Manipulation Language (язык изменения данных)

func (q tQuery) DMLSelect() (result string, err error) {
	// - It's almost done
	sOperation := "internal -> analyzers -> sql -> DML -> DMLSelect"
	defer func() { e.Wrapper(sOperation, err) }()

	var (
		stResultRow []gtypes.TResponseRow
		isOkSelect  bool
		stRes       gtypes.TResponse
		stResSelect gtypes.TResponseSelect
		stSelectIn  = gtypes.TSelectStruct{
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

	sLogin, sDB, stAccess, sNewTicket, err := preChecker(q.Ticket)
	if err != nil {
		stRes.State = "error"
		stRes.Result = err.Error()
		return ecowriter.EncodeJSON(stRes), err
	}

	if sNewTicket != "" {
		stResSelect.Ticket = sNewTicket
		stRes.Ticket = sNewTicket
	}

	// Parsing an expression - Begin

	sInstruction := vqlexp.RegExpCollection["SelectWord"].ReplaceAllLiteralString(q.Instruction, "")

	sOrderBy := ""
	sGroupBy := ""
	isOrder := false
	isGroup := false

	if vqlexp.RegExpCollection["OrderbyToEnd"].MatchString(sInstruction) {
		sOrderBy = vqlexp.RegExpCollection["OrderbyToEnd"].FindString(sInstruction)
		sInstruction = vqlexp.RegExpCollection["OrderbyToEnd"].ReplaceAllLiteralString(sInstruction, "")
		sOrderBy = vqlexp.RegExpCollection["Orderby"].ReplaceAllLiteralString(sOrderBy, "")
		isOrder = true
	}

	if vqlexp.RegExpCollection["GroupbyToEnd"].MatchString(sInstruction) {
		sGroupBy = vqlexp.RegExpCollection["GroupbyToEnd"].FindString(sInstruction)
		sInstruction = vqlexp.RegExpCollection["GroupbyToEnd"].ReplaceAllLiteralString(sInstruction, "")
		sGroupBy = vqlexp.RegExpCollection["Groupby"].ReplaceAllLiteralString(sGroupBy, "")
		isGroup = true
	}

	if vqlexp.RegExpCollection["WhereToEnd"].MatchString(sInstruction) {
		sWhere := vqlexp.RegExpCollection["WhereToEnd"].FindString(sInstruction)
		sInstruction = vqlexp.RegExpCollection["WhereToEnd"].ReplaceAllLiteralString(sInstruction, "")
		sWhere = vqlexp.RegExpCollection["Where"].ReplaceAllLiteralString(sWhere, "")
		expression, err := parseWhere(sWhere)
		if err != nil {
			return `{"state":"error", "result":"condition error"}`, errors.New("condition error")
		}
		stSelectIn.Where = append(stSelectIn.Where, expression...)
		stSelectIn.IsWhere = true
	}

	sTable := vqlexp.RegExpCollection["SelectFromToEnd"].FindString(sInstruction)
	sInstruction = vqlexp.RegExpCollection["SelectFromToEnd"].ReplaceAllLiteralString(sInstruction, "")
	sTable = vqlexp.RegExpCollection["SelectFromWord"].ReplaceAllLiteralString(sTable, "")
	sTable = vqlexp.RegExpCollection["Spaces"].ReplaceAllLiteralString(sTable, "")
	sTable = trimQuotationMarks(sTable)
	if sTable == "" {
		return `{"state":"error", "result":"invalid table name"}`, errors.New("invalid table name")
	}

	if !core.IfExistTable(sDB, sTable) {
		return `{"state":"error", "result":"invalid table name"}`, errors.New("invalid table name")
	}

	isDistinct := vqlexp.RegExpCollection["SelectDistinctWord"].MatchString(sInstruction)
	if isDistinct {
		sInstruction = vqlexp.RegExpCollection["SelectDistinctWord"].ReplaceAllLiteralString(sInstruction, "")
	}
	stSelectIn.Distinct = isDistinct

	sColumns := strings.TrimSpace(sInstruction)
	slColumns := vqlexp.RegExpCollection["Comma"].Split(sColumns, -1)
	for _, sCol := range slColumns {
		sCol = vqlexp.RegExpCollection["Spaces"].ReplaceAllLiteralString(sCol, "")
		sCol = trimQuotationMarks(sCol)
		if sCol != "" {
			stSelectIn.Columns = append(stSelectIn.Columns, sCol)
		}
	}
	if len(stSelectIn.Columns) < 1 {
		return `{"state":"error", "result":"no columns"}`, errors.New("no columns")
	}

	if isOrder {
		stOrderByExp, err := parseOrderBy(sOrderBy, stSelectIn.Columns)
		if err != nil {
			return `{"state":"error", "result":"condition error"}`, errors.New("condition error")
		}
		stSelectIn.Orderby = stOrderByExp
		stSelectIn.IsOrder = isOrder
	}

	if isGroup {
		slSGroupByCols, err := parseGroupBy(sGroupBy, stSelectIn.Columns)
		if err != nil {
			return `{"state":"error", "result":"condition error"}`, errors.New("condition error")
		}
		stSelectIn.Groupby = append(stSelectIn.Groupby, slSGroupByCols...)
		stSelectIn.IsGroup = isGroup
	}

	// Parsing an expression - End

	// Post checking

	isLuxUser, stFlagsAcs, err := dourPostChecker(sDB, sTable, sLogin, stAccess)
	if err != nil {
		stRes.State = "error"
		stRes.Result = err.Error()
		return ecowriter.EncodeJSON(stRes), err
	}

	if !isLuxUser && !stFlagsAcs.Select {
		return `{"state":"error", "result":"not enough rights"}`, errors.New("not enough rights")
	}

	// Request execution

	// TODO: Make an implementation in the kernel
	stResultRow, isOkSelect = core.SelectRows(sDB, sTable, stSelectIn)
	if !isOkSelect {
		return `{"state":"error", "result":"the record(s) cannot be updated"}`, errors.New("the record cannot be updated")
	}

	stResSelect.State = "ok"
	stResSelect.Result = stResultRow
	return ecowriter.EncodeJSON(stResSelect), nil
}

func (q tQuery) DMLInsert() (result string, err error) {
	// This method is complete
	sOperation := "internal -> analyzers -> sql -> DML -> DMLInsert"
	defer func() { e.Wrapper(sOperation, err) }()

	var (
		slResultIDs = make([]uint64, 0)
		isOkInsert  bool
		stRes       gtypes.TResponse
		stResArr    gtypes.TResponseUints
		slColumnsIn = make([]string, 0)
	)

	// Pre checking

	sLogin, sDB, stAccess, sNewTicket, err := preChecker(q.Ticket)
	if err != nil {
		stRes.State = "error"
		stRes.Result = err.Error()
		return ecowriter.EncodeJSON(stRes), err
	}

	if sNewTicket != "" {
		stResArr.Ticket = sNewTicket
		stRes.Ticket = sNewTicket
	}

	// Parsing an expression - Begin

	sInstruction := vqlexp.RegExpCollection["InsertWord"].ReplaceAllLiteralString(q.Instruction, "")
	sValues := vqlexp.RegExpCollection["InsertValuesToEnd"].FindString(sInstruction)
	sInstruction = vqlexp.RegExpCollection["InsertValuesToEnd"].ReplaceAllLiteralString(sInstruction, "")

	sColumns := vqlexp.RegExpCollection["InsertColParenthesis"].FindString(sInstruction)
	sColumns = vqlexp.RegExpCollection["InsertParenthesis"].ReplaceAllLiteralString(sColumns, "")
	sColumns = vqlexp.RegExpCollection["Spaces"].ReplaceAllLiteralString(sColumns, "")
	sColumns = trimQuotationMarks(sColumns)
	slColumnsIn = vqlexp.RegExpCollection["Comma"].Split(sColumns, -1)

	sTable := vqlexp.RegExpCollection["InsertColParenthesis"].ReplaceAllLiteralString(sInstruction, "")
	sTable = vqlexp.RegExpCollection["Spaces"].ReplaceAllLiteralString(sTable, "")
	sTable = trimQuotationMarks(sTable)

	if !core.IfExistTable(sDB, sTable) {
		return `{"state":"error", "result":"invalid table name"}`, errors.New("invalid table name")
	}

	var slRowsIn [][]string
	sValues = vqlexp.RegExpCollection["InsertValuesWord"].ReplaceAllLiteralString(sValues, "")
	slValues := vqlexp.RegExpCollection["InsertSplitParenthesis"].Split(sValues, -1)
	for _, sValue := range slValues {
		sValue = vqlexp.RegExpCollection["InsertParenthesis"].ReplaceAllLiteralString(sValue, "")
		slValueIn := vqlexp.RegExpCollection["Comma"].Split(sValue, -1)
		var slRowIn []string
		for _, sVal := range slValueIn {
			sVal = strings.TrimSpace(sVal)
			sVal = strings.TrimRight(sVal, `"'`)
			sVal = strings.TrimRight(sVal, "`")
			sVal = strings.TrimLeft(sVal, `"'`)
			sVal = strings.TrimLeft(sVal, "`")
			slRowIn = append(slRowIn, sVal)
		}
		slRowsIn = append(slRowsIn, slRowIn)
	}

	if len(slColumnsIn) == 0 || slColumnsIn[0] == "" {
		stDBInfo, isOkDB := core.GetDBInfo(sDB)
		if isOkDB {
			slColumnsIn = stDBInfo.Tables[sTable].Order
		} else {
			return `{"state":"error", "result":"invalid database name"}`, errors.New("invalid database name")
		}
	}

	// Parsing an expression - End

	// Post checking

	isLuxUser, stFlagsAcs, err := friendlyPostChecker(sDB, sTable, sLogin, stAccess)
	if err != nil {
		stRes.State = "error"
		stRes.Result = err.Error()
		return ecowriter.EncodeJSON(stRes), err
	}

	if !isLuxUser && !stFlagsAcs.Insert {
		return `{"state":"error", "result":"not enough rights"}`, errors.New("not enough rights")
	}

	// Request execution

	slResultIDs, isOkInsert = core.InsertRows(sDB, sTable, slColumnsIn, slRowsIn)
	if !isOkInsert {
		return `{"state":"error", "result":"the record(s) cannot be inserted"}`, errors.New("the record cannot be inserted")
	}

	stResArr.State = "ok"
	stResArr.Result = slResultIDs
	return ecowriter.EncodeJSON(stResArr), nil
}

func (q tQuery) DMLUpdate() (result string, err error) {
	// This function is complete
	sOperation := "internal -> analyzers -> sql -> DML -> DMLUpdate"
	defer func() { e.Wrapper(sOperation, err) }()

	var (
		slResultIDs = make([]uint64, 0)
		isOkUpdate  bool
		stRes       gtypes.TResponse
		stResArr    gtypes.TResponseUints
		stUpdateIn  = gtypes.TUpdaateStruct{
			Where:   make([]gtypes.TConditions, 0, 4),
			Couples: make(map[string]string),
		}
	)

	// Pre checking

	sLogin, sDB, stAccess, sNewTicket, err := preChecker(q.Ticket)
	if err != nil {
		stRes.State = "error"
		stRes.Result = err.Error()
		return ecowriter.EncodeJSON(stRes), err
	}

	if sNewTicket != "" {
		stResArr.Ticket = sNewTicket
		stRes.Ticket = sNewTicket
	}

	// Parsing an expression - Begin

	sInstruction := vqlexp.RegExpCollection["UpdateWord"].ReplaceAllLiteralString(q.Instruction, "")
	sWhere := vqlexp.RegExpCollection["WhereToEnd"].FindString(sInstruction)
	sWhere = vqlexp.RegExpCollection["Where"].ReplaceAllLiteralString(sWhere, "")
	// columnsValuesIn.Where = sWhere

	slExpression, err := parseWhere(sWhere)
	if err != nil {
		return `{"state":"error", "result":"condition error"}`, err
	}
	if len(slExpression) == 0 {
		return `{"state":"error", "result":"condition error"}`, errors.New("condition error")
	}
	stUpdateIn.Where = append(stUpdateIn.Where, slExpression...)

	sInstruction = vqlexp.RegExpCollection["WhereToEnd"].ReplaceAllLiteralString(sInstruction, "")

	sColumnsValues := vqlexp.RegExpCollection["UpdateSetToEnd"].FindString(sInstruction)
	sColumnsValues = vqlexp.RegExpCollection["UpdateSetWord"].ReplaceAllLiteralString(sColumnsValues, "")
	slColumnsValues := vqlexp.RegExpCollection["Comma"].Split(sColumnsValues, -1)

	if len(slColumnsValues) == 0 || slColumnsValues[0] == "" {
		return `{"state":"error", "result":"incorrect command syntax"}`, errors.New("incorrect command syntax")
	}

	for _, sColVal := range slColumnsValues {
		slColVal := vqlexp.RegExpCollection["SignEqual"].Split(sColVal, -1)
		sCol := slColVal[0]
		sVal := slColVal[1]

		sCol = vqlexp.RegExpCollection["Spaces"].ReplaceAllLiteralString(sCol, "")
		sCol = trimQuotationMarks(sCol)

		if len(sCol) == 0 {
			return `{"state":"error", "result":"incorrect syntax"}`, errors.New("incorrect syntax")
		}

		sVal = strings.TrimSpace(sVal)
		sVal = trimQuotationMarks(sVal)

		stUpdateIn.Couples[sCol] = sVal
	}

	sTable := vqlexp.RegExpCollection["UpdateSetToEnd"].ReplaceAllLiteralString(sInstruction, "")
	sTable = vqlexp.RegExpCollection["Spaces"].ReplaceAllLiteralString(sTable, "")
	sTable = trimQuotationMarks(sTable)

	if !core.IfExistTable(sDB, sTable) {
		return `{"state":"error", "result":"invalid table name"}`, errors.New("invalid table name")
	}

	// Parsing an expression - End

	// Post checking

	isLuxUser, stFlagsAcs, err := friendlyPostChecker(sDB, sTable, sLogin, stAccess)
	if err != nil {
		stRes.State = "error"
		stRes.Result = err.Error()
		return ecowriter.EncodeJSON(stRes), err
	}

	if !isLuxUser && !stFlagsAcs.Update {
		return `{"state":"error", "result":"not enough rights"}`, errors.New("not enough rights")
	}

	// Request execution

	slResultIDs, isOkUpdate = core.UpdateRows(sDB, sTable, stUpdateIn)
	if !isOkUpdate {
		return `{"state":"error", "result":"the record(s) cannot be updated"}`, errors.New("the record cannot be updated")
	}

	stResArr.State = "ok"
	stResArr.Result = slResultIDs
	return ecowriter.EncodeJSON(stResArr), nil
}

func (q tQuery) DMLDelete() (result string, err error) {
	// This function is complete
	sOperation := "internal -> analyzers -> sql -> DML -> DMLDelete"
	defer func() { e.Wrapper(sOperation, err) }()

	var (
		slResultIDs = make([]uint64, 0)
		isOkDel     bool
		stRes       gtypes.TResponse
		stResArr    gtypes.TResponseUints
		stDeleteIn  = gtypes.TDeleteStruct{
			Where:   make([]gtypes.TConditions, 0, 4),
			IsWhere: false,
		}
	)

	// Pre checking

	sLogin, sDB, stAccess, sNewTicket, err := preChecker(q.Ticket)
	if err != nil {
		stRes.State = "error"
		stRes.Result = err.Error()
		return ecowriter.EncodeJSON(stRes), err
	}

	if sNewTicket != "" {
		stResArr.Ticket = sNewTicket
		stRes.Ticket = sNewTicket
	}

	// Parsing an expression - Begin

	sInstruction := vqlexp.RegExpCollection["DeleteWord"].ReplaceAllLiteralString(q.Instruction, "")

	if vqlexp.RegExpCollection["WhereToEnd"].MatchString(sInstruction) {
		sWhere := vqlexp.RegExpCollection["WhereToEnd"].FindString(sInstruction)
		sInstruction = vqlexp.RegExpCollection["WhereToEnd"].ReplaceAllLiteralString(sInstruction, "")
		sWhere = vqlexp.RegExpCollection["Where"].ReplaceAllLiteralString(sWhere, "")
		slExpression, err := parseWhere(sWhere)
		if err != nil {
			return `{"state":"error", "result":"condition error"}`, errors.New("condition error")
		}
		stDeleteIn.Where = append(stDeleteIn.Where, slExpression...)
		stDeleteIn.IsWhere = true
	}

	sTable := vqlexp.RegExpCollection["Spaces"].ReplaceAllLiteralString(sInstruction, "")
	sTable = trimQuotationMarks(sTable)
	if sTable == "" {
		return `{"state":"error", "result":"invalid table name"}`, errors.New("invalid table name")
	}

	if !core.IfExistTable(sDB, sTable) {
		return `{"state":"error", "result":"invalid table name"}`, errors.New("invalid table name")
	}

	// Parsing an expression - End

	// Post checking

	isLuxUser, stFlagsAcs, err := dourPostChecker(sDB, sTable, sLogin, stAccess)
	if err != nil {
		stRes.State = "error"
		stRes.Result = err.Error()
		return ecowriter.EncodeJSON(stRes), err
	}

	if !isLuxUser && !stFlagsAcs.Delete {
		return `{"state":"error", "result":"not enough rights"}`, errors.New("not enough rights")
	}

	// Request execution

	slResultIDs, isOkDel = core.DeleteRows(sDB, sTable, stDeleteIn)
	if !isOkDel {
		return `{"state":"error", "result":"the record(s) cannot be updated"}`, errors.New("the record cannot be updated")
	}

	stResArr.State = "ok"
	stResArr.Result = slResultIDs
	return ecowriter.EncodeJSON(stResArr), nil
}

func (q tQuery) DMLCommit() (result string, err error) {
	// -
	sOperation := "internal -> analyzers -> sql -> DML -> DMLCommit"
	defer func() { e.Wrapper(sOperation, err) }()

	return "DMLCommit", nil
}

func (q tQuery) DMLRollback() (result string, err error) {
	// -
	sOperation := "internal -> analyzers -> sql -> DML -> DMLRollback"
	defer func() { e.Wrapper(sOperation, err) }()

	return "DMLRollback", nil
}

func (q tQuery) DMLTruncateTable() (result string, err error) {
	// This method is complete
	sOperation := "internal -> analyzers -> sql -> DML -> DMLTruncate"
	defer func() { e.Wrapper(sOperation, err) }()

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

	sTable := vqlexp.RegExpCollection["TruncateTableWord"].ReplaceAllLiteralString(q.Instruction, "")
	sTable = strings.TrimSpace(sTable)
	sTable = trimQuotationMarks(sTable)

	if !core.IfExistTable(sDB, sTable) {
		return `{"state":"error", "result":"invalid table name"}`, errors.New("invalid table name")
	}

	// Parsing an expression - End

	// Post checking

	isLuxUser, stFlagsAcs, err := friendlyPostChecker(sDB, sTable, sLogin, stAccess)
	if err != nil {
		stRes.State = "error"
		stRes.Result = err.Error()
		return ecowriter.EncodeJSON(stRes), err
	}

	if !isLuxUser && !stFlagsAcs.Delete {
		return `{"state":"error", "result":"not enough rights"}`, errors.New("not enough rights")
	}

	// Request execution

	if !core.TruncateTable(sDB, sTable) {
		return `{"state":"error", "result":"the table cannot be truncated"}`, errors.New("the table cannot be truncated")
	}

	stRes.State = "ok"
	return ecowriter.EncodeJSON(stRes), nil
}
