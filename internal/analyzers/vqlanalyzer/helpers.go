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

func parseOrderBy(orderbyStr string, columns []string) (gtypes.TOrderBy, error) {
	var obCols = gtypes.TOrderBy{
		Cols: make([]string, 0, 2),
		Sort: make([]uint8, 0, 2),
	}

	orderbyArr := vqlexp.RegExpCollection["Comma"].Split(orderbyStr, -1)
	for _, obCol := range orderbyArr {
		// разобрать ...
		col := ""
		uad := uint8(0)

		if vqlexp.RegExpCollection["ASC"].MatchString(obCol) {
			col = vqlexp.RegExpCollection["ASC"].ReplaceAllLiteralString(obCol, "")
			uad = 1
		} else if vqlexp.RegExpCollection["DESC"].MatchString(obCol) {
			col = vqlexp.RegExpCollection["DESC"].ReplaceAllLiteralString(obCol, "")
			uad = 2
		} else {
			col = obCol
			uad = 0
		}

		col = vqlexp.RegExpCollection["Spaces"].ReplaceAllLiteralString(col, "")
		col = trimQuotationMarks(col)
		if col != "" {
			obCols.Cols = append(obCols.Cols, col)
			obCols.Sort = append(obCols.Sort, uad)
		}
	}

	if len(obCols.Cols) < 1 {
		return obCols, errors.New("group-by error")
	}

	for _, obCol := range obCols.Cols {
		if !slices.Contains(columns, obCol) {
			return obCols, errors.New("group-by error")
		}
	}

	return obCols, nil
}

func parseGroupBy(groupbyStr string, columns []string) ([]string, error) {
	var gbCols = make([]string, 0, 4)
	groupbyArr := vqlexp.RegExpCollection["Comma"].Split(groupbyStr, -1)
	for _, gbCol := range groupbyArr {
		gbCol = vqlexp.RegExpCollection["Spaces"].ReplaceAllLiteralString(gbCol, "")
		gbCol = trimQuotationMarks(gbCol)
		if gbCol != "" {
			gbCols = append(gbCols, gbCol)
		}
	}
	if len(gbCols) < 1 {
		return gbCols, errors.New("group-by error")
	}
	for _, gbCol := range gbCols {
		if !slices.Contains(columns, gbCol) {
			return gbCols, errors.New("group-by error")
		}
	}
	return gbCols, nil
}

func parseWhere(whereStr string) ([]gtypes.TConditions, error) {
	var expression = make([]gtypes.TConditions, 0, 4)
	for {
		headCond := vqlexp.RegExpCollection["WhereExpression"].ReplaceAllLiteralString(whereStr, "")
		condition := vqlexp.RegExpCollection["WhereOperationConditions"].Split(headCond, -1)
		keyIn := condition[0]
		valueIn := condition[1]

		keyIn = vqlexp.RegExpCollection["Spaces"].ReplaceAllLiteralString(keyIn, "")
		keyIn = trimQuotationMarks(keyIn)

		valueIn = strings.TrimSpace(valueIn)
		valueIn = trimQuotationMarks(valueIn)

		if keyIn == "" {
			return []gtypes.TConditions{}, errors.New("condition error")
		}
		if valueIn == "" {
			return []gtypes.TConditions{}, errors.New("condition error")
		} // null value, maybe delete a condition

		exp := gtypes.TConditions{
			Type:  "operation",
			Key:   keyIn,
			Value: valueIn,
		}

		if vqlexp.RegExpCollection["WhereOperation_<="].MatchString(headCond) {
			exp.Operation = "<="
		} else if vqlexp.RegExpCollection["WhereOperation_>="].MatchString(headCond) {
			exp.Operation = ">="
		} else if vqlexp.RegExpCollection["WhereOperation_<"].MatchString(headCond) {
			exp.Operation = "<"
		} else if vqlexp.RegExpCollection["WhereOperation_>"].MatchString(headCond) {
			exp.Operation = ">"
		} else if vqlexp.RegExpCollection["WhereOperation_="].MatchString(headCond) {
			exp.Operation = "="
		} else if vqlexp.RegExpCollection["WhereOperation_LIKE"].MatchString(headCond) {
			exp.Operation = "like"
		} else if vqlexp.RegExpCollection["WhereOperation_REGEXP"].MatchString(headCond) {
			exp.Operation = "regexp"
		} else {
			return []gtypes.TConditions{}, errors.New("condition error")
		}
		expression = append(expression, exp)

		whereStr = vqlexp.RegExpCollection["WhereExpression"].FindString(whereStr)
		logicOper := vqlexp.RegExpCollection["WhereExpression_And_Or_Word"].FindString(whereStr)
		// logicOper = strings.TrimSpace(logicOper)

		if vqlexp.RegExpCollection["OR"].MatchString(logicOper) {
			expression = append(expression, gtypes.TConditions{
				Type: "or",
			})
		} else if vqlexp.RegExpCollection["AND"].MatchString(logicOper) {
			expression = append(expression, gtypes.TConditions{
				Type: "and",
			})
		} else {
			break
		}

		whereStr = vqlexp.RegExpCollection["WhereExpression_And_Or_Word"].ReplaceAllLiteralString(whereStr, "")
	}
	return expression, nil
}

func trimQuotationMarks(input string) string {
	if vqlexp.RegExpCollection["QuotationMarks"].MatchString(input) {
		input = vqlexp.RegExpCollection["QuotationMarks"].ReplaceAllLiteralString(input, "")
		return input
	}

	if vqlexp.RegExpCollection["SpecQuotationMark"].MatchString(input) {
		input = vqlexp.RegExpCollection["SpecQuotationMark"].ReplaceAllLiteralString(input, "")
		return input
	}

	return input
}

func preChecker(ticket string) (login string, db string, access gauth.TProfile, newticket string, err error) {
	if ticket == "" {
		return login, db, access, newticket, errors.New("an empty ticket")
	}

	login, access, newticket, err = gauth.CheckTicket(ticket)
	if err != nil {
		return login, db, access, newticket, err
	}

	if access.Status.IsBad() {
		return login, db, access, newticket, errors.New("auth error")
	}

	state, ok := core.States[ticket]
	if !ok {
		return login, db, access, newticket, errors.New("unknown database")
	}
	db = state.CurrentDB
	if db == "" {
		return login, db, access, newticket, errors.New("no database selected")
	}

	return login, db, access, newticket, nil
}

func angryPostChecker(db, table, login string, access gauth.TProfile) (luxUser bool, flagsAcs gtypes.TAccessFlags, err error) {
	dbInfo, okDB := core.GetDBInfo(db)
	if okDB {
		var okFlags bool = false
		flagsAcs = gtypes.TAccessFlags{}
		luxUser = false

		_, okTable := dbInfo.Tables[table]
		if !okTable {
			return luxUser, flagsAcs, errors.New("invalid table name")
		}

		dbAccess, okAccess := core.GetDBAccess(db)
		if okAccess {
			flagsAcs, okFlags = dbAccess.Flags[login]
			if dbAccess.Owner != login {
				for role := range access.Roles {
					if role == int(gauth.ADMIN) || role == int(gauth.ENGINEER) {
						luxUser = true
						break
					}
				}
				if !luxUser {
					if !okFlags {
						return luxUser, flagsAcs, errors.New("not enough rights")
					}
				}
			} else {
				luxUser = true
			}
		} else {
			return luxUser, flagsAcs, errors.New("internal error")
		}

		return luxUser, flagsAcs, nil

	} else {
		return luxUser, flagsAcs, errors.New("invalid database name")
	}
}

func friendlyPostChecker(db, table, login string, access gauth.TProfile) (luxUser bool, flagsAcs gtypes.TAccessFlags, err error) {
LabelCheck:
	dbInfo, okDB := core.GetDBInfo(db)
	if okDB {
		var okFlags bool = false
		flagsAcs = gtypes.TAccessFlags{}
		luxUser = false

		_, okTable := dbInfo.Tables[table]
		if !okTable {
			if core.LocalCoreSettings.FriendlyMode {
				if !core.CreateTable(db, table, true) {
					return luxUser, flagsAcs, errors.New("invalid table name")
				}
				goto LabelCheck
			}
			return luxUser, flagsAcs, errors.New("invalid table name")
		}

		dbAccess, okAccess := core.GetDBAccess(db)
		if okAccess {
			flagsAcs, okFlags = dbAccess.Flags[login]
			if dbAccess.Owner != login {
				for role := range access.Roles {
					if role == int(gauth.ADMIN) || role == int(gauth.ENGINEER) {
						luxUser = true
						break
					}
				}
				if !luxUser {
					if !okFlags {
						return luxUser, flagsAcs, errors.New("not enough rights")
					}
				}
			} else {
				luxUser = true
			}
		} else {
			return luxUser, flagsAcs, errors.New("internal error")
		}

		return luxUser, flagsAcs, nil

	} else {
		if core.LocalCoreSettings.FriendlyMode {
			if !core.CreateDB(db, login, true) {
				return luxUser, flagsAcs, errors.New("invalid database name")
			}
			goto LabelCheck
		}
		return luxUser, flagsAcs, errors.New("internal error")
	}
}
