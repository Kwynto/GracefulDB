package gtypes

import (
	"reflect"
	"testing"
)

func Test_DefaultData(t *testing.T) {
	t.Run("DefaultData() function testing", func(t *testing.T) {
		res := DefaultData()
		if reflect.TypeOf(res) != reflect.TypeOf(VData{}) {
			t.Error("DefaultData() error = The function returns the wrong type")
		}
	})
}
