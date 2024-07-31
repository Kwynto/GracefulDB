package vqlanalyzer

import (
	"fmt"
	"strings"

	"github.com/Kwynto/GracefulDB/internal/engine/basicsystem/gtypes"
	"github.com/Kwynto/GracefulDB/internal/engine/basicsystem/vqlexp"
)

// Directives and reserved words

func parseWhere(sWhere string) ([]gtypes.TConditions, bool) {
	// This functions is complete
	var slExpression = make([]gtypes.TConditions, 0, 4)

	sWhere = vqlexp.MRegExpCollection["WhereWord"].ReplaceAllLiteralString(sWhere, "")

	for {
		sHeadCond := vqlexp.MRegExpCollection["WhereExpression"].ReplaceAllLiteralString(sWhere, "")
		slCondition := vqlexp.MRegExpCollection["WhereOperationConditions"].Split(sHeadCond, -1)
		sKeyIn := slCondition[0]
		sValueIn := slCondition[1]

		sKeyIn = vqlexp.MRegExpCollection["Spaces"].ReplaceAllLiteralString(sKeyIn, "")
		sKeyIn = trimQuotationMarks(sKeyIn)

		sValueIn = strings.TrimSpace(sValueIn)
		sValueIn = trimQuotationMarks(sValueIn)

		if sKeyIn == "" {
			return []gtypes.TConditions{}, false
		}
		if sValueIn == "" {
			return []gtypes.TConditions{}, false
		} // null value, maybe delete a condition

		stExp := gtypes.TConditions{
			Type:  "operation",
			Key:   sKeyIn,
			Value: sValueIn,
		}

		if vqlexp.MRegExpCollection["WhereOperation_<="].MatchString(sHeadCond) {
			stExp.Operation = "<="
		} else if vqlexp.MRegExpCollection["WhereOperation_>="].MatchString(sHeadCond) {
			stExp.Operation = ">="
		} else if vqlexp.MRegExpCollection["WhereOperation_<"].MatchString(sHeadCond) {
			stExp.Operation = "<"
		} else if vqlexp.MRegExpCollection["WhereOperation_>"].MatchString(sHeadCond) {
			stExp.Operation = ">"
		} else if vqlexp.MRegExpCollection["WhereOperation_=="].MatchString(sHeadCond) {
			stExp.Operation = "=="
		} else if vqlexp.MRegExpCollection["WhereOperation_LIKE"].MatchString(sHeadCond) {
			stExp.Operation = "like"
		} else if vqlexp.MRegExpCollection["WhereOperation_REGEXP"].MatchString(sHeadCond) {
			stExp.Operation = "regexp"
		} else {
			return []gtypes.TConditions{}, false
		}
		slExpression = append(slExpression, stExp)

		sWhere = vqlexp.MRegExpCollection["WhereExpression"].FindString(sWhere)
		sLogicOper := vqlexp.MRegExpCollection["WhereExpression_And_Or_Word"].FindString(sWhere)

		if vqlexp.MRegExpCollection["OR"].MatchString(sLogicOper) {
			slExpression = append(slExpression, gtypes.TConditions{
				Type: "or",
			})
		} else if vqlexp.MRegExpCollection["AND"].MatchString(sLogicOper) {
			slExpression = append(slExpression, gtypes.TConditions{
				Type: "and",
			})
		} else {
			break
		}

		sWhere = vqlexp.MRegExpCollection["WhereExpression_And_Or_Word"].ReplaceAllLiteralString(sWhere, "")
	}
	return slExpression, true
}

func (q tQuery) DirectWhere(lineInd int) (result string, ok bool) {
	// -
	sLine := q.QueryCode[lineInd]

	sLeft := vqlexp.MRegExpCollection["WhereRight"].ReplaceAllLiteralString(sLine, "")
	sLeft = strings.TrimRight(sLeft, "= ")
	if sLeft == "" {
		sLeft = "$result"
	}
	if !vqlexp.MRegExpCollection["VariableWholeString"].MatchString(sLeft) {
		return "", false
	}

	sRight := vqlexp.MRegExpCollection["WhereRight"].FindAllString(sLine, -1)[0]
	// TODO: orderby and limit
	stExpression, ok := parseWhere(sRight)
	if !ok {
		return "", false
	}

	_ = stExpression // FIXME: --

	sOut := fmt.Sprintf("{\"%s\": %s}", sLeft, result)
	return sOut, true
}
