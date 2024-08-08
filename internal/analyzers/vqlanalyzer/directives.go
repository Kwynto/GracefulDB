package vqlanalyzer

import (
	"strconv"
	"strings"
	"time"

	"github.com/Kwynto/GracefulDB/internal/engine/basicsystem/gtypes"
	"github.com/Kwynto/GracefulDB/internal/engine/basicsystem/vqlexp"
	"github.com/Kwynto/GracefulDB/internal/engine/core"
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
		iStart = 0
		iOffset, err = strconv.Atoi(slLimitAfter[0])
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

// func (q tQuery) DirectWhere(lineInd int) (result string, ok bool) {
func (q tQuery) DirectWhere(lineInd int) (map[string]any, bool) {
	// This function is complete
	var stOrderByExp gtypes.TOrderBy
	var stLimitExp gtypes.TLimit

	result := make(map[string]any)

	sLine := q.QueryCode[lineInd]

	sLeft := vqlexp.MRegExpCollection["WhereRight"].ReplaceAllLiteralString(sLine, "")
	sLeft = strings.TrimRight(sLeft, "= ")
	if sLeft == "" {
		sLeft = "$result"
	}
	if !vqlexp.MRegExpCollection["VariableWholeString"].MatchString(sLeft) {
		return result, false
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
		stOrderByExp.Is = true
	}

	sWhere = strings.TrimSpace(sWhere)
	stExpression, okW := parseWhere(sWhere)
	if !okW {
		return result, false
	}

	if len(stExpression) == 0 {
		stBaseCond := gtypes.TConditions{
			Type:      "operation",
			Key:       "_id",
			Operation: ">",
			Value:     "0",
		}
		stExpression = append(stExpression, stBaseCond)
	}

	if q.DB == "" || q.Table == "" {
		return result, false
	}

	stDBInfo, isOkDB := core.GetDBInfo(q.DB)
	if !isOkDB {
		return result, false
	}
	stTableInfo, isOkTable := stDBInfo.Tables[q.Table]
	if !isOkTable {
		return result, false
	}

	// chacking keys
	for _, stWhereElem := range stExpression {
		if stWhereElem.Type == "operation" {
			if stWhereElem.Key != "_id" && stWhereElem.Key != "_time" && stWhereElem.Key != "_status" && stWhereElem.Key != "_shape" {
				_, isOk := stTableInfo.Columns[stWhereElem.Key]
				if !isOk {
					return result, false
				}
			}
		}
	}

	stAdditionalData := gtypes.TAdditionalData{
		Db:    q.DB,
		Table: q.Table,
		Stamp: time.Now().Unix(),
	}

	slUWhereIds := core.WhereSelection(stExpression, stAdditionalData)

	slUWhereIds = core.OrderByVQL(slUWhereIds, stOrderByExp, stAdditionalData)

	if stLimitExp.Is {
		iLenIds := len(slUWhereIds)
		iEnd := stLimitExp.Start + stLimitExp.Offset
		if iEnd < iLenIds {
			slUWhereIds = slUWhereIds[stLimitExp.Start:iEnd]
		} else {
			slUWhereIds = slUWhereIds[stLimitExp.Start:]
		}
	}

	// bResult, err := json.Marshal(slUWhereIds)
	// if err != nil {
	// 	return "", false
	// }
	// result = string(bResult)
	// sOut := fmt.Sprintf("{\"%s\": %s}", sLeft, result)
	result[sLeft] = slUWhereIds
	return result, true
}
