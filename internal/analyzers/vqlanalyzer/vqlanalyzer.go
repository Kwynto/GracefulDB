package vqlanalyzer

import (
	"fmt"
	"strings"

	"github.com/Kwynto/GracefulDB/internal/engine/basicsystem/gauth"
	"github.com/Kwynto/GracefulDB/internal/engine/basicsystem/gtypes"
	"github.com/Kwynto/GracefulDB/internal/engine/basicsystem/sqlexp"
	"github.com/Kwynto/GracefulDB/internal/engine/basicsystem/vqlexp"
	"github.com/Kwynto/GracefulDB/pkg/lib/ecowriter"
)

type tInVariables struct {
	Name string
	Type string
}

type tStFuncCode struct {
	Name         string
	InVariables  []tInVariables
	OutVariables []string
	Lines        []string
}

type tQuery struct {
	Login          string
	Access         gauth.TProfile
	Ticket         string
	QueryCode      []string
	LocalFunctions map[string]tStFuncCode
	Variables      map[string]any
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

func prepareLocalFunctions(sSlIn []string) ([]string, map[string]tStFuncCode, error) {
	// This functions is complete
	var slPrepLines []string
	var sNameFunc string = ""
	mFuncsCode := make(map[string]tStFuncCode)
	countCodeBlocks := 0

	iLenIn := len(sSlIn)
	for i := 0; i < iLenIn; i++ {
		sLine := sSlIn[i]

		if vqlexp.MRegExpCollection["FuncSignature"].MatchString(sLine) {
			if countCodeBlocks > 0 {
				return slPrepLines, mFuncsCode, fmt.Errorf("sintax error in \"%s\"", sLine)
			}

			var stFuncCode tStFuncCode

			sNameFunc = sqlexp.MRegExpCollection["FuncWord"].ReplaceAllLiteralString(sLine, "")
			sNameFunc = sqlexp.MRegExpCollection["FuncDesc"].ReplaceAllLiteralString(sNameFunc, "")

			sLine = sqlexp.MRegExpCollection["BeginBlock"].ReplaceAllLiteralString(sLine, "")
			sLine = sqlexp.MRegExpCollection["FuncWordAndName"].ReplaceAllLiteralString(sLine, "")

			sInVars := sqlexp.MRegExpCollection["FuncInVarString"].FindAllString(sLine, -1)[0]
			sInVars = strings.TrimLeft(sInVars, " (")
			sInVars = strings.TrimRight(sInVars, ") ")
			slSInVars := strings.Split(sInVars, ",")

			var slInVariables []tInVariables
			for _, v := range slSInVars {
				v = strings.TrimSpace(v)
				slV := sqlexp.MRegExpCollection["Spaces"].Split(v, -1)

				slRNameVar := []rune(slV[0])
				if slRNameVar[0] != rune('$') {
					return slPrepLines, mFuncsCode, fmt.Errorf("sintax error in functione's declaration")
				}
				stVar := tInVariables{
					Name: slV[0],
					Type: slV[1],
				}
				slInVariables = append(slInVariables, stVar)
			}

			sOutVars := sqlexp.MRegExpCollection["FuncInVarString"].ReplaceAllLiteralString(sLine, "")
			sOutVars = strings.TrimLeft(sOutVars, " (")
			sOutVars = strings.TrimRight(sOutVars, ") ")
			slSOutVars := strings.Split(sOutVars, ",")

			var slOutVariables []string
			for _, v := range slSOutVars {
				v = strings.TrimSpace(v)
				slOutVariables = append(slOutVariables, v)
			}

			stFuncCode.Name = sNameFunc
			stFuncCode.InVariables = slInVariables
			stFuncCode.OutVariables = slOutVariables

			mFuncsCode[sNameFunc] = stFuncCode
			countCodeBlocks = 1
			continue
		}

		if countCodeBlocks == 0 {
			slPrepLines = append(slPrepLines, sLine)
			continue
		}

		if sNameFunc != "" {
			stFuncCode, ok := mFuncsCode[sNameFunc]
			if ok {
				stFuncCode.Lines = append(stFuncCode.Lines, sLine)

				if vqlexp.MRegExpCollection["BeginBlock"].MatchString(sLine) {
					countCodeBlocks += 1
				} else if vqlexp.MRegExpCollection["EndBlock"].MatchString(sLine) {
					countCodeBlocks -= 1
					if countCodeBlocks == 0 {
						sNameFunc = ""
					}
				}
			}
		}

		if countCodeBlocks > 0 && i == (iLenIn-1) {
			return slPrepLines, mFuncsCode, fmt.Errorf("sintax error")
		}
	}

	return slPrepLines, mFuncsCode, nil
}

func preparation(sIn string) ([]string, map[string]tStFuncCode, error) {
	// This functions is complete
	slPrepLines := vqlexp.MRegExpCollection["LineBreak"].Split(sIn, -1)
	slPrepLines = prepareSpacesInLine(slPrepLines)
	slPrepLines = prepareRemoveComments(slPrepLines)
	slPrepLines = preparePipelineInLine(slPrepLines)
	slPrepLines, mLocalFunctions, err := prepareLocalFunctions(slPrepLines)
	return slPrepLines, mLocalFunctions, err
}

func execution(query tQuery) (gtypes.TResponse, error) {
	// -

	for lineInd, sLine := range query.QueryCode {
		for _, sExpName := range vqlexp.ArParsingOrder {
			if vqlexp.MRegExpCollection[sExpName].MatchString(sLine) {
				switch sExpName {
				case "Where":
					query.DirectWhere(lineInd)
				}
			}
		}
	}

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
	slQryLines, mLocalFunctions, errP := preparation(sOriginalCode)
	if errP != nil {
		stRes.State = "error"
		stRes.Result = errP.Error()
		return ecowriter.EncodeJSON(stRes)
	}

	mVariables, errU := ecowriter.DecodeJSONMap(sVariables)
	if errU != nil {
		stRes.State = "error"
		stRes.Result = errU.Error()
		return ecowriter.EncodeJSON(stRes)
	}

	var query tQuery = tQuery{
		Login:          sLogin,
		Access:         stAccess,
		Ticket:         sTicket,
		QueryCode:      slQryLines,
		LocalFunctions: mLocalFunctions,
		Variables:      mVariables,
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
