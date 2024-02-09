package sqlanalyzer

import (
	"fmt"
	"strings"
	"unicode/utf8"

	"github.com/Kwynto/GracefulDB/internal/engine/core"
)

type tQuery struct {
	Ticket      string
	Instruction string
	Placeholder []string
}

// TODO: Request
func Request(ticket *string, instruction *string, placeholder *[]string) *string {
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
		Ticket:      *ticket,
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
			case "SearchAuth":
				res, _ = query.DCLAuth()
			}
		}
	}

	return &res
}
