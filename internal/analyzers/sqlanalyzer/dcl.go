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

// DCL — язык управления данными (Data Control Language)

func (q tQuery) DCLGrant() (result string, err error) {
	// This method is complete
	op := "internal -> analyzers -> sql -> DCL -> DCLGrant"
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

	var (
		dbs   []string
		users []string
	)

	privilegesStr := core.RegExpCollection["GrantPrivileges"].FindString(q.Instruction)
	privilegesStr = core.RegExpCollection["GrantWord"].ReplaceAllLiteralString(privilegesStr, "")
	privilegesStr = core.RegExpCollection["ON"].ReplaceAllLiteralString(privilegesStr, "")
	privileges := core.RegExpCollection["GrantPrivilegesList"].FindAllString(privilegesStr, -1)

	if len(privileges) == 0 {
		return `{"state":"error", "result":"privileges are not specified"}`, errors.New("privileges are not specified")
	}

	dbsStr := core.RegExpCollection["GrantOnTo"].FindString(q.Instruction)
	dbsStr = core.RegExpCollection["ON"].ReplaceAllLiteralString(dbsStr, "")
	dbsStr = core.RegExpCollection["TO"].ReplaceAllLiteralString(dbsStr, "")
	dbsStr = core.RegExpCollection["Spaces"].ReplaceAllLiteralString(dbsStr, "")
	dbsStr = core.RegExpCollection["QuotationMarks"].ReplaceAllLiteralString(dbsStr, "")
	dbsStr = core.RegExpCollection["SpecQuotationMark"].ReplaceAllLiteralString(dbsStr, "")
	dbsIn := core.RegExpCollection["Comma"].Split(dbsStr, -1)
	for _, db := range dbsIn {
		if _, ok := core.StorageInfo.DBs[db]; ok {
			dbs = append(dbs, db)
		}
	}
	if len(dbs) == 0 {
		return `{"state":"error", "result":"databases are not specified"}`, errors.New("databases are not specified")
	}

	usersStr := core.RegExpCollection["GrantToEnd"].FindString(q.Instruction)
	usersStr = core.RegExpCollection["TO"].ReplaceAllLiteralString(usersStr, "")
	usersStr = core.RegExpCollection["Spaces"].ReplaceAllLiteralString(usersStr, "")
	usersStr = core.RegExpCollection["QuotationMarks"].ReplaceAllLiteralString(usersStr, "")
	usersStr = core.RegExpCollection["SpecQuotationMark"].ReplaceAllLiteralString(usersStr, "")
	usersIn := core.RegExpCollection["Comma"].Split(usersStr, -1)
	for _, user := range usersIn {
		if _, err := gauth.GetProfile(user); err == nil {
			users = append(users, user)
		}
	}
	if len(users) == 0 {
		return `{"state":"error", "result":"users are not specified"}`, errors.New("users are not specified")
	}

	for _, db := range dbs {
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
					return `{"state":"error", "result":"auth error"}`, errors.New("auth error")
				}
			}
			for _, user := range users {
				var aFlags gtypes.TAccessFlags
				aFlags, ok := core.StorageInfo.Access[db].Flags[user]
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

				core.StorageInfo.Access[db].Flags[user] = aFlags
			}
		}
	}

	core.StorageInfo.Save()

	res.State = "ok"
	return ecowriter.EncodeString(res), nil
}

func (q tQuery) DCLRevoke() (result string, err error) {
	// This method is complete
	op := "internal -> analyzers -> sql -> DCL -> DCLRevoke"
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

	var (
		dbs   []string
		users []string
	)

	privilegesStr := core.RegExpCollection["RevokePrivileges"].FindString(q.Instruction)
	privilegesStr = core.RegExpCollection["RevokeWord"].ReplaceAllLiteralString(privilegesStr, "")
	privilegesStr = core.RegExpCollection["ON"].ReplaceAllLiteralString(privilegesStr, "")
	privileges := core.RegExpCollection["RevokePrivilegesList"].FindAllString(privilegesStr, -1)

	if len(privileges) == 0 {
		return `{"state":"error", "result":"privileges are not specified"}`, errors.New("privileges are not specified")
	}

	dbsStr := core.RegExpCollection["RevokeOnTo"].FindString(q.Instruction)
	dbsStr = core.RegExpCollection["ON"].ReplaceAllLiteralString(dbsStr, "")
	dbsStr = core.RegExpCollection["TO"].ReplaceAllLiteralString(dbsStr, "")
	dbsStr = core.RegExpCollection["Spaces"].ReplaceAllLiteralString(dbsStr, "")
	dbsStr = core.RegExpCollection["QuotationMarks"].ReplaceAllLiteralString(dbsStr, "")
	dbsStr = core.RegExpCollection["SpecQuotationMark"].ReplaceAllLiteralString(dbsStr, "")
	dbsIn := core.RegExpCollection["Comma"].Split(dbsStr, -1)
	for _, db := range dbsIn {
		if _, ok := core.StorageInfo.DBs[db]; ok {
			dbs = append(dbs, db)
		}
	}
	if len(dbs) == 0 {
		return `{"state":"error", "result":"databases are not specified"}`, errors.New("databases are not specified")
	}

	usersStr := core.RegExpCollection["RevokeToEnd"].FindString(q.Instruction)
	usersStr = core.RegExpCollection["TO"].ReplaceAllLiteralString(usersStr, "")
	usersStr = core.RegExpCollection["Spaces"].ReplaceAllLiteralString(usersStr, "")
	usersStr = core.RegExpCollection["QuotationMarks"].ReplaceAllLiteralString(usersStr, "")
	usersStr = core.RegExpCollection["SpecQuotationMark"].ReplaceAllLiteralString(usersStr, "")
	usersIn := core.RegExpCollection["Comma"].Split(usersStr, -1)
	for _, user := range usersIn {
		if _, err := gauth.GetProfile(user); err == nil {
			users = append(users, user)
		}
	}
	if len(users) == 0 {
		return `{"state":"error", "result":"users are not specified"}`, errors.New("users are not specified")
	}

	for _, db := range dbs {
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
					return `{"state":"error", "result":"auth error"}`, errors.New("auth error")
				}
			}
			for _, user := range users {
				var aFlags gtypes.TAccessFlags
				aFlags, ok := core.StorageInfo.Access[db].Flags[user]
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

				core.StorageInfo.Access[db].Flags[user] = aFlags
			}
		}
	}

	core.StorageInfo.Save()

	res.State = "ok"
	return ecowriter.EncodeString(res), nil
}

func (q tQuery) DCLUse() (result string, err error) {
	// This method is complete
	op := "internal -> analyzers -> sql -> DCL -> DCLUse"
	defer func() { e.Wrapper(op, err) }()

	var ticket string
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
		ticket = newticket
		res.Ticket = newticket
	} else {
		ticket = q.Ticket
	}

	db := core.RegExpCollection["UseWord"].ReplaceAllLiteralString(q.Instruction, "")
	db = strings.TrimSpace(db)
	db = core.RegExpCollection["QuotationMarks"].ReplaceAllLiteralString(db, "")
	db = core.RegExpCollection["SpecQuotationMark"].ReplaceAllLiteralString(db, "")

	if !core.RegExpCollection["EntityName"].MatchString(db) {
		return `{"state":"error", "result":"invalid database name"}`, errors.New("invalid database name")
	}

	if core.LocalCoreSettings.FreezeMode {
		if _, ok := core.StorageInfo.DBs[db]; !ok {
			return `{"state":"error", "result":"the database does not exist"}`, errors.New("the database does not exist")
		}
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

	core.States[ticket] = core.TState{
		CurrentDB: db,
	}

	res.State = "ok"
	res.Result = db
	return ecowriter.EncodeString(res), nil
}

func (q tQuery) DCLShow() (result string, err error) {
	// This method is complete
	op := "internal -> analyzers -> sql -> DCL -> DCLShow"
	defer func() { e.Wrapper(op, err) }()

	var res gtypes.Response
	var resArr gtypes.ResponseArray

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

	isDBs := core.RegExpCollection["ShowDatabasesWord"].MatchString(q.Instruction)
	isTables := core.RegExpCollection["ShowTablesWord"].MatchString(q.Instruction)

	if isDBs {
		var namesDBs []string = []string{}
		for nameDB := range core.StorageInfo.DBs {
			namesDBs = append(namesDBs, nameDB)
		}

		resArr.State = "ok"
		resArr.Ticket = res.Ticket
		resArr.Result = namesDBs
		return ecowriter.EncodeString(resArr), nil
	} else if isTables {
		var namesTables []string = []string{}

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

		dbInfo, ok := core.StorageInfo.DBs[db]
		if !ok {
			res.State = "error"
			res.Result = "incorrect database"
			return ecowriter.EncodeString(res), errors.New("incorrect database")
		}

		for nameTable := range dbInfo.Tables {
			namesTables = append(namesTables, nameTable)
		}

		resArr.State = "ok"
		resArr.Ticket = res.Ticket
		resArr.Result = namesTables
		return ecowriter.EncodeString(resArr), nil
	}

	res.State = "error"
	res.Result = "unknown command"
	return ecowriter.EncodeString(res), nil
}

func (q tQuery) DCLAuth() (result string, err error) {
	// This method is complete
	op := "internal -> analyzers -> sql -> DCL -> DCLAuth"
	defer func() { e.Wrapper(op, err) }()

	var roles []gauth.TRole

	isNew := core.RegExpCollection["AuthNew"].MatchString(q.Instruction)
	isChange := core.RegExpCollection["AuthChange"].MatchString(q.Instruction)
	isRemove := core.RegExpCollection["AuthRemove"].MatchString(q.Instruction)

	login := core.RegExpCollection["Login"].FindString(q.Instruction)
	login = core.RegExpCollection["LoginWord"].ReplaceAllLiteralString(login, " ")
	login = strings.TrimSpace(login)
	login = core.RegExpCollection["QuotationMarks"].ReplaceAllLiteralString(login, "")
	login = core.RegExpCollection["SpecQuotationMark"].ReplaceAllLiteralString(login, "")

	password := core.RegExpCollection["Password"].FindString(q.Instruction)
	password = core.RegExpCollection["PasswordWord"].ReplaceAllLiteralString(password, " ")
	password = strings.TrimSpace(password)
	password = core.RegExpCollection["QuotationMarks"].ReplaceAllLiteralString(password, "")
	password = core.RegExpCollection["SpecQuotationMark"].ReplaceAllLiteralString(password, "")

	hash := core.RegExpCollection["Hash"].FindString(q.Instruction)
	hash = core.RegExpCollection["HashWord"].ReplaceAllLiteralString(hash, " ")
	hash = strings.TrimSpace(hash)
	hash = core.RegExpCollection["QuotationMarks"].ReplaceAllLiteralString(hash, "")
	hash = core.RegExpCollection["SpecQuotationMark"].ReplaceAllLiteralString(hash, "")

	isRole := core.RegExpCollection["Role"].MatchString(q.Instruction)
	if isRole {
		roleStr := core.RegExpCollection["Role"].FindString(q.Instruction)
		roleStr = core.RegExpCollection["RoleWord"].ReplaceAllLiteralString(roleStr, "")
		roleStr = core.RegExpCollection["Spaces"].ReplaceAllLiteralString(roleStr, "")
		roleStr = core.RegExpCollection["QuotationMarks"].ReplaceAllLiteralString(roleStr, "")
		roleStr = core.RegExpCollection["SpecQuotationMark"].ReplaceAllLiteralString(roleStr, "")
		roleIn := core.RegExpCollection["Comma"].Split(roleStr, -1)
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
			return ecowriter.EncodeString(res), errors.New("auth error")
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
				return ecowriter.EncodeString(gtypes.Response{
					State:  "error",
					Result: err.Error(),
				}), err
			}
		}

		if isChange {
			access, err := gauth.GetProfile(login)
			if err != nil {
				return ecowriter.EncodeString(gtypes.Response{
					State:  "error",
					Result: err.Error(),
				}), err
			}

			if isRole {
				access.Roles = roles
			}

			err = gauth.UpdateUser(login, password, access)
			if err != nil {
				return ecowriter.EncodeString(gtypes.Response{
					State:  "error",
					Result: err.Error(),
				}), err
			}
		}

		if isRemove {
			err := gauth.DeleteUser(login)
			if err != nil {
				return ecowriter.EncodeString(gtypes.Response{
					State:  "error",
					Result: err.Error(),
				}), err
			}
		}

		res.State = "ok"
		return ecowriter.EncodeString(res), nil
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

	return ecowriter.EncodeString(gtypes.Response{
		State:  "ok",
		Ticket: ticket,
	}), nil
}
