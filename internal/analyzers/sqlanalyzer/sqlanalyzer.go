package sqlanalyzer

import (
	"fmt"
	"strings"
	"unicode/utf8"

	"github.com/Kwynto/GracefulDB/internal/engine/core"
)

type tQuery struct {
	Instruction string
	Placeholder []string
	// QueryLine   []string
}

// func (q *tQuery) Decomposition() (err error) {
// 	// -
// 	op := "internal -> analyzers -> sql -> Decomposition"
// 	defer func() { e.Wrapper(op, err) }()

// 	var lastLen int = 0
// 	for {
// 		if curentLen := len(q.Instruction); curentLen <= 0 || curentLen == lastLen {
// 			break
// 		}
// 		lastLen = len(q.Instruction)
// 		q.SearchExp("AnyCommand")
// 	}

// 	return nil
// }

// func (q *tQuery) HeadCleaner() bool {
// 	// This method is completes
// 	re := core.RegExpCollection["HeadCleaner"]

// 	location := re.FindStringIndex(q.Instruction)
// 	if len(location) > 0 && location[0] == 0 {
// 		res := re.FindString(q.Instruction)
// 		q.Instruction = strings.Replace(q.Instruction, res, "", 1)
// 		return true
// 	}

// 	return false
// }

// func (q *tQuery) SearchExp(nameExp string) bool {
// 	// This method is complete
// 	re := core.RegExpCollection[nameExp]

// 	location := re.FindStringIndex(q.Instruction)
// 	if len(location) > 0 && location[0] == 0 {
// 		res := re.FindString(q.Instruction)
// 		q.QueryLine = append(q.QueryLine, res)
// 		q.Instruction = strings.Replace(q.Instruction, res, "", 1)
// 		q.HeadCleaner()
// 		return true
// 	}

// 	return false
// }

// func (q *tQuery) SearchAllExp(nameExps []string) (err error) {
// 	// This method is complete
// 	op := "internal -> analyzers -> sql -> SearchAllExp"
// 	defer func() { e.Wrapper(op, err) }()

// 	for _, v := range nameExps {
// 		if q.SearchExp(v) {
// 			break
// 		}
// 	}

// 	return nil
// }

// TODO: Request
func Request(instruction *string, placeholder *[]string) *string {
	// -
	var res string

	// Prep
	inst := *instruction
	// inst = core.RegExpCollection["LineBreak"].ReplaceAllLiteralString(inst, " ")
	inst = strings.TrimSpace(inst)
	if r, _ := utf8.DecodeLastRuneInString(inst); r != ';' {
		inst = fmt.Sprintf("%s;", inst)
	}

	var query tQuery = tQuery{
		Instruction: inst,
		Placeholder: *placeholder,
	}

	for _, expName := range core.ParsingOrder {
		re := core.RegExpCollection[expName]

		location := re.FindStringIndex(query.Instruction)
		if len(location) == 1 && location[0] == 0 {
			switch expName {
			case "SearchCreate":
				res, _ = query.DDLCreate()
			case "SearchAlter":
				res, _ = query.DDLAlter()
			case "SearchDrop":
				res, _ = query.DDLDrop()
			case "SearchSelect":
				res, _ = query.DMLSelect()
			case "SearchInsert":
				res, _ = query.DMLInsert()
			case "SearchUpdate":
				res, _ = query.DMLUpdate()
			case "SearchDelete":
				res, _ = query.DMLDelete()
			case "SearchTruncate":
				res, _ = query.DMLTruncate()
			case "SearchCommit":
				res, _ = query.DMLCommit()
			case "SearchRollback":
				res, _ = query.DMLRollback()
			case "SearchUse":
				res, _ = query.DCLUse()
			case "SearchGrant":
				res, _ = query.DCLGrant()
			case "SearchRevoke":
				res, _ = query.DCLRevoke()
			}
		}
	}

	return &res
}
