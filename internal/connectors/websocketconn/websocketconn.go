package websocketconn

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/Kwynto/GracefulDB/internal/analyzers/vqlanalyzer"
	"github.com/Kwynto/GracefulDB/internal/config"
	"github.com/Kwynto/GracefulDB/pkg/lib/closer"
	"github.com/Kwynto/GracefulDB/pkg/lib/ordinarylogger"

	"github.com/gorilla/websocket"
)

type tVQuery struct {
	Ticket      string   `json:"ticket"`
	Instruction string   `json:"instruction"`
	Placeholder []string `json:"placeholder"`
}

var sAddress string
var stMuxWS *http.ServeMux

var stSrvWS *http.Server

var stConf config.TWebSocketConnector

func home(w http.ResponseWriter, r *http.Request) {
	http.NotFound(w, r)
}

func query(w http.ResponseWriter, r *http.Request) {
	var stMsgSQuery *tVQuery

	var stUpgrader = websocket.Upgrader{
		ReadBufferSize:  stConf.BufferSize.Read,
		WriteBufferSize: stConf.BufferSize.Write,
	}

	stWebSocket, err := stUpgrader.Upgrade(w, r, nil)
	if err != nil {
		slog.Error("Failed to create connection", slog.String("err", err.Error()))
		return
	}
	slog.Debug("Websocket Connected! - SQuery")

	for {
		// read a message
		iMessageType, slBMessageContent, err := stWebSocket.ReadMessage()
		if err != nil {
			slog.Debug("Error reading the message", slog.String("err", err.Error()))
			return
		}

		// Data processing

		if err := json.Unmarshal(slBMessageContent, &stMsgSQuery); err != nil {
			slog.Debug("Query error", slog.String("err", err.Error()))
			stWebSocket.WriteMessage(iMessageType, []byte("Bad request - query error."))
			stWebSocket.Close()
			return
		}

		// reponse message
		sMessageResponse := vqlanalyzer.Request(stMsgSQuery.Ticket, stMsgSQuery.Instruction, stMsgSQuery.Placeholder)

		if err := stWebSocket.WriteMessage(iMessageType, []byte(sMessageResponse)); err != nil {
			slog.Debug("Error sending response", slog.String("err", err.Error()))
			stWebSocket.Close()
			return
		}
	}
}

func routes() *http.ServeMux {
	stMux := http.NewServeMux()
	stMux.HandleFunc("/", home)
	stMux.HandleFunc("/query", query)

	return stMux
}

func Start(cfg *config.TConfig) {
	stConf = cfg.WebSocketConnector

	sAddress = fmt.Sprintf("%s:%s", stConf.Address, stConf.Port)
	stMuxWS = routes()

	stSrvWS = &http.Server{
		Addr:     sAddress,
		ErrorLog: ordinarylogger.LogServerError,
		Handler:  stMuxWS,
	}

	slog.Info("WebSocket server is running", slog.String("address", sAddress))
	if err := stSrvWS.ListenAndServe(); err != nil {
		slog.Debug(err.Error())
		return
	}
}

func Shutdown(ctx context.Context, c *closer.TCloser) {
	if err := stSrvWS.Shutdown(ctx); err != nil {
		sMsg := fmt.Sprintf("There was a problem with stopping the WebSocket-server: %s", err.Error())
		c.AddMsg(sMsg)
	}
	slog.Info("WebSocket server stopped")
	c.Done()
}
