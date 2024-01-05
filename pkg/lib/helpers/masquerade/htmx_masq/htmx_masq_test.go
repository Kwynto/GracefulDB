package htmx_masq

import "testing"

const (
	CLOSER_TESTING_ITER_MIN = 5
)

func Test_DefaultBlank(t *testing.T) {
	t.Run("DefaultBlank() function testing - positive", func(t *testing.T) {
		if _, err := DefaultBlank(); err != nil {
			t.Errorf("DefaultBlank() error: %v.", err)
		}
	})

	t.Run("DefaultBlank() function testing - negative", func(t *testing.T) {
		Default = " "
		if _, err := DefaultBlank(); err == nil {
			t.Error("DefaultBlank() do not working error.")
		}
	})
}
