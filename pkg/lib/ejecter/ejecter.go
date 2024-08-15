package ejecter

import "regexp"

type tRegExpCollection map[string]*regexp.Regexp

var MRegExpCollection tRegExpCollection

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

	return mRECol
}

func init() {
	MRegExpCollection = CompileRegExpCollection()
}
