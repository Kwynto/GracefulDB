package webmanage

import (
	"context"
	"net/http"
	"reflect"
	"testing"
	"time"

	"github.com/Kwynto/GracefulDB/internal/config"
	"github.com/Kwynto/GracefulDB/pkg/lib/closer"
)

func Test_routes(t *testing.T) {
	t.Run("routes() function testing", func(t *testing.T) {
		res := routes() // calling the tested function

		if reflect.TypeOf(res) != reflect.TypeOf(&http.ServeMux{}) {
			t.Error("routes() error = The function returns the wrong type")
		}
	})
}

func Test_Start_and_Shutdown(t *testing.T) {
	t.Run("Start() and Shutdown() function testing", func(t *testing.T) {
		tf := "../../../../config/develop.yaml"
		config.MustLoad(tf)
		go Start(&config.StDefaultConfig) // calling the tested function
		closer.AddHandler(Shutdown)
		time.Sleep(2 * time.Second)
		// srvRest.Shutdown(context.Background())
		Shutdown(context.Background(), closer.StCloseProcs)

		if reflect.TypeOf(muxWeb) != reflect.TypeOf(&http.ServeMux{}) {
			t.Error("Start() error = The function has created an incorrect dependency.")
		}

		if reflect.TypeOf(srvWeb) != reflect.TypeOf(&http.Server{}) {
			t.Error("Start() error = The function has created an incorrect dependency.")
		}

		if closer.StCloseProcs.Counter != 0 {
			t.Errorf("Shutdown() error: %v.", closer.StCloseProcs.Counter)
		}
	})

	t.Run("Shutdown() function testing - positive", func(t *testing.T) {
		Shutdown(context.Background(), closer.StCloseProcs)

		if len(closer.StCloseProcs.Msgs) > 0 {
			t.Errorf("Shutdown() error.")
		}
	})
}
