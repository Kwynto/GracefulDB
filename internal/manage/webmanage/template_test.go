package webmanage

import (
	"testing"

	"github.com/Kwynto/GracefulDB/pkg/lib/helpers/masquerade/home_masq"
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
	t.Run("loadTemplateFromVar() function testing - positive", func(t *testing.T) {
		err := LoadTemplateFromString(HOME_TEMP_NAME, home_masq.HtmlHome) // calling the tested function
		if err != nil {
			t.Error("loadTemplateFromVar() error.")
		}
	})

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
