package rest

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/Kwynto/GracefulDB/internal/analyzers/sqlanalyzer"
	"github.com/Kwynto/GracefulDB/internal/config"
	"github.com/Kwynto/GracefulDB/pkg/lib/closer"
	"github.com/Kwynto/GracefulDB/pkg/lib/prettylogger"
)

var address string
var muxRest *http.ServeMux

var srvRest *http.Server

func home(w http.ResponseWriter, r *http.Request) {
	http.NotFound(w, r)
}

func squery(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		slog.Debug("The method is prohibited!", slog.String("method", http.MethodPost))
		w.Header().Set("Allow", http.MethodPost)
		http.Error(w, "The method is prohibited!", http.StatusMethodNotAllowed)
		return
	}

	r.ParseForm()
	var placeholder []string
	ticket := r.PostForm.Get("ticket")
	instruction := r.PostForm.Get("instruction")
	placeholderJSONArray := r.PostForm.Get("placeholder")
	if err := json.Unmarshal([]byte(placeholderJSONArray), &placeholder); err != nil {
		slog.Debug("Placeholder error", slog.String("err", err.Error()))
		http.Error(w, "Bad request - placeholder error (The placeholder must be in JSON format, in the form of an array of strings).", http.StatusBadRequest)
		return
	}

	response := sqlanalyzer.Request(ticket, instruction, placeholder)

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(response))
}

func routes() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", home)
	mux.HandleFunc("/squery", squery)

	return mux
}

func Start(cfg *config.Config) {
	address = fmt.Sprintf("%s:%s", cfg.RestConnector.Address, cfg.RestConnector.Port)
	muxRest = routes()

	srvRest = &http.Server{
		Addr:     address,
		ErrorLog: prettylogger.LogServerError,
		Handler:  muxRest,
	}

	slog.Info("REST server is running", slog.String("address", address))
	if err := srvRest.ListenAndServe(); err != nil {
		slog.Debug(err.Error())
		return
	}
}

func Shutdown(ctx context.Context, c *closer.Closer) {
	if err := srvRest.Shutdown(ctx); err != nil {
		msg := fmt.Sprintf("There was a problem with stopping the REST-server: %s", err.Error())
		c.AddMsg(msg)
	}
	slog.Info("REST server stopped")
	c.Done()
}
