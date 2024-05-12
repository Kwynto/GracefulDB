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

var sAddress string
var stMuxWeb *http.ServeMux

var stSrvWeb *http.Server

func routes() *http.ServeMux {
	// Main routes
	stMux := http.NewServeMux()
	stMux.HandleFunc("/", fnHome)
	stMux.HandleFunc("/log.out", fnLogout)

	// HTMX routes
	stMux.HandleFunc("/hx/", fnNavDefault)
	stMux.HandleFunc("/hx/nav/logout", nav_logout)

	stMux.HandleFunc("/hx/nav/dashboard", nav_dashboard)

	stMux.HandleFunc("/hx/nav/databases", nav_databases)
	// mux.HandleFunc("/hx/databases/request", database_request)

	stMux.HandleFunc("/hx/nav/console", nav_console)
	stMux.HandleFunc("/hx/console/request", console_request)

	stMux.HandleFunc("/hx/nav/accounts", nav_accounts)
	stMux.HandleFunc("/hx/accounts/create_load_form", account_create_load_form)
	stMux.HandleFunc("/hx/accounts/create_ok", account_create_ok)
	stMux.HandleFunc("/hx/accounts/edit_load_form", account_edit_load_form)
	stMux.HandleFunc("/hx/accounts/edit_ok", account_edit_ok)
	stMux.HandleFunc("/hx/accounts/ban_load_form", account_ban_load_form)
	stMux.HandleFunc("/hx/accounts/ban_ok", account_ban_ok)
	stMux.HandleFunc("/hx/accounts/unban_load_form", account_unban_load_form)
	stMux.HandleFunc("/hx/accounts/unban_ok", account_unban_ok)
	stMux.HandleFunc("/hx/accounts/del_load_form", account_del_load_form)
	stMux.HandleFunc("/hx/accounts/del_ok", account_del_ok)
	stMux.HandleFunc("/hx/accounts/selfedit_load_form", selfedit_load_form)
	stMux.HandleFunc("/hx/accounts/selfedit_ok", selfedit_ok)

	stMux.HandleFunc("/hx/nav/settings", nav_settings)
	stMux.HandleFunc("/hx/settings/core_friendly_change_sw", settings_core_friendly_change_sw)
	stMux.HandleFunc("/hx/settings/wsc_change_sw", settings_wsc_change_sw)
	stMux.HandleFunc("/hx/settings/rest_change_sw", settings_rest_change_sw)
	stMux.HandleFunc("/hx/settings/grpc_change_sw", settings_grpc_change_sw)
	stMux.HandleFunc("/hx/settings/web_change_sw", settings_web_change_sw)

	// Isolation of static files
	// fileServer := http.FileServer(IsolatedFS{http.Dir("./ui/static/")})
	// mux.Handle("/static", http.NotFoundHandler())
	// mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	// Embed of static file
	inFileServer := http.FileServer(http.FS(emStaticDir))
	stMux.Handle("/ui/static", http.NotFoundHandler())
	stMux.Handle("/ui/static/", inFileServer)

	return stMux
}

func Start(stCfg *config.TConfig) {
	// This function is completes
	parseTemplates()

	sAddress = fmt.Sprintf("%s:%s", stCfg.WebServer.Address, stCfg.WebServer.Port)
	stMuxWeb = routes()

	stSrvWeb = &http.Server{
		Addr:     sAddress,
		ErrorLog: ordinarylogger.LogServerError,
		Handler:  stMuxWeb,
	}

	slog.Info("Web manager is running", slog.String("address", sAddress))
	if err := stSrvWeb.ListenAndServe(); err != nil {
		slog.Debug(err.Error())
		return
	}
}

func Shutdown(ctx context.Context, c *closer.TCloser) {
	// This function is complete
	if err := stSrvWeb.Shutdown(ctx); err != nil {
		sMsg := fmt.Sprintf("There was a problem with stopping the Web manager: %s", err.Error())
		c.AddMsg(sMsg)
	}
	slog.Info("Web manager stopped")
	c.Done()
}
