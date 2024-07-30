package vqlexp

import "regexp"

type tRegExpCollection map[string]*regexp.Regexp

var MRegExpCollection tRegExpCollection

var ArParsingOrder = [...]string{
	"Where",
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
	mRECol = mRECol.CompileExp("Comment", `(?m)^\/\/`)

	// recol = recol.CompileExp("HeadCleaner", `(?m)^\s*\n*\s*`)
	// recol = recol.CompileExp("AnyCommand", `(?m)^[a-zA-Z].*;\s*`)
	mRECol = mRECol.CompileExp("EntityName", `(?m)^[a-zA-Z][a-zA-Z0-9_-]*$`) // protection of technical names
	mRECol = mRECol.CompileExp("QuotationMarks", `(?m)^[\'\"]|[\'\"]$`)
	mRECol = mRECol.CompileExp("SpecQuotationMark", "(?m)^[`]|[`]$")
	mRECol = mRECol.CompileExp("Comma", `(?m),`)
	mRECol = mRECol.CompileExp("SignEqual", `=`) // FIXME: there may be problems with the equality symbol inside the values

	mRECol = mRECol.CompileExp("Spaces", `(?m)\s+`)
	mRECol = mRECol.CompileExp("BeginBlock", `(?m)\s*\{$`)
	mRECol = mRECol.CompileExp("EndBlock", `(?m)^\}$`)

	// Directives and reserved words TODO: Разработать шаблоны
	mRECol = mRECol.CompileExp("FuncSignature", `(?m)^func\s+[a-zA-Z][a-zA-Z0-9_\-]*\([a-zA-Z0-9_\-\$\s\,\[\]\"\'\{\}]*\)\s+\(*[a-zA-Z_\-\[\]\,\s\{\}]*\)*\s*\{$`)
	mRECol = mRECol.CompileExp("FuncWord", `(?m)^func\s+`)
	mRECol = mRECol.CompileExp("FuncWordAndName", `(?m)^func\s+[a-zA-Z][a-zA-Z0-9_\-]*`)
	mRECol = mRECol.CompileExp("FuncInVarString", `(?m)^\([a-zA-Z0-9_\-\$\s\,\[\]\"\'\{\}]*\)\s*`)
	mRECol = mRECol.CompileExp("FuncDesc", `(?m)\([a-zA-Z0-9_\-\$\s\,\[\]\"\'\{\}]*\)\s+[a-zA-Z_\-\[\]\(\)\,\s\{\}]*\s*\{$`)

	mRECol = mRECol.CompileExp("Where", `(?m)^\/\/`) // FIXME:

	// GPF - General Purpose Functions TODO: Разработать шаблоны

	// DDF - Data Definition Functions TODO: Разработать шаблоны

	// DMF - Data Manipulation Functions TODO: Разработать шаблоны

	// DCF - Data Control Functions TODO: Разработать шаблоны

	return mRECol
}

func init() {
	MRegExpCollection = CompileRegExpCollection()
}
