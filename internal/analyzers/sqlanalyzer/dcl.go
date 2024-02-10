package sqlanalyzer

import (
	"fmt"
	"strings"

	"github.com/Kwynto/GracefulDB/internal/engine/core"
	"github.com/Kwynto/GracefulDB/pkg/lib/e"
)

// DCL — язык управления данными (Data Control Language)

func (q *tQuery) DCLGrant() (result string, err error) {
	// -
	op := "internal -> analyzers -> sql -> DCL -> DCLGrant"
	defer func() { e.Wrapper(op, err) }()

	return "DCLGrant", nil
}

func (q *tQuery) DCLRevoke() (result string, err error) {
	// -
	op := "internal -> analyzers -> sql -> DCL -> DCLRevoke"
	defer func() { e.Wrapper(op, err) }()

	return "DCLRevoke", nil
}

func (q *tQuery) DCLUse() (result string, err error) {
	// -
	op := "internal -> analyzers -> sql -> DCL -> DCLUse"
	defer func() { e.Wrapper(op, err) }()

	return "DCLUse", nil
}

func (q *tQuery) DCLAuth() (result string, err error) {
	// -
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

	fmt.Println(login)
	fmt.Println(password)

	return login, nil
}
