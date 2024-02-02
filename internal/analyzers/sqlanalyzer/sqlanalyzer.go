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

	// q.DCLSearchUse()
	// q.DCLSearchGrant()
	// q.DCLSearchRevoke()

	nameExps := []string{
		"SearchUse",
		"SearchGrant",
		"SearchRevoke",
	}

	var lastLen int = 0

	for {
		if curentLen := len(q.Instruction); curentLen <= 0 || curentLen == lastLen {
			break
		}
		lastLen = len(q.Instruction)
		q.SearchAllExp(nameExps)
	}

	return nil
}

func (q *tQuery) HeadCleaner() bool {
	// This method is completes

	re := core.RegExpCollection["HeadCleaner"]

	location := re.FindStringIndex(q.Instruction)
	if len(location) > 0 && location[0] == 0 {
		res := re.FindString(q.Instruction)
		q.Instruction = strings.Replace(q.Instruction, res, "", 1)
		return true
	}

	return false
}

func (q *tQuery) SearchExp(nameExp string) bool {
	// This method is complete
	re := core.RegExpCollection[nameExp]

	location := re.FindStringIndex(q.Instruction)
	if len(location) > 0 && location[0] == 0 {
		res := re.FindString(q.Instruction)
		q.QueryLine = append(q.QueryLine, res)
		q.Instruction = strings.Replace(q.Instruction, res, "", 1)
		q.HeadCleaner()
		return true
	}

	return false
}

func (q *tQuery) SearchAllExp(nameExps []string) (err error) {
	// This method is complete
	op := "internal -> analyzers -> sql -> SearchAllExp"
	defer func() { e.Wrapper(op, err) }()

	for _, v := range nameExps {
		if q.SearchExp(v) {
			break
		}
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

	query.HeadCleaner()

	query.Decomposition()

	res = fmt.Sprint(query.QueryLine) // FIXME: Temporarily for tests

	return &res
}
