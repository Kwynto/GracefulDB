package rest

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
)

var sAddress string
var stMuxRest *http.ServeMux

var stSrvRest *http.Server

func home(w http.ResponseWriter, r *http.Request) {
	http.NotFound(w, r)
}

func query(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		slog.Debug("The method is prohibited!", slog.String("method", http.MethodPost))
		w.Header().Set("Allow", http.MethodPost)
		http.Error(w, "The method is prohibited!", http.StatusMethodNotAllowed)
		return
	}

	r.ParseForm()
	var slPlaceholder []string
	sTicket := r.PostForm.Get("ticket")
	sInstruction := r.PostForm.Get("instruction")
	jPlaceholderJSONArray := r.PostForm.Get("placeholder")
	if err := json.Unmarshal([]byte(jPlaceholderJSONArray), &slPlaceholder); err != nil {
		slog.Debug("Placeholder error", slog.String("err", err.Error()))
		http.Error(w, "Bad request - placeholder error (The placeholder must be in JSON format, in the form of an array of strings).", http.StatusBadRequest)
		return
	}

	sResponse := vqlanalyzer.Request(sTicket, sInstruction, slPlaceholder)

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(sResponse))
}

func routes() *http.ServeMux {
	stMux := http.NewServeMux()
	stMux.HandleFunc("/", home)
	stMux.HandleFunc("/query", query)

	return stMux
}

func Start(cfg *config.TConfig) {
	sAddress = fmt.Sprintf("%s:%s", cfg.RestConnector.Address, cfg.RestConnector.Port)
	stMuxRest = routes()

	stSrvRest = &http.Server{
		Addr:     sAddress,
		ErrorLog: ordinarylogger.LogServerError,
		Handler:  stMuxRest,
	}

	slog.Info("REST server is running", slog.String("address", sAddress))
	if err := stSrvRest.ListenAndServe(); err != nil {
		slog.Debug(err.Error())
		return
	}
}

func Shutdown(ctx context.Context, c *closer.TCloser) {
	if err := stSrvRest.Shutdown(ctx); err != nil {
		sMsg := fmt.Sprintf("There was a problem with stopping the REST-server: %s", err.Error())
		c.AddMsg(sMsg)
	}
	slog.Info("REST server stopped")
	c.Done()
}
