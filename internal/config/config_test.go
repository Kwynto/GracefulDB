package config

import (
	"os"
	"reflect"
	"testing"
)

func Test_defaultConfig(t *testing.T) {
	res := defaultConfig()
	if reflect.TypeOf(res) != reflect.TypeOf(Config{}) {
		t.Error("defaultConfig() error = The function returns the wrong type")
	}
}

func Test_MustLoad(t *testing.T) {
	res := MustLoad("")
	if reflect.TypeOf(res) != reflect.TypeOf(&Config{}) {
		t.Error("MustLoad() error = The function returns the wrong type")
	}

	res = MustLoad("./../../config/default.yaml")
	if reflect.TypeOf(res) != reflect.TypeOf(&Config{}) {
		t.Error("MustLoad() error = The function returns the wrong type")
	}

	f, _ := os.Create("./test.yaml")
	f.Write([]byte("bla-bla-bla"))
	f.Close()

	res = MustLoad("./test.yaml")
	if reflect.TypeOf(res) != reflect.TypeOf(&Config{}) {
		t.Error("MustLoad() error = The function returns the wrong type")
	}

	os.Remove("./test.yaml")
}
