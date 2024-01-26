package core

import (
	"reflect"
	"testing"

	"github.com/Kwynto/GracefulDB/internal/config"
)

func Test_LoadLocalCoreSettings(t *testing.T) {
	t.Run("LoadLocalCoreSettings() function testing", func(t *testing.T) {
		res := LoadLocalCoreSettings(&config.DefaultConfig)
		if reflect.TypeOf(res) != reflect.TypeOf(tCoreSettings{}) {
			t.Error("LoadLocalCoreSettings() error = The function returns the wrong type")
		}
	})
}

func Test_Start(t *testing.T) {
	t.Run("Start() function testing", func(t *testing.T) {
		Start(&config.DefaultConfig)
		if reflect.TypeOf(LocalCoreSettings) != reflect.TypeOf(tCoreSettings{}) {
			t.Error("Start() error = The function returns the wrong type")
		}
	})
}
