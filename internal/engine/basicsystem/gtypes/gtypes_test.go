package gtypes

import (
	"reflect"
	"testing"
)

func Test_DefaultSecret(t *testing.T) {
	t.Run("DefaultSecret() function testing", func(t *testing.T) {
		res := DefaultSecret()
		if reflect.TypeOf(res) != reflect.TypeOf(Secret{}) {
			t.Error("DefaultSecret() error = The function returns the wrong type")
		}
	})
}
