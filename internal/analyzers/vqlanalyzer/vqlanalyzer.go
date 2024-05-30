package vqlanalyzer

import (
	"strings"

	"github.com/Kwynto/GracefulDB/internal/engine/basicsystem/vqlexp"
)

type tQuery struct {
	Ticket      string
	Instruction string
	Placeholder []string
}

// TODO: Request
func Request(sTicket string, sInstruction string, slPlaceholder []string) string {
	// -
	// Prep
	sInstruction = vqlexp.MRegExpCollection["LineBreak"].ReplaceAllLiteralString(sInstruction, " ")
	sInstruction = strings.TrimRight(sInstruction, "; ")
	sInstruction = strings.TrimLeft(sInstruction, " ")

	var query tQuery = tQuery{
		Ticket:      sTicket,
		Instruction: sInstruction,
		Placeholder: slPlaceholder,
	}

	for _, sExpName := range vqlexp.ArParsingOrder {
		if vqlexp.MRegExpCollection[sExpName].MatchString(query.Instruction) {
			switch sExpName {
			case "SearchCreate":
				sResult, _ := query.DDLCreate()
				return sResult
			case "SearchAlter":
				sResult, _ := query.DDLAlter()
				return sResult
			case "SearchDrop":
				sResult, _ := query.DDLDrop()
				return sResult
			case "SearchSelect":
				sResult, _ := query.DMLSelect()
				return sResult
			case "SearchInsert":
				sResult, _ := query.DMLInsert()
				return sResult
			case "SearchUpdate":
				sResult, _ := query.DMLUpdate()
				return sResult
			case "SearchDelete":
				sResult, _ := query.DMLDelete()
				return sResult
			case "SearchTruncateTable":
				sResult, _ := query.DMLTruncateTable()
				return sResult
			case "SearchCommit":
				sResult, _ := query.DMLCommit()
				return sResult
			case "SearchRollback":
				sResult, _ := query.DMLRollback()
				return sResult
			case "SearchUse":
				sResult, _ := query.DCLUse()
				return sResult
			case "SearchShow":
				sResult, _ := query.DCLShow()
				return sResult
			case "SearchDesc", "SearchDescribe", "SearchExplain":
				sResult, _ := query.DCLDesc()
				return sResult
			case "SearchGrant":
				sResult, _ := query.DCLGrant()
				return sResult
			case "SearchRevoke":
				sResult, _ := query.DCLRevoke()
				return sResult
			case "SearchAuth":
				sResult, _ := query.DCLAuth()
				return sResult
			}
		}
	}

	sResult := "{\"state\":\"error\",\"result\":\"unknown command\"}"
	return sResult
}
