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

// DDL — язык определения данных (Data Definition Language)

func (q tQuery) DDLCreate() (result string, err error) {
	// -
	op := "internal -> analyzers -> sql -> DDL -> DDLCreate"
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

	isDB := core.RegExpCollection["CreateDatabaseWord"].MatchString(q.Instruction)
	isTable := core.RegExpCollection["CreateTableWord"].MatchString(q.Instruction)
	isINE := core.RegExpCollection["IfNotExistsWord"].MatchString(q.Instruction)

	if isDB {
		db := core.RegExpCollection["CreateDatabaseWord"].ReplaceAllLiteralString(q.Instruction, "")
		if isINE {
			db = core.RegExpCollection["IfNotExistsWord"].ReplaceAllLiteralString(db, "")
		}
		db = strings.TrimSpace(db)
		db = core.RegExpCollection["QuotationMarks"].ReplaceAllLiteralString(db, "")
		db = core.RegExpCollection["SpecQuotationMark"].ReplaceAllLiteralString(db, "")
		if !core.CreateDB(db, login, true) {
			res.State = "error"
			res.Result = "invalid database name"
			return ecowriter.EncodeString(res), errors.New("invalid database name")
		}
	} else if isTable {
		table := core.RegExpCollection["CreateTableWord"].ReplaceAllLiteralString(q.Instruction, "")
		if isINE {
			table = core.RegExpCollection["IfNotExistsWord"].ReplaceAllLiteralString(table, "")
		}
		table = strings.TrimSpace(table)
		table = core.RegExpCollection["QuotationMarks"].ReplaceAllLiteralString(table, "")
		table = core.RegExpCollection["SpecQuotationMark"].ReplaceAllLiteralString(table, "")

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

		if !core.CreateTable(db, table, true) {
			res.State = "error"
			res.Result = "invalid database name"
			return ecowriter.EncodeString(res), errors.New("invalid database name")
		}
	}

	res.State = "ok"
	return ecowriter.EncodeString(res), nil
}

func (q tQuery) DDLAlter() (result string, err error) {
	// -
	op := "internal -> analyzers -> sql -> DDL -> DDLAlter"
	defer func() { e.Wrapper(op, err) }()

	return "DDLAlter", nil
}

func (q tQuery) DDLDrop() (result string, err error) {
	// -
	op := "internal -> analyzers -> sql -> DDL -> DDLDrop"
	defer func() { e.Wrapper(op, err) }()

	return "DDLDrop", nil
}
