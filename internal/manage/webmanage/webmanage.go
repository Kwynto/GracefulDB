package webmanage

import (
	"context"
	"fmt"
	"html/template"
	"log/slog"
	"net/http"

	"github.com/Kwynto/gosession"

	"github.com/Kwynto/GracefulDB/internal/config"
	"github.com/Kwynto/GracefulDB/pkg/lib/closer"
	"github.com/Kwynto/GracefulDB/pkg/lib/prettylogger"
)

var address string
var muxWeb *http.ServeMux

var srvWeb *http.Server

func homeDefault(w http.ResponseWriter, r *http.Request, login string) {
	ts, err := template.ParseFiles("./ui/html/home.html")
	if err != nil {
		slog.Debug("Internal Server Error", slog.String("err", err.Error()))
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	err = ts.Execute(w, nil)
	if err != nil {
		slog.Debug("Internal Server Error", slog.String("err", err.Error()))
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func homeAuth(w http.ResponseWriter, r *http.Request) {
	ts, err := template.ParseFiles("./ui/html/auth.html")
	if err != nil {
		slog.Debug("Internal Server Error", slog.String("err", err.Error()))
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	err = ts.Execute(w, nil)
	if err != nil {
		slog.Debug("Internal Server Error", slog.String("err", err.Error()))
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	sesID := gosession.Start(&w, r)
	auth := sesID.Get("auth")
	if auth == nil {
		homeAuth(w, r)
	} else {
		login := fmt.Sprint(auth)
		homeDefault(w, r, login)
	}
}

func routes() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", home)

	return mux
}

func Start(cfg *config.Config) {
	address = fmt.Sprintf("%s:%s", cfg.WebServer.Address, cfg.WebServer.Port)
	muxWeb = routes()

	srvWeb = &http.Server{
		Addr:     address,
		ErrorLog: prettylogger.LogServerError,
		Handler:  muxWeb,
	}

	slog.Info("Web manager is running", slog.String("address", address))
	if err := srvWeb.ListenAndServe(); err != nil {
		slog.Debug(err.Error())
		return
	}
}

func Shutdown(ctx context.Context, c *closer.Closer) {
	if err := srvWeb.Shutdown(ctx); err != nil {
		slog.Error("There was a problem with stopping the Web manager", slog.String("err", err.Error()))
	}
	slog.Info("Web manager stopped")
	c.Done()
}
