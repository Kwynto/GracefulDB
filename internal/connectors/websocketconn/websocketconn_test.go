package websocketconn

import (
	"context"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/Kwynto/GracefulDB/internal/config"
	"github.com/Kwynto/GracefulDB/pkg/lib/closer"
	"github.com/gorilla/websocket"
)

func Test_home(t *testing.T) {
	t.Run("home() function testing - positive", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/", nil)

		home(w, r) // calling the tested function
		status := w.Code
		if status != http.StatusNotFound {
			t.Error("home() error.")
		}
	})

	t.Run("home() function testing - negative", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/", nil)

		home(w, r) // calling the tested function
		status := w.Code
		if status == http.StatusOK {
			t.Error("home() error.")
		}
	})
}

func Test_squery(t *testing.T) {
	t.Run("query() function testing - negative", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/", nil)

		query(w, r) // calling the tested function
		status := w.Code
		if status != http.StatusBadRequest {
			t.Errorf("query() error. Status: %v", status)
		}
	})

	t.Run("query() function testing #1", func(t *testing.T) {
		// Create test server with the echo handler.
		s := httptest.NewServer(http.HandlerFunc(query))
		defer s.Close()
		// Convert http://127.0.0.1 to ws://127.0.0.
		u := "ws" + strings.TrimPrefix(s.URL, "http")
		// Connect to the server
		ws, _, err := websocket.DefaultDialer.Dial(u, nil)
		if err != nil {
			t.Fatalf("%v", err)
		}
		defer ws.Close()

		// Send message to server, read response and check to see if it's what we expect.

		if err := ws.WriteMessage(websocket.TextMessage, []byte("hello")); err != nil {
			t.Errorf("query() error: %v", err)
		}
		_, p, err := ws.ReadMessage()
		if err != nil {
			t.Errorf("query() error: %v", err)
		}
		if len(p) == 0 {
			t.Error("query() error.")
		}
	})

	t.Run("query() function testing #2", func(t *testing.T) {
		// Create test server with the echo handler.
		s := httptest.NewServer(http.HandlerFunc(query))
		defer s.Close()
		// Convert http://127.0.0.1 to ws://127.0.0.
		u := "ws" + strings.TrimPrefix(s.URL, "http")
		// Connect to the server
		ws, _, err := websocket.DefaultDialer.Dial(u, nil)
		if err != nil {
			t.Fatalf("%v", err)
		}
		defer ws.Close()

		// Send message to server, read response and check to see if it's what we expect.

		if err := ws.WriteMessage(websocket.TextMessage, []byte("{\"hello\": \"hello\"}")); err != nil {
			t.Errorf("query() error: %v", err)
		}
		_, p, err := ws.ReadMessage()
		if err != nil {
			t.Errorf("query() error: %v", err)
		}
		if len(p) == 0 {
			t.Error("query() error.")
		}
	})

	t.Run("query() function testing #3", func(t *testing.T) {
		// Create test server with the echo handler.
		s := httptest.NewServer(http.HandlerFunc(query))
		defer s.Close()
		// Convert http://127.0.0.1 to ws://127.0.0.
		u := "ws" + strings.TrimPrefix(s.URL, "http")
		// Connect to the server
		ws, _, err := websocket.DefaultDialer.Dial(u, nil)
		if err != nil {
			t.Fatalf("%v", err)
		}
		// defer ws.Close()

		// Send message to server, read response and check to see if it's what we expect.

		if err := ws.WriteMessage(websocket.TextMessage, []byte("{\"hello\": \"hello\"}")); err != nil {
			t.Errorf("query() error: %v", err)
		}
		ws.Close()
		_, p, err := ws.ReadMessage()
		if err == nil {
			t.Errorf("query() error: %v", err)
		}
		if len(p) != 0 {
			t.Error("query() error.")
		}
	})
}

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

		if reflect.TypeOf(stMuxWS) != reflect.TypeOf(&http.ServeMux{}) {
			t.Error("Start() error = The function has created an incorrect dependency.")
		}

		if reflect.TypeOf(stSrvWS) != reflect.TypeOf(&http.Server{}) {
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
