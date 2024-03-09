package vqlexp

import "regexp"

type tRegExpCollection map[string]*regexp.Regexp

var RegExpCollection tRegExpCollection

var ParsingOrder = [...]string{
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

func (r tRegExpCollection) CompileExp(name string, expr string) tRegExpCollection {
	// This method is complete
	re, err := regexp.Compile(expr)
	if err != nil {
		return r
	}
	r[name] = re

	return r
}

func CompileRegExpCollection() tRegExpCollection {
	// -
	var recol tRegExpCollection = make(tRegExpCollection)

	recol = recol.CompileExp("LineBreak", `(?m)\n`)
	// recol = recol.CompileExp("HeadCleaner", `(?m)^\s*\n*\s*`)
	// recol = recol.CompileExp("AnyCommand", `(?m)^[a-zA-Z].*;\s*`)
	recol = recol.CompileExp("EntityName", `(?m)^[a-zA-Z][a-zA-Z0-9_-]*$`) // protection of technical names
	recol = recol.CompileExp("QuotationMarks", `(?m)[\'\"]`)
	recol = recol.CompileExp("SpecQuotationMark", "(?m)[`]")
	recol = recol.CompileExp("Spaces", `(?m)\s*`)
	recol = recol.CompileExp("Comma", `(?m),`)
	recol = recol.CompileExp("SignEqual", `(?m)=`)

	recol = recol.CompileExp("IfNotExistsWord", `(?m)[iI][fF]\s*[nN][oO][tT]\s*[eE][xX][iI][sS][tT][sS]`)
	recol = recol.CompileExp("IfExistsWord", `(?m)[iI][fF]\s*[eE][xX][iI][sS][tT][sS]`)

	recol = recol.CompileExp("ON", `(?m)[oO][nN]`)
	recol = recol.CompileExp("TO", `(?m)[tT][oO]`)
	recol = recol.CompileExp("FROM", `(?m)[fF][rR][oO][mM]`)
	recol = recol.CompileExp("ADD", `(?m)[aA][dD][dD]`)
	recol = recol.CompileExp("DROP", `(?m)[dD][rR][oO][pP]`)
	recol = recol.CompileExp("MODIFY", `(?m)[mM][oO][dD][iI][fF][yY]`)
	recol = recol.CompileExp("RenameTo", `(?m)[rR][eE][nN][aA][mM][eE]\s*[tT][oO]`)

	recol = recol.CompileExp("WhereToEnd", `(?m)[wW][hH][eE][rR][eE].*`)
	recol = recol.CompileExp("Where", `(?m)[wW][hH][eE][rR][eE]`)
	recol = recol.CompileExp("WhereExpression", `(?m)(\s+[aA][nN][dD]\s+|\s+[oO][rR]\s+).*`)
	recol = recol.CompileExp("WhereExpression_And_Or_Word", `(?m)^(\s+[aA][nN][dD]\s+|\s+[oO][rR]\s+)`)
	recol = recol.CompileExp("AND", `(?m)[aA][nN][dD]`)
	recol = recol.CompileExp("OR", `(?m)[oO][rR]`)
	recol = recol.CompileExp("WhereOperationConditions", `(?m)(<|>|<=|>=|=|[lL][iI][kK][eE])`)
	recol = recol.CompileExp("WhereOperation_<=", `(?m)<=`)
	recol = recol.CompileExp("WhereOperation_>=", `(?m)>=`)
	recol = recol.CompileExp("WhereOperation_<", `(?m)<`)
	recol = recol.CompileExp("WhereOperation_>", `(?m)>`)
	recol = recol.CompileExp("WhereOperation_=", `(?m)=`)
	recol = recol.CompileExp("WhereOperation_LIKE", `(?m)[lL][iI][kK][eE]`)

	// DDL
	recol = recol.CompileExp("ColumnUnique", `(?m)[uU][nN][iI][qQ][uU][eE]\s*$`)
	recol = recol.CompileExp("ColumnNotNull", `(?m)[nN][oO][tT]\s*[nN][uU][lL][lL]\s*$`)
	recol = recol.CompileExp("ColumnDefault", `(?m)[dD][eE][fF][aA][uU][lL][tT]:.+`)
	recol = recol.CompileExp("ColumnDefaultWord", `(?m)[dD][eE][fF][aA][uU][lL][tT]:`)
	recol = recol.CompileExp("TableParenthesis", `(?m)[\(\)]`)

	recol = recol.CompileExp("SearchCreate", `(?m)^[cC][rR][eE][aA][tT][eE].*`)
	recol = recol.CompileExp("CreateDatabaseWord", `(?m)^[cC][rR][eE][aA][tT][eE]\s*[dD][aA][tT][aA][bB][aA][sS][eE]`)
	recol = recol.CompileExp("CreateTableWord", `(?m)^[cC][rR][eE][aA][tT][eE]\s*[tT][aA][bB][lL][eE]`)
	recol = recol.CompileExp("TableColumns", `(?m)\(.*\)`)

	recol = recol.CompileExp("SearchAlter", `(?m)^[aA][lL][tT][eE][rR].*`)
	recol = recol.CompileExp("AlterDatabaseWord", `(?m)^[aA][lL][tT][eE][rR]\s*[dD][aA][tT][aA][bB][aA][sS][eE]`)
	recol = recol.CompileExp("AlterDatabaseRenameTo", `(?m)^[aA][lL][tT][eE][rR]\s*[dD][aA][tT][aA][bB][aA][sS][eE].*[rR][eE][nN][aA][mM][eE]\s*[tT][oO]`)
	recol = recol.CompileExp("AlterTableWord", `(?m)^[aA][lL][tT][eE][rR]\s*[tT][aA][bB][lL][eE]`)
	recol = recol.CompileExp("AlterTableAdd", `(?m)^[aA][lL][tT][eE][rR]\s*[tT][aA][bB][lL][eE].*[aA][dD][dD]`)
	recol = recol.CompileExp("AlterTableDrop", `(?m)^[aA][lL][tT][eE][rR]\s*[tT][aA][bB][lL][eE].*[dD][rR][oO][pP]`)
	recol = recol.CompileExp("AlterTableModify", `(?m)^[aA][lL][tT][eE][rR]\s*[tT][aA][bB][lL][eE].*[mM][oO][dD][iI][fF][yY]`)
	recol = recol.CompileExp("AlterTableRenameTo", `(?m)^[aA][lL][tT][eE][rR]\s*[tT][aA][bB][lL][eE].*[rR][eE][nN][aA][mM][eE]\s*[tT][oO]`)

	recol = recol.CompileExp("SearchDrop", `(?m)^[dD][rR][oO][pP].*`)
	recol = recol.CompileExp("DropDatabaseWord", `(?m)^[dD][rR][oO][pP]\s*[dD][aA][tT][aA][bB][aA][sS][eE]`)
	recol = recol.CompileExp("DropTableWord", `(?m)^[dD][rR][oO][pP]\s*[tT][aA][bB][lL][eE]`)

	// DML TODO: Разработать шаблоны
	recol = recol.CompileExp("SearchSelect", `(?m)^;`)

	recol = recol.CompileExp("SearchInsert", `(?m)^[iI][nN][sS][eE][rR][tT]\s*[iI][nN][tT][oO].*`)
	recol = recol.CompileExp("InsertWord", `(?m)^[iI][nN][sS][eE][rR][tT]\s*[iI][nN][tT][oO]`)
	recol = recol.CompileExp("InsertValuesToEnd", `(?m)[vV][aA][lL][uU][eE][sS].*`)
	recol = recol.CompileExp("InsertColParenthesis", `(?m)\(.*\)`)
	recol = recol.CompileExp("InsertParenthesis", `(?m)[\(\)]`)
	recol = recol.CompileExp("InsertValuesWord", `(?m)[vV][aA][lL][uU][eE][sS]`)
	recol = recol.CompileExp("InsertSplitParenthesis", `(?m)\),\s*\(`)

	recol = recol.CompileExp("SearchUpdate", `(?m)^[uU][pP][dD][aA][tT][eE].*[sS][sE][tT].*`)
	recol = recol.CompileExp("UpdateWord", `(?m)^[uU][pP][dD][aA][tT][eE]`)
	recol = recol.CompileExp("UpdateSetToEnd", `(?m)\s*[sS][eE][tT]\s.*`)
	recol = recol.CompileExp("UpdateSetWord", `(?m)\s*[sS][eE][tT]\s`)

	recol = recol.CompileExp("SearchDelete", `(?m)^;`)
	recol = recol.CompileExp("SearchCommit", `(?m)^;`)
	recol = recol.CompileExp("SearchRollback", `(?m)^;`)

	recol = recol.CompileExp("SearchTruncateTable", `(?m)^[tT][rR][uU][nN][cC][aA][tT][eE]\s*[tT][aA][bB][lL][eE].*`)
	recol = recol.CompileExp("TruncateTableWord", `(?m)^[tT][rR][uU][nN][cC][aA][tT][eE]\s*[tT][aA][bB][lL][eE]`)

	// DCL
	// recol = recol.CompileExp("SearchUse", `(?m)^[uU][sS][eE]\s*[a-zA-Z][a-zA-Z0-9_-]+\s*`)
	recol = recol.CompileExp("SearchUse", "(?m)^[uU][sS][eE] *[\"'`]?[a-zA-Z][a-zA-Z0-9_-]+[\"'`]?")
	recol = recol.CompileExp("UseWord", `(?m)^[uU][sS][eE]`)

	recol = recol.CompileExp("SearchShow", `(?m)^[sS][hH][oO][wW].*`)
	recol = recol.CompileExp("ShowDatabasesWord", `(?m)[sS][hH][oO][wW]\s*[dD][aA][tT][aA][bB][aA][sS][eE][sS]`)
	recol = recol.CompileExp("ShowTablesWord", `(?m)[sS][hH][oO][wW]\s*[tT][aA][bB][lL][eE][sS]`)

	recol = recol.CompileExp("SearchExplain", `(?m)^[eE][xX][pP][lL][aA][iI][nN]\s+.*`)
	recol = recol.CompileExp("SearchDescribe", `(?m)^[dD][eE][sS][cC][rR][iI][bB][eE]\s+.*`)
	recol = recol.CompileExp("SearchDesc", `(?m)^[dD][eE][sS][cC]\s+.*`)

	recol = recol.CompileExp("ExplainWord", `(?m)^[eE][xX][pP][lL][aA][iI][nN]`)
	recol = recol.CompileExp("DescribeWord", `(?m)^[dD][eE][sS][cC][rR][iI][bB][eE]`)
	recol = recol.CompileExp("DescWord", `(?m)^[dD][eE][sS][cC]`)

	recol = recol.CompileExp("SearchGrant", `(?m)^[gG][rR][aA][nN][tT].*`)
	recol = recol.CompileExp("GrantWord", `(?m)^[gG][rR][aA][nN][tT]`)
	recol = recol.CompileExp("GrantPrivileges", `(?m)^[gG][rR][aA][nN][tT].*[oO][nN]`)
	recol = recol.CompileExp("GrantPrivilegesList", `(?m)[cC][rR][eE][aA][tT][eE]|[sS][eE][lL][eE][cC][tT]|[iI][nN][sS][eE][rR][tT]|[uU][pP][dD][aA][tT][eE]|[dD][eE][lL][eE][tT][eE]`)
	recol = recol.CompileExp("GrantOnTo", `(?m)[oO][nN].*[tT][oO]`)
	recol = recol.CompileExp("GrantToEnd", `(?m)[tT][oO].*`)

	recol = recol.CompileExp("SearchRevoke", `(?m)^[rR][eE][vV][oO][kK][eE].*`)
	recol = recol.CompileExp("RevokeWord", `(?m)^[rR][eE][vV][oO][kK][eE]`)
	recol = recol.CompileExp("RevokePrivileges", `(?m)^[rR][eE][vV][oO][kK][eE].*[oO][nN]`)
	recol = recol.CompileExp("RevokePrivilegesList", `(?m)[cC][rR][eE][aA][tT][eE]|[sS][eE][lL][eE][cC][tT]|[iI][nN][sS][eE][rR][tT]|[uU][pP][dD][aA][tT][eE]|[dD][eE][lL][eE][tT][eE]`)
	recol = recol.CompileExp("RevokeOnTo", `(?m)[oO][nN].*[tT][oO]`)
	recol = recol.CompileExp("RevokeToEnd", `(?m)[tT][oO].*`)

	recol = recol.CompileExp("SearchAuth", `(?m)^[aA][uU][tT][hH].+`)
	recol = recol.CompileExp("AuthNew", `(?m)^[aA][uU][tT][hH]\s*[nN][eE][wW]`)
	recol = recol.CompileExp("AuthChange", `(?m)^[aA][uU][tT][hH]\s*[cC][hH][aA][nN][gG][eE]`)
	recol = recol.CompileExp("AuthRemove", `(?m)^[aA][uU][tT][hH]\s*[rR][eE][mM][oO][vV][eE]`)
	recol = recol.CompileExp("Login", `(?m)[lL][oO][gG][iI][nN]\s+\S+(\s+|$)`)
	recol = recol.CompileExp("LoginWord", `(?m)[lL][oO][gG][iI][nN]`)
	recol = recol.CompileExp("Password", `(?m)[pP][aA][sS][sS][wW][oO][rR][dD]\s+\S+(\s+|$)`)
	recol = recol.CompileExp("PasswordWord", `(?m)[pP][aA][sS][sS][wW][oO][rR][dD]`)
	recol = recol.CompileExp("Hash", `(?m)[hH][aA][sS][hH]\s+\S+(\s+|$)`)
	recol = recol.CompileExp("HashWord", `(?m)[hH][aA][sS][hH]`)
	recol = recol.CompileExp("Role", `(?m)[rR][oO][lL][eE].*`)
	recol = recol.CompileExp("RoleWord", `(?m)[rR][oO][lL][eE]`)

	return recol
}

func init() {
	RegExpCollection = CompileRegExpCollection()
}
