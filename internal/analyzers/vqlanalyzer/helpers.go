package vqlanalyzer

import (
	"errors"
	"slices"
	"strings"

	"github.com/Kwynto/GracefulDB/internal/engine/basicsystem/gauth"
	"github.com/Kwynto/GracefulDB/internal/engine/basicsystem/gtypes"
	"github.com/Kwynto/GracefulDB/internal/engine/basicsystem/vqlexp"
	"github.com/Kwynto/GracefulDB/internal/engine/core"
)

// Helpers for VQLAnalyzer

func parseOrderBy(sOrderBy string, slColumns []string) (gtypes.TOrderBy, error) {
	var stOBCols = gtypes.TOrderBy{
		Cols: make([]string, 0, 2),
		Sort: make([]uint8, 0, 2),
	}

	slOrderBy := vqlexp.MRegExpCollection["Comma"].Split(sOrderBy, -1)
	for _, sOBCol := range slOrderBy {
		// разобрать ...
		sCol := ""
		uAD := uint8(0)

		if vqlexp.MRegExpCollection["ASC"].MatchString(sOBCol) {
			sCol = vqlexp.MRegExpCollection["ASC"].ReplaceAllLiteralString(sOBCol, "")
			uAD = 1
		} else if vqlexp.MRegExpCollection["DESC"].MatchString(sOBCol) {
			sCol = vqlexp.MRegExpCollection["DESC"].ReplaceAllLiteralString(sOBCol, "")
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

	if len(stOBCols.Cols) < 1 {
		return stOBCols, errors.New("group-by error")
	}

	for _, obCol := range stOBCols.Cols {
		if !slices.Contains(slColumns, obCol) {
			return stOBCols, errors.New("group-by error")
		}
	}

	return stOBCols, nil
}

func parseGroupBy(sGroupBy string, slColumns []string) ([]string, error) {
	var slGBCols = make([]string, 0, 4)
	slGroupBy := vqlexp.MRegExpCollection["Comma"].Split(sGroupBy, -1)
	for _, sGBCol := range slGroupBy {
		sGBCol = vqlexp.MRegExpCollection["Spaces"].ReplaceAllLiteralString(sGBCol, "")
		sGBCol = trimQuotationMarks(sGBCol)
		if sGBCol != "" {
			slGBCols = append(slGBCols, sGBCol)
		}
	}
	if len(slGBCols) < 1 {
		return slGBCols, errors.New("group-by error")
	}
	for _, sGBCol := range slGBCols {
		if !slices.Contains(slColumns, sGBCol) {
			return slGBCols, errors.New("group-by error")
		}
	}
	return slGBCols, nil
}

func parseWhere(sWhere string) ([]gtypes.TConditions, error) {
	var slExpression = make([]gtypes.TConditions, 0, 4)
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
			return []gtypes.TConditions{}, errors.New("condition error")
		}
		if sValueIn == "" {
			return []gtypes.TConditions{}, errors.New("condition error")
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
		} else if vqlexp.MRegExpCollection["WhereOperation_="].MatchString(sHeadCond) {
			stExp.Operation = "="
		} else if vqlexp.MRegExpCollection["WhereOperation_LIKE"].MatchString(sHeadCond) {
			stExp.Operation = "like"
		} else if vqlexp.MRegExpCollection["WhereOperation_REGEXP"].MatchString(sHeadCond) {
			stExp.Operation = "regexp"
		} else {
			return []gtypes.TConditions{}, errors.New("condition error")
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
	return slExpression, nil
}

func trimQuotationMarks(input string) string {
	if vqlexp.MRegExpCollection["QuotationMarks"].MatchString(input) {
		input = vqlexp.MRegExpCollection["QuotationMarks"].ReplaceAllLiteralString(input, "")
		return input
	}

	if vqlexp.MRegExpCollection["SpecQuotationMark"].MatchString(input) {
		input = vqlexp.MRegExpCollection["SpecQuotationMark"].ReplaceAllLiteralString(input, "")
		return input
	}

	return input
}

func preChecker(sTicket string) (sLogin string, sDB string, stAccess gauth.TProfile, sNewTicket string, err error) {
	if sTicket == "" {
		return sLogin, sDB, stAccess, sNewTicket, errors.New("an empty ticket")
	}

	sLogin, stAccess, sNewTicket, err = gauth.CheckTicket(sTicket)
	if err != nil {
		return sLogin, sDB, stAccess, sNewTicket, err
	}

	if stAccess.Status.IsBad() {
		return sLogin, sDB, stAccess, sNewTicket, errors.New("auth error")
	}

	stState, isOk := core.MStates[sTicket]
	if !isOk {
		return sLogin, sDB, stAccess, sNewTicket, errors.New("unknown database")
	}
	sDB = stState.CurrentDB
	if sDB == "" {
		return sLogin, sDB, stAccess, sNewTicket, errors.New("no database selected")
	}

	return sLogin, sDB, stAccess, sNewTicket, nil
}

func dourPostChecker(sDB, sTable, sLogin string, stAccess gauth.TProfile) (isLuxUser bool, stFlagsAcs gtypes.TAccessFlags, err error) {
	stDBInfo, isOkDB := core.GetDBInfo(sDB)
	if isOkDB {
		var isOkFlags bool = false
		stFlagsAcs = gtypes.TAccessFlags{}
		isLuxUser = false

		_, isOkTable := stDBInfo.Tables[sTable]
		if !isOkTable {
			return isLuxUser, stFlagsAcs, errors.New("invalid table name")
		}

		stDBAccess, isOkAccess := core.GetDBAccess(sDB)
		if isOkAccess {
			stFlagsAcs, isOkFlags = stDBAccess.Flags[sLogin]
			if stDBAccess.Owner != sLogin {
				for iRole := range stAccess.Roles {
					if iRole == int(gauth.ADMIN) || iRole == int(gauth.ENGINEER) {
						isLuxUser = true
						break
					}
				}
				if !isLuxUser {
					if !isOkFlags {
						return isLuxUser, stFlagsAcs, errors.New("not enough rights")
					}
				}
			} else {
				isLuxUser = true
			}
		} else {
			return isLuxUser, stFlagsAcs, errors.New("internal error")
		}

		return isLuxUser, stFlagsAcs, nil

	} else {
		return isLuxUser, stFlagsAcs, errors.New("invalid database name")
	}
}

func friendlyPostChecker(sDB, sTable, sLogin string, stAccess gauth.TProfile) (isLuxUser bool, stFlagsAcs gtypes.TAccessFlags, err error) {
labelCheck:
	stDBInfo, isOkDB := core.GetDBInfo(sDB)
	if isOkDB {
		var isOkFlags bool = false
		stFlagsAcs = gtypes.TAccessFlags{}
		isLuxUser = false

		_, isOkTable := stDBInfo.Tables[sTable]
		if !isOkTable {
			if core.StLocalCoreSettings.FriendlyMode {
				if !core.CreateTable(sDB, sTable, true) {
					return isLuxUser, stFlagsAcs, errors.New("invalid table name")
				}
				goto labelCheck
			}
			return isLuxUser, stFlagsAcs, errors.New("invalid table name")
		}

		stDBAccess, isOkAccess := core.GetDBAccess(sDB)
		if isOkAccess {
			stFlagsAcs, isOkFlags = stDBAccess.Flags[sLogin]
			if stDBAccess.Owner != sLogin {
				for iRole := range stAccess.Roles {
					if iRole == int(gauth.ADMIN) || iRole == int(gauth.ENGINEER) {
						isLuxUser = true
						break
					}
				}
				if !isLuxUser {
					if !isOkFlags {
						return isLuxUser, stFlagsAcs, errors.New("not enough rights")
					}
				}
			} else {
				isLuxUser = true
			}
		} else {
			return isLuxUser, stFlagsAcs, errors.New("internal error")
		}

		return isLuxUser, stFlagsAcs, nil

	} else {
		if core.StLocalCoreSettings.FriendlyMode {
			if !core.CreateDB(sDB, sLogin, true) {
				return isLuxUser, stFlagsAcs, errors.New("invalid database name")
			}
			goto labelCheck
		}
		return isLuxUser, stFlagsAcs, errors.New("internal error")
	}
}
