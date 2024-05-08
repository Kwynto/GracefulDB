package webmanage

import (
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/Kwynto/gosession"

	"github.com/Kwynto/GracefulDB/internal/analyzers/vqlanalyzer"
	"github.com/Kwynto/GracefulDB/internal/config"
	"github.com/Kwynto/GracefulDB/internal/connectors/grpc"
	"github.com/Kwynto/GracefulDB/internal/connectors/rest"
	"github.com/Kwynto/GracefulDB/internal/connectors/websocketconn"
	"github.com/Kwynto/GracefulDB/internal/engine/basicsystem/gauth"
	"github.com/Kwynto/GracefulDB/internal/engine/basicsystem/gtypes"
	"github.com/Kwynto/GracefulDB/internal/engine/core"

	"github.com/Kwynto/GracefulDB/pkg/lib/closer"
)

type TViewAccountsTable struct {
	System      bool
	Superuser   bool
	Baned       bool
	Login       string
	Status      string
	Roles       string
	Description string
}

/*
The main section
*/

// Handler after authorization
func homeDefault(w http.ResponseWriter, r *http.Request) {
	// This function is complete
	if IsolatedAuth(w, r, []gauth.TRole{gauth.MANAGER, gauth.ENGINEER, gauth.USER}) {
		logout(w, r)
		return
	}

	sesID := gosession.Start(&w, r)
	auth := sesID.Get("auth")
	login := fmt.Sprint(auth)
	profile, _ := gauth.GetProfile(login) // There is no point in checking the error, since erroneous data acquisition is eliminated at the isolation stage.
	// profile, err := gauth.GetProfile(login)
	// if err != nil {
	// 	logout(w, r)
	// 	return
	// }

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

	err := TemplatesMap[HOME_TEMP_NAME].Execute(w, data)
	if err != nil {
		slog.Debug("Internal Server Error", slog.String("err", err.Error()))
		// http.Error(w, "Internal Server Error", http.StatusInternalServerError)
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

			secret := gtypes.TSecret{
				Login:    username,
				Password: password,
			}
			ticket, err2 := gauth.NewAuth(&secret)
			if err2 == nil {
				core.States[ticket] = core.TState{
					CurrentDB: "",
				}
			}
		}
		http.Redirect(w, r, "/", http.StatusFound)
	} else {
		TemplatesMap[AUTH_TEMP_NAME].Execute(w, nil)
		// err := TemplatesMap[AUTH_TEMP_NAME].Execute(w, nil)
		// if err != nil {
		// 	slog.Debug("Internal Server Error", slog.String("err", err.Error()))
		// 	// http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		// }
	}
}

// Handler for the main route
func home(w http.ResponseWriter, r *http.Request) {
	// This function is complete
	if r.URL.Path != "/" {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	sesID := gosession.Start(&w, r)
	auth := sesID.Get("auth")
	if auth == nil {
		homeAuth(w, r)
	} else {
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
Profile section
*/

func selfedit_load_form(w http.ResponseWriter, r *http.Request) {
	if IsolatedAuth(w, r, []gauth.TRole{gauth.MANAGER, gauth.ENGINEER, gauth.USER}) {
		TemplatesMap[BLOCK_TEMP_ACCESS_DENIED].Execute(w, nil)
		return
	}

	sesID := gosession.Start(&w, r)
	auth := sesID.Get("auth")
	login := fmt.Sprint(auth)
	profile, _ := gauth.GetProfile(login) // There is no point in checking the error, since erroneous data acquisition is eliminated at the isolation stage.
	// profile, err := gauth.GetProfile(login)
	// if err != nil {
	// 	TemplatesMap[BLOCK_TEMP_ACCOUNT_SELFEDIT_ERROR].Execute(w, nil)
	// 	return
	// }

	data := struct {
		Login string
		Desc  string
	}{
		Login: login,
		Desc:  profile.Description,
	}

	TemplatesMap[BLOCK_TEMP_ACCOUNT_SELFEDIT_LOAD].Execute(w, data)
}

func selfedit_ok(w http.ResponseWriter, r *http.Request) {
	if IsolatedAuth(w, r, []gauth.TRole{gauth.MANAGER, gauth.ENGINEER, gauth.USER}) {
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
		nav_default(w, r)
		return
	}

	var data = struct {
		MsgErr string
	}{
		MsgErr: "",
	}

	sesID := gosession.Start(&w, r)
	auth := sesID.Get("auth")
	login := fmt.Sprint(auth)
	profile, _ := gauth.GetProfile(login) // There is no point in checking the error, since erroneous data acquisition is eliminated at the isolation stage.
	// profile, err := gauth.GetProfile(login)
	// if err != nil {
	// 	data.MsgErr = "Unknown user."
	// 	TemplatesMap[BLOCK_TEMP_ACCOUNT_SELFEDIT_ERROR].Execute(w, data)
	// 	return
	// }

	password := strings.TrimSpace(r.PostForm.Get("password"))
	if password == "" {
		slog.Debug("Update user", slog.String("err", "an empty password"))
		data.MsgErr = "The password cannot be empty."
		TemplatesMap[BLOCK_TEMP_ACCOUNT_SELFEDIT_ERROR].Execute(w, data)
		return
	}

	desc := strings.TrimSpace(r.PostForm.Get("desc"))
	profile.Description = desc

	gauth.UpdateUser(login, password, profile) // An error is not possible, since all fields have already been checked.
	// err = gauth.UpdateUser(login, password, profile)
	// if err != nil {
	// 	slog.Debug("Update user", slog.String("err", err.Error()))
	// 	data.MsgErr = "The user could not be updated."
	// 	TemplatesMap[BLOCK_TEMP_ACCOUNT_SELFEDIT_ERROR].Execute(w, data)
	// 	return
	// }

	TemplatesMap[BLOCK_TEMP_ACCOUNT_SELFEDIT_OK].Execute(w, nil)
}

/*
Dashboard section
*/

func nav_dashboard(w http.ResponseWriter, r *http.Request) {
	if IsolatedAuth(w, r, []gauth.TRole{gauth.MANAGER, gauth.ENGINEER, gauth.USER}) {
		TemplatesMap[BLOCK_TEMP_ACCESS_DENIED].Execute(w, nil)
		return
	}

	TemplatesMap[BLOCK_TEMP_DASHBOARD].Execute(w, nil)
}

/*
Databases section
*/

func nav_databases(w http.ResponseWriter, r *http.Request) {
	if IsolatedAuth(w, r, []gauth.TRole{gauth.ENGINEER}) {
		TemplatesMap[BLOCK_TEMP_ACCESS_DENIED].Execute(w, nil)
		return
	}

	TemplatesMap[BLOCK_TEMP_DATABASES].Execute(w, nil)
}

/*
Console section
*/

func nav_console(w http.ResponseWriter, r *http.Request) {
	if IsolatedAuth(w, r, []gauth.TRole{gauth.ENGINEER}) {
		TemplatesMap[BLOCK_TEMP_ACCESS_DENIED].Execute(w, nil)
		return
	}

	TemplatesMap[BLOCK_TEMP_CONSOLE].Execute(w, nil)
}

func console_request(w http.ResponseWriter, r *http.Request) {
	timeR := time.Now().Format(CONSOLE_TIME_FORMAT)

	if IsolatedAuth(w, r, []gauth.TRole{gauth.ENGINEER}) {
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
		nav_default(w, r)
		return
	}

	request := strings.TrimSpace(r.PostForm.Get("request"))

	sesID := gosession.Start(&w, r)
	auth := sesID.Get("auth")
	login := fmt.Sprint(auth)

	ticket, err := gauth.GetTicket(login)
	if err != nil {
		TemplatesMap[BLOCK_TEMP_ACCESS_DENIED].Execute(w, nil)
		return
	}

	answer := vqlanalyzer.Request(ticket, request, []string{})

	timeA := time.Now().Format(CONSOLE_TIME_FORMAT)

	data := struct {
		From    string
		Request string
		Answer  string
		TimeR   string
		TimeA   string
	}{
		From:    login,
		Request: request,
		Answer:  answer,
		TimeR:   timeR,
		TimeA:   timeA,
	}
	TemplatesMap[BLOCK_TEMP_CONSOLE_REQUEST_ANSWER].Execute(w, data)
}

/*
Accounts section
*/

func nav_accounts(w http.ResponseWriter, r *http.Request) {
	if IsolatedAuth(w, r, []gauth.TRole{gauth.MANAGER}) {
		TemplatesMap[BLOCK_TEMP_ACCESS_DENIED].Execute(w, nil)
		return
	}

	var table = make([]TViewAccountsTable, 0, 10)
	for key := range gauth.MHash {
		element := TViewAccountsTable{
			System:      false,
			Superuser:   false,
			Baned:       false,
			Login:       key,
			Status:      gauth.MAccess[key].Status.String(),
			Roles:       "",
			Description: gauth.MAccess[key].Description,
		}

		for _, role := range gauth.MAccess[key].Roles {
			if role == gauth.SYSTEM {
				element.System = true
			}
			element.Roles = fmt.Sprintf("%s %s", element.Roles, role.String())
		}

		if key == "root" {
			element.Superuser = true
		}
		if gauth.MAccess[key].Status == gauth.BANED {
			element.Baned = true
		}

		table = append(table, element)
	}

	TemplatesMap[BLOCK_TEMP_ACCOUNTS].Execute(w, table)
}

func account_create_load_form(w http.ResponseWriter, r *http.Request) {
	if IsolatedAuth(w, r, []gauth.TRole{gauth.MANAGER}) {
		TemplatesMap[BLOCK_TEMP_ACCESS_DENIED].Execute(w, nil)
		return
	}
	TemplatesMap[BLOCK_TEMP_ACCOUNT_CREATE_FORM_LOAD].Execute(w, nil)
}

func account_create_ok(w http.ResponseWriter, r *http.Request) {
	if IsolatedAuth(w, r, []gauth.TRole{gauth.MANAGER}) {
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
	if IsolatedAuth(w, r, []gauth.TRole{gauth.MANAGER}) {
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
	}{
		System: false,
		Login:  user,
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

	TemplatesMap[BLOCK_TEMP_ACCOUNT_EDIT_FORM_LOAD].Execute(w, data)
}

func account_edit_ok(w http.ResponseWriter, r *http.Request) {
	if IsolatedAuth(w, r, []gauth.TRole{gauth.MANAGER}) {
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

	var roles []gauth.TRole
	if Login != "root" {
		rolesIn := r.Form["role_names"]
		for _, role := range rolesIn {
			switch role {
			case "SYSTEM":
				roles = append(roles, gauth.SYSTEM)
			case "ADMIN":
				roles = append(roles, gauth.ADMIN)
			case "MANAGER":
				roles = append(roles, gauth.MANAGER)
			case "ENGINEER":
				roles = append(roles, gauth.ENGINEER)
			case "USER":
				roles = append(roles, gauth.USER)
			default:
				roles = append(roles, gauth.USER)
			}
		}
	}

	if Login == "root" {
		desc = ""
		status = 2
		roles = append(roles, gauth.ADMIN)
	}

	access := gauth.TProfile{
		Description: desc,
		Status:      gauth.TStatus(status),
		Roles:       roles,
	}

	gauth.UpdateUser(Login, password, access) // An error is not possible, since all fields have already been checked.
	// err = gauth.UpdateUser(Login, password, access)
	// if err != nil {
	// 	slog.Debug("Update user", slog.String("err", err.Error()))
	// 	data.MsgErr = "The user could not be updated."
	// 	TemplatesMap[BLOCK_TEMP_ACCOUNT_EDIT_FORM_ERROR].Execute(w, data)
	// 	return
	// }

	TemplatesMap[BLOCK_TEMP_ACCOUNT_EDIT_FORM_OK].Execute(w, data)
}

func account_ban_load_form(w http.ResponseWriter, r *http.Request) {
	if IsolatedAuth(w, r, []gauth.TRole{gauth.MANAGER}) {
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
	if IsolatedAuth(w, r, []gauth.TRole{gauth.MANAGER}) {
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
	if IsolatedAuth(w, r, []gauth.TRole{gauth.MANAGER}) {
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
	if IsolatedAuth(w, r, []gauth.TRole{gauth.MANAGER}) {
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
	if IsolatedAuth(w, r, []gauth.TRole{gauth.MANAGER}) {
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
	if IsolatedAuth(w, r, []gauth.TRole{gauth.MANAGER}) {
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
Settings section
*/

func nav_settings(w http.ResponseWriter, r *http.Request) {
	if IsolatedAuth(w, r, []gauth.TRole{gauth.ADMIN}) {
		TemplatesMap[BLOCK_TEMP_ACCESS_DENIED].Execute(w, nil)
		return
	}

	data := config.StDefaultConfig
	TemplatesMap[BLOCK_TEMP_SETTINGS].Execute(w, data)
}

func settings_core_friendly_change_sw(w http.ResponseWriter, r *http.Request) {
	if IsolatedAuth(w, r, []gauth.TRole{gauth.ADMIN}) {
		TemplatesMap[BLOCK_TEMP_ACCESS_DENIED].Execute(w, nil)
		return
	}

	if config.StDefaultConfig.CoreSettings.FriendlyMode {
		config.StDefaultConfig.CoreSettings.FriendlyMode = false
	} else {
		config.StDefaultConfig.CoreSettings.FriendlyMode = true
	}
	core.LocalCoreSettings = core.LoadLocalCoreSettings(&config.StDefaultConfig)
	msg := "The friendly mode has been switched."
	slog.Warn(msg, slog.String("FriendlyMode", fmt.Sprintf("%v", core.LocalCoreSettings.FriendlyMode)))

	nav_settings(w, r)
}

func settings_wsc_change_sw(w http.ResponseWriter, r *http.Request) {
	if IsolatedAuth(w, r, []gauth.TRole{gauth.ADMIN}) {
		TemplatesMap[BLOCK_TEMP_ACCESS_DENIED].Execute(w, nil)
		return
	}

	if config.StDefaultConfig.WebSocketConnector.Enable {
		config.StDefaultConfig.WebSocketConnector.Enable = false
		closer.RunAndDelHandler(websocketconn.Shutdown)
	} else {
		config.StDefaultConfig.WebSocketConnector.Enable = true
		go websocketconn.Start(&config.StDefaultConfig)
		closer.AddHandler(websocketconn.Shutdown) // Register a shutdown handler.
	}
	slog.Warn("The service has been switched.", slog.String("service", "WebSocketConnector"))

	nav_settings(w, r)
}

func settings_rest_change_sw(w http.ResponseWriter, r *http.Request) {
	if IsolatedAuth(w, r, []gauth.TRole{gauth.ADMIN}) {
		TemplatesMap[BLOCK_TEMP_ACCESS_DENIED].Execute(w, nil)
		return
	}

	if config.StDefaultConfig.RestConnector.Enable {
		config.StDefaultConfig.RestConnector.Enable = false
		closer.RunAndDelHandler(rest.Shutdown)
	} else {
		config.StDefaultConfig.RestConnector.Enable = true
		go rest.Start(&config.StDefaultConfig)
		closer.AddHandler(rest.Shutdown) // Register a shutdown handler.
	}
	slog.Warn("The service has been switched.", slog.String("service", "RestConnector"))

	nav_settings(w, r)
}

func settings_grpc_change_sw(w http.ResponseWriter, r *http.Request) {
	if IsolatedAuth(w, r, []gauth.TRole{gauth.ADMIN}) {
		TemplatesMap[BLOCK_TEMP_ACCESS_DENIED].Execute(w, nil)
		return
	}

	if config.StDefaultConfig.GrpcConnector.Enable {
		config.StDefaultConfig.GrpcConnector.Enable = false
		closer.RunAndDelHandler(grpc.Shutdown)
	} else {
		config.StDefaultConfig.GrpcConnector.Enable = true
		go grpc.Start(&config.StDefaultConfig)
		closer.AddHandler(grpc.Shutdown) // Register a shutdown handler.
	}
	slog.Warn("The service has been switched.", slog.String("service", "GrpcConnector"))

	nav_settings(w, r)
}

func settings_web_change_sw(w http.ResponseWriter, r *http.Request) {
	if IsolatedAuth(w, r, []gauth.TRole{gauth.ADMIN}) {
		TemplatesMap[BLOCK_TEMP_ACCESS_DENIED].Execute(w, nil)
		return
	}

	slog.Warn("This service cannot be disabled.", slog.String("service", "WebServer"))

	nav_settings(w, r)
}
