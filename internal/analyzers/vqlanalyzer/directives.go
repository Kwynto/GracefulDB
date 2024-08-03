package vqlanalyzer

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/Kwynto/GracefulDB/internal/engine/basicsystem/gtypes"
	"github.com/Kwynto/GracefulDB/internal/engine/basicsystem/vqlexp"
)

// Directives and reserved words

func parseLimit(sLimit string) gtypes.TLimit {
	// -
	var iStart, iOffset int
	var err error
	var slLimitAfter []string

	slLimitBefore := vqlexp.MRegExpCollection["Comma"].Split(sLimit, -1)
	for _, sValue := range slLimitBefore {
		sValue = vqlexp.MRegExpCollection["Spaces"].ReplaceAllLiteralString(sValue, "")
		// sValue = trimQuotationMarks(sValue)
		slLimitAfter = append(slLimitAfter, sValue)
	}
	iLenLimit := len(slLimitAfter)

	switch iLenLimit {
	case 1:
		iStart, err = strconv.Atoi(slLimitAfter[0])
		if err != nil {
			return gtypes.TLimit{
				Is: false,
			}
		}
	case 2:
		iStart, err = strconv.Atoi(slLimitAfter[0])
		if err != nil {
			return gtypes.TLimit{
				Is: false,
			}
		}
		iOffset, err = strconv.Atoi(slLimitAfter[1])
		if err != nil {
			return gtypes.TLimit{
				Is: false,
			}
		}
	default:
		return gtypes.TLimit{
			Is: false,
		}
	}

	return gtypes.TLimit{
		Is:     true,
		Start:  iStart,
		Offset: iOffset,
	}
}

func parseOrderBy(sOrderBy string) gtypes.TOrderBy {
	// -
	var stOBCols = gtypes.TOrderBy{
		Cols: make([]string, 0, 2),
		Sort: make([]uint8, 0, 2),
	}

	slOrderBy := vqlexp.MRegExpCollection["Comma"].Split(sOrderBy, -1)
	for _, sOBCol := range slOrderBy {
		sCol := ""
		uAD := uint8(0)

		if vqlexp.MRegExpCollection["asc"].MatchString(sOBCol) {
			sCol = vqlexp.MRegExpCollection["asc"].ReplaceAllLiteralString(sOBCol, "")
			uAD = 1
		} else if vqlexp.MRegExpCollection["desc"].MatchString(sOBCol) {
			sCol = vqlexp.MRegExpCollection["desc"].ReplaceAllLiteralString(sOBCol, "")
			uAD = 2
		} else {
			sCol = sOBCol
			uAD = 0
		}

		sCol = vqlexp.MRegExpCollection["Spaces"].ReplaceAllLiteralString(sCol, "")
		sCol = trimQuotationMarks(sCol)
		if sCol != "" {
			stOBCols.Cols = append(stOBCols.Cols, sCol)
			stOBCols.Sort = append(stOBCols.Sort, uAD)
		}
	}

	return stOBCols
}

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
	var stOrderByExp gtypes.TOrderBy
	var stLimitExp gtypes.TLimit

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
	sWhere := strings.TrimLeft(sRight, "= ")

	if vqlexp.MRegExpCollection["LimitToEnd"].MatchString(sWhere) {
		sLimit := vqlexp.MRegExpCollection["LimitToEnd"].FindString(sWhere)
		sWhere = vqlexp.MRegExpCollection["LimitToEnd"].ReplaceAllLiteralString(sWhere, "")
		sLimit = vqlexp.MRegExpCollection["Limit"].ReplaceAllLiteralString(sLimit, "")
		stLimitExp = parseLimit(sLimit)
	}

	if vqlexp.MRegExpCollection["OrderbyToEnd"].MatchString(sWhere) {
		sOrderBy := vqlexp.MRegExpCollection["OrderbyToEnd"].FindString(sWhere)
		sWhere = vqlexp.MRegExpCollection["OrderbyToEnd"].ReplaceAllLiteralString(sWhere, "")
		sOrderBy = vqlexp.MRegExpCollection["Orderby"].ReplaceAllLiteralString(sOrderBy, "")
		stOrderByExp = parseOrderBy(sOrderBy)
	}

	sWhere = strings.TrimSpace(sWhere)
	stExpression, okW := parseWhere(sWhere)
	if !okW {
		return "", false
	}

	// FIXME: --
	_ = stExpression
	_ = stOrderByExp
	_ = stLimitExp

	sOut := fmt.Sprintf("{\"%s\": %s}", sLeft, result)
	return sOut, true
}
