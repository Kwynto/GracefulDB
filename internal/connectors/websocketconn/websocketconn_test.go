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
	t.Run("squery() function testing - negative", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/", nil)

		squery(w, r) // calling the tested function
		status := w.Code
		if status != http.StatusBadRequest {
			t.Errorf("squery() error. Status: %v", status)
		}
	})

	t.Run("squery() function testing #1", func(t *testing.T) {
		// Create test server with the echo handler.
		s := httptest.NewServer(http.HandlerFunc(squery))
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
			t.Errorf("squery() error: %v", err)
		}
		_, p, err := ws.ReadMessage()
		if err != nil {
			t.Errorf("squery() error: %v", err)
		}
		if len(p) == 0 {
			t.Error("squery() error.")
		}
	})

	t.Run("squery() function testing #2", func(t *testing.T) {
		// Create test server with the echo handler.
		s := httptest.NewServer(http.HandlerFunc(squery))
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
			t.Errorf("squery() error: %v", err)
		}
		_, p, err := ws.ReadMessage()
		if err != nil {
			t.Errorf("squery() error: %v", err)
		}
		if len(p) == 0 {
			t.Error("squery() error.")
		}
	})

	t.Run("squery() function testing #3", func(t *testing.T) {
		// Create test server with the echo handler.
		s := httptest.NewServer(http.HandlerFunc(squery))
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
			t.Errorf("squery() error: %v", err)
		}
		ws.Close()
		_, p, err := ws.ReadMessage()
		if err == nil {
			t.Errorf("squery() error: %v", err)
		}
		if len(p) != 0 {
			t.Error("squery() error.")
		}
	})
}

func Test_vquery(t *testing.T) {
	t.Run("vquery() function testing - negative", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/", nil)

		vquery(w, r) // calling the tested function
		status := w.Code
		if status != http.StatusBadRequest {
			t.Errorf("vquery() error. Status: %v", status)
		}
	})

	t.Run("vquery() function testing #1", func(t *testing.T) {
		// Create test server with the echo handler.
		s := httptest.NewServer(http.HandlerFunc(vquery))
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
			t.Errorf("vquery() error: %v", err)
		}
		_, p, err := ws.ReadMessage()
		if err != nil {
			t.Errorf("vquery() error: %v", err)
		}
		if len(p) == 0 {
			t.Error("vquery() error.")
		}
	})

	t.Run("vquery() function testing #2", func(t *testing.T) {
		// Create test server with the echo handler.
		s := httptest.NewServer(http.HandlerFunc(vquery))
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
			t.Errorf("vquery() error: %v", err)
		}
		_, p, err := ws.ReadMessage()
		if err != nil {
			t.Errorf("vquery() error: %v", err)
		}
		if len(p) == 0 {
			t.Error("vquery() error.")
		}
	})

	t.Run("vquery() function testing #3", func(t *testing.T) {
		// Create test server with the echo handler.
		s := httptest.NewServer(http.HandlerFunc(vquery))
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
			t.Errorf("vquery() error: %v", err)
		}
		ws.Close()
		_, p, err := ws.ReadMessage()
		if err == nil {
			t.Errorf("vquery() error: %v", err)
		}
		if len(p) != 0 {
			t.Error("vquery() error.")
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
		go Start(&config.DefaultConfig) // calling the tested function
		closer.AddHandler(Shutdown)
		time.Sleep(2 * time.Second)
		// srvRest.Shutdown(context.Background())
		Shutdown(context.Background(), closer.CloseProcs)

		if reflect.TypeOf(muxWS) != reflect.TypeOf(&http.ServeMux{}) {
			t.Error("Start() error = The function has created an incorrect dependency.")
		}

		if reflect.TypeOf(srvWS) != reflect.TypeOf(&http.Server{}) {
			t.Error("Start() error = The function has created an incorrect dependency.")
		}

		if closer.CloseProcs.Counter != 0 {
			t.Errorf("Shutdown() error: %v.", closer.CloseProcs.Counter)
		}
	})

	t.Run("Shutdown() function testing - positive", func(t *testing.T) {
		Shutdown(context.Background(), closer.CloseProcs)

		if len(closer.CloseProcs.Msgs) > 0 {
			t.Errorf("Shutdown() error.")
		}
	})
}
