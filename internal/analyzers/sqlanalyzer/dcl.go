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

	var (
		dbs   []string
		users []string
	)

	privilegesStr := core.RegExpCollection["GrantPrivileges"].FindString(q.Instruction)
	privilegesStr = core.RegExpCollection["GrantWord"].ReplaceAllLiteralString(privilegesStr, "")
	privilegesStr = core.RegExpCollection["ON"].ReplaceAllLiteralString(privilegesStr, "")
	privileges := core.RegExpCollection["GrantPrivilegesList"].FindAllString(privilegesStr, -1)

	if len(privileges) == 0 {
		return ecowriter.EncodeString(gtypes.Response{
			State:  "error",
			Result: "privileges are not specified",
		}), errors.New("privileges are not specified")
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
		return ecowriter.EncodeString(gtypes.Response{
			State:  "error",
			Result: "databases are not specified",
		}), errors.New("databases are not specified")
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
		return ecowriter.EncodeString(gtypes.Response{
			State:  "error",
			Result: "users are not specified",
		}), errors.New("users are not specified")
	}

	for _, db := range dbs {
		dbAccess, ok := core.StorageInfo.Access[db]
		if ok {
			if dbAccess.Owner != login {
				var luxUser bool = false
				for role := range access.Roles {
					if role == 1 || role == 3 {
						luxUser = true
						break
					}
				}
				if !luxUser {
					return ecowriter.EncodeString(gtypes.Response{
						State:  "error",
						Result: "auth error",
					}), errors.New("auth error")
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
	// -
	op := "internal -> analyzers -> sql -> DCL -> DCLRevoke"
	defer func() { e.Wrapper(op, err) }()

	return "DCLRevoke", nil
}

func (q tQuery) DCLUse() (result string, err error) {
	// This method is complete
	op := "internal -> analyzers -> sql -> DCL -> DCLUse"
	defer func() { e.Wrapper(op, err) }()

	var ticket string
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
		return ecowriter.EncodeString(gtypes.Response{
			State:  "error",
			Result: "invalid database name",
		}), errors.New("invalid database name")
	}

	if core.LocalCoreSettings.FreezeMode {
		if _, ok := core.StorageInfo.DBs[db]; !ok {
			return ecowriter.EncodeString(gtypes.Response{
				State:  "error",
				Result: "the database does not exist",
			}), errors.New("the database does not exist")
		}
	}

	dbAccess, ok := core.StorageInfo.Access[db]
	if ok {
		if dbAccess.Owner != login {
			var luxUser bool = false
			for role := range access.Roles {
				if role == 1 || role == 3 {
					luxUser = true
					break
				}
			}

			if !luxUser {
				flags, ok := dbAccess.Flags[login]
				if !ok {
					return ecowriter.EncodeString(gtypes.Response{
						State:  "error",
						Result: "auth error",
					}), errors.New("auth error")
				}
				if !flags.AnyTrue() {
					return ecowriter.EncodeString(gtypes.Response{
						State:  "error",
						Result: "auth error",
					}), errors.New("auth error")
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

func (q tQuery) DCLAuth() (result string, err error) {
	// This method is complete
	op := "internal -> analyzers -> sql -> DCL -> DCLAuth"
	defer func() { e.Wrapper(op, err) }()

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

	profile, err := gauth.GetProfile(login)
	if err != nil {
		return ecowriter.EncodeString(gtypes.Response{
			State:  "error",
			Result: "auth error",
		}), errors.New("auth error")
	}

	if profile.Status.IsBad() {
		return ecowriter.EncodeString(gtypes.Response{
			State:  "error",
			Result: "auth error",
		}), errors.New("auth error")
	}

	secret := gtypes.Secret{
		Login:    login,
		Password: password,
		Hash:     hash,
	}
	ticket, err := gauth.NewAuth(&secret)
	if err != nil {
		return ecowriter.EncodeString(gtypes.Response{
			State:  "error",
			Result: "auth error",
		}), errors.New("auth error")
	}

	return ecowriter.EncodeString(gtypes.Response{
		State:  "ok",
		Ticket: ticket,
	}), nil
}
