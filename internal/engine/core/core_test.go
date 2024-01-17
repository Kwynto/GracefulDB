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

func Test_Engine(t *testing.T) {
	t.Run("Engine() function testing", func(t *testing.T) {
		Engine(&config.DefaultConfig)
		if reflect.TypeOf(LocalCoreSettings) != reflect.TypeOf(tCoreSettings{}) {
			t.Error("Engine() error = The function returns the wrong type")
		}
	})
}
