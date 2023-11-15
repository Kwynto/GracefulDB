package rest

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/Kwynto/GracefulDB/internal/config"
	"github.com/Kwynto/GracefulDB/pkg/lib/closer"
)

var address string
var muxRest *http.ServeMux

func home(w http.ResponseWriter, r *http.Request) {
	http.NotFound(w, r)
}

func squery(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("The method is prohibited!"))
		return
	}

	w.Write([]byte("This should be a respons for the client."))
}

func vquery(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("The method is prohibited!"))
		return
	}

	w.Write([]byte("This should be a respons for the client."))
}

func Start(cfg *config.Config) {
	address = fmt.Sprintf("%s:%s", cfg.RestConnector.Address, cfg.RestConnector.Port)

	muxRest = http.NewServeMux()
	muxRest.HandleFunc("/", home)
	muxRest.HandleFunc("/squery", squery)
	muxRest.HandleFunc("/vquery", vquery)

	slog.Info("REST server is running", slog.String("address", address))
	if err := http.ListenAndServe(address, muxRest); err != nil {
		slog.Error("Failed to start REST-listener", slog.String("err", err.Error()))
		return
	}
}

func Shutdown(ctx context.Context, c *closer.Closer) {
	slog.Info("REST server stopped")

	c.Done()
}
