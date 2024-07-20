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

	_ = query

	sResult := "{\"state\":\"error\",\"result\":\"unknown command\"}"
	return sResult
}
