package core

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
	"SearchTruncate",
	"SearchCommit",
	"SearchRollback",

	"SearchShow",
	"SearchCreate",
	"SearchAlter",
	"SearchDrop",

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
	recol = recol.CompileExp("EntityName", `(?m)^[a-zA-Z][a-zA-Z0-9_-]*$`)
	recol = recol.CompileExp("QuotationMarks", `(?m)[\'\"]`)
	recol = recol.CompileExp("SpecQuotationMark", "(?m)[`]")
	recol = recol.CompileExp("Spaces", `(?m)\s*`)
	recol = recol.CompileExp("Comma", `(?m),`)

	recol = recol.CompileExp("IfNotExistsWord", `(?m)[iI][fF]\s*[nN][oO][tT]\s*[eE][xX][iI][sS][tT][sS]`)
	recol = recol.CompileExp("IfExistsWord", `(?m)[iI][fF]\s*[eE][xX][iI][sS][tT][sS]`)

	recol = recol.CompileExp("ON", `(?m)[oO][nN]`)
	recol = recol.CompileExp("TO", `(?m)[tT][oO]`)
	recol = recol.CompileExp("FROM", `(?m)[fF][rR][oO][mM]`)

	// DDL TODO: Разработать шаблоны
	recol = recol.CompileExp("SearchCreate", `(?m)^[cC][rR][eE][aA][tT][eE].*`)
	recol = recol.CompileExp("CreateDatabaseWord", `(?m)^[cC][rR][eE][aA][tT][eE]\s*[dD][aA][tT][aA][bB][aA][sS][eE]`)
	recol = recol.CompileExp("CreateTableWord", `(?m)^[cC][rR][eE][aA][tT][eE]\s*[tT][aA][bB][lL][eE]`)
	recol = recol.CompileExp("TableColumns", `(?m)\(.*\)`)
	recol = recol.CompileExp("TableParenthesis", `(?m)[\(\)]`)
	recol = recol.CompileExp("ColumnUnique", `(?m)[uU][nN][iI][qQ][uU][eE]`)
	recol = recol.CompileExp("ColumnNotNull", `(?m)[nN][oO][tT]\s*[nN][uU][lL][lL]`)
	recol = recol.CompileExp("ColumnDefault", `(?m)[dD][eE][fF][aA][uU][lL][tT]:.+`)
	recol = recol.CompileExp("ColumnDefaultWord", `(?m)[dD][eE][fF][aA][uU][lL][tT]:`)

	recol = recol.CompileExp("SearchAlter", `(?m)^[aA][lL][tT][eE][rR].*`)
	recol = recol.CompileExp("AlterDatabaseWord", `(?m)^[aA][lL][tT][eE][rR]\s*[dD][aA][tT][aA][bB][aA][sS][eE]`)
	recol = recol.CompileExp("AlterTableWord", `(?m)^[aA][lL][tT][eE][rR]\s*[tT][aA][bB][lL][eE]`)
	recol = recol.CompileExp("AlterRenameTo", `(?m)[rR][eE][nN][aA][mM][eE]\s*[tT][oO]`)

	recol = recol.CompileExp("SearchDrop", `(?m)^[dD][rR][oO][pP].*`)
	recol = recol.CompileExp("DropDatabaseWord", `(?m)^[dD][rR][oO][pP]\s*[dD][aA][tT][aA][bB][aA][sS][eE]`)
	recol = recol.CompileExp("DropTableWord", `(?m)^[dD][rR][oO][pP]\s*[tT][aA][bB][lL][eE]`)

	// DML TODO: Разработать шаблоны
	recol = recol.CompileExp("SearchSelect", `(?m)^;`)
	recol = recol.CompileExp("SearchInsert", `(?m)^;`)
	recol = recol.CompileExp("SearchUpdate", `(?m)^;`)
	recol = recol.CompileExp("SearchDelete", `(?m)^;`)
	recol = recol.CompileExp("SearchTruncate", `(?m)^;`)
	recol = recol.CompileExp("SearchCommit", `(?m)^;`)
	recol = recol.CompileExp("SearchRollback", `(?m)^;`)

	// DCL
	// recol = recol.CompileExp("SearchUse", `(?m)^[uU][sS][eE]\s*[a-zA-Z][a-zA-Z0-9_-]+\s*`)
	recol = recol.CompileExp("SearchUse", "(?m)^[uU][sS][eE] *[\"'`]?[a-zA-Z][a-zA-Z0-9_-]+[\"'`]?")
	recol = recol.CompileExp("UseWord", `(?m)^[uU][sS][eE]`)

	recol = recol.CompileExp("SearchShow", `(?m)^[sS][hH][oO][wW].*`)
	recol = recol.CompileExp("ShowDatabasesWord", `(?m)[sS][hH][oO][wW]\s*[dD][aA][tT][aA][bB][aA][sS][eE][sS]`)
	recol = recol.CompileExp("ShowTablesWord", `(?m)[sS][hH][oO][wW]\s*[tT][aA][bB][lL][eE][sS]`)

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
