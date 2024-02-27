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

// DML — язык изменения данных (Data Manipulation Language)

func (q tQuery) DMLSelect() (result string, err error) {
	// -
	op := "internal -> analyzers -> sql -> DML -> DMLSelect"
	defer func() { e.Wrapper(op, err) }()

	return "DMLSelect", nil
}

func (q tQuery) DMLInsert() (result string, err error) {
	// -
	op := "internal -> analyzers -> sql -> DML -> DMLInsert"
	defer func() { e.Wrapper(op, err) }()

	return "DMLInsert", nil
}

func (q tQuery) DMLUpdate() (result string, err error) {
	// -
	op := "internal -> analyzers -> sql -> DML -> DMLUpdate"
	defer func() { e.Wrapper(op, err) }()

	return "DMLUpdate", nil
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

	_, okDB := core.GetDBInfo(db)
	if okDB {
		var flagsAcs gtypes.TAccessFlags
		var okFlags bool = false
		var luxUser bool = false

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
			res.State = "error"
			res.Result = "internal error"
			return ecowriter.EncodeJSON(res), errors.New("internal error")
		}

		if !luxUser && !flagsAcs.Delete {
			return `{"state":"error", "result":"not enough rights"}`, errors.New("not enough rights")
		}

		if !core.TruncateTable(db, table) {
			return `{"state":"error", "result":"the table cannot be truncated"}`, errors.New("the table cannot be truncated")
		}
	} else {
		res.State = "error"
		res.Result = "internal error"
		return ecowriter.EncodeJSON(res), errors.New("internal error")
	}

	res.State = "ok"
	return ecowriter.EncodeJSON(res), nil
}
