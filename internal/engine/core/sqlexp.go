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

	"SearchCreate",
	"SearchAlter",
	"SearchDrop",

	"SearchGrant",
	"SearchRevoke",
}

func (r tRegExpCollection) CompileExp(name string, expr string) tRegExpCollection {
	// This method is completes
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
	// recol = recol.CompileExp("LineBreak", `(?m)\n`)
	// recol = recol.CompileExp("HeadCleaner", `(?m)^\s*\n*\s*`)
	// recol = recol.CompileExp("AnyCommand", `(?m)^[a-zA-Z].*;\s*`)
	recol = recol.CompileExp("EntityName", `(?m)^[a-zA-Z][a-zA-Z0-9_-]*$`)
	recol = recol.CompileExp("QuotationMarks", `(?m)[\'\"]`)
	recol = recol.CompileExp("SpecQuotationMark", "(?m)[`]")

	// DDL TODO: Разработать шаблоны
	recol = recol.CompileExp("SearchCreate", `(?m)^;`)
	recol = recol.CompileExp("SearchAlter", `(?m)^;`)
	recol = recol.CompileExp("SearchDrop", `(?m)^;`)

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
	recol = recol.CompileExp("UseWord", `(?m)[uU][sS][eE]`)

	recol = recol.CompileExp("SearchGrant", `(?m)^[gG][rR][aA][nN][tT].*`)
	recol = recol.CompileExp("SearchRevoke", `(?m)^[rR][eE][vV][oO][kK][eE].*`)

	recol = recol.CompileExp("SearchAuth", `(?m)^[aA][uU][tT][hH].+`)
	// recol = recol.CompileExp("Auth", `(?m)^[aA][uU][tT][hH]`)
	recol = recol.CompileExp("Login", `(?m)[lL][oO][gG][iI][nN]\s+\S+(\s+|$)`)
	recol = recol.CompileExp("LoginWord", `(?m)[lL][oO][gG][iI][nN]`)
	recol = recol.CompileExp("Password", `(?m)[pP][aA][sS][sS][wW][oO][rR][dD]\s+\S+(\s+|$)`)
	recol = recol.CompileExp("PasswordWord", `(?m)[pP][aA][sS][sS][wW][oO][rR][dD]`)
	recol = recol.CompileExp("Hash", `(?m)[hH][aA][sS][hH]\s+\S+(\s+|$)`)
	recol = recol.CompileExp("HashWord", `(?m)[hH][aA][sS][hH]`)

	return recol
}
