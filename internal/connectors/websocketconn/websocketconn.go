package websocketconn

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/Kwynto/GracefulDB/internal/config"
	"github.com/Kwynto/GracefulDB/pkg/lib/closer"
	"github.com/Kwynto/GracefulDB/pkg/lib/prettylogger"

	"github.com/gorilla/websocket"
)

var address string
var muxWS *http.ServeMux

var srvWS *http.Server

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func listenHome(conn *websocket.Conn) {
	for {
		// read a message
		messageType, messageContent, err := conn.ReadMessage()
		if err != nil {
			slog.Debug("Error reading the message", slog.String("err", err.Error()))
			return
		}

		timeReceive := time.Now()

		// Data processing

		// reponse message
		messageResponse := fmt.Sprintf("Your message is: %s. Time received : %v", messageContent, timeReceive)

		if err := conn.WriteMessage(messageType, []byte(messageResponse)); err != nil {
			slog.Debug("Error sending response", slog.String("err", err.Error()))
			return
		}
	}
}

func home(w http.ResponseWriter, r *http.Request) {
	websocket, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		slog.Error("Failed to create connection", slog.String("err", err.Error()))
		return
	}
	slog.Debug("Websocket Connected!")
	listenHome(websocket)
}

func squery(w http.ResponseWriter, r *http.Request) {
	websocket, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		slog.Error("Failed to create connection", slog.String("err", err.Error()))
		return
	}
	slog.Debug("Websocket Connected!")
	// FIXME: листен
	listenHome(websocket)
}

func vquery(w http.ResponseWriter, r *http.Request) {
	websocket, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		slog.Error("Failed to create connection", slog.String("err", err.Error()))
		return
	}
	slog.Debug("Websocket Connected!")
	// FIXME: листен
	listenHome(websocket)
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
