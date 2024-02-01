package sqlanalyzer

import (
	"fmt"
	"strings"
	"unicode/utf8"

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
	// This method is completes
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

func TailSign(inst string) string {
	// This function is complete
	r, _ := utf8.DecodeLastRuneInString(inst)
	if r != ';' {
		return fmt.Sprintf("%s;", inst)
	}
	return inst
}

// TODO: Request
func Request(instruction *string, placeholder *[]string) *string {
	// -
	var res string

	inst := *instruction
	inst = strings.TrimSpace(inst)
	inst = TailSign(inst)

	var query tQuery = tQuery{
		Instruction: inst,
		Placeholder: *placeholder,
		QueryLine:   make([]string, 5),
	}

	query.Decomposition()

	res = fmt.Sprint(query.QueryLine) // FIXME: Temporarily for tests

	return &res
}
