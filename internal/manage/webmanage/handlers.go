package webmanage

import (
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
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
	System      bool
	Superuser   bool
	Baned       bool
	Login       string
	Status      string
	Role        string
	Description string
}

/*
The main block
*/

// Handler after authorization
func homeDefault(w http.ResponseWriter, r *http.Request) {
	// This function is complete
	if IsolatedAuth(w, r, gauth.ENGINEER) {
		logout(w, r)
		return
	}

	sesID := gosession.Start(&w, r)
	auth := sesID.Get("auth")
	login := fmt.Sprint(auth)
	profile, err := gauth.GetProfile(login)
	if err != nil {
		logout(w, r)
		return
	}

	var data = struct {
		Login string
		Roles string
	}{
		Login: login,
		Roles: "",
	}

	for _, role := range profile.Roles {
		data.Roles = fmt.Sprintf("%s %s", data.Roles, role.String())
	}

	err = TemplatesMap[HOME_TEMP_NAME].Execute(w, data)
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
		// login := fmt.Sprint(auth)
		// homeDefault(w, r, login)
		homeDefault(w, r)
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
	if IsolatedAuth(w, r, gauth.ENGINEER) {
		TemplatesMap[BLOCK_TEMP_ACCESS_DENIED].Execute(w, nil)
		return
	}

	TemplatesMap[BLOCK_TEMP_DASHBOARD].Execute(w, nil)
}

/*
Databases block
*/

func nav_databases(w http.ResponseWriter, r *http.Request) {
	if IsolatedAuth(w, r, gauth.ENGINEER) {
		TemplatesMap[BLOCK_TEMP_ACCESS_DENIED].Execute(w, nil)
		return
	}

	TemplatesMap[BLOCK_TEMP_DATABASES].Execute(w, nil)
}

/*
Accounts block
*/

func nav_accounts(w http.ResponseWriter, r *http.Request) {
	if IsolatedAuth(w, r, gauth.MANAGER) {
		TemplatesMap[BLOCK_TEMP_ACCESS_DENIED].Execute(w, nil)
		return
	}

	var table = make([]TViewAccountsTable, 0, 10)
	for key := range gauth.HashMap {
		element := TViewAccountsTable{
			System:      false,
			Superuser:   false,
			Baned:       false,
			Login:       key,
			Status:      gauth.AccessMap[key].Status.String(),
			Role:        "",
			Description: gauth.AccessMap[key].Description,
		}

		for _, role := range gauth.AccessMap[key].Roles {
			if role == gauth.SYSTEM {
				element.System = true
			}
			element.Role = fmt.Sprintf("%s %s", element.Role, role.String())
		}

		// if gauth.AccessMap[key].Role == gauth.SYSTEM {
		// 	element.System = true
		// }

		if key == "root" {
			element.Superuser = true
		}
		if gauth.AccessMap[key].Status == gauth.BANED {
			element.Baned = true
		}

		table = append(table, element)
	}

	// view := table
	// TemplatesMap[BLOCK_TEMP_ACCOUNTS].Execute(w, view)
	TemplatesMap[BLOCK_TEMP_ACCOUNTS].Execute(w, table)
}

func account_create_load_form(w http.ResponseWriter, r *http.Request) {
	if IsolatedAuth(w, r, gauth.MANAGER) {
		TemplatesMap[BLOCK_TEMP_ACCESS_DENIED].Execute(w, nil)
		return
	}
	TemplatesMap[BLOCK_TEMP_ACCOUNT_CREATE_FORM_LOAD].Execute(w, nil)
}

func account_create_ok(w http.ResponseWriter, r *http.Request) {
	if IsolatedAuth(w, r, gauth.MANAGER) {
		TemplatesMap[BLOCK_TEMP_ACCESS_DENIED].Execute(w, nil)
		return
	}

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

	if len(Login) == 0 || len(password) == 0 {
		TemplatesMap[BLOCK_TEMP_ACCOUNT_CREATE_FORM_ERROR].Execute(w, data)
		return
	}

	access := gauth.TProfile{
		Description: desc,
		Status:      gauth.NEW,
		Roles:       []gauth.TRole{gauth.USER},
		Rules:       []string{""},
	}

	err = gauth.AddUser(Login, password, access)
	if err != nil {
		TemplatesMap[BLOCK_TEMP_ACCOUNT_CREATE_FORM_ERROR].Execute(w, data)
		return
	}

	slog.Info("The user has been created", slog.String("user", Login))
	TemplatesMap[BLOCK_TEMP_ACCOUNT_CREATE_FORM_OK].Execute(w, data)
}

func account_edit_load_form(w http.ResponseWriter, r *http.Request) {
	if IsolatedAuth(w, r, gauth.MANAGER) {
		TemplatesMap[BLOCK_TEMP_ACCESS_DENIED].Execute(w, nil)
		return
	}

	user := strings.TrimSpace(r.URL.Query().Get("user"))
	data := struct {
		System      bool
		Login       string
		Description string
		Status      gauth.TStatus
		Roles       []string
		Rules       string
	}{
		System: false,
		Login:  user,
		Rules:  "",
	}

	profile, err := gauth.GetProfile(user)
	if err != nil {
		TemplatesMap[BLOCK_TEMP_ACCOUNT_EDIT_FORM_ERROR].Execute(w, data)
		return
	}
	data.Description = profile.Description
	data.Status = profile.Status

	for _, role := range profile.Roles {
		if role == gauth.SYSTEM {
			data.System = true
		}
		data.Roles = append(data.Roles, role.String())
	}

	// data.Roles = profile.Roles

	for _, v := range profile.Rules {
		data.Rules = fmt.Sprintf("%s\n%s", data.Rules, v)
	}

	TemplatesMap[BLOCK_TEMP_ACCOUNT_EDIT_FORM_LOAD].Execute(w, data)
}

func account_edit_ok(w http.ResponseWriter, r *http.Request) {
	if IsolatedAuth(w, r, gauth.MANAGER) {
		TemplatesMap[BLOCK_TEMP_ACCESS_DENIED].Execute(w, nil)
		return
	}

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

	data := struct {
		Login  string
		MsgErr string
	}{
		Login:  "",
		MsgErr: "",
	}

	Login := strings.TrimSpace(r.PostForm.Get("login"))
	if Login == "" {
		slog.Debug("Update user", slog.String("err", "invalid username"))
		data.MsgErr = "Invalid username."
		TemplatesMap[BLOCK_TEMP_ACCOUNT_EDIT_FORM_ERROR].Execute(w, data)
		return
	}
	data.Login = Login

	password := strings.TrimSpace(r.PostForm.Get("password"))
	if password == "" {
		slog.Debug("Update user", slog.String("err", "an empty password"))
		data.MsgErr = "The password cannot be empty."
		TemplatesMap[BLOCK_TEMP_ACCOUNT_EDIT_FORM_ERROR].Execute(w, data)
		return
	}

	desc := strings.TrimSpace(r.PostForm.Get("desc"))

	status, err := strconv.Atoi(strings.TrimSpace(r.PostForm.Get("status")))
	if (err != nil || status < 1) && Login != "root" {
		slog.Debug("Update user", slog.String("err", "incorrect status"))
		data.MsgErr = "Incorrect status."
		TemplatesMap[BLOCK_TEMP_ACCOUNT_EDIT_FORM_ERROR].Execute(w, data)
		return
	}

	role, err := strconv.Atoi(strings.TrimSpace(r.PostForm.Get("role")))
	if (err != nil || role == 0) && Login != "root" {
		slog.Debug("Update user", slog.String("err", "incorrect role"))
		data.MsgErr = "Incorrect role."
		TemplatesMap[BLOCK_TEMP_ACCOUNT_EDIT_FORM_ERROR].Execute(w, data)
		return
	}

	rulesIn := strings.TrimSpace(r.PostForm.Get("rules"))
	rules := strings.Split(rulesIn, "\n")

	if Login == "root" {
		desc = ""
		status = 2
		role = 1
		rules = []string{""}
	}

	access := gauth.TProfile{
		Description: desc,
		Status:      gauth.TStatus(status),
		Roles:       []gauth.TRole{gauth.TRole(role)},
		Rules:       rules,
	}

	err = gauth.UpdateUser(Login, password, access)
	if err != nil {
		slog.Debug("Update user", slog.String("err", err.Error()))
		data.MsgErr = "The user could not be updated."
		TemplatesMap[BLOCK_TEMP_ACCOUNT_EDIT_FORM_ERROR].Execute(w, data)
		return
	}

	TemplatesMap[BLOCK_TEMP_ACCOUNT_EDIT_FORM_OK].Execute(w, data)
}

func account_ban_load_form(w http.ResponseWriter, r *http.Request) {
	if IsolatedAuth(w, r, gauth.MANAGER) {
		TemplatesMap[BLOCK_TEMP_ACCESS_DENIED].Execute(w, nil)
		return
	}

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
	if IsolatedAuth(w, r, gauth.MANAGER) {
		TemplatesMap[BLOCK_TEMP_ACCESS_DENIED].Execute(w, nil)
		return
	}

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

	slog.Info("The user has been blocked", slog.String("user", Login))
	TemplatesMap[BLOCK_TEMP_ACCOUNT_BAN_FORM_OK].Execute(w, data)
}

func account_unban_load_form(w http.ResponseWriter, r *http.Request) {
	if IsolatedAuth(w, r, gauth.MANAGER) {
		TemplatesMap[BLOCK_TEMP_ACCESS_DENIED].Execute(w, nil)
		return
	}

	user := strings.TrimSpace(r.URL.Query().Get("user"))
	data := struct {
		Login string
	}{
		Login: user,
	}

	if user == "" || user == "root" {
		TemplatesMap[BLOCK_TEMP_ACCOUNT_UNBAN_FORM_ERROR].Execute(w, data)
		return
	}

	TemplatesMap[BLOCK_TEMP_ACCOUNT_UNBAN_FORM_LOAD].Execute(w, data)
}

func account_unban_ok(w http.ResponseWriter, r *http.Request) {
	if IsolatedAuth(w, r, gauth.MANAGER) {
		TemplatesMap[BLOCK_TEMP_ACCESS_DENIED].Execute(w, nil)
		return
	}

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
		TemplatesMap[BLOCK_TEMP_ACCOUNT_UNBAN_FORM_ERROR].Execute(w, data)
		return
	}

	err = gauth.UnblockUser(Login)
	if err != nil {
		TemplatesMap[BLOCK_TEMP_ACCOUNT_UNBAN_FORM_ERROR].Execute(w, data)
		return
	}

	slog.Info("The user has been unblocked", slog.String("user", Login))
	TemplatesMap[BLOCK_TEMP_ACCOUNT_UNBAN_FORM_OK].Execute(w, data)
}

func account_del_load_form(w http.ResponseWriter, r *http.Request) {
	if IsolatedAuth(w, r, gauth.MANAGER) {
		TemplatesMap[BLOCK_TEMP_ACCESS_DENIED].Execute(w, nil)
		return
	}

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
	if IsolatedAuth(w, r, gauth.MANAGER) {
		TemplatesMap[BLOCK_TEMP_ACCESS_DENIED].Execute(w, nil)
		return
	}

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

	slog.Info("The user has been removed", slog.String("user", Login))
	TemplatesMap[BLOCK_TEMP_ACCOUNT_DEL_FORM_OK].Execute(w, data)
}

/*
Settings block
*/

func nav_settings(w http.ResponseWriter, r *http.Request) {
	if IsolatedAuth(w, r, gauth.ADMIN) {
		TemplatesMap[BLOCK_TEMP_ACCESS_DENIED].Execute(w, nil)
		return
	}

	data := config.DefaultConfig
	TemplatesMap[BLOCK_TEMP_SETTINGS].Execute(w, data)
}

func settings_wsc_change_sw(w http.ResponseWriter, r *http.Request) {
	if IsolatedAuth(w, r, gauth.ADMIN) {
		TemplatesMap[BLOCK_TEMP_ACCESS_DENIED].Execute(w, nil)
		return
	}

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
	if IsolatedAuth(w, r, gauth.ADMIN) {
		TemplatesMap[BLOCK_TEMP_ACCESS_DENIED].Execute(w, nil)
		return
	}

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
	if IsolatedAuth(w, r, gauth.ADMIN) {
		TemplatesMap[BLOCK_TEMP_ACCESS_DENIED].Execute(w, nil)
		return
	}

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
	if IsolatedAuth(w, r, gauth.ADMIN) {
		TemplatesMap[BLOCK_TEMP_ACCESS_DENIED].Execute(w, nil)
		return
	}

	slog.Warn("This service cannot be disabled.", slog.String("service", "WebServer"))

	nav_settings(w, r)
}
