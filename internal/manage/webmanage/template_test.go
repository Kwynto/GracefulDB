package webmanage

import (
	"testing"
)

func Test_parseTemplates(t *testing.T) {
	t.Run("parseTemplates() function testing", func(t *testing.T) {
		parseTemplates() // calling the tested function

		iCount := len(MTemplates)
		if iCount == 0 {
			t.Errorf("parseTemplates() error. Count: %d", iCount)
		}
	})
}

func Test_loadTemplateFromVar(t *testing.T) {

	t.Run("loadTemplateFromVar() function testing - negative", func(t *testing.T) {
		sWrong := `
		<html>
			{{ errorexp }}
		</html>
		`

		err := LoadTemplateFromString("wrongStr", sWrong) // calling the tested function
		if err == nil {
			t.Error("LoadTemplateFromString() error.")
		}
	})
}
