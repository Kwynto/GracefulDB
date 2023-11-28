package webmanage

import (
	"context"
	"fmt"
	"html/template"
	"log/slog"
	"net/http"

	"github.com/Kwynto/GracefulDB/internal/config"

	"github.com/Kwynto/GracefulDB/pkg/lib/closer"
	"github.com/Kwynto/GracefulDB/pkg/lib/prettylogger"

	"github.com/Kwynto/GracefulDB/pkg/lib/helpers/masquerade/auth_masq"
	"github.com/Kwynto/GracefulDB/pkg/lib/helpers/masquerade/home_masq"
)

const (
	HOME_TEMP_NAME = "home.html"
	AUTH_TEMP_NAME = "auth.html"
)

var address string
var muxWeb *http.ServeMux

var srvWeb *http.Server

var templatesMap = make(map[string]*template.Template)

func parseTemplates() {
	// ts, err := template.ParseFiles("./ui/html/home.html")
	ts, err := template.New(HOME_TEMP_NAME).Parse(home_masq.HtmlHome)
	if err != nil {
		slog.Debug("Internal Server Error", slog.String("err", err.Error()))
		// http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	templatesMap[HOME_TEMP_NAME] = ts

	// ts, err := template.ParseFiles("./ui/html/auth.html")
	ts, err = template.New(AUTH_TEMP_NAME).Parse(auth_masq.HtmlAuth)
	if err != nil {
		slog.Debug("Internal Server Error", slog.String("err", err.Error()))
		// http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	templatesMap[AUTH_TEMP_NAME] = ts
}

func routes() *http.ServeMux {
	// Main routes
	mux := http.NewServeMux()
	mux.HandleFunc("/", home)
	mux.HandleFunc("/log.out", logout)

	// HTMX routes
	mux.HandleFunc("/hx/firstmsg", firstmsg)
	mux.HandleFunc("/hx/mainunit", mainunit)

	// Isolation of static files
	fileServer := http.FileServer(isolatedFS{http.Dir("./ui/static/")})
	mux.Handle("/static", http.NotFoundHandler())
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	return mux
}

func Start(cfg *config.Config) {
	// This function is completes
	parseTemplates()

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
	// This function is complete
	if err := srvWeb.Shutdown(ctx); err != nil {
		// slog.Error("There was a problem with stopping the Web manager", slog.String("err", err.Error()))
		msg := fmt.Sprintf("There was a problem with stopping the Web manager: %s", err.Error())
		c.AddMsg(msg)
	}
	slog.Info("Web manager stopped")
	c.Done()
}
