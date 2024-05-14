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
	stMux.HandleFunc("/hx/nav/logout", fnNavLogout)

	stMux.HandleFunc("/hx/nav/dashboard", fnNavDashboard)

	stMux.HandleFunc("/hx/nav/databases", fnNavDatabases)
	// mux.HandleFunc("/hx/databases/request", database_request)

	stMux.HandleFunc("/hx/nav/console", fnNavConsole)
	stMux.HandleFunc("/hx/console/request", fnConsoleRequest)

	stMux.HandleFunc("/hx/nav/accounts", fnNavAccounts)
	stMux.HandleFunc("/hx/accounts/create_load_form", fnAccountCreateLoadForm)
	stMux.HandleFunc("/hx/accounts/create_ok", fnAccountCreateOk)
	stMux.HandleFunc("/hx/accounts/edit_load_form", fnAccountEditLoadForm)
	stMux.HandleFunc("/hx/accounts/edit_ok", fnAccountEditOk)
	stMux.HandleFunc("/hx/accounts/ban_load_form", fnAccountBanLoadForm)
	stMux.HandleFunc("/hx/accounts/ban_ok", fnAccountBanOk)
	stMux.HandleFunc("/hx/accounts/unban_load_form", fnAccountUnbanLoadForm)
	stMux.HandleFunc("/hx/accounts/unban_ok", fnAccountUnbanOk)
	stMux.HandleFunc("/hx/accounts/del_load_form", fnAccountDelLoadForm)
	stMux.HandleFunc("/hx/accounts/del_ok", fnAccountDelOk)
	stMux.HandleFunc("/hx/accounts/selfedit_load_form", fnSelfeditLoadForm)
	stMux.HandleFunc("/hx/accounts/selfedit_ok", fnSelfeditOk)

	stMux.HandleFunc("/hx/nav/settings", fnNavSettings)
	stMux.HandleFunc("/hx/settings/core_friendly_change_sw", fnSettingsCoreFriendlyChangeSw)
	stMux.HandleFunc("/hx/settings/wsc_change_sw", fnSettingsWScChangeSw)
	stMux.HandleFunc("/hx/settings/rest_change_sw", fnSettingsRestChangeSw)
	stMux.HandleFunc("/hx/settings/grpc_change_sw", fnSettingsGrpcChangeSw)
	stMux.HandleFunc("/hx/settings/web_change_sw", fnSettingsWebChangeSw)

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
