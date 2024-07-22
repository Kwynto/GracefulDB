package vqlanalyzer

import (
	"fmt"
	"strings"

	"github.com/Kwynto/GracefulDB/internal/engine/basicsystem/vqlexp"
)

type tQuery struct {
	Ticket      string
	Instruction string
	Placeholder []string
}

func prepSpacesInLine(sSlIn []string) []string {
	// This functions is complete
	var slPrepLines []string
	for _, v := range sSlIn {
		v = strings.TrimRight(v, "; \t")
		v = strings.TrimLeft(v, " \t")
		slPrepLines = append(slPrepLines, v)
	}
	return slPrepLines
}

func prepPipelineInLine(sSlIn []string) []string {
	// This functions is complete
	var slPrepLines []string
	var isEndlySimbolPL bool = false

	for iQryInd, sLine := range sSlIn {
		rSlCurLine := []rune(sLine)
		rCurStartSimbol := rSlCurLine[0]
		rCurFinishSimbol := rSlCurLine[len(rSlCurLine)-1]

		if iQryInd == 0 {
			sLine = strings.TrimLeft(sLine, "| ")

			if string(rCurFinishSimbol) == "|" {
				sLine = strings.TrimRight(sLine, " |")
				sLine = fmt.Sprintf("%s|", sLine)
				isEndlySimbolPL = true
			} else {
				isEndlySimbolPL = false
			}

			slPrepLines = append(slPrepLines, sLine)
			continue
		}

		if isEndlySimbolPL {
			sLine = strings.TrimLeft(sLine, "| ")

			if string(rCurFinishSimbol) == "|" {
				sLine = strings.TrimRight(sLine, " |")
				sLine = fmt.Sprintf("%s|", sLine)
				isEndlySimbolPL = true
			} else {
				isEndlySimbolPL = false
			}

			tempLine := slPrepLines[len(slPrepLines)-1]
			slPrepLines = slPrepLines[:len(slPrepLines)-1]
			tempLine = fmt.Sprintf("%s%s", tempLine, sLine)
			slPrepLines = append(slPrepLines, tempLine)
			continue
		}

		if string(rCurStartSimbol) == "|" {
			sLine = strings.TrimLeft(sLine, "| ")
			sLine = fmt.Sprintf("|%s", sLine)

			if string(rCurFinishSimbol) == "|" {
				sLine = strings.TrimRight(sLine, " |")
				sLine = fmt.Sprintf("%s|", sLine)
				isEndlySimbolPL = true
			} else {
				isEndlySimbolPL = false
			}

			tempLine := slPrepLines[len(slPrepLines)-1]
			slPrepLines = slPrepLines[:len(slPrepLines)-1]
			tempLine = fmt.Sprintf("%s%s", tempLine, sLine)
			slPrepLines = append(slPrepLines, tempLine)
		} else {
			if string(rCurFinishSimbol) == "|" {
				sLine = strings.TrimRight(sLine, " |")
				sLine = fmt.Sprintf("%s%s", sLine, "|")
				isEndlySimbolPL = true
			} else {
				isEndlySimbolPL = false
			}

			slPrepLines = append(slPrepLines, sLine)
		}
	}

	return slPrepLines
}

// TODO: Request
func Request(sTicket string, sInstruction string, slPlaceholder []string) string {
	// -

	// Prep
	slQryLines := vqlexp.MRegExpCollection["LineBreak"].Split(sInstruction, -1)
	slQryLines = prepSpacesInLine(slQryLines)
	slQryLines = prepPipelineInLine(slQryLines)

	var query tQuery = tQuery{
		Ticket:      sTicket,
		Instruction: sInstruction,
		Placeholder: slPlaceholder,
	}

	_ = query
	_ = slQryLines

	sResult := "{\"state\":\"error\",\"result\":\"unknown command\"}"
	return sResult
}
