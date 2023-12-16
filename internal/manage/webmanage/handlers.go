package webmanage

import (
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	"github.com/Kwynto/gosession"

	"github.com/Kwynto/GracefulDB/internal/base/basicsystem/gauth"
	"github.com/Kwynto/GracefulDB/internal/config"
	"github.com/Kwynto/GracefulDB/internal/connectors/grpc"
	"github.com/Kwynto/GracefulDB/internal/connectors/rest"
	"github.com/Kwynto/GracefulDB/internal/connectors/websocketconn"

	"github.com/Kwynto/GracefulDB/pkg/lib/closer"
)

type TViewAccountsTable struct {
	Superuser   bool
	Login       string
	Status      string
	Role        string
	Description string
}

/*
The main block
*/

// Handler after authorization
func homeDefault(w http.ResponseWriter, r *http.Request, login string) {
	// This function is complete
	err := TemplatesMap[HOME_TEMP_NAME].Execute(w, nil)
	if err != nil {
		slog.Debug("Internal Server Error", slog.String("err", err.Error()))
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

// Authorization Handler
func homeAuth(w http.ResponseWriter, r *http.Request) {
	// This function is complete
	if r.Method == http.MethodPost {
		err := r.ParseForm()
		if err != nil {
			slog.Debug("Bad request", slog.String("err", err.Error()))
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}

		username := r.PostForm.Get("username")
		password := r.PostForm.Get("password")
		isAuth := gauth.CheckUser(username, password)
		if isAuth {
			sesID := gosession.Start(&w, r)
			sesID.Set("auth", username)
		}
		http.Redirect(w, r, "/", http.StatusFound)
	} else {
		err := TemplatesMap[AUTH_TEMP_NAME].Execute(w, nil)
		if err != nil {
			slog.Debug("Internal Server Error", slog.String("err", err.Error()))
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
	}
}

// Handler for the main route
func home(w http.ResponseWriter, r *http.Request) {
	// This function is complete
	if r.URL.Path != "/" {
		// http.NotFound(w, r)
		http.Redirect(w, r, "/", http.StatusFound)
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

// Exit handler
func logout(w http.ResponseWriter, r *http.Request) {
	// This function is complete
	sesID := gosession.Start(&w, r)
	sesID.Remove("auth")
	http.Redirect(w, r, "/", http.StatusFound)
}

/*
Nav Menu Handlers
*/

func nav_default(w http.ResponseWriter, r *http.Request) {
	// This function is complete
	TemplatesMap[BLOCK_TEMP_DEFAULT].Execute(w, nil)
}

func nav_logout(w http.ResponseWriter, r *http.Request) {
	// This function is complete
	w.Header().Set("HX-Redirect", "/log.out")
}

/*
Dashboard block
*/

func nav_dashboard(w http.ResponseWriter, r *http.Request) {
	TemplatesMap[BLOCK_TEMP_DASHBOARD].Execute(w, nil)
}

/*
Databases block
*/

func nav_databases(w http.ResponseWriter, r *http.Request) {
	TemplatesMap[BLOCK_TEMP_DATABASES].Execute(w, nil)
}

/*
Accounts block
*/

func nav_accounts(w http.ResponseWriter, r *http.Request) {
	var table = make([]TViewAccountsTable, 0, 10)
	for key := range gauth.HashMap {
		element := TViewAccountsTable{
			Superuser:   false,
			Login:       key,
			Status:      gauth.AccessMap[key].Status.String(),
			Role:        gauth.AccessMap[key].Role.String(),
			Description: gauth.AccessMap[key].Description,
		}
		if key == "root" {
			element.Superuser = true
		}
		table = append(table, element)
	}

	// view := table
	// TemplatesMap[BLOCK_TEMP_ACCOUNTS].Execute(w, view)
	TemplatesMap[BLOCK_TEMP_ACCOUNTS].Execute(w, table)
}

func account_create_load_form(w http.ResponseWriter, r *http.Request) {
	TemplatesMap[BLOCK_TEMP_ACCOUNT_CREATE_FORM_LOAD].Execute(w, nil)
}

func account_create_ok(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		nav_default(w, r)
		return
	}

	err := r.ParseForm()
	if err != nil {
		slog.Debug("Bad request", slog.String("err", err.Error()))
		// http.Error(w, "Bad request", http.StatusBadRequest)
		nav_default(w, r)
		return
	}

	Login := strings.TrimSpace(r.PostForm.Get("login"))
	password := strings.TrimSpace(r.PostForm.Get("password"))
	desc := strings.TrimSpace(r.PostForm.Get("desc"))

	var data = struct {
		Login string
	}{
		Login,
	}

	if len(Login) == 0 || len(password) == 0 || len(desc) == 0 {
		TemplatesMap[BLOCK_TEMP_ACCOUNT_CREATE_FORM_ERROR].Execute(w, data)
		return
	}

	access := gauth.TRights{
		Description: desc,
		Status:      gauth.NEW,
		Role:        gauth.USER,
		Rules:       []string{},
	}

	err = gauth.AddUser(Login, password, access)
	if err != nil {
		TemplatesMap[BLOCK_TEMP_ACCOUNT_CREATE_FORM_ERROR].Execute(w, data)
		return
	}

	TemplatesMap[BLOCK_TEMP_ACCOUNT_CREATE_FORM_OK].Execute(w, data)
}

func account_edit_form(w http.ResponseWriter, r *http.Request) {
	TemplatesMap[BLOCK_TEMP_ACCOUNT_EDIT_FORM].Execute(w, nil)
}

func account_ban_load_form(w http.ResponseWriter, r *http.Request) {
	user := strings.TrimSpace(r.URL.Query().Get("user"))
	data := struct {
		Login string
	}{
		Login: user,
	}

	if user == "" || user == "root" {
		TemplatesMap[BLOCK_TEMP_ACCOUNT_BAN_FORM_ERROR].Execute(w, data)
		return
	}

	TemplatesMap[BLOCK_TEMP_ACCOUNT_BAN_FORM_LOAD].Execute(w, data)
}

func account_ban_ok(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		nav_default(w, r)
		return
	}

	err := r.ParseForm()
	if err != nil {
		slog.Debug("Bad request", slog.String("err", err.Error()))
		// http.Error(w, "Bad request", http.StatusBadRequest)
		nav_default(w, r)
		return
	}

	Login := strings.TrimSpace(r.PostForm.Get("login"))

	var data = struct {
		Login string
	}{
		Login,
	}

	if len(Login) == 0 {
		TemplatesMap[BLOCK_TEMP_ACCOUNT_BAN_FORM_ERROR].Execute(w, data)
		return
	}

	err = gauth.BlockUser(Login)
	if err != nil {
		TemplatesMap[BLOCK_TEMP_ACCOUNT_BAN_FORM_ERROR].Execute(w, data)
		return
	}

	TemplatesMap[BLOCK_TEMP_ACCOUNT_BAN_FORM_OK].Execute(w, data)
}

func account_del_load_form(w http.ResponseWriter, r *http.Request) {
	user := strings.TrimSpace(r.URL.Query().Get("user"))
	data := struct {
		Login string
	}{
		Login: user,
	}

	if user == "" || user == "root" {
		TemplatesMap[BLOCK_TEMP_ACCOUNT_DEL_FORM_ERROR].Execute(w, data)
		return
	}

	TemplatesMap[BLOCK_TEMP_ACCOUNT_DEL_FORM_LOAD].Execute(w, data)
}

func account_del_ok(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		nav_default(w, r)
		return
	}

	err := r.ParseForm()
	if err != nil {
		slog.Debug("Bad request", slog.String("err", err.Error()))
		// http.Error(w, "Bad request", http.StatusBadRequest)
		nav_default(w, r)
		return
	}

	Login := strings.TrimSpace(r.PostForm.Get("login"))

	var data = struct {
		Login string
	}{
		Login,
	}

	if len(Login) == 0 {
		TemplatesMap[BLOCK_TEMP_ACCOUNT_DEL_FORM_ERROR].Execute(w, data)
		return
	}

	err = gauth.DeleteUser(Login)
	if err != nil {
		TemplatesMap[BLOCK_TEMP_ACCOUNT_DEL_FORM_ERROR].Execute(w, data)
		return
	}

	TemplatesMap[BLOCK_TEMP_ACCOUNT_DEL_FORM_OK].Execute(w, data)
}

/*
Settings block
*/

func nav_settings(w http.ResponseWriter, r *http.Request) {
	data := config.DefaultConfig
	TemplatesMap[BLOCK_TEMP_SETTINGS].Execute(w, data)
}

func settings_wsc_change_sw(w http.ResponseWriter, r *http.Request) {
	if config.DefaultConfig.WebSocketConnector.Enable {
		config.DefaultConfig.WebSocketConnector.Enable = false
		closer.RunAndDelHandler(websocketconn.Shutdown)
	} else {
		config.DefaultConfig.WebSocketConnector.Enable = true
		go websocketconn.Start(&config.DefaultConfig)
		closer.AddHandler(websocketconn.Shutdown) // Register a shutdown handler.
	}
	slog.Warn("The service has been switched.", slog.String("service", "WebSocketConnector"))

	nav_settings(w, r)
}

func settings_rest_change_sw(w http.ResponseWriter, r *http.Request) {
	if config.DefaultConfig.RestConnector.Enable {
		config.DefaultConfig.RestConnector.Enable = false
		closer.RunAndDelHandler(rest.Shutdown)
	} else {
		config.DefaultConfig.RestConnector.Enable = true
		go rest.Start(&config.DefaultConfig)
		closer.AddHandler(rest.Shutdown) // Register a shutdown handler.
	}
	slog.Warn("The service has been switched.", slog.String("service", "RestConnector"))

	nav_settings(w, r)
}

func settings_grpc_change_sw(w http.ResponseWriter, r *http.Request) {
	if config.DefaultConfig.GrpcConnector.Enable {
		config.DefaultConfig.GrpcConnector.Enable = false
		closer.RunAndDelHandler(grpc.Shutdown)
	} else {
		config.DefaultConfig.GrpcConnector.Enable = true
		go grpc.Start(&config.DefaultConfig)
		closer.AddHandler(grpc.Shutdown) // Register a shutdown handler.
	}
	slog.Warn("The service has been switched.", slog.String("service", "GrpcConnector"))

	nav_settings(w, r)
}

func settings_web_change_sw(w http.ResponseWriter, r *http.Request) {
	slog.Warn("This service cannot be disabled.", slog.String("service", "WebServer"))

	nav_settings(w, r)
}
