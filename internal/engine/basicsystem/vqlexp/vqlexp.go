package vqlexp

import "regexp"

type tRegExpCollection map[string]*regexp.Regexp

var MRegExpCollection tRegExpCollection

// var ArParsingOrder = [...]string{

// }

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
	mRECol = mRECol.CompileExp("Comment", `(?m)^\/\/`)

	// recol = recol.CompileExp("HeadCleaner", `(?m)^\s*\n*\s*`)
	// recol = recol.CompileExp("AnyCommand", `(?m)^[a-zA-Z].*;\s*`)
	mRECol = mRECol.CompileExp("EntityName", `(?m)^[a-zA-Z][a-zA-Z0-9_-]*$`) // protection of technical names
	mRECol = mRECol.CompileExp("QuotationMarks", `(?m)^[\'\"]|[\'\"]$`)
	mRECol = mRECol.CompileExp("SpecQuotationMark", "(?m)^[`]|[`]$")
	mRECol = mRECol.CompileExp("Spaces", `(?m)\s*`)
	mRECol = mRECol.CompileExp("Comma", `(?m),`)
	mRECol = mRECol.CompileExp("SignEqual", `=`) // FIXME: there may be problems with the equality symbol inside the values

	// DDF TODO: Разработать шаблоны

	// DMF TODO: Разработать шаблоны

	// DCF TODO: Разработать шаблоны

	return mRECol
}

func init() {
	MRegExpCollection = CompileRegExpCollection()
}
