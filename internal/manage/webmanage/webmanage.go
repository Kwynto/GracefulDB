package webmanage

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/Kwynto/GracefulDB/internal/config"

	"github.com/Kwynto/GracefulDB/pkg/lib/closer"
	"github.com/Kwynto/GracefulDB/pkg/lib/ordinarylogger"
)

const (
	CONSOLE_TIME_FORMAT = "2006-01-02 15:04:05"
)

var address string
var muxWeb *http.ServeMux

var srvWeb *http.Server

func routes() *http.ServeMux {
	// Main routes
	mux := http.NewServeMux()
	mux.HandleFunc("/", home)
	mux.HandleFunc("/log.out", logout)

	// HTMX routes
	mux.HandleFunc("/hx/", nav_default)
	mux.HandleFunc("/hx/nav/logout", nav_logout)

	mux.HandleFunc("/hx/nav/dashboard", nav_dashboard)

	mux.HandleFunc("/hx/nav/databases", nav_databases)
	mux.HandleFunc("/hx/databases/request", database_request)

	mux.HandleFunc("/hx/nav/accounts", nav_accounts)
	mux.HandleFunc("/hx/accounts/create_load_form", account_create_load_form)
	mux.HandleFunc("/hx/accounts/create_ok", account_create_ok)
	mux.HandleFunc("/hx/accounts/edit_load_form", account_edit_load_form)
	mux.HandleFunc("/hx/accounts/edit_ok", account_edit_ok)
	mux.HandleFunc("/hx/accounts/ban_load_form", account_ban_load_form)
	mux.HandleFunc("/hx/accounts/ban_ok", account_ban_ok)
	mux.HandleFunc("/hx/accounts/unban_load_form", account_unban_load_form)
	mux.HandleFunc("/hx/accounts/unban_ok", account_unban_ok)
	mux.HandleFunc("/hx/accounts/del_load_form", account_del_load_form)
	mux.HandleFunc("/hx/accounts/del_ok", account_del_ok)
	mux.HandleFunc("/hx/accounts/selfedit_load_form", selfedit_load_form)
	mux.HandleFunc("/hx/accounts/selfedit_ok", selfedit_ok)

	mux.HandleFunc("/hx/nav/settings", nav_settings)
	mux.HandleFunc("/hx/settings/core_friendly_change_sw", settings_core_friendly_change_sw)
	mux.HandleFunc("/hx/settings/wsc_change_sw", settings_wsc_change_sw)
	mux.HandleFunc("/hx/settings/rest_change_sw", settings_rest_change_sw)
	mux.HandleFunc("/hx/settings/grpc_change_sw", settings_grpc_change_sw)
	mux.HandleFunc("/hx/settings/web_change_sw", settings_web_change_sw)

	// Isolation of static files
	// fileServer := http.FileServer(IsolatedFS{http.Dir("./ui/static/")})
	// mux.Handle("/static", http.NotFoundHandler())
	// mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	// Embed of static file
	fileServer := http.FileServer(http.FS(uiStaticDir))
	mux.Handle("/ui/static", http.NotFoundHandler())
	mux.Handle("/ui/static/", fileServer)

	return mux
}

func Start(cfg *config.Config) {
	// This function is completes
	parseTemplates()

	address = fmt.Sprintf("%s:%s", cfg.WebServer.Address, cfg.WebServer.Port)
	muxWeb = routes()

	srvWeb = &http.Server{
		Addr:     address,
		ErrorLog: ordinarylogger.LogServerError,
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
		msg := fmt.Sprintf("There was a problem with stopping the Web manager: %s", err.Error())
		c.AddMsg(msg)
	}
	slog.Info("Web manager stopped")
	c.Done()
}
