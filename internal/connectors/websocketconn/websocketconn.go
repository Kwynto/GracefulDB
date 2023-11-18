package websocketconn

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/Kwynto/GracefulDB/internal/analyzers/sqlanalyzer"
	"github.com/Kwynto/GracefulDB/internal/analyzers/vqlanalyzer"
	"github.com/Kwynto/GracefulDB/internal/config"
	"github.com/Kwynto/GracefulDB/pkg/lib/closer"
	"github.com/Kwynto/GracefulDB/pkg/lib/prettylogger"

	"github.com/gorilla/websocket"
)

type tSQuery struct {
	Instruction string   `json:"instruction"`
	Placeholder []string `json:"placeholder"`
}

var address string
var muxWS *http.ServeMux

var srvWS *http.Server

func home(w http.ResponseWriter, r *http.Request) {
	http.NotFound(w, r)
}

func squery(w http.ResponseWriter, r *http.Request) {
	var msgSQuery *tSQuery

	var upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	websocket, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		slog.Error("Failed to create connection", slog.String("err", err.Error()))
		return
	}
	slog.Debug("Websocket Connected! - SQuery")

	for {
		// read a message
		messageType, messageContent, err := websocket.ReadMessage()
		if err != nil {
			slog.Debug("Error reading the message", slog.String("err", err.Error()))
			return
		}

		// Data processing
		slog.Debug(string(messageContent))

		if err := json.Unmarshal(messageContent, &msgSQuery); err != nil {
			slog.Debug("Query error", slog.String("err", err.Error()))
			websocket.WriteMessage(messageType, []byte("Bad request - query error."))
			websocket.Close()
			return
		}

		// reponse message
		messageResponse := sqlanalyzer.Request(msgSQuery.Instruction, msgSQuery.Placeholder)

		if err := websocket.WriteMessage(messageType, []byte(messageResponse)); err != nil {
			slog.Debug("Error sending response", slog.String("err", err.Error()))
			return
		}
	}
}

func vquery(w http.ResponseWriter, r *http.Request) {
	var upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	websocket, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		slog.Error("Failed to create connection", slog.String("err", err.Error()))
		return
	}
	slog.Debug("Websocket Connected! - VQuery")

	for {
		// read a message
		messageType, messageContent, err := websocket.ReadMessage()
		if err != nil {
			slog.Debug("Error reading the message", slog.String("err", err.Error()))
			return
		}

		// Data processing
		slog.Debug(string(messageContent))

		// reponse message
		messageResponse := vqlanalyzer.Request(string(messageContent))

		if err := websocket.WriteMessage(messageType, []byte(messageResponse)); err != nil {
			slog.Debug("Error sending response", slog.String("err", err.Error()))
			return
		}
	}
}

func routes() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", home)
	mux.HandleFunc("/squery", squery)
	mux.HandleFunc("/vquery", vquery)

	return mux
}

func Start(cfg *config.Config) {
	address = fmt.Sprintf("%s:%s", cfg.WebSocketConnector.Address, cfg.WebSocketConnector.Port)
	muxWS = routes()

	srvWS = &http.Server{
		Addr:     address,
		ErrorLog: prettylogger.LogServerError,
		Handler:  muxWS,
	}

	slog.Info("WebSocket server is running", slog.String("address", address))
	if err := srvWS.ListenAndServe(); err != nil {
		slog.Debug(err.Error())
		return
	}
}

func Shutdown(ctx context.Context, c *closer.Closer) {
	if err := srvWS.Shutdown(ctx); err != nil {
		slog.Error("There was a problem with stopping the WebSocket-server", slog.String("err", err.Error()))
	}
	slog.Info("WebSocket server stopped")
	c.Done()
}
