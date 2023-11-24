package webmanage

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/Kwynto/GracefulDB/internal/config"

	"github.com/Kwynto/GracefulDB/pkg/lib/closer"
	"github.com/Kwynto/GracefulDB/pkg/lib/prettylogger"
)

var address string
var muxWeb *http.ServeMux

var srvWeb *http.Server

func routes() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", home)

	fileServer := http.FileServer(isolatedFS{http.Dir("./ui/static/")})
	mux.Handle("/static", http.NotFoundHandler())
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	// fileServer := http.FileServer(http.Dir("./ui/static/"))
	// mux.Handle("/static/", http.StripPrefix("/static", fileServer))

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
		// slog.Error("There was a problem with stopping the Web manager", slog.String("err", err.Error()))
		msg := fmt.Sprintf("There was a problem with stopping the Web manager: %s", err.Error())
		c.AddMsg(msg)
	}
	slog.Info("Web manager stopped")
	c.Done()
}
