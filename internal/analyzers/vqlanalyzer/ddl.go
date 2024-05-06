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
	sOperation := "internal -> analyzers -> sql -> DDL -> DDLCreate"
	defer func() { e.Wrapper(sOperation, err) }()

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

	isRT := vqlexp.RegExpCollection["AlterDatabaseRenameTo"].MatchString(q.Instruction)
	if !isRT {
		return `{"state":"error", "result":"invalid command format"}`, errors.New("invalid command format")
	}

	sOldDBName := vqlexp.RegExpCollection["AlterDatabaseRenameTo"].FindString(q.Instruction)
	sOldDBName = vqlexp.RegExpCollection["AlterDatabaseWord"].ReplaceAllLiteralString(sOldDBName, "")
	sOldDBName = vqlexp.RegExpCollection["RenameTo"].ReplaceAllLiteralString(sOldDBName, "")
	sOldDBName = trimQuotationMarks(sOldDBName)
	sOldDBName = strings.TrimSpace(sOldDBName)

	sNewDBName := vqlexp.RegExpCollection["AlterDatabaseRenameTo"].ReplaceAllLiteralString(q.Instruction, "")
	sNewDBName = trimQuotationMarks(sNewDBName)
	sNewDBName = strings.TrimSpace(sNewDBName)

	if sOldDBName == "" || sNewDBName == "" {
		return `{"state":"error", "result":"invalid command format"}`, errors.New("invalid command format")
	}

	// Parsing an expression - End

	// Post checking and execution

	_, isOk := core.GetDBInfo(sOldDBName)
	if isOk {
		stDBAccess, isOk := core.GetDBAccess(sOldDBName)
		if isOk {
			stFlagsAcs, isOkFlags := stDBAccess.Flags[sLogin]
			if stDBAccess.Owner != sLogin {
				var isLuxUser bool = false
				for iRole := range stAccess.Roles {
					if iRole == int(gauth.ADMIN) || iRole == int(gauth.ENGINEER) {
						isLuxUser = true
						break
					}
				}
				if !isLuxUser {
					if isOkFlags {
						if !stFlagsAcs.Update {
							return `{"state":"error", "result":"not enough rights"}`, errors.New("not enough rights")
						}
					} else {
						return `{"state":"error", "result":"not enough rights"}`, errors.New("not enough rights")
					}
				}
			}
			if !core.RenameDB(sOldDBName, sNewDBName, true) {
				stRes.State = "error"
				stRes.Result = "the database cannot be renamed"
				return ecowriter.EncodeJSON(stRes), errors.New("the database cannot be renamed")
			}
		} else {
			stRes.State = "error"
			stRes.Result = "internal error"
			return ecowriter.EncodeJSON(stRes), errors.New("internal error")
		}
	} else {
		stRes.State = "error"
		stRes.Result = "invalid database name"
		return ecowriter.EncodeJSON(stRes), errors.New("invalid database name")
	}

	stRes.State = "ok"
	return ecowriter.EncodeJSON(stRes), nil
}

func (q tQuery) DDLAlterTableAdd() (result string, err error) {
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

	sTableName := vqlexp.RegExpCollection["AlterTableAdd"].FindString(q.Instruction)
	sTableName = vqlexp.RegExpCollection["AlterTableWord"].ReplaceAllLiteralString(sTableName, "")
	sTableName = vqlexp.RegExpCollection["ADD"].ReplaceAllLiteralString(sTableName, "")
	sTableName = trimQuotationMarks(sTableName)
	sTableName = strings.TrimSpace(sTableName)

	sColumns := vqlexp.RegExpCollection["AlterTableAdd"].ReplaceAllLiteralString(q.Instruction, "")
	sColumns = vqlexp.RegExpCollection["TableParenthesis"].ReplaceAllLiteralString(sColumns, "")
	slColumnsIn := vqlexp.RegExpCollection["Comma"].Split(sColumns, -1)

	var slStColumns = []gtypes.TColumnForWrite{}

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

		slStColumns = append(slStColumns, stCol)
	}

	if len(slStColumns) < 1 {
		return `{"state":"error", "result":"invalid command format"}`, errors.New("invalid command format")
	}

	// Parsing an expression - End

	// Post checking

	isLuxUser, stFlagsAcs, err := dourPostChecker(sDB, sTableName, sLogin, stAccess)
	if err != nil {
		stRes.State = "error"
		stRes.Result = err.Error()
		return ecowriter.EncodeJSON(stRes), err
	}

	if !isLuxUser {
		if !(stFlagsAcs.Alter && stFlagsAcs.Create) {
			return `{"state":"error", "result":"not enough rights"}`, errors.New("not enough rights")
		}
	}

	// Request execution

	for _, stColName := range slStColumns {
		if !core.CreateColumn(sDB, sTableName, stColName.Name, true, stColName.Spec) {
			stRes.State = "error"
			stRes.Result = "the column cannot be added"
			return ecowriter.EncodeJSON(stRes), errors.New("the column cannot be added")
		}
	}

	stRes.State = "ok"
	return ecowriter.EncodeJSON(stRes), nil
}

func (q tQuery) DDLAlterTableDrop() (result string, err error) {
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

	sTableName := vqlexp.RegExpCollection["AlterTableDrop"].FindString(q.Instruction)
	sTableName = vqlexp.RegExpCollection["AlterTableWord"].ReplaceAllLiteralString(sTableName, "")
	sTableName = vqlexp.RegExpCollection["DROP"].ReplaceAllLiteralString(sTableName, "")
	sTableName = trimQuotationMarks(sTableName)
	sTableName = strings.TrimSpace(sTableName)

	if sTableName == "" {
		return `{"state":"error", "result":"invalid command format"}`, errors.New("invalid command format")
	}

	sColumns := vqlexp.RegExpCollection["AlterTableDrop"].ReplaceAllLiteralString(q.Instruction, "")
	sColumns = vqlexp.RegExpCollection["TableParenthesis"].ReplaceAllLiteralString(sColumns, "")
	slColumnsIn := vqlexp.RegExpCollection["Comma"].Split(sColumns, -1)

	var slSColumns = []string{}

	for _, sColumn := range slColumnsIn {
		sColumn = vqlexp.RegExpCollection["Spaces"].ReplaceAllLiteralString(sColumn, "")
		sColumn = trimQuotationMarks(sColumn)

		slSColumns = append(slSColumns, sColumn)
	}

	if len(slSColumns) < 1 {
		return `{"state":"error", "result":"invalid command format"}`, errors.New("invalid command format")
	}

	// Parsing an expression - End

	// Post checking

	sLuxUser, stFlagsAcs, err := dourPostChecker(sDB, sTableName, sLogin, stAccess)
	if err != nil {
		stRes.State = "error"
		stRes.Result = err.Error()
		return ecowriter.EncodeJSON(stRes), err
	}

	if !sLuxUser {
		if !(stFlagsAcs.Alter && stFlagsAcs.Drop) {
			return `{"state":"error", "result":"not enough rights"}`, errors.New("not enough rights")
		}
	}

	// Request execution
	for _, sColName := range slSColumns {
		if !core.RemoveColumn(sDB, sTableName, sColName) {
			stRes.State = "error"
			stRes.Result = "the column cannot be deleted"
			return ecowriter.EncodeJSON(stRes), errors.New("the column cannot be deleted")
		}
	}

	stRes.State = "ok"
	return ecowriter.EncodeJSON(stRes), nil
}

func (q tQuery) DDLAlterTableModify() (result string, err error) {
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

	sTableName := vqlexp.RegExpCollection["AlterTableModify"].FindString(q.Instruction)
	sTableName = vqlexp.RegExpCollection["AlterTableWord"].ReplaceAllLiteralString(sTableName, "")
	sTableName = vqlexp.RegExpCollection["MODIFY"].ReplaceAllLiteralString(sTableName, "")
	sTableName = trimQuotationMarks(sTableName)
	sTableName = strings.TrimSpace(sTableName)

	sColumns := vqlexp.RegExpCollection["AlterTableModify"].ReplaceAllLiteralString(q.Instruction, "")
	sColumns = vqlexp.RegExpCollection["TableParenthesis"].ReplaceAllLiteralString(sColumns, "")
	slColumnsIn := vqlexp.RegExpCollection["Comma"].Split(sColumns, -1)

	var slStColumns = []gtypes.TColumnForWrite{}

	for _, sColumn := range slColumnsIn {
		stCol := gtypes.TColumnForWrite{
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

		if vqlexp.RegExpCollection["ColumnUnique"].MatchString(sColumn) {
			sColumn = vqlexp.RegExpCollection["ColumnUnique"].ReplaceAllLiteralString(sColumn, "")
			stCol.Spec.Unique = true
			// col.IsChUniqut = true
		}
		if vqlexp.RegExpCollection["ColumnNotNull"].MatchString(sColumn) {
			sColumn = vqlexp.RegExpCollection["ColumnNotNull"].ReplaceAllLiteralString(sColumn, "")
			stCol.Spec.NotNull = true
			// col.IsChNotNull = true
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
			// col.IsChDefault = true
		}

		if vqlexp.RegExpCollection["RenameTo"].MatchString(sColumn) {
			slNames := vqlexp.RegExpCollection["RenameTo"].Split(sColumn, -1)
			sOldName := slNames[0]
			sNewName := slNames[1]

			sOldName = vqlexp.RegExpCollection["Spaces"].ReplaceAllLiteralString(sOldName, "")
			sOldName = trimQuotationMarks(sOldName)

			sNewName = vqlexp.RegExpCollection["Spaces"].ReplaceAllLiteralString(sNewName, "")
			sNewName = trimQuotationMarks(sNewName)

			if sNewName != sOldName {
				stCol.Name = sNewName
				stCol.OldName = sOldName
				stCol.IsChName = true
			} else {
				stCol.Name = sOldName
			}
		} else {
			sColumn = vqlexp.RegExpCollection["Spaces"].ReplaceAllLiteralString(sColumn, "")
			sColumn = trimQuotationMarks(sColumn)

			stCol.Name = sColumn
		}

		if stCol.Name != "" {
			slStColumns = append(slStColumns, stCol)
		}
	}

	if len(slStColumns) < 1 {
		return `{"state":"error", "result":"invalid command format"}`, errors.New("invalid command format")
	}

	// Parsing an expression - End

	// Post checking

	isLuxUser, stFlagsAcs, err := dourPostChecker(sDB, sTableName, sLogin, stAccess)
	if err != nil {
		stRes.State = "error"
		stRes.Result = err.Error()
		return ecowriter.EncodeJSON(stRes), err
	}

	if !isLuxUser {
		if !stFlagsAcs.Alter {
			return `{"state":"error", "result":"not enough rights"}`, errors.New("not enough rights")
		}
	}

	// Request execution

	for _, stColumn := range slStColumns {
		if !core.ChangeColumn(sDB, sTableName, stColumn, true) {
			stRes.State = "error"
			stRes.Result = "the column cannot be changed"
			return ecowriter.EncodeJSON(stRes), errors.New("the column cannot be changed")
		}
	}

	stRes.State = "ok"
	return ecowriter.EncodeJSON(stRes), nil
}

func (q tQuery) DDLAlterTableRenameTo() (result string, err error) {
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

	isRT := vqlexp.RegExpCollection["AlterTableRenameTo"].MatchString(q.Instruction)
	if !isRT {
		return `{"state":"error", "result":"invalid command format"}`, errors.New("invalid command format")
	}

	sOldTableName := vqlexp.RegExpCollection["AlterTableRenameTo"].FindString(q.Instruction)
	sOldTableName = vqlexp.RegExpCollection["AlterTableWord"].ReplaceAllLiteralString(sOldTableName, "")
	sOldTableName = vqlexp.RegExpCollection["RenameTo"].ReplaceAllLiteralString(sOldTableName, "")
	sOldTableName = trimQuotationMarks(sOldTableName)
	sOldTableName = strings.TrimSpace(sOldTableName)

	sNewTableName := vqlexp.RegExpCollection["AlterTableRenameTo"].ReplaceAllLiteralString(q.Instruction, "")
	sNewTableName = trimQuotationMarks(sNewTableName)
	sNewTableName = strings.TrimSpace(sNewTableName)

	if sOldTableName == "" || sNewTableName == "" {
		return `{"state":"error", "result":"invalid command format"}`, errors.New("invalid command format")
	}

	// Parsing an expression - End

	// Post checking

	isLuxUser, stFlagsAcs, err := dourPostChecker(sDB, sOldTableName, sLogin, stAccess)
	if err != nil {
		stRes.State = "error"
		stRes.Result = err.Error()
		return ecowriter.EncodeJSON(stRes), err
	}

	if !isLuxUser {
		if !stFlagsAcs.Alter {
			return `{"state":"error", "result":"not enough rights"}`, errors.New("not enough rights")
		}
	}

	// Request execution

	if !core.RenameTable(sDB, sOldTableName, sNewTableName, true) {
		stRes.State = "error"
		stRes.Result = "the database cannot be renamed"
		return ecowriter.EncodeJSON(stRes), errors.New("the database cannot be renamed")
	}

	stRes.State = "ok"
	return ecowriter.EncodeJSON(stRes), nil
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
	// This method is complete
	sOperation := "internal -> analyzers -> sql -> DDL -> DDLAlter"
	defer func() { e.Wrapper(sOperation, err) }()

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

	isIE := vqlexp.RegExpCollection["IfExistsWord"].MatchString(q.Instruction)

	sDB := vqlexp.RegExpCollection["DropDatabaseWord"].ReplaceAllLiteralString(q.Instruction, "")
	if isIE {
		sDB = vqlexp.RegExpCollection["IfExistsWord"].ReplaceAllLiteralString(sDB, "")
	}
	sDB = strings.TrimSpace(sDB)
	sDB = trimQuotationMarks(sDB)

	// Parsing an expression - End

	// Post checking

	_, isOk := core.GetDBInfo(sDB)
	if !isOk {
		if isIE {
			stRes.State = "error"
			stRes.Result = "the database not exists"
			return ecowriter.EncodeJSON(stRes), errors.New("the database not exists")
		}

		stRes.State = "ok"
		return ecowriter.EncodeJSON(stRes), nil
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

	// Request execution

	if !core.RemoveDB(sDB) {
		stRes.State = "error"
		stRes.Result = "the database cannot be deleted"
		return ecowriter.EncodeJSON(stRes), errors.New("the database cannot be deleted")
	}

	stRes.State = "ok"
	return ecowriter.EncodeJSON(stRes), nil
}

func (q tQuery) DDLDropTable() (result string, err error) {
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

	isIE := vqlexp.RegExpCollection["IfExistsWord"].MatchString(q.Instruction)

	sTable := vqlexp.RegExpCollection["DropTableWord"].ReplaceAllLiteralString(q.Instruction, "")
	if isIE {
		sTable = vqlexp.RegExpCollection["IfExistsWord"].ReplaceAllLiteralString(sTable, "")
	}
	sTable = strings.TrimSpace(sTable)
	sTable = trimQuotationMarks(sTable)

	// Parsing an expression - End

	// Post checking

	isLuxUser, stFlagsAcs, err := dourPostChecker(sDB, sTable, sLogin, stAccess)
	if err != nil {
		stRes.State = "error"
		stRes.Result = err.Error()
		return ecowriter.EncodeJSON(stRes), err
	}

	if !isLuxUser && !stFlagsAcs.Drop {
		return `{"state":"error", "result":"not enough rights"}`, errors.New("not enough rights")
	}

	// Request execution

	if !core.RemoveTable(sDB, sTable) {
		return `{"state":"error", "result":"the table cannot be deleted"}`, errors.New("the table cannot be deleted")
	}

	stRes.State = "ok"
	return ecowriter.EncodeJSON(stRes), nil
}

func (q tQuery) DDLDrop() (result string, err error) {
	// This method is complete
	sOperation := "internal -> analyzers -> sql -> DDL -> DDLDrop"
	defer func() { e.Wrapper(sOperation, err) }()

	isDB := vqlexp.RegExpCollection["DropDatabaseWord"].MatchString(q.Instruction)
	isTable := vqlexp.RegExpCollection["DropTableWord"].MatchString(q.Instruction)

	if isDB {
		return q.DDLDropDB()
	} else if isTable {
		return q.DDLDropTable()
	}

	return `{"state":"error", "result":"unknown command"}`, errors.New("unknown command")
}
