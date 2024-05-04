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

// DCL — Data Control Language (язык управления данными)

func (q tQuery) DCLGrant() (result string, err error) {
	// This method is complete
	op := "internal -> analyzers -> sql -> DCL -> DCLGrant"
	defer func() { e.Wrapper(op, err) }()

	var stRes gtypes.TResponse
	var (
		slDBs   []string
		slUsers []string
	)

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

	sPrivileges := vqlexp.RegExpCollection["GrantPrivileges"].FindString(q.Instruction)
	sPrivileges = vqlexp.RegExpCollection["GrantWord"].ReplaceAllLiteralString(sPrivileges, "")
	sPrivileges = vqlexp.RegExpCollection["ON"].ReplaceAllLiteralString(sPrivileges, "")
	slPrivileges := vqlexp.RegExpCollection["GrantPrivilegesList"].FindAllString(sPrivileges, -1)

	if len(slPrivileges) == 0 {
		return `{"state":"error", "result":"privileges are not specified"}`, errors.New("privileges are not specified")
	}

	sDBs := vqlexp.RegExpCollection["GrantOnTo"].FindString(q.Instruction)
	sDBs = vqlexp.RegExpCollection["ON"].ReplaceAllLiteralString(sDBs, "")
	sDBs = vqlexp.RegExpCollection["TO"].ReplaceAllLiteralString(sDBs, "")
	sDBs = vqlexp.RegExpCollection["Spaces"].ReplaceAllLiteralString(sDBs, "")
	sDBs = trimQuotationMarks(sDBs)
	slDBsIn := vqlexp.RegExpCollection["Comma"].Split(sDBs, -1)
	for _, sDB := range slDBsIn {
		if _, ok := core.GetDBInfo(sDB); ok {
			slDBs = append(slDBs, sDB)
		}
	}
	if len(slDBs) == 0 {
		return `{"state":"error", "result":"databases are not specified"}`, errors.New("databases are not specified")
	}

	sUsers := vqlexp.RegExpCollection["GrantToEnd"].FindString(q.Instruction)
	sUsers = vqlexp.RegExpCollection["TO"].ReplaceAllLiteralString(sUsers, "")
	sUsers = vqlexp.RegExpCollection["Spaces"].ReplaceAllLiteralString(sUsers, "")
	sUsers = trimQuotationMarks(sUsers)
	slUsersIn := vqlexp.RegExpCollection["Comma"].Split(sUsers, -1)
	for _, sUser := range slUsersIn {
		if _, err := gauth.GetProfile(sUser); err == nil {
			slUsers = append(slUsers, sUser)
		}
	}
	if len(slUsers) == 0 {
		return `{"state":"error", "result":"users are not specified"}`, errors.New("users are not specified")
	}

	// Parsing an expression - End

	// Post checking and execution

	for _, sDB := range slDBs {
		stDBAccess, isOk := core.GetDBAccess(sDB)
		if isOk {
			if stDBAccess.Owner != sLogin {
				var isLuxUser bool = false
				for role := range stAccess.Roles {
					if role == int(gauth.ADMIN) || role == int(gauth.ENGINEER) {
						isLuxUser = true
						break
					}
				}
				if !isLuxUser {
					return `{"state":"error", "result":"auth error"}`, errors.New("auth error")
				}
			}
			for _, sUser := range slUsers {
				// var stAccessFlags gtypes.TAccessFlags
				stAccessFlags, isOk := stDBAccess.Flags[sUser]
				if !isOk {
					stAccessFlags = gtypes.TAccessFlags{}
				}

				for _, sPrivilege := range slPrivileges {
					switch strings.ToLower(sPrivilege) {
					case "create":
						stAccessFlags.Create = true
					case "select":
						stAccessFlags.Select = true
					case "insert":
						stAccessFlags.Insert = true
					case "update":
						stAccessFlags.Update = true
					case "delete":
						stAccessFlags.Delete = true
					}
				}

				// core.StorageInfo.Access[db].Flags[user] = aFlags
				core.SetAccessFlags(sDB, sUser, stAccessFlags)
			}
		}
	}

	stRes.State = "ok"
	return ecowriter.EncodeJSON(stRes), nil
}

func (q tQuery) DCLRevoke() (result string, err error) {
	// This method is complete
	op := "internal -> analyzers -> sql -> DCL -> DCLRevoke"
	defer func() { e.Wrapper(op, err) }()

	var stRes gtypes.TResponse
	var (
		slDBs   []string
		slUsers []string
	)

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

	sPrivileges := vqlexp.RegExpCollection["RevokePrivileges"].FindString(q.Instruction)
	sPrivileges = vqlexp.RegExpCollection["RevokeWord"].ReplaceAllLiteralString(sPrivileges, "")
	sPrivileges = vqlexp.RegExpCollection["ON"].ReplaceAllLiteralString(sPrivileges, "")
	slPrivileges := vqlexp.RegExpCollection["RevokePrivilegesList"].FindAllString(sPrivileges, -1)

	if len(slPrivileges) == 0 {
		return `{"state":"error", "result":"privileges are not specified"}`, errors.New("privileges are not specified")
	}

	sDBs := vqlexp.RegExpCollection["RevokeOnTo"].FindString(q.Instruction)
	sDBs = vqlexp.RegExpCollection["ON"].ReplaceAllLiteralString(sDBs, "")
	sDBs = vqlexp.RegExpCollection["TO"].ReplaceAllLiteralString(sDBs, "")
	sDBs = vqlexp.RegExpCollection["Spaces"].ReplaceAllLiteralString(sDBs, "")
	sDBs = trimQuotationMarks(sDBs)
	slDBsIn := vqlexp.RegExpCollection["Comma"].Split(sDBs, -1)
	for _, sDB := range slDBsIn {
		if _, isOk := core.GetDBInfo(sDB); isOk {
			slDBs = append(slDBs, sDB)
		}
	}
	if len(slDBs) == 0 {
		return `{"state":"error", "result":"databases are not specified"}`, errors.New("databases are not specified")
	}

	sUsers := vqlexp.RegExpCollection["RevokeToEnd"].FindString(q.Instruction)
	sUsers = vqlexp.RegExpCollection["TO"].ReplaceAllLiteralString(sUsers, "")
	sUsers = vqlexp.RegExpCollection["Spaces"].ReplaceAllLiteralString(sUsers, "")
	sUsers = trimQuotationMarks(sUsers)
	slUsersIn := vqlexp.RegExpCollection["Comma"].Split(sUsers, -1)
	for _, sUser := range slUsersIn {
		if _, err := gauth.GetProfile(sUser); err == nil {
			slUsers = append(slUsers, sUser)
		}
	}
	if len(slUsers) == 0 {
		return `{"state":"error", "result":"users are not specified"}`, errors.New("users are not specified")
	}

	// Parsing an expression - End

	// Post checking and execution

	for _, sDB := range slDBs {
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
					return `{"state":"error", "result":"auth error"}`, errors.New("auth error")
				}
			}
			for _, sUser := range slUsers {
				// var aFlags gtypes.TAccessFlags
				stAccessFlags, isOk := stDBAccess.Flags[sUser]
				if !isOk {
					stAccessFlags = gtypes.TAccessFlags{}
				}

				for _, sPrivilege := range slPrivileges {
					switch strings.ToLower(sPrivilege) {
					case "create":
						stAccessFlags.Create = false
					case "select":
						stAccessFlags.Select = false
					case "insert":
						stAccessFlags.Insert = false
					case "update":
						stAccessFlags.Update = false
					case "delete":
						stAccessFlags.Delete = false
					}
				}

				// core.StorageInfo.Access[db].Flags[user] = aFlags
				core.SetAccessFlags(sDB, sUser, stAccessFlags)
			}
		}
	}

	stRes.State = "ok"
	return ecowriter.EncodeJSON(stRes), nil
}

func (q tQuery) DCLUse() (result string, err error) {
	// This method is complete
	op := "internal -> analyzers -> sql -> DCL -> DCLUse"
	defer func() { e.Wrapper(op, err) }()

	var sTicket string
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
		sTicket = sNewTicket
		stRes.Ticket = sNewTicket
	} else {
		sTicket = q.Ticket
	}

	// Parsing an expression - Begin

	sDB := vqlexp.RegExpCollection["UseWord"].ReplaceAllLiteralString(q.Instruction, "")
	sDB = strings.TrimSpace(sDB)
	sDB = trimQuotationMarks(sDB)

	if !vqlexp.RegExpCollection["EntityName"].MatchString(sDB) {
		return `{"state":"error", "result":"invalid database name"}`, errors.New("invalid database name")
	}

	if !core.LocalCoreSettings.FriendlyMode {
		if _, isOk := core.GetDBInfo(sDB); !isOk {
			return `{"state":"error", "result":"the database does not exist"}`, errors.New("the database does not exist")
		}
	}

	// Parsing an expression - End

	// Post checking

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
				stAccessFlags, isOk := stDBAccess.Flags[sLogin]
				if !isOk {
					return `{"state":"error", "result":"auth error"}`, errors.New("auth error")
				}
				if !stAccessFlags.AnyTrue() {
					return `{"state":"error", "result":"auth error"}`, errors.New("auth error")
				}
			}
		}
	}

	// Request execution

	core.States[sTicket] = core.TState{
		CurrentDB: sDB,
	}

	stRes.State = "ok"
	stRes.Result = sDB
	return ecowriter.EncodeJSON(stRes), nil
}

func (q tQuery) DCLShow() (result string, err error) {
	// This method is complete
	op := "internal -> analyzers -> sql -> DCL -> DCLShow"
	defer func() { e.Wrapper(op, err) }()

	var (
		stRes      gtypes.TResponse
		stResArray gtypes.TResponseStrings
	)

	// Pre checking

	if q.Ticket == "" {
		return `{"state":"error", "result":"an empty ticket"}`, errors.New("an empty ticket")
	}

	_, stAccess, sNewTicket, err := gauth.CheckTicket(q.Ticket)
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

	isDBs := vqlexp.RegExpCollection["ShowDatabasesWord"].MatchString(q.Instruction)
	isTables := vqlexp.RegExpCollection["ShowTablesWord"].MatchString(q.Instruction)

	// Parsing an expression - End

	// Post checking and execution

	if isDBs {
		var slNamesDBs []string = []string{}
		for sNameDB := range core.StorageInfo.DBs {
			slNamesDBs = append(slNamesDBs, sNameDB)
		}

		stResArray.State = "ok"
		stResArray.Ticket = stRes.Ticket
		stResArray.Result = slNamesDBs
		return ecowriter.EncodeJSON(stResArray), nil
	} else if isTables {
		var slNamesTables []string = []string{}

		stState, isOk := core.States[q.Ticket]
		if !isOk {
			stRes.State = "error"
			stRes.Result = "unknown database"
			return ecowriter.EncodeJSON(stRes), errors.New("unknown database")
		}
		sDB := stState.CurrentDB
		if sDB == "" {
			stRes.State = "error"
			stRes.Result = "no database selected"
			return ecowriter.EncodeJSON(stRes), errors.New("no database selected")
		}

		stDBInfo, isOk := core.GetDBInfo(sDB)
		if !isOk {
			stRes.State = "error"
			stRes.Result = "incorrect database"
			return ecowriter.EncodeJSON(stRes), errors.New("incorrect database")
		}

		for sNameTable := range stDBInfo.Tables {
			slNamesTables = append(slNamesTables, sNameTable)
		}

		stResArray.State = "ok"
		stResArray.Ticket = stRes.Ticket
		stResArray.Result = slNamesTables
		return ecowriter.EncodeJSON(stResArray), nil
	}

	stRes.State = "error"
	stRes.Result = "unknown command"
	return ecowriter.EncodeJSON(stRes), nil
}

func (q tQuery) DCLDesc() (result string, err error) {
	// This method is complete
	op := "internal -> analyzers -> sql -> DCL -> DCLDesc"
	defer func() { e.Wrapper(op, err) }()

	var stRes gtypes.TResponse
	var stResArray gtypes.TResponseColumns

	var sTable string

	// Pre checking

	sLogin, sDB, stAccess, sNewTicket, err := preChecker(q.Ticket)
	if err != nil {
		stRes.State = "error"
		stRes.Result = err.Error()
		return ecowriter.EncodeJSON(stRes), err
	}

	if sNewTicket != "" {
		stResArray.Ticket = sNewTicket
		stRes.Ticket = sNewTicket
	}

	// Parsing an expression - Begin

	if vqlexp.RegExpCollection["SearchExplain"].MatchString(q.Instruction) {
		sTable = vqlexp.RegExpCollection["ExplainWord"].ReplaceAllLiteralString(q.Instruction, "")
	} else if vqlexp.RegExpCollection["SearchDescribe"].MatchString(q.Instruction) {
		sTable = vqlexp.RegExpCollection["DescribeWord"].ReplaceAllLiteralString(q.Instruction, "")
	} else if vqlexp.RegExpCollection["SearchDesc"].MatchString(q.Instruction) {
		sTable = vqlexp.RegExpCollection["DescWord"].ReplaceAllLiteralString(q.Instruction, "")
	}

	sTable = vqlexp.RegExpCollection["Spaces"].ReplaceAllLiteralString(sTable, "")
	sTable = trimQuotationMarks(sTable)

	// Parsing an expression - End

	// Post checking

	isLuxUser, stAccessFlags, err := dourPostChecker(sDB, sTable, sLogin, stAccess)
	if err != nil {
		stRes.State = "error"
		stRes.Result = err.Error()
		return ecowriter.EncodeJSON(stRes), err
	}

	if !isLuxUser && !stAccessFlags.Select {
		return `{"state":"error", "result":"not enough rights"}`, errors.New("not enough rights")
	}

	// Request execution

	stDBInfo, isOkDB := core.GetDBInfo(sDB)
	if !isOkDB {
		stRes.State = "error"
		stRes.Result = "invalid database name"
		return ecowriter.EncodeJSON(stRes), errors.New("invalid database name")
	}

	stTableInfo, isOk := stDBInfo.Tables[sTable]
	if !isOk {
		stRes.State = "error"
		stRes.Result = "unknown table"
		return ecowriter.EncodeJSON(stRes), errors.New("unknown table")
	}

	if len(stTableInfo.Order) < 1 {
		stRes.State = "error"
		stRes.Result = "there are no columns"
		return ecowriter.EncodeJSON(stRes), errors.New("there are no columns")
	}

	for _, sColName := range stTableInfo.Order {
		stColumn, isOkCol := stTableInfo.Columns[sColName]
		if isOkCol {
			var stResColumn gtypes.TResultColumn

			stResColumn.Field = stColumn.Name
			stResColumn.Default = stColumn.Specification.Default
			stResColumn.NotNull = stColumn.Specification.NotNull
			stResColumn.Unique = stColumn.Specification.Unique
			stResColumn.LastUpdate = stColumn.LastUpdate

			stResArray.Result = append(stResArray.Result, stResColumn)
		}
	}
	stResArray.State = "ok"
	return ecowriter.EncodeJSON(stResArray), nil
}

func (q tQuery) DCLAuth() (result string, err error) {
	// This method is complete
	op := "internal -> analyzers -> sql -> DCL -> DCLAuth"
	defer func() { e.Wrapper(op, err) }()

	var slRoles []gauth.TRole

	// Parsing an expression - Begin

	isNew := vqlexp.RegExpCollection["AuthNew"].MatchString(q.Instruction)
	isChange := vqlexp.RegExpCollection["AuthChange"].MatchString(q.Instruction)
	isRemove := vqlexp.RegExpCollection["AuthRemove"].MatchString(q.Instruction)

	sLogin := vqlexp.RegExpCollection["Login"].FindString(q.Instruction)
	sLogin = vqlexp.RegExpCollection["LoginWord"].ReplaceAllLiteralString(sLogin, " ")
	sLogin = strings.TrimSpace(sLogin)
	sLogin = trimQuotationMarks(sLogin)

	sPassword := vqlexp.RegExpCollection["Password"].FindString(q.Instruction)
	sPassword = vqlexp.RegExpCollection["PasswordWord"].ReplaceAllLiteralString(sPassword, " ")
	sPassword = strings.TrimSpace(sPassword)
	sPassword = trimQuotationMarks(sPassword)

	sHash := vqlexp.RegExpCollection["Hash"].FindString(q.Instruction)
	sHash = vqlexp.RegExpCollection["HashWord"].ReplaceAllLiteralString(sHash, " ")
	sHash = strings.TrimSpace(sHash)
	sHash = trimQuotationMarks(sHash)

	isRole := vqlexp.RegExpCollection["Role"].MatchString(q.Instruction)
	if isRole {
		sRole := vqlexp.RegExpCollection["Role"].FindString(q.Instruction)
		sRole = vqlexp.RegExpCollection["RoleWord"].ReplaceAllLiteralString(sRole, "")
		sRole = vqlexp.RegExpCollection["Spaces"].ReplaceAllLiteralString(sRole, "")
		sRole = trimQuotationMarks(sRole)
		slRoleIn := vqlexp.RegExpCollection["Comma"].Split(sRole, -1)
		if len(slRoleIn) == 0 {
			return `{"state":"error", "result":"incorrect roles"}`, errors.New("incorrect roles")
		}
		for _, sRoleIt := range slRoleIn {
			switch strings.ToUpper(sRoleIt) {
			case "SYSTEM":
				slRoles = append(slRoles, gauth.SYSTEM)
			case "ADMIN":
				slRoles = append(slRoles, gauth.ADMIN)
			case "MANAGER":
				slRoles = append(slRoles, gauth.MANAGER)
			case "ENGINEER":
				slRoles = append(slRoles, gauth.ENGINEER)
			case "USER":
				slRoles = append(slRoles, gauth.USER)
			}
		}
	}

	// Parsing an expression - End

	// Request execution

	if isNew || isChange || isRemove {
		var stRes gtypes.TResponse

		if q.Ticket == "" {
			return `{"state":"error", "result":"an empty ticket"}`, errors.New("an empty ticket")
		}

		_, stCurentAccess, sNewTicket, err := gauth.CheckTicket(q.Ticket)
		if err != nil {
			return `{"state":"error", "result":"authorization failed"}`, err
		}

		if stCurentAccess.Status.IsBad() {
			return `{"state":"error", "result":"auth error"}`, errors.New("auth error")
		}

		if sNewTicket != "" {
			stRes.Ticket = sNewTicket
		}

		var isLuxUser bool = false
		for iRole := range stCurentAccess.Roles {
			if iRole == int(gauth.ADMIN) || iRole == int(gauth.MANAGER) {
				isLuxUser = true
				break
			}
		}

		if !isLuxUser {
			stRes.State = "error"
			stRes.Result = "auth error"
			return ecowriter.EncodeJSON(stRes), errors.New("auth error")
		}

		if isNew {
			stAccess := gauth.TProfile{
				Description: "",
				Status:      gauth.NEW,
			}

			if isRole {
				stAccess.Roles = slRoles
			} else {
				stAccess.Roles = []gauth.TRole{gauth.USER}
			}

			err := gauth.AddUser(sLogin, sPassword, stAccess)
			if err != nil {
				return ecowriter.EncodeJSON(gtypes.TResponse{
					State:  "error",
					Result: err.Error(),
				}), err
			}
		}

		if isChange {
			stAccess, err := gauth.GetProfile(sLogin)
			if err != nil {
				return ecowriter.EncodeJSON(gtypes.TResponse{
					State:  "error",
					Result: err.Error(),
				}), err
			}

			if isRole {
				stAccess.Roles = slRoles
			}

			err = gauth.UpdateUser(sLogin, sPassword, stAccess)
			if err != nil {
				return ecowriter.EncodeJSON(gtypes.TResponse{
					State:  "error",
					Result: err.Error(),
				}), err
			}
		}

		if isRemove {
			err := gauth.DeleteUser(sLogin)
			if err != nil {
				return ecowriter.EncodeJSON(gtypes.TResponse{
					State:  "error",
					Result: err.Error(),
				}), err
			}
		}

		stRes.State = "ok"
		return ecowriter.EncodeJSON(stRes), nil
	}

	stProfile, err := gauth.GetProfile(sLogin)
	if err != nil {
		return `{"state":"error", "result":"auth error"}`, errors.New("auth error")
	}

	if stProfile.Status.IsBad() {
		return `{"state":"error", "result":"auth error"}`, errors.New("auth error")
	}

	stSecret := gtypes.TSecret{
		Login:    sLogin,
		Password: sPassword,
		Hash:     sHash,
	}
	sTicket, err := gauth.NewAuth(&stSecret)
	if err != nil {
		return `{"state":"error", "result":"auth error"}`, errors.New("auth error")
	}

	return ecowriter.EncodeJSON(gtypes.TResponse{
		State:  "ok",
		Ticket: sTicket,
	}), nil
}
