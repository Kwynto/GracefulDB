package vqlanalyzer

import (
	"fmt"
	"strings"

	"github.com/Kwynto/GracefulDB/internal/engine/basicsystem/gauth"
	"github.com/Kwynto/GracefulDB/internal/engine/basicsystem/gtypes"
	"github.com/Kwynto/GracefulDB/internal/engine/basicsystem/vqlexp"
	"github.com/Kwynto/GracefulDB/pkg/lib/ecowriter"
)

type tQuery struct {
	Login     string
	Access    gauth.TProfile
	Ticket    string
	QueryCode []string
	Variables map[string]any
}

func prepareSpacesInLine(sSlIn []string) []string {
	// This functions is complete
	var slPrepLines []string
	for _, sLine := range sSlIn {
		sLine = strings.TrimRight(sLine, "; \t")
		sLine = strings.TrimLeft(sLine, " \t")
		slPrepLines = append(slPrepLines, sLine)
	}
	return slPrepLines
}

func preparePipelineInLine(sSlIn []string) []string {
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

func prepareRemoveComments(sSlIn []string) []string {
	// This functions is complete
	var slPrepLines []string
	for _, sLine := range sSlIn {
		if !vqlexp.MRegExpCollection["Comment"].MatchString(sLine) {
			slPrepLines = append(slPrepLines, sLine)
		}
	}
	return slPrepLines
}

func preparation(sIn string) []string {
	// This functions is complete
	slPrepLines := vqlexp.MRegExpCollection["LineBreak"].Split(sIn, -1)
	slPrepLines = prepareSpacesInLine(slPrepLines)
	slPrepLines = prepareRemoveComments(slPrepLines)
	slPrepLines = preparePipelineInLine(slPrepLines)
	return slPrepLines
}

func execution(query tQuery) (gtypes.TResponse, error) {
	// -

	_ = query

	return gtypes.TResponse{}, nil
}

// TODO: Request
func Request(sTicket string, sOriginalCode string, sVariables string) string {
	// -
	var stRes gtypes.TResponse
	// mVariables := make(map[string]any)

	// Pre checking
	sLogin, stAccess, sNewTicket, errC := preCheckerVQL(sTicket)
	if errC != nil {
		stRes.State = "error"
		stRes.Result = errC.Error()
		return ecowriter.EncodeJSON(stRes)
	}

	if sNewTicket != "" {
		stRes.Ticket = sNewTicket
		sTicket = sNewTicket
	}

	// Preparation query
	slQryLines := preparation(sOriginalCode)

	mVariables, errU := ecowriter.DecodeJSONMap(sVariables)
	if errU != nil {
		stRes.State = "error"
		stRes.Result = errU.Error()
		return ecowriter.EncodeJSON(stRes)
	}

	var query tQuery = tQuery{
		Login:     sLogin,
		Access:    stAccess,
		Ticket:    sTicket,
		QueryCode: slQryLines,
		Variables: mVariables,
	}

	// Execution query
	stRes, errEx := execution(query)
	if errEx != nil {
		stRes.State = "error"
		stRes.Result = errEx.Error()
		return ecowriter.EncodeJSON(stRes)
	}

	return ecowriter.EncodeJSON(stRes)
}
