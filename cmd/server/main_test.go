package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

var appName string = "testapp.exe"

func Test_main(t *testing.T) {
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	exPath := filepath.Dir(ex)
	appName = fmt.Sprintf("%s%s", exPath, appName)

	fmt.Println("-> Building...")
	build := exec.Command("go", "build", "-o", appName)
	if err := build.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error building %s: %s", appName, err)
		os.Exit(1)
	}

	// Running
	fmt.Println("-> Running...")
	t.Run("Testing main() function", func(t *testing.T) {
		command := exec.Command(appName)
		if err := command.Run(); err != nil {
			t.Errorf("main() error: %v", err)
		}
	})

	fmt.Println("-> Getting done...")

	os.Remove(appName)
}
