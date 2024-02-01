package sqlanalyzer

import (
	"regexp"
	"strings"

	"github.com/Kwynto/GracefulDB/pkg/lib/e"
)

// DCL — язык управления данными (Data Control Language)

func (q *tQuery) DCLGrant() (err error) {
	// -
	op := "internal -> analyzers -> sql -> DCL -> DCLGrant"
	defer func() { e.Wrapper(op, err) }()

	return nil
}

func (q *tQuery) DCLRevoke() (err error) {
	// -
	op := "internal -> analyzers -> sql -> DCL -> DCLRevoke"
	defer func() { e.Wrapper(op, err) }()

	return nil
}

func (q *tQuery) DCLUse() (err error) {
	// -
	op := "internal -> analyzers -> sql -> DCL -> DCLUse"
	defer func() { e.Wrapper(op, err) }()

	re, err := regexp.Compile(`(?m)^[uU][sS][eE]\s*[a-zA-Z][a-zA-Z0-1]+;`)
	if err != nil {
		return err
	}

	location := re.FindStringIndex(q.Instruction)
	if len(location) > 0 && location[0] == 0 {
		res := re.FindString(q.Instruction)
		q.QueryLine = append(q.QueryLine, res)
		q.Instruction = strings.Replace(q.Instruction, res, "", 1)
		q.HeadCleaner()
	}

	return nil
}
