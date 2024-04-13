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

// DCL — язык управления данными (Data Control Language)

func (q tQuery) DCLGrant() (result string, err error) {
	// This method is complete
	op := "internal -> analyzers -> sql -> DCL -> DCLGrant"
	defer func() { e.Wrapper(op, err) }()

	var res gtypes.Response
	var (
		dbs   []string
		users []string
	)

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

	privilegesStr := vqlexp.RegExpCollection["GrantPrivileges"].FindString(q.Instruction)
	privilegesStr = vqlexp.RegExpCollection["GrantWord"].ReplaceAllLiteralString(privilegesStr, "")
	privilegesStr = vqlexp.RegExpCollection["ON"].ReplaceAllLiteralString(privilegesStr, "")
	privileges := vqlexp.RegExpCollection["GrantPrivilegesList"].FindAllString(privilegesStr, -1)

	if len(privileges) == 0 {
		return `{"state":"error", "result":"privileges are not specified"}`, errors.New("privileges are not specified")
	}

	dbsStr := vqlexp.RegExpCollection["GrantOnTo"].FindString(q.Instruction)
	dbsStr = vqlexp.RegExpCollection["ON"].ReplaceAllLiteralString(dbsStr, "")
	dbsStr = vqlexp.RegExpCollection["TO"].ReplaceAllLiteralString(dbsStr, "")
	dbsStr = vqlexp.RegExpCollection["Spaces"].ReplaceAllLiteralString(dbsStr, "")
	dbsStr = trimQuotationMarks(dbsStr)
	dbsIn := vqlexp.RegExpCollection["Comma"].Split(dbsStr, -1)
	for _, db := range dbsIn {
		if _, ok := core.GetDBInfo(db); ok {
			dbs = append(dbs, db)
		}
	}
	if len(dbs) == 0 {
		return `{"state":"error", "result":"databases are not specified"}`, errors.New("databases are not specified")
	}

	usersStr := vqlexp.RegExpCollection["GrantToEnd"].FindString(q.Instruction)
	usersStr = vqlexp.RegExpCollection["TO"].ReplaceAllLiteralString(usersStr, "")
	usersStr = vqlexp.RegExpCollection["Spaces"].ReplaceAllLiteralString(usersStr, "")
	usersStr = trimQuotationMarks(usersStr)
	usersIn := vqlexp.RegExpCollection["Comma"].Split(usersStr, -1)
	for _, user := range usersIn {
		if _, err := gauth.GetProfile(user); err == nil {
			users = append(users, user)
		}
	}
	if len(users) == 0 {
		return `{"state":"error", "result":"users are not specified"}`, errors.New("users are not specified")
	}

	// Parsing an expression - End

	// Post checking and execution

	for _, db := range dbs {
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
					return `{"state":"error", "result":"auth error"}`, errors.New("auth error")
				}
			}
			for _, user := range users {
				var aFlags gtypes.TAccessFlags
				aFlags, ok := dbAccess.Flags[user]
				if !ok {
					aFlags = gtypes.TAccessFlags{}
				}

				for _, privilege := range privileges {
					switch strings.ToLower(privilege) {
					case "create":
						aFlags.Create = true
					case "select":
						aFlags.Select = true
					case "insert":
						aFlags.Insert = true
					case "update":
						aFlags.Update = true
					case "delete":
						aFlags.Delete = true
					}
				}

				// core.StorageInfo.Access[db].Flags[user] = aFlags
				core.SetAccessFlags(db, user, aFlags)
			}
		}
	}

	res.State = "ok"
	return ecowriter.EncodeJSON(res), nil
}

func (q tQuery) DCLRevoke() (result string, err error) {
	// This method is complete
	op := "internal -> analyzers -> sql -> DCL -> DCLRevoke"
	defer func() { e.Wrapper(op, err) }()

	var res gtypes.Response
	var (
		dbs   []string
		users []string
	)

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

	privilegesStr := vqlexp.RegExpCollection["RevokePrivileges"].FindString(q.Instruction)
	privilegesStr = vqlexp.RegExpCollection["RevokeWord"].ReplaceAllLiteralString(privilegesStr, "")
	privilegesStr = vqlexp.RegExpCollection["ON"].ReplaceAllLiteralString(privilegesStr, "")
	privileges := vqlexp.RegExpCollection["RevokePrivilegesList"].FindAllString(privilegesStr, -1)

	if len(privileges) == 0 {
		return `{"state":"error", "result":"privileges are not specified"}`, errors.New("privileges are not specified")
	}

	dbsStr := vqlexp.RegExpCollection["RevokeOnTo"].FindString(q.Instruction)
	dbsStr = vqlexp.RegExpCollection["ON"].ReplaceAllLiteralString(dbsStr, "")
	dbsStr = vqlexp.RegExpCollection["TO"].ReplaceAllLiteralString(dbsStr, "")
	dbsStr = vqlexp.RegExpCollection["Spaces"].ReplaceAllLiteralString(dbsStr, "")
	dbsStr = trimQuotationMarks(dbsStr)
	dbsIn := vqlexp.RegExpCollection["Comma"].Split(dbsStr, -1)
	for _, db := range dbsIn {
		if _, ok := core.GetDBInfo(db); ok {
			dbs = append(dbs, db)
		}
	}
	if len(dbs) == 0 {
		return `{"state":"error", "result":"databases are not specified"}`, errors.New("databases are not specified")
	}

	usersStr := vqlexp.RegExpCollection["RevokeToEnd"].FindString(q.Instruction)
	usersStr = vqlexp.RegExpCollection["TO"].ReplaceAllLiteralString(usersStr, "")
	usersStr = vqlexp.RegExpCollection["Spaces"].ReplaceAllLiteralString(usersStr, "")
	usersStr = trimQuotationMarks(usersStr)
	usersIn := vqlexp.RegExpCollection["Comma"].Split(usersStr, -1)
	for _, user := range usersIn {
		if _, err := gauth.GetProfile(user); err == nil {
			users = append(users, user)
		}
	}
	if len(users) == 0 {
		return `{"state":"error", "result":"users are not specified"}`, errors.New("users are not specified")
	}

	// Parsing an expression - End

	// Post checking and execution

	for _, db := range dbs {
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
					return `{"state":"error", "result":"auth error"}`, errors.New("auth error")
				}
			}
			for _, user := range users {
				var aFlags gtypes.TAccessFlags
				aFlags, ok := dbAccess.Flags[user]
				if !ok {
					aFlags = gtypes.TAccessFlags{}
				}

				for _, privilege := range privileges {
					switch strings.ToLower(privilege) {
					case "create":
						aFlags.Create = false
					case "select":
						aFlags.Select = false
					case "insert":
						aFlags.Insert = false
					case "update":
						aFlags.Update = false
					case "delete":
						aFlags.Delete = false
					}
				}

				// core.StorageInfo.Access[db].Flags[user] = aFlags
				core.SetAccessFlags(db, user, aFlags)
			}
		}
	}

	res.State = "ok"
	return ecowriter.EncodeJSON(res), nil
}

func (q tQuery) DCLUse() (result string, err error) {
	// This method is complete
	op := "internal -> analyzers -> sql -> DCL -> DCLUse"
	defer func() { e.Wrapper(op, err) }()

	var ticket string
	var res gtypes.Response

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
		ticket = newticket
		res.Ticket = newticket
	} else {
		ticket = q.Ticket
	}

	// Parsing an expression - Begin

	db := vqlexp.RegExpCollection["UseWord"].ReplaceAllLiteralString(q.Instruction, "")
	db = strings.TrimSpace(db)
	db = trimQuotationMarks(db)

	if !vqlexp.RegExpCollection["EntityName"].MatchString(db) {
		return `{"state":"error", "result":"invalid database name"}`, errors.New("invalid database name")
	}

	if !core.LocalCoreSettings.FriendlyMode {
		if _, ok := core.GetDBInfo(db); !ok {
			return `{"state":"error", "result":"the database does not exist"}`, errors.New("the database does not exist")
		}
	}

	// Parsing an expression - End

	// Post checking

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
				flags, ok := dbAccess.Flags[login]
				if !ok {
					return `{"state":"error", "result":"auth error"}`, errors.New("auth error")
				}
				if !flags.AnyTrue() {
					return `{"state":"error", "result":"auth error"}`, errors.New("auth error")
				}
			}
		}
	}

	// Request execution

	core.States[ticket] = core.TState{
		CurrentDB: db,
	}

	res.State = "ok"
	res.Result = db
	return ecowriter.EncodeJSON(res), nil
}

func (q tQuery) DCLShow() (result string, err error) {
	// This method is complete
	op := "internal -> analyzers -> sql -> DCL -> DCLShow"
	defer func() { e.Wrapper(op, err) }()

	var (
		res    gtypes.Response
		resArr gtypes.ResponseStrings
	)

	// Pre checking

	if q.Ticket == "" {
		return `{"state":"error", "result":"an empty ticket"}`, errors.New("an empty ticket")
	}

	_, access, newticket, err := gauth.CheckTicket(q.Ticket)
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

	isDBs := vqlexp.RegExpCollection["ShowDatabasesWord"].MatchString(q.Instruction)
	isTables := vqlexp.RegExpCollection["ShowTablesWord"].MatchString(q.Instruction)

	// Parsing an expression - End

	// Post checking and execution

	if isDBs {
		var namesDBs []string = []string{}
		for nameDB := range core.StorageInfo.DBs {
			namesDBs = append(namesDBs, nameDB)
		}

		resArr.State = "ok"
		resArr.Ticket = res.Ticket
		resArr.Result = namesDBs
		return ecowriter.EncodeJSON(resArr), nil
	} else if isTables {
		var namesTables []string = []string{}

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

		dbInfo, ok := core.GetDBInfo(db)
		if !ok {
			res.State = "error"
			res.Result = "incorrect database"
			return ecowriter.EncodeJSON(res), errors.New("incorrect database")
		}

		for nameTable := range dbInfo.Tables {
			namesTables = append(namesTables, nameTable)
		}

		resArr.State = "ok"
		resArr.Ticket = res.Ticket
		resArr.Result = namesTables
		return ecowriter.EncodeJSON(resArr), nil
	}

	res.State = "error"
	res.Result = "unknown command"
	return ecowriter.EncodeJSON(res), nil
}

func (q tQuery) DCLDesc() (result string, err error) {
	// This method is complete
	op := "internal -> analyzers -> sql -> DCL -> DCLDesc"
	defer func() { e.Wrapper(op, err) }()

	var res gtypes.Response
	var resArr gtypes.ResponseColumns

	var table string

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

	if vqlexp.RegExpCollection["SearchExplain"].MatchString(q.Instruction) {
		table = vqlexp.RegExpCollection["ExplainWord"].ReplaceAllLiteralString(q.Instruction, "")
	} else if vqlexp.RegExpCollection["SearchDescribe"].MatchString(q.Instruction) {
		table = vqlexp.RegExpCollection["DescribeWord"].ReplaceAllLiteralString(q.Instruction, "")
	} else if vqlexp.RegExpCollection["SearchDesc"].MatchString(q.Instruction) {
		table = vqlexp.RegExpCollection["DescWord"].ReplaceAllLiteralString(q.Instruction, "")
	}

	table = vqlexp.RegExpCollection["Spaces"].ReplaceAllLiteralString(table, "")
	table = trimQuotationMarks(table)

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

	dbInfo, okDB := core.GetDBInfo(db)
	if !okDB {
		res.State = "error"
		res.Result = "invalid database name"
		return ecowriter.EncodeJSON(res), errors.New("invalid database name")
	}

	tableInfo, ok := dbInfo.Tables[table]
	if !ok {
		res.State = "error"
		res.Result = "unknown table"
		return ecowriter.EncodeJSON(res), errors.New("unknown table")
	}

	if len(tableInfo.Order) < 1 {
		res.State = "error"
		res.Result = "there are no columns"
		return ecowriter.EncodeJSON(res), errors.New("there are no columns")
	}

	for _, colName := range tableInfo.Order {
		column, okCol := tableInfo.Columns[colName]
		if okCol {
			var resColumn gtypes.ResultColumn

			resColumn.Field = column.Name
			resColumn.Default = column.Specification.Default
			resColumn.NotNull = column.Specification.NotNull
			resColumn.Unique = column.Specification.Unique
			resColumn.LastUpdate = column.LastUpdate

			resArr.Result = append(resArr.Result, resColumn)
		}
	}
	resArr.State = "ok"
	return ecowriter.EncodeJSON(resArr), nil
}

func (q tQuery) DCLAuth() (result string, err error) {
	// This method is complete
	op := "internal -> analyzers -> sql -> DCL -> DCLAuth"
	defer func() { e.Wrapper(op, err) }()

	var roles []gauth.TRole

	// Parsing an expression - Begin

	isNew := vqlexp.RegExpCollection["AuthNew"].MatchString(q.Instruction)
	isChange := vqlexp.RegExpCollection["AuthChange"].MatchString(q.Instruction)
	isRemove := vqlexp.RegExpCollection["AuthRemove"].MatchString(q.Instruction)

	login := vqlexp.RegExpCollection["Login"].FindString(q.Instruction)
	login = vqlexp.RegExpCollection["LoginWord"].ReplaceAllLiteralString(login, " ")
	login = strings.TrimSpace(login)
	login = trimQuotationMarks(login)

	password := vqlexp.RegExpCollection["Password"].FindString(q.Instruction)
	password = vqlexp.RegExpCollection["PasswordWord"].ReplaceAllLiteralString(password, " ")
	password = strings.TrimSpace(password)
	password = trimQuotationMarks(password)

	hash := vqlexp.RegExpCollection["Hash"].FindString(q.Instruction)
	hash = vqlexp.RegExpCollection["HashWord"].ReplaceAllLiteralString(hash, " ")
	hash = strings.TrimSpace(hash)
	hash = trimQuotationMarks(hash)

	isRole := vqlexp.RegExpCollection["Role"].MatchString(q.Instruction)
	if isRole {
		roleStr := vqlexp.RegExpCollection["Role"].FindString(q.Instruction)
		roleStr = vqlexp.RegExpCollection["RoleWord"].ReplaceAllLiteralString(roleStr, "")
		roleStr = vqlexp.RegExpCollection["Spaces"].ReplaceAllLiteralString(roleStr, "")
		roleStr = trimQuotationMarks(roleStr)
		roleIn := vqlexp.RegExpCollection["Comma"].Split(roleStr, -1)
		if len(roleIn) == 0 {
			return `{"state":"error", "result":"incorrect roles"}`, errors.New("incorrect roles")
		}
		for _, role := range roleIn {
			switch strings.ToUpper(role) {
			case "SYSTEM":
				roles = append(roles, gauth.SYSTEM)
			case "ADMIN":
				roles = append(roles, gauth.ADMIN)
			case "MANAGER":
				roles = append(roles, gauth.MANAGER)
			case "ENGINEER":
				roles = append(roles, gauth.ENGINEER)
			case "USER":
				roles = append(roles, gauth.USER)
			}
		}
	}

	// Parsing an expression - End

	// Request execution

	if isNew || isChange || isRemove {
		var res gtypes.Response

		if q.Ticket == "" {
			return `{"state":"error", "result":"an empty ticket"}`, errors.New("an empty ticket")
		}

		_, curaccess, newticket, err := gauth.CheckTicket(q.Ticket)
		if err != nil {
			return `{"state":"error", "result":"authorization failed"}`, err
		}

		if curaccess.Status.IsBad() {
			return `{"state":"error", "result":"auth error"}`, errors.New("auth error")
		}

		if newticket != "" {
			res.Ticket = newticket
		}

		var luxUser bool = false
		for role := range curaccess.Roles {
			if role == int(gauth.ADMIN) || role == int(gauth.MANAGER) {
				luxUser = true
				break
			}
		}

		if !luxUser {
			res.State = "error"
			res.Result = "auth error"
			return ecowriter.EncodeJSON(res), errors.New("auth error")
		}

		if isNew {
			access := gauth.TProfile{
				Description: "",
				Status:      gauth.NEW,
			}

			if isRole {
				access.Roles = roles
			} else {
				access.Roles = []gauth.TRole{gauth.USER}
			}

			err := gauth.AddUser(login, password, access)
			if err != nil {
				return ecowriter.EncodeJSON(gtypes.Response{
					State:  "error",
					Result: err.Error(),
				}), err
			}
		}

		if isChange {
			access, err := gauth.GetProfile(login)
			if err != nil {
				return ecowriter.EncodeJSON(gtypes.Response{
					State:  "error",
					Result: err.Error(),
				}), err
			}

			if isRole {
				access.Roles = roles
			}

			err = gauth.UpdateUser(login, password, access)
			if err != nil {
				return ecowriter.EncodeJSON(gtypes.Response{
					State:  "error",
					Result: err.Error(),
				}), err
			}
		}

		if isRemove {
			err := gauth.DeleteUser(login)
			if err != nil {
				return ecowriter.EncodeJSON(gtypes.Response{
					State:  "error",
					Result: err.Error(),
				}), err
			}
		}

		res.State = "ok"
		return ecowriter.EncodeJSON(res), nil
	}

	profile, err := gauth.GetProfile(login)
	if err != nil {
		return `{"state":"error", "result":"auth error"}`, errors.New("auth error")
	}

	if profile.Status.IsBad() {
		return `{"state":"error", "result":"auth error"}`, errors.New("auth error")
	}

	secret := gtypes.Secret{
		Login:    login,
		Password: password,
		Hash:     hash,
	}
	ticket, err := gauth.NewAuth(&secret)
	if err != nil {
		return `{"state":"error", "result":"auth error"}`, errors.New("auth error")
	}

	return ecowriter.EncodeJSON(gtypes.Response{
		State:  "ok",
		Ticket: ticket,
	}), nil
}
