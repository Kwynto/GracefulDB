package webmanage

import (
	"testing"
)

func Test_parseTemplates(t *testing.T) {
	t.Run("parseTemplates() function testing", func(t *testing.T) {
		parseTemplates() // calling the tested function

		count := len(TemplatesMap)
		if count == 0 {
			t.Errorf("parseTemplates() error. Count: %d", count)
		}
	})
}

func Test_loadTemplateFromVar(t *testing.T) {

	t.Run("loadTemplateFromVar() function testing - negative", func(t *testing.T) {
		wrongStr := `
		<html>
			{{ errorexp }}
		</html>
		`

		err := LoadTemplateFromString("wrongStr", wrongStr) // calling the tested function
		if err == nil {
			t.Error("loadTemplateFromVar() error.")
		}
	})
}
