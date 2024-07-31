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

	mRECol = mRECol.CompileExp("EntityName", `(?m)^[a-zA-Z][a-zA-Z0-9_-]*$`) // protection of technical names
	mRECol = mRECol.CompileExp("Comma", `(?m),`)
	mRECol = mRECol.CompileExp("SignEqual", `=`) // FIXME: there may be problems with the equality symbol inside the values

	mRECol = mRECol.CompileExp("Spaces", `(?m)\s+`)
	mRECol = mRECol.CompileExp("QuotationMarks", `(?m)^[\'\"]|[\'\"]$`)
	mRECol = mRECol.CompileExp("SpecQuotationMark", "(?m)^[`]|[`]$")
	mRECol = mRECol.CompileExp("BeginBlock", `(?m)\s*\{$`)
	mRECol = mRECol.CompileExp("EndBlock", `(?m)^\}$`)
	mRECol = mRECol.CompileExp("VariableWholeString", `(?m)^\$[a-zA-Z0-9]*$`)
	mRECol = mRECol.CompileExp("Variable", `(?m)\$[a-zA-Z0-9]*`)

	// Directives and reserved words TODO: Разработать шаблоны
	mRECol = mRECol.CompileExp("FuncSignature", `(?m)^func\s+[a-zA-Z][a-zA-Z0-9_\-]*\([a-zA-Z0-9_\-\$\s\,\[\]\"\'\{\}]*\)\s+\(*[a-zA-Z_\-\[\]\,\s\{\}]*\)*\s*\{$`)
	mRECol = mRECol.CompileExp("FuncWord", `(?m)^func\s+`)
	mRECol = mRECol.CompileExp("FuncWordAndName", `(?m)^func\s+[a-zA-Z][a-zA-Z0-9_\-]*`)
	mRECol = mRECol.CompileExp("FuncInVarString", `(?m)^\([a-zA-Z0-9_\-\$\s\,\[\]\"\'\{\}]*\)\s*`)
	mRECol = mRECol.CompileExp("FuncDesc", `(?m)\([a-zA-Z0-9_\-\$\s\,\[\]\"\'\{\}]*\)\s+[a-zA-Z_\-\[\]\(\)\,\s\{\}]*\s*\{$`)

	mRECol = mRECol.CompileExp("Where", `(?m)^[\$a-zA-Z0-9\s\=]*where.*$`)
	mRECol = mRECol.CompileExp("WhereWord", `(?m)^\s*where\s*`)
	mRECol = mRECol.CompileExp("WhereRight", `(?m)\s+where.*`)
	mRECol = mRECol.CompileExp("WhereExpression", `(?m)(\s+[aA][nN][dD]\s+|\s+[oO][rR]\s+).*`)
	mRECol = mRECol.CompileExp("WhereExpression_And_Or_Word", `(?m)^(\s+[aA][nN][dD]\s+|\s+[oO][rR]\s+)`)
	mRECol = mRECol.CompileExp("AND", `(?m)[aA][nN][dD]`)
	mRECol = mRECol.CompileExp("OR", `(?m)[oO][rR]`)
	mRECol = mRECol.CompileExp("WhereOperationConditions", `(?m)(<|>|<=|>=|==|[lL][iI][kK][eE]|[rR][eE][gG][eE][xX][pP])`)
	mRECol = mRECol.CompileExp("WhereOperation_<=", `(?m)<=`)
	mRECol = mRECol.CompileExp("WhereOperation_>=", `(?m)>=`)
	mRECol = mRECol.CompileExp("WhereOperation_<", `(?m)<`)
	mRECol = mRECol.CompileExp("WhereOperation_>", `(?m)>`)
	mRECol = mRECol.CompileExp("WhereOperation_==", `(?m)==`)
	mRECol = mRECol.CompileExp("WhereOperation_LIKE", `(?m)[lL][iI][kK][eE]`)
	mRECol = mRECol.CompileExp("WhereOperation_REGEXP", `(?m)[rR][eE][gG][eE][xX][pP]`)

	// GPF - General Purpose Functions TODO: Разработать шаблоны

	// DDF - Data Definition Functions TODO: Разработать шаблоны

	// DMF - Data Manipulation Functions TODO: Разработать шаблоны

	// DCF - Data Control Functions TODO: Разработать шаблоны

	return mRECol
}

func init() {
	MRegExpCollection = CompileRegExpCollection()
}
