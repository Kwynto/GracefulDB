package sqlanalyzer

import (
	"fmt"
	"strings"

	"github.com/Kwynto/GracefulDB/internal/engine/core"
	"github.com/Kwynto/GracefulDB/pkg/lib/e"
)

type tQuery struct {
	Instruction string
	Placeholder []string
	QueryLine   []string
}

func (q *tQuery) Decomposition() (err error) {
	// -
	op := "internal -> analyzers -> sql -> Decomposition"
	defer func() { e.Wrapper(op, err) }()

	q.DCLSearchUse()

	return nil
}

func (q *tQuery) HeadCleaner() (err error) {
	// -
	op := "internal -> analyzers -> sql -> HeadCleaner"
	defer func() { e.Wrapper(op, err) }()

	re := core.RegExpCollection["HeadCleaner"]

	location := re.FindStringIndex(q.Instruction)
	if len(location) > 0 && location[0] == 0 {
		res := re.FindString(q.Instruction)
		q.Instruction = strings.Replace(q.Instruction, res, "", 1)
	}

	return nil
}

// TODO: Request
func Request(instruction *string, placeholder *[]string) *string {
	// -
	var res string

	var query tQuery = tQuery{
		Instruction: *instruction,
		Placeholder: *placeholder,
		QueryLine:   make([]string, 1),
	}

	query.Decomposition()

	res = fmt.Sprint(query.QueryLine) // FIXME: Temporarily for tests

	return &res
}
