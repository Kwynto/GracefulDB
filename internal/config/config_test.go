package config

import (
	"os"
	"reflect"
	"testing"
)

func Test_defaultConfig(t *testing.T) {
	stRes := defaultConfig()
	if reflect.TypeOf(stRes) != reflect.TypeOf(TConfig{}) {
		t.Error("defaultConfig() error = The function returns the wrong type")
	}
}

func Test_MustLoad(t *testing.T) {
	stRes := MustLoad("")
	if reflect.TypeOf(stRes) != reflect.TypeOf(&TConfig{}) {
		t.Error("MustLoad() error = The function returns the wrong type")
	}

	stRes = MustLoad("./../../config/default.yaml")
	if reflect.TypeOf(stRes) != reflect.TypeOf(&TConfig{}) {
		t.Error("MustLoad() error = The function returns the wrong type")
	}

	f, _ := os.Create("./test.yaml")
	f.Write([]byte("bla-bla-bla"))
	f.Close()

	stRes = MustLoad("./test.yaml")
	if reflect.TypeOf(stRes) != reflect.TypeOf(&TConfig{}) {
		t.Error("MustLoad() error = The function returns the wrong type")
	}

	os.Remove("./test.yaml")
}
