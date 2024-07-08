package sqlexp

import "regexp"

type tRegExpCollection map[string]*regexp.Regexp

var MRegExpCollection tRegExpCollection

var ArParsingOrder = [...]string{
	"SearchSelect",
	"SearchInsert",
	"SearchUpdate",

	"SearchUse",
	"SearchAuth",

	"SearchDelete",
	"SearchCommit",
	"SearchRollback",

	"SearchShow",
	"SearchCreate",
	"SearchExplain",
	"SearchDescribe",
	"SearchDesc",
	"SearchDrop",
	"SearchAlter",
	"SearchTruncateTable",

	"SearchGrant",
	"SearchRevoke",
}

func (r tRegExpCollection) CompileExp(sName string, sExpr string) tRegExpCollection {
	// This method is complete
	re, err := regexp.Compile(sExpr)
	if err != nil {
		return r
	}
	r[sName] = re

	return r
}

func CompileRegExpCollection() tRegExpCollection {
	// -
	var mRECol tRegExpCollection = make(tRegExpCollection)

	mRECol = mRECol.CompileExp("LineBreak", `(?m)\n`)
	// recol = recol.CompileExp("HeadCleaner", `(?m)^\s*\n*\s*`)
	// recol = recol.CompileExp("AnyCommand", `(?m)^[a-zA-Z].*;\s*`)
	mRECol = mRECol.CompileExp("EntityName", `(?m)^[a-zA-Z][a-zA-Z0-9_-]*$`) // protection of technical names
	mRECol = mRECol.CompileExp("QuotationMarks", `(?m)^[\'\"]|[\'\"]$`)
	mRECol = mRECol.CompileExp("SpecQuotationMark", "(?m)^[`]|[`]$")
	mRECol = mRECol.CompileExp("Spaces", `(?m)\s*`)
	mRECol = mRECol.CompileExp("Comma", `(?m),`)
	mRECol = mRECol.CompileExp("SignEqual", `=`) // FIXME: there may be problems with the equality symbol inside the values

	mRECol = mRECol.CompileExp("IfNotExistsWord", `(?m)[iI][fF]\s*[nN][oO][tT]\s*[eE][xX][iI][sS][tT][sS]`)
	mRECol = mRECol.CompileExp("IfExistsWord", `(?m)[iI][fF]\s*[eE][xX][iI][sS][tT][sS]`)

	mRECol = mRECol.CompileExp("ON", `(?m)[oO][nN]`)
	mRECol = mRECol.CompileExp("TO", `(?m)[tT][oO]`)
	mRECol = mRECol.CompileExp("FROM", `(?m)[fF][rR][oO][mM]`)
	mRECol = mRECol.CompileExp("ADD", `(?m)[aA][dD][dD]`)
	mRECol = mRECol.CompileExp("DROP", `(?m)[dD][rR][oO][pP]`)
	mRECol = mRECol.CompileExp("MODIFY", `(?m)[mM][oO][dD][iI][fF][yY]`)
	mRECol = mRECol.CompileExp("RenameTo", `(?m)[rR][eE][nN][aA][mM][eE]\s*[tT][oO]`)

	mRECol = mRECol.CompileExp("GroupbyToEnd", `(?m)\s+[gG][rR][oO][uU][pP]\s+[bB][yY].*`)
	mRECol = mRECol.CompileExp("Groupby", `(?m)^\s*[gG][rR][oO][uU][pP]\s+[bB][yY]`)

	mRECol = mRECol.CompileExp("OrderbyToEnd", `(?m)\s+[oO][rR][dD][eE][rR]\s+[bB][yY].*`)
	mRECol = mRECol.CompileExp("Orderby", `(?m)^\s*[oO][rR][dD][eE][rR]\s+[bB][yY]`)
	mRECol = mRECol.CompileExp("ASC", `(?m)\s*[aA][sS][cC]\s*`)
	mRECol = mRECol.CompileExp("DESC", `(?m)\s*[dD][eE][sS][cC]\s*`)

	mRECol = mRECol.CompileExp("WhereToEnd", `(?m)\s+[wW][hH][eE][rR][eE].*`)
	mRECol = mRECol.CompileExp("Where", `(?m)^\s*[wW][hH][eE][rR][eE]\s+`)
	mRECol = mRECol.CompileExp("WhereExpression", `(?m)(\s+[aA][nN][dD]\s+|\s+[oO][rR]\s+).*`)
	mRECol = mRECol.CompileExp("WhereExpression_And_Or_Word", `(?m)^(\s+[aA][nN][dD]\s+|\s+[oO][rR]\s+)`)
	mRECol = mRECol.CompileExp("AND", `(?m)[aA][nN][dD]`)
	mRECol = mRECol.CompileExp("OR", `(?m)[oO][rR]`)
	mRECol = mRECol.CompileExp("WhereOperationConditions", `(?m)(<|>|<=|>=|=|[lL][iI][kK][eE])`)
	mRECol = mRECol.CompileExp("WhereOperation_<=", `(?m)<=`)
	mRECol = mRECol.CompileExp("WhereOperation_>=", `(?m)>=`)
	mRECol = mRECol.CompileExp("WhereOperation_<", `(?m)<`)
	mRECol = mRECol.CompileExp("WhereOperation_>", `(?m)>`)
	mRECol = mRECol.CompileExp("WhereOperation_=", `(?m)=`)
	mRECol = mRECol.CompileExp("WhereOperation_LIKE", `(?m)[lL][iI][kK][eE]`)
	mRECol = mRECol.CompileExp("WhereOperation_REGEXP", `(?m)[rR][eE][gG][eE][xX][pP]`)

	// DDL
	mRECol = mRECol.CompileExp("ColumnUnique", `(?m)[uU][nN][iI][qQ][uU][eE]\s*$`)
	mRECol = mRECol.CompileExp("ColumnNotNull", `(?m)[nN][oO][tT]\s*[nN][uU][lL][lL]\s*$`)
	mRECol = mRECol.CompileExp("ColumnDefault", `(?m)[dD][eE][fF][aA][uU][lL][tT]:.+`)
	mRECol = mRECol.CompileExp("ColumnDefaultWord", `(?m)[dD][eE][fF][aA][uU][lL][tT]:`)
	mRECol = mRECol.CompileExp("TableParenthesis", `(?m)[\(\)]`)

	mRECol = mRECol.CompileExp("SearchCreate", `(?m)^[cC][rR][eE][aA][tT][eE].*`)
	mRECol = mRECol.CompileExp("CreateDatabaseWord", `(?m)^[cC][rR][eE][aA][tT][eE]\s*[dD][aA][tT][aA][bB][aA][sS][eE]`)
	mRECol = mRECol.CompileExp("CreateTableWord", `(?m)^[cC][rR][eE][aA][tT][eE]\s*[tT][aA][bB][lL][eE]`)
	mRECol = mRECol.CompileExp("TableColumns", `(?m)\(.*\)`)

	mRECol = mRECol.CompileExp("SearchAlter", `(?m)^[aA][lL][tT][eE][rR].*`)
	mRECol = mRECol.CompileExp("AlterDatabaseWord", `(?m)^[aA][lL][tT][eE][rR]\s*[dD][aA][tT][aA][bB][aA][sS][eE]`)
	mRECol = mRECol.CompileExp("AlterDatabaseRenameTo", `(?m)^[aA][lL][tT][eE][rR]\s*[dD][aA][tT][aA][bB][aA][sS][eE].*[rR][eE][nN][aA][mM][eE]\s*[tT][oO]`)
	mRECol = mRECol.CompileExp("AlterTableWord", `(?m)^[aA][lL][tT][eE][rR]\s*[tT][aA][bB][lL][eE]`)
	mRECol = mRECol.CompileExp("AlterTableAdd", `(?m)^[aA][lL][tT][eE][rR]\s*[tT][aA][bB][lL][eE].*[aA][dD][dD]`)
	mRECol = mRECol.CompileExp("AlterTableDrop", `(?m)^[aA][lL][tT][eE][rR]\s*[tT][aA][bB][lL][eE].*[dD][rR][oO][pP]`)
	mRECol = mRECol.CompileExp("AlterTableModify", `(?m)^[aA][lL][tT][eE][rR]\s*[tT][aA][bB][lL][eE].*[mM][oO][dD][iI][fF][yY]`)
	mRECol = mRECol.CompileExp("AlterTableRenameTo", `(?m)^[aA][lL][tT][eE][rR]\s*[tT][aA][bB][lL][eE].*[rR][eE][nN][aA][mM][eE]\s*[tT][oO]`)

	mRECol = mRECol.CompileExp("SearchDrop", `(?m)^[dD][rR][oO][pP].*`)
	mRECol = mRECol.CompileExp("DropDatabaseWord", `(?m)^[dD][rR][oO][pP]\s*[dD][aA][tT][aA][bB][aA][sS][eE]`)
	mRECol = mRECol.CompileExp("DropTableWord", `(?m)^[dD][rR][oO][pP]\s*[tT][aA][bB][lL][eE]`)

	// DML TODO: Разработать шаблоны
	mRECol = mRECol.CompileExp("SearchSelect", `(?m)^[sS][eE][lL][eE][cC][tT]\s*.*`)
	mRECol = mRECol.CompileExp("SelectWord", `(?m)^[sS][eE][lL][eE][cC][tT]`)
	mRECol = mRECol.CompileExp("SelectFromToEnd", `(?m)\s+[fF][rR][oO][mM].*`)
	mRECol = mRECol.CompileExp("SelectFromWord", `(?m)^\s*[fF][rR][oO][mM]`)
	mRECol = mRECol.CompileExp("SelectDistinctWord", `(?m)^\s*[dD][iI][sS][tT][iI][nN][cC][tT]`)

	mRECol = mRECol.CompileExp("SearchInsert", `(?m)^[iI][nN][sS][eE][rR][tT]\s*[iI][nN][tT][oO].*`)
	mRECol = mRECol.CompileExp("InsertWord", `(?m)^[iI][nN][sS][eE][rR][tT]\s*[iI][nN][tT][oO]`)
	mRECol = mRECol.CompileExp("InsertValuesToEnd", `(?m)[vV][aA][lL][uU][eE][sS].*`)
	mRECol = mRECol.CompileExp("InsertColParenthesis", `(?m)\(.*\)`)
	mRECol = mRECol.CompileExp("InsertParenthesis", `(?m)[\(\)]`)
	mRECol = mRECol.CompileExp("InsertValuesWord", `(?m)[vV][aA][lL][uU][eE][sS]`)
	mRECol = mRECol.CompileExp("InsertSplitParenthesis", `(?m)\),\s*\(`)

	mRECol = mRECol.CompileExp("SearchUpdate", `(?m)^[uU][pP][dD][aA][tT][eE].*[sS][eE][tT].*`)
	mRECol = mRECol.CompileExp("UpdateWord", `(?m)^[uU][pP][dD][aA][tT][eE]`)
	mRECol = mRECol.CompileExp("UpdateSetToEnd", `(?m)\s*[sS][eE][tT]\s.*`)
	mRECol = mRECol.CompileExp("UpdateSetWord", `(?m)\s*[sS][eE][tT]\s`)

	mRECol = mRECol.CompileExp("SearchDelete", `(?m)^[dD][eE][lL][eE][tT][eE]\s+[fF][rR][oO][mM].*`)
	mRECol = mRECol.CompileExp("DeleteWord", `(?m)^[dD][eE][lL][eE][tT][eE]\s+[fF][rR][oO][mM]`)

	mRECol = mRECol.CompileExp("SearchCommit", `(?m)^;`)
	mRECol = mRECol.CompileExp("SearchRollback", `(?m)^;`)

	mRECol = mRECol.CompileExp("SearchTruncateTable", `(?m)^[tT][rR][uU][nN][cC][aA][tT][eE]\s*[tT][aA][bB][lL][eE].*`)
	mRECol = mRECol.CompileExp("TruncateTableWord", `(?m)^[tT][rR][uU][nN][cC][aA][tT][eE]\s*[tT][aA][bB][lL][eE]`)

	// DCL
	// recol = recol.CompileExp("SearchUse", `(?m)^[uU][sS][eE]\s*[a-zA-Z][a-zA-Z0-9_-]+\s*`)
	mRECol = mRECol.CompileExp("SearchUse", "(?m)^[uU][sS][eE] *[\"'`]?[a-zA-Z][a-zA-Z0-9_-]+[\"'`]?")
	mRECol = mRECol.CompileExp("UseWord", `(?m)^[uU][sS][eE]`)

	mRECol = mRECol.CompileExp("SearchShow", `(?m)^[sS][hH][oO][wW].*`)
	mRECol = mRECol.CompileExp("ShowDatabasesWord", `(?m)[sS][hH][oO][wW]\s*[dD][aA][tT][aA][bB][aA][sS][eE][sS]`)
	mRECol = mRECol.CompileExp("ShowTablesWord", `(?m)[sS][hH][oO][wW]\s*[tT][aA][bB][lL][eE][sS]`)

	mRECol = mRECol.CompileExp("SearchExplain", `(?m)^[eE][xX][pP][lL][aA][iI][nN]\s+.*`)
	mRECol = mRECol.CompileExp("SearchDescribe", `(?m)^[dD][eE][sS][cC][rR][iI][bB][eE]\s+.*`)
	mRECol = mRECol.CompileExp("SearchDesc", `(?m)^[dD][eE][sS][cC]\s+.*`)

	mRECol = mRECol.CompileExp("ExplainWord", `(?m)^[eE][xX][pP][lL][aA][iI][nN]`)
	mRECol = mRECol.CompileExp("DescribeWord", `(?m)^[dD][eE][sS][cC][rR][iI][bB][eE]`)
	mRECol = mRECol.CompileExp("DescWord", `(?m)^[dD][eE][sS][cC]`)

	mRECol = mRECol.CompileExp("SearchGrant", `(?m)^[gG][rR][aA][nN][tT].*`)
	mRECol = mRECol.CompileExp("GrantWord", `(?m)^[gG][rR][aA][nN][tT]`)
	mRECol = mRECol.CompileExp("GrantPrivileges", `(?m)^[gG][rR][aA][nN][tT].*[oO][nN]`)
	mRECol = mRECol.CompileExp("GrantPrivilegesList", `(?m)[cC][rR][eE][aA][tT][eE]|[sS][eE][lL][eE][cC][tT]|[iI][nN][sS][eE][rR][tT]|[uU][pP][dD][aA][tT][eE]|[dD][eE][lL][eE][tT][eE]`)
	mRECol = mRECol.CompileExp("GrantOnTo", `(?m)[oO][nN].*[tT][oO]`)
	mRECol = mRECol.CompileExp("GrantToEnd", `(?m)[tT][oO].*`)

	mRECol = mRECol.CompileExp("SearchRevoke", `(?m)^[rR][eE][vV][oO][kK][eE].*`)
	mRECol = mRECol.CompileExp("RevokeWord", `(?m)^[rR][eE][vV][oO][kK][eE]`)
	mRECol = mRECol.CompileExp("RevokePrivileges", `(?m)^[rR][eE][vV][oO][kK][eE].*[oO][nN]`)
	mRECol = mRECol.CompileExp("RevokePrivilegesList", `(?m)[cC][rR][eE][aA][tT][eE]|[sS][eE][lL][eE][cC][tT]|[iI][nN][sS][eE][rR][tT]|[uU][pP][dD][aA][tT][eE]|[dD][eE][lL][eE][tT][eE]`)
	mRECol = mRECol.CompileExp("RevokeOnTo", `(?m)[oO][nN].*[tT][oO]`)
	mRECol = mRECol.CompileExp("RevokeToEnd", `(?m)[tT][oO].*`)

	mRECol = mRECol.CompileExp("SearchAuth", `(?m)^[aA][uU][tT][hH].+`)
	mRECol = mRECol.CompileExp("AuthNew", `(?m)^[aA][uU][tT][hH]\s*[nN][eE][wW]`)
	mRECol = mRECol.CompileExp("AuthChange", `(?m)^[aA][uU][tT][hH]\s*[cC][hH][aA][nN][gG][eE]`)
	mRECol = mRECol.CompileExp("AuthRemove", `(?m)^[aA][uU][tT][hH]\s*[rR][eE][mM][oO][vV][eE]`)
	mRECol = mRECol.CompileExp("Login", `(?m)[lL][oO][gG][iI][nN]\s+\S+(\s+|$)`)
	mRECol = mRECol.CompileExp("LoginWord", `(?m)[lL][oO][gG][iI][nN]`)
	mRECol = mRECol.CompileExp("Password", `(?m)[pP][aA][sS][sS][wW][oO][rR][dD]\s+\S+(\s+|$)`)
	mRECol = mRECol.CompileExp("PasswordWord", `(?m)[pP][aA][sS][sS][wW][oO][rR][dD]`)
	mRECol = mRECol.CompileExp("Hash", `(?m)[hH][aA][sS][hH]\s+\S+(\s+|$)`)
	mRECol = mRECol.CompileExp("HashWord", `(?m)[hH][aA][sS][hH]`)
	mRECol = mRECol.CompileExp("Role", `(?m)[rR][oO][lL][eE].*`)
	mRECol = mRECol.CompileExp("RoleWord", `(?m)[rR][oO][lL][eE]`)

	return mRECol
}

func init() {
	MRegExpCollection = CompileRegExpCollection()
}
