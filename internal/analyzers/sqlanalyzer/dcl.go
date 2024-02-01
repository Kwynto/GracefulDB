package sqlanalyzer

import (
	"strings"

	"github.com/Kwynto/GracefulDB/internal/engine/core"
	"github.com/Kwynto/GracefulDB/pkg/lib/e"
)

// DCL — язык управления данными (Data Control Language)

func (q *tQuery) DCLSearchGrant() (err error) {
	// -
	op := "internal -> analyzers -> sql -> DCL -> DCLSearchGrant"
	defer func() { e.Wrapper(op, err) }()

	return nil
}

func (q *tQuery) DCLSearchRevoke() (err error) {
	// -
	op := "internal -> analyzers -> sql -> DCL -> DCLSearchRevoke"
	defer func() { e.Wrapper(op, err) }()

	return nil
}

func (q *tQuery) DCLSearchUse() (err error) {
	// This method is complete
	op := "internal -> analyzers -> sql -> DCL -> DCLSearchUse"
	defer func() { e.Wrapper(op, err) }()

	re := core.RegExpCollection["SearchUse"]

	location := re.FindStringIndex(q.Instruction)
	if len(location) > 0 && location[0] == 0 {
		res := re.FindString(q.Instruction)
		q.QueryLine = append(q.QueryLine, res)
		q.Instruction = strings.Replace(q.Instruction, res, "", 1)
		q.HeadCleaner()
	}

	return nil
}
