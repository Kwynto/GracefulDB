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
	"github.com/Kwynto/GracefulDB/pkg/lib/helpers/masquerade/htmx_masq"
)

const (
	// The names of the templates in the cache
	HOME_TEMP_NAME = "home.html"
	AUTH_TEMP_NAME = "auth.html"

	BLOCK_TEMP_DEFAULT                   = "Default"
	BLOCK_TEMP_DASHBOARD                 = "Dashboard"
	BLOCK_TEMP_DATABASES                 = "Databases"
	BLOCK_TEMP_ACCOUNTS                  = "Accounts"
	BLOCK_TEMP_ACCOUNT_CREATE_FORM_OK    = "AccountCreateFormOk"
	BLOCK_TEMP_ACCOUNT_CREATE_FORM_LOAD  = "AccountCreateFormLoad"
	BLOCK_TEMP_ACCOUNT_CREATE_FORM_ERROR = "AccountCreateFormError"
	BLOCK_TEMP_ACCOUNT_EDIT_FORM         = "AccountEditForm"
	BLOCK_TEMP_ACCOUNT_BAN_FORM_OK       = "AccountBanFormOk"
	BLOCK_TEMP_ACCOUNT_BAN_FORM_LOAD     = "AccountBanFormLoad"
	BLOCK_TEMP_ACCOUNT_BAN_FORM_ERROR    = "AccountBanFormError"
	BLOCK_TEMP_ACCOUNT_UNBAN_FORM_OK     = "AccountUnBanFormOk"
	BLOCK_TEMP_ACCOUNT_UNBAN_FORM_LOAD   = "AccountUnBanFormLoad"
	BLOCK_TEMP_ACCOUNT_UNBAN_FORM_ERROR  = "AccountUnBanFormError"
	BLOCK_TEMP_ACCOUNT_DEL_FORM_OK       = "AccountDelFormOk"
	BLOCK_TEMP_ACCOUNT_DEL_FORM_LOAD     = "AccountDelFormLoad"
	BLOCK_TEMP_ACCOUNT_DEL_FORM_ERROR    = "AccountDelFormError"
	BLOCK_TEMP_SETTINGS                  = "Settings"
)

var address string
var muxWeb *http.ServeMux

var srvWeb *http.Server

var TemplatesMap = make(map[string]*template.Template)

func parseTemplates() {
	ts, err := template.New(HOME_TEMP_NAME).Parse(home_masq.HtmlHome)
	if err != nil {
		slog.Debug("Error reading the template", slog.String("err", err.Error()))
		return
	}
	TemplatesMap[HOME_TEMP_NAME] = ts

	ts, err = template.New(AUTH_TEMP_NAME).Parse(auth_masq.HtmlAuth)
	if err != nil {
		slog.Debug("Error reading the template", slog.String("err", err.Error()))
		return
	}
	TemplatesMap[AUTH_TEMP_NAME] = ts

	ts, err = template.New(BLOCK_TEMP_DEFAULT).Parse(htmx_masq.Default)
	if err != nil {
		slog.Debug("Error reading the template", slog.String("err", err.Error()))
		return
	}
	TemplatesMap[BLOCK_TEMP_DEFAULT] = ts

	TemplatesMap[BLOCK_TEMP_DASHBOARD] = ts
	ts, err = template.New(BLOCK_TEMP_DASHBOARD).Parse(htmx_masq.Dashboard)
	if err != nil {
		slog.Debug("Error reading the template", slog.String("err", err.Error()))
		return
	}
	TemplatesMap[BLOCK_TEMP_DASHBOARD] = ts

	ts, err = template.New(BLOCK_TEMP_DATABASES).Parse(htmx_masq.Databases)
	if err != nil {
		slog.Debug("Error reading the template", slog.String("err", err.Error()))
		return
	}
	TemplatesMap[BLOCK_TEMP_DATABASES] = ts

	ts, err = template.New(BLOCK_TEMP_ACCOUNTS).Parse(htmx_masq.Accounts)
	if err != nil {
		slog.Debug("Error reading the template", slog.String("err", err.Error()))
		return
	}
	TemplatesMap[BLOCK_TEMP_ACCOUNTS] = ts

	ts, err = template.New(BLOCK_TEMP_ACCOUNT_CREATE_FORM_OK).Parse(htmx_masq.AccountCreateFormOk)
	if err != nil {
		slog.Debug("Error reading the template", slog.String("err", err.Error()))
		return
	}
	TemplatesMap[BLOCK_TEMP_ACCOUNT_CREATE_FORM_OK] = ts

	ts, err = template.New(BLOCK_TEMP_ACCOUNT_CREATE_FORM_LOAD).Parse(htmx_masq.AccountCreateFormLoad)
	if err != nil {
		slog.Debug("Error reading the template", slog.String("err", err.Error()))
		return
	}
	TemplatesMap[BLOCK_TEMP_ACCOUNT_CREATE_FORM_LOAD] = ts

	ts, err = template.New(BLOCK_TEMP_ACCOUNT_CREATE_FORM_ERROR).Parse(htmx_masq.AccountCreateFormError)
	if err != nil {
		slog.Debug("Error reading the template", slog.String("err", err.Error()))
		return
	}
	TemplatesMap[BLOCK_TEMP_ACCOUNT_CREATE_FORM_ERROR] = ts

	ts, err = template.New(BLOCK_TEMP_ACCOUNT_EDIT_FORM).Parse(htmx_masq.AccountEditForm)
	if err != nil {
		slog.Debug("Error reading the template", slog.String("err", err.Error()))
		return
	}
	TemplatesMap[BLOCK_TEMP_ACCOUNT_EDIT_FORM] = ts

	ts, err = template.New(BLOCK_TEMP_ACCOUNT_BAN_FORM_OK).Parse(htmx_masq.AccountBanFormOk)
	if err != nil {
		slog.Debug("Error reading the template", slog.String("err", err.Error()))
		return
	}
	TemplatesMap[BLOCK_TEMP_ACCOUNT_BAN_FORM_OK] = ts

	ts, err = template.New(BLOCK_TEMP_ACCOUNT_BAN_FORM_LOAD).Parse(htmx_masq.AccountBanFormLoad)
	if err != nil {
		slog.Debug("Error reading the template", slog.String("err", err.Error()))
		return
	}
	TemplatesMap[BLOCK_TEMP_ACCOUNT_BAN_FORM_LOAD] = ts

	ts, err = template.New(BLOCK_TEMP_ACCOUNT_BAN_FORM_ERROR).Parse(htmx_masq.AccountBanFormError)
	if err != nil {
		slog.Debug("Error reading the template", slog.String("err", err.Error()))
		return
	}
	TemplatesMap[BLOCK_TEMP_ACCOUNT_BAN_FORM_ERROR] = ts

	ts, err = template.New(BLOCK_TEMP_ACCOUNT_UNBAN_FORM_OK).Parse(htmx_masq.AccountUnBanFormOk)
	if err != nil {
		slog.Debug("Error reading the template", slog.String("err", err.Error()))
		return
	}
	TemplatesMap[BLOCK_TEMP_ACCOUNT_UNBAN_FORM_OK] = ts

	ts, err = template.New(BLOCK_TEMP_ACCOUNT_UNBAN_FORM_LOAD).Parse(htmx_masq.AccountUnBanFormLoad)
	if err != nil {
		slog.Debug("Error reading the template", slog.String("err", err.Error()))
		return
	}
	TemplatesMap[BLOCK_TEMP_ACCOUNT_UNBAN_FORM_LOAD] = ts

	ts, err = template.New(BLOCK_TEMP_ACCOUNT_UNBAN_FORM_ERROR).Parse(htmx_masq.AccountUnBanFormError)
	if err != nil {
		slog.Debug("Error reading the template", slog.String("err", err.Error()))
		return
	}
	TemplatesMap[BLOCK_TEMP_ACCOUNT_UNBAN_FORM_ERROR] = ts

	ts, err = template.New(BLOCK_TEMP_ACCOUNT_DEL_FORM_OK).Parse(htmx_masq.AccountDelFormOk)
	if err != nil {
		slog.Debug("Error reading the template", slog.String("err", err.Error()))
		return
	}
	TemplatesMap[BLOCK_TEMP_ACCOUNT_DEL_FORM_OK] = ts

	ts, err = template.New(BLOCK_TEMP_ACCOUNT_DEL_FORM_LOAD).Parse(htmx_masq.AccountDelFormLoad)
	if err != nil {
		slog.Debug("Error reading the template", slog.String("err", err.Error()))
		return
	}
	TemplatesMap[BLOCK_TEMP_ACCOUNT_DEL_FORM_LOAD] = ts

	ts, err = template.New(BLOCK_TEMP_ACCOUNT_DEL_FORM_ERROR).Parse(htmx_masq.AccountDelFormError)
	if err != nil {
		slog.Debug("Error reading the template", slog.String("err", err.Error()))
		return
	}
	TemplatesMap[BLOCK_TEMP_ACCOUNT_DEL_FORM_ERROR] = ts

	ts, err = template.New(BLOCK_TEMP_SETTINGS).Parse(htmx_masq.Settings)
	if err != nil {
		slog.Debug("Error reading the template", slog.String("err", err.Error()))
		return
	}
	TemplatesMap[BLOCK_TEMP_SETTINGS] = ts

}

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

	mux.HandleFunc("/hx/nav/accounts", nav_accounts)
	mux.HandleFunc("/hx/accounts/create_load_form", account_create_load_form)
	mux.HandleFunc("/hx/accounts/create_ok", account_create_ok)
	mux.HandleFunc("/hx/accounts/edit_form", account_edit_form)
	mux.HandleFunc("/hx/accounts/ban_load_form", account_ban_load_form)
	mux.HandleFunc("/hx/accounts/ban_ok", account_ban_ok)
	mux.HandleFunc("/hx/accounts/unban_load_form", account_unban_load_form)
	mux.HandleFunc("/hx/accounts/unban_ok", account_unban_ok)
	mux.HandleFunc("/hx/accounts/del_load_form", account_del_load_form)
	mux.HandleFunc("/hx/accounts/del_ok", account_del_ok)

	mux.HandleFunc("/hx/nav/settings", nav_settings)
	mux.HandleFunc("/hx/settings/wsc_change_sw", settings_wsc_change_sw)
	mux.HandleFunc("/hx/settings/rest_change_sw", settings_rest_change_sw)
	mux.HandleFunc("/hx/settings/grpc_change_sw", settings_grpc_change_sw)
	mux.HandleFunc("/hx/settings/web_change_sw", settings_web_change_sw)

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
