package home_masq

import (
	"testing"
)

const (
	CLOSER_TESTING_ITER_MIN = 5
)

func Test_Default(t *testing.T) {
	t.Run("Default() function testing - positive", func(t *testing.T) {
		if _, err := Default(); err != nil {
			t.Errorf("Default() error: %v.", err)
		}
	})

	t.Run("Default() function testing - negative", func(t *testing.T) {
		HtmlHome = " "
		if _, err := Default(); err == nil {
			t.Error("Default() do not working error.")
		}
	})
}
