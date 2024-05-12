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
func fnHomeDefault(w http.ResponseWriter, r *http.Request) {
	// This function is complete
	if IsolatedAuth(w, r, []gauth.TRole{gauth.MANAGER, gauth.ENGINEER, gauth.USER}) {
		fnLogout(w, r)
		return
	}

	sSesID := gosession.Start(&w, r)
	inAuth := sSesID.Get("auth")
	sLogin := fmt.Sprint(inAuth)
	stProfile, _ := gauth.GetProfile(sLogin) // There is no point in checking the error, since erroneous data acquisition is eliminated at the isolation stage.

	var stData = struct {
		Login string
		Roles string
	}{
		Login: sLogin,
		Roles: "",
	}

	for _, iRole := range stProfile.Roles {
		stData.Roles = fmt.Sprintf("%s %s", stData.Roles, iRole.String())
	}

	err := MTemplates[HOME_TEMP_NAME].Execute(w, stData)
	if err != nil {
		slog.Debug("Internal Server Error", slog.String("err", err.Error()))
	}
}

// Authorization Handler
func fnHomeAuth(w http.ResponseWriter, r *http.Request) {
	// This function is complete
	if r.Method == http.MethodPost {
		err := r.ParseForm()
		if err != nil {
			slog.Debug("Bad request", slog.String("err", err.Error()))
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}

		sUsername := r.PostForm.Get("username")
		sPassword := r.PostForm.Get("password")
		isAuth := gauth.CheckUser(sUsername, sPassword)
		if isAuth {
			sSesID := gosession.Start(&w, r)
			sSesID.Set("auth", sUsername)

			stSecret := gtypes.TSecret{
				Login:    sUsername,
				Password: sPassword,
			}
			sTicket, err2 := gauth.NewAuth(&stSecret)
			if err2 == nil {
				core.MStates[sTicket] = core.TState{
					CurrentDB: "",
				}
			}
		}
		http.Redirect(w, r, "/", http.StatusFound)
	} else {
		MTemplates[AUTH_TEMP_NAME].Execute(w, nil)
	}
}

// Handler for the main route
func fnHome(w http.ResponseWriter, r *http.Request) {
	// This function is complete
	if r.URL.Path != "/" {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	sSesID := gosession.Start(&w, r)
	inAuth := sSesID.Get("auth")
	if inAuth == nil {
		fnHomeAuth(w, r)
	} else {
		fnHomeDefault(w, r)
	}
}

// Exit handler
func fnLogout(w http.ResponseWriter, r *http.Request) {
	// This function is complete
	sSesID := gosession.Start(&w, r)
	sSesID.Remove("auth")
	http.Redirect(w, r, "/", http.StatusFound)
}

/*
Nav Menu Handlers
*/

func fnNavDefault(w http.ResponseWriter, r *http.Request) {
	// This function is complete
	MTemplates[BLOCK_TEMP_DEFAULT].Execute(w, nil)
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
		MTemplates[BLOCK_TEMP_ACCESS_DENIED].Execute(w, nil)
		return
	}

	sSesID := gosession.Start(&w, r)
	inAuth := sSesID.Get("auth")
	sLogin := fmt.Sprint(inAuth)
	stProfile, _ := gauth.GetProfile(sLogin) // There is no point in checking the error, since erroneous data acquisition is eliminated at the isolation stage.

	stData := struct {
		Login string
		Desc  string
	}{
		Login: sLogin,
		Desc:  stProfile.Description,
	}

	MTemplates[BLOCK_TEMP_ACCOUNT_SELFEDIT_LOAD].Execute(w, stData)
}

func selfedit_ok(w http.ResponseWriter, r *http.Request) {
	if IsolatedAuth(w, r, []gauth.TRole{gauth.MANAGER, gauth.ENGINEER, gauth.USER}) {
		MTemplates[BLOCK_TEMP_ACCESS_DENIED].Execute(w, nil)
		return
	}

	if r.Method != http.MethodPost {
		fnNavDefault(w, r)
		return
	}

	err := r.ParseForm()
	if err != nil {
		slog.Debug("Bad request", slog.String("err", err.Error()))
		fnNavDefault(w, r)
		return
	}

	var stData = struct {
		MsgErr string
	}{
		MsgErr: "",
	}

	sSesID := gosession.Start(&w, r)
	inAuth := sSesID.Get("auth")
	sLogin := fmt.Sprint(inAuth)
	stProfile, _ := gauth.GetProfile(sLogin) // There is no point in checking the error, since erroneous data acquisition is eliminated at the isolation stage.

	sPassword := strings.TrimSpace(r.PostForm.Get("password"))
	if sPassword == "" {
		slog.Debug("Update user", slog.String("err", "an empty password"))
		stData.MsgErr = "The password cannot be empty."
		MTemplates[BLOCK_TEMP_ACCOUNT_SELFEDIT_ERROR].Execute(w, stData)
		return
	}

	sDesc := strings.TrimSpace(r.PostForm.Get("desc"))
	stProfile.Description = sDesc

	gauth.UpdateUser(sLogin, sPassword, stProfile) // An error is not possible, since all fields have already been checked.

	MTemplates[BLOCK_TEMP_ACCOUNT_SELFEDIT_OK].Execute(w, nil)
}

/*
Dashboard section
*/

func nav_dashboard(w http.ResponseWriter, r *http.Request) {
	if IsolatedAuth(w, r, []gauth.TRole{gauth.MANAGER, gauth.ENGINEER, gauth.USER}) {
		MTemplates[BLOCK_TEMP_ACCESS_DENIED].Execute(w, nil)
		return
	}

	MTemplates[BLOCK_TEMP_DASHBOARD].Execute(w, nil)
}

/*
Databases section
*/

func nav_databases(w http.ResponseWriter, r *http.Request) {
	if IsolatedAuth(w, r, []gauth.TRole{gauth.ENGINEER}) {
		MTemplates[BLOCK_TEMP_ACCESS_DENIED].Execute(w, nil)
		return
	}

	MTemplates[BLOCK_TEMP_DATABASES].Execute(w, nil)
}

/*
Console section
*/

func nav_console(w http.ResponseWriter, r *http.Request) {
	if IsolatedAuth(w, r, []gauth.TRole{gauth.ENGINEER}) {
		MTemplates[BLOCK_TEMP_ACCESS_DENIED].Execute(w, nil)
		return
	}

	MTemplates[BLOCK_TEMP_CONSOLE].Execute(w, nil)
}

func console_request(w http.ResponseWriter, r *http.Request) {
	sTimeR := time.Now().Format(CONSOLE_TIME_FORMAT)

	if IsolatedAuth(w, r, []gauth.TRole{gauth.ENGINEER}) {
		MTemplates[BLOCK_TEMP_ACCESS_DENIED].Execute(w, nil)
		return
	}

	if r.Method != http.MethodPost {
		fnNavDefault(w, r)
		return
	}

	err := r.ParseForm()
	if err != nil {
		slog.Debug("Bad request", slog.String("err", err.Error()))
		fnNavDefault(w, r)
		return
	}

	sRequest := strings.TrimSpace(r.PostForm.Get("request"))

	sSesID := gosession.Start(&w, r)
	inAuth := sSesID.Get("auth")
	sLogin := fmt.Sprint(inAuth)

	sTicket, err := gauth.GetTicket(sLogin)
	if err != nil {
		MTemplates[BLOCK_TEMP_ACCESS_DENIED].Execute(w, nil)
		return
	}

	sAnswer := vqlanalyzer.Request(sTicket, sRequest, []string{})

	sTimeA := time.Now().Format(CONSOLE_TIME_FORMAT)

	stData := struct {
		From    string
		Request string
		Answer  string
		TimeR   string
		TimeA   string
	}{
		From:    sLogin,
		Request: sRequest,
		Answer:  sAnswer,
		TimeR:   sTimeR,
		TimeA:   sTimeA,
	}
	MTemplates[BLOCK_TEMP_CONSOLE_REQUEST_ANSWER].Execute(w, stData)
}

/*
Accounts section
*/

func nav_accounts(w http.ResponseWriter, r *http.Request) {
	if IsolatedAuth(w, r, []gauth.TRole{gauth.MANAGER}) {
		MTemplates[BLOCK_TEMP_ACCESS_DENIED].Execute(w, nil)
		return
	}

	var stTable = make([]TViewAccountsTable, 0, 10)
	for sKey := range gauth.MHash {
		stElement := TViewAccountsTable{
			System:      false,
			Superuser:   false,
			Baned:       false,
			Login:       sKey,
			Status:      gauth.MAccess[sKey].Status.String(),
			Roles:       "",
			Description: gauth.MAccess[sKey].Description,
		}

		for _, iRole := range gauth.MAccess[sKey].Roles {
			if iRole == gauth.SYSTEM {
				stElement.System = true
			}
			stElement.Roles = fmt.Sprintf("%s %s", stElement.Roles, iRole.String())
		}

		if sKey == "root" {
			stElement.Superuser = true
		}
		if gauth.MAccess[sKey].Status == gauth.BANED {
			stElement.Baned = true
		}

		stTable = append(stTable, stElement)
	}

	MTemplates[BLOCK_TEMP_ACCOUNTS].Execute(w, stTable)
}

func account_create_load_form(w http.ResponseWriter, r *http.Request) {
	if IsolatedAuth(w, r, []gauth.TRole{gauth.MANAGER}) {
		MTemplates[BLOCK_TEMP_ACCESS_DENIED].Execute(w, nil)
		return
	}
	MTemplates[BLOCK_TEMP_ACCOUNT_CREATE_FORM_LOAD].Execute(w, nil)
}

func account_create_ok(w http.ResponseWriter, r *http.Request) {
	if IsolatedAuth(w, r, []gauth.TRole{gauth.MANAGER}) {
		MTemplates[BLOCK_TEMP_ACCESS_DENIED].Execute(w, nil)
		return
	}

	if r.Method != http.MethodPost {
		fnNavDefault(w, r)
		return
	}

	err := r.ParseForm()
	if err != nil {
		slog.Debug("Bad request", slog.String("err", err.Error()))
		fnNavDefault(w, r)
		return
	}

	sLogin := strings.TrimSpace(r.PostForm.Get("login"))
	sPassword := strings.TrimSpace(r.PostForm.Get("password"))
	sDesc := strings.TrimSpace(r.PostForm.Get("desc"))

	var stData = struct {
		Login string
	}{
		Login: sLogin,
	}

	if len(sLogin) == 0 || len(sPassword) == 0 {
		MTemplates[BLOCK_TEMP_ACCOUNT_CREATE_FORM_ERROR].Execute(w, stData)
		return
	}

	stAccess := gauth.TProfile{
		Description: sDesc,
		Status:      gauth.NEW,
		Roles:       []gauth.TRole{gauth.USER},
	}

	err = gauth.AddUser(sLogin, sPassword, stAccess)
	if err != nil {
		MTemplates[BLOCK_TEMP_ACCOUNT_CREATE_FORM_ERROR].Execute(w, stData)
		return
	}

	slog.Info("The user has been created", slog.String("user", sLogin))
	MTemplates[BLOCK_TEMP_ACCOUNT_CREATE_FORM_OK].Execute(w, stData)
}

func account_edit_load_form(w http.ResponseWriter, r *http.Request) {
	if IsolatedAuth(w, r, []gauth.TRole{gauth.MANAGER}) {
		MTemplates[BLOCK_TEMP_ACCESS_DENIED].Execute(w, nil)
		return
	}

	sUser := strings.TrimSpace(r.URL.Query().Get("user"))
	stData := struct {
		System      bool
		Login       string
		Description string
		Status      gauth.TStatus
		Roles       []string
	}{
		System: false,
		Login:  sUser,
	}

	stProfile, err := gauth.GetProfile(sUser)
	if err != nil {
		MTemplates[BLOCK_TEMP_ACCOUNT_EDIT_FORM_ERROR].Execute(w, stData)
		return
	}
	stData.Description = stProfile.Description
	stData.Status = stProfile.Status

	for _, iRole := range stProfile.Roles {
		if iRole == gauth.SYSTEM {
			stData.System = true
		}
		stData.Roles = append(stData.Roles, iRole.String())
	}

	MTemplates[BLOCK_TEMP_ACCOUNT_EDIT_FORM_LOAD].Execute(w, stData)
}

func account_edit_ok(w http.ResponseWriter, r *http.Request) {
	if IsolatedAuth(w, r, []gauth.TRole{gauth.MANAGER}) {
		MTemplates[BLOCK_TEMP_ACCESS_DENIED].Execute(w, nil)
		return
	}

	if r.Method != http.MethodPost {
		fnNavDefault(w, r)
		return
	}

	err := r.ParseForm()
	if err != nil {
		slog.Debug("Bad request", slog.String("err", err.Error()))
		fnNavDefault(w, r)
		return
	}

	stData := struct {
		Login  string
		MsgErr string
	}{
		Login:  "",
		MsgErr: "",
	}

	sLogin := strings.TrimSpace(r.PostForm.Get("login"))
	if sLogin == "" {
		slog.Debug("Update user", slog.String("err", "invalid username"))
		stData.MsgErr = "Invalid username."
		MTemplates[BLOCK_TEMP_ACCOUNT_EDIT_FORM_ERROR].Execute(w, stData)
		return
	}
	stData.Login = sLogin

	sPassword := strings.TrimSpace(r.PostForm.Get("password"))
	if sPassword == "" {
		slog.Debug("Update user", slog.String("err", "an empty password"))
		stData.MsgErr = "The password cannot be empty."
		MTemplates[BLOCK_TEMP_ACCOUNT_EDIT_FORM_ERROR].Execute(w, stData)
		return
	}

	sDesc := strings.TrimSpace(r.PostForm.Get("desc"))

	iStatus, err := strconv.Atoi(strings.TrimSpace(r.PostForm.Get("status")))
	if (err != nil || iStatus < 1) && sLogin != "root" {
		slog.Debug("Update user", slog.String("err", "incorrect status"))
		stData.MsgErr = "Incorrect status."
		MTemplates[BLOCK_TEMP_ACCOUNT_EDIT_FORM_ERROR].Execute(w, stData)
		return
	}

	var slIRoles []gauth.TRole
	if sLogin != "root" {
		slSRolesIn := r.Form["role_names"]
		for _, sRole := range slSRolesIn {
			switch sRole {
			case "SYSTEM":
				slIRoles = append(slIRoles, gauth.SYSTEM)
			case "ADMIN":
				slIRoles = append(slIRoles, gauth.ADMIN)
			case "MANAGER":
				slIRoles = append(slIRoles, gauth.MANAGER)
			case "ENGINEER":
				slIRoles = append(slIRoles, gauth.ENGINEER)
			case "USER":
				slIRoles = append(slIRoles, gauth.USER)
			default:
				slIRoles = append(slIRoles, gauth.USER)
			}
		}
	}

	if sLogin == "root" {
		sDesc = ""
		iStatus = 2
		slIRoles = append(slIRoles, gauth.ADMIN)
	}

	stAccess := gauth.TProfile{
		Description: sDesc,
		Status:      gauth.TStatus(iStatus),
		Roles:       slIRoles,
	}

	gauth.UpdateUser(sLogin, sPassword, stAccess) // An error is not possible, since all fields have already been checked.

	MTemplates[BLOCK_TEMP_ACCOUNT_EDIT_FORM_OK].Execute(w, stData)
}

func account_ban_load_form(w http.ResponseWriter, r *http.Request) {
	if IsolatedAuth(w, r, []gauth.TRole{gauth.MANAGER}) {
		MTemplates[BLOCK_TEMP_ACCESS_DENIED].Execute(w, nil)
		return
	}

	sUser := strings.TrimSpace(r.URL.Query().Get("user"))
	stData := struct {
		Login string
	}{
		Login: sUser,
	}

	if sUser == "" || sUser == "root" {
		MTemplates[BLOCK_TEMP_ACCOUNT_BAN_FORM_ERROR].Execute(w, stData)
		return
	}

	MTemplates[BLOCK_TEMP_ACCOUNT_BAN_FORM_LOAD].Execute(w, stData)
}

func account_ban_ok(w http.ResponseWriter, r *http.Request) {
	if IsolatedAuth(w, r, []gauth.TRole{gauth.MANAGER}) {
		MTemplates[BLOCK_TEMP_ACCESS_DENIED].Execute(w, nil)
		return
	}

	if r.Method != http.MethodPost {
		fnNavDefault(w, r)
		return
	}

	err := r.ParseForm()
	if err != nil {
		slog.Debug("Bad request", slog.String("err", err.Error()))
		fnNavDefault(w, r)
		return
	}

	sLogin := strings.TrimSpace(r.PostForm.Get("login"))

	var stData = struct {
		Login string
	}{
		Login: sLogin,
	}

	if len(sLogin) == 0 {
		MTemplates[BLOCK_TEMP_ACCOUNT_BAN_FORM_ERROR].Execute(w, stData)
		return
	}

	err = gauth.BlockUser(sLogin)
	if err != nil {
		MTemplates[BLOCK_TEMP_ACCOUNT_BAN_FORM_ERROR].Execute(w, stData)
		return
	}

	slog.Info("The user has been blocked", slog.String("user", sLogin))
	MTemplates[BLOCK_TEMP_ACCOUNT_BAN_FORM_OK].Execute(w, stData)
}

func account_unban_load_form(w http.ResponseWriter, r *http.Request) {
	if IsolatedAuth(w, r, []gauth.TRole{gauth.MANAGER}) {
		MTemplates[BLOCK_TEMP_ACCESS_DENIED].Execute(w, nil)
		return
	}

	sUser := strings.TrimSpace(r.URL.Query().Get("user"))
	stData := struct {
		Login string
	}{
		Login: sUser,
	}

	if sUser == "" || sUser == "root" {
		MTemplates[BLOCK_TEMP_ACCOUNT_UNBAN_FORM_ERROR].Execute(w, stData)
		return
	}

	MTemplates[BLOCK_TEMP_ACCOUNT_UNBAN_FORM_LOAD].Execute(w, stData)
}

func account_unban_ok(w http.ResponseWriter, r *http.Request) {
	if IsolatedAuth(w, r, []gauth.TRole{gauth.MANAGER}) {
		MTemplates[BLOCK_TEMP_ACCESS_DENIED].Execute(w, nil)
		return
	}

	if r.Method != http.MethodPost {
		fnNavDefault(w, r)
		return
	}

	err := r.ParseForm()
	if err != nil {
		slog.Debug("Bad request", slog.String("err", err.Error()))
		fnNavDefault(w, r)
		return
	}

	sLogin := strings.TrimSpace(r.PostForm.Get("login"))

	var stData = struct {
		Login string
	}{
		Login: sLogin,
	}

	if len(sLogin) == 0 {
		MTemplates[BLOCK_TEMP_ACCOUNT_UNBAN_FORM_ERROR].Execute(w, stData)
		return
	}

	err = gauth.UnblockUser(sLogin)
	if err != nil {
		MTemplates[BLOCK_TEMP_ACCOUNT_UNBAN_FORM_ERROR].Execute(w, stData)
		return
	}

	slog.Info("The user has been unblocked", slog.String("user", sLogin))
	MTemplates[BLOCK_TEMP_ACCOUNT_UNBAN_FORM_OK].Execute(w, stData)
}

func account_del_load_form(w http.ResponseWriter, r *http.Request) {
	if IsolatedAuth(w, r, []gauth.TRole{gauth.MANAGER}) {
		MTemplates[BLOCK_TEMP_ACCESS_DENIED].Execute(w, nil)
		return
	}

	sUser := strings.TrimSpace(r.URL.Query().Get("user"))
	stData := struct {
		Login string
	}{
		Login: sUser,
	}

	if sUser == "" || sUser == "root" {
		MTemplates[BLOCK_TEMP_ACCOUNT_DEL_FORM_ERROR].Execute(w, stData)
		return
	}

	MTemplates[BLOCK_TEMP_ACCOUNT_DEL_FORM_LOAD].Execute(w, stData)
}

func account_del_ok(w http.ResponseWriter, r *http.Request) {
	if IsolatedAuth(w, r, []gauth.TRole{gauth.MANAGER}) {
		MTemplates[BLOCK_TEMP_ACCESS_DENIED].Execute(w, nil)
		return
	}

	if r.Method != http.MethodPost {
		fnNavDefault(w, r)
		return
	}

	err := r.ParseForm()
	if err != nil {
		slog.Debug("Bad request", slog.String("err", err.Error()))
		fnNavDefault(w, r)
		return
	}

	sLogin := strings.TrimSpace(r.PostForm.Get("login"))

	var stData = struct {
		Login string
	}{
		Login: sLogin,
	}

	if len(sLogin) == 0 {
		MTemplates[BLOCK_TEMP_ACCOUNT_DEL_FORM_ERROR].Execute(w, stData)
		return
	}

	err = gauth.DeleteUser(sLogin)
	if err != nil {
		MTemplates[BLOCK_TEMP_ACCOUNT_DEL_FORM_ERROR].Execute(w, stData)
		return
	}

	slog.Info("The user has been removed", slog.String("user", sLogin))
	MTemplates[BLOCK_TEMP_ACCOUNT_DEL_FORM_OK].Execute(w, stData)
}

/*
Settings section
*/

func nav_settings(w http.ResponseWriter, r *http.Request) {
	if IsolatedAuth(w, r, []gauth.TRole{gauth.ADMIN}) {
		MTemplates[BLOCK_TEMP_ACCESS_DENIED].Execute(w, nil)
		return
	}

	stData := config.StDefaultConfig
	MTemplates[BLOCK_TEMP_SETTINGS].Execute(w, stData)
}

func settings_core_friendly_change_sw(w http.ResponseWriter, r *http.Request) {
	if IsolatedAuth(w, r, []gauth.TRole{gauth.ADMIN}) {
		MTemplates[BLOCK_TEMP_ACCESS_DENIED].Execute(w, nil)
		return
	}

	if config.StDefaultConfig.CoreSettings.FriendlyMode {
		config.StDefaultConfig.CoreSettings.FriendlyMode = false
	} else {
		config.StDefaultConfig.CoreSettings.FriendlyMode = true
	}
	core.StLocalCoreSettings = core.LoadLocalCoreSettings(&config.StDefaultConfig)
	msg := "The friendly mode has been switched."
	slog.Warn(msg, slog.String("FriendlyMode", fmt.Sprintf("%v", core.StLocalCoreSettings.FriendlyMode)))

	nav_settings(w, r)
}

func settings_wsc_change_sw(w http.ResponseWriter, r *http.Request) {
	if IsolatedAuth(w, r, []gauth.TRole{gauth.ADMIN}) {
		MTemplates[BLOCK_TEMP_ACCESS_DENIED].Execute(w, nil)
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
		MTemplates[BLOCK_TEMP_ACCESS_DENIED].Execute(w, nil)
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
		MTemplates[BLOCK_TEMP_ACCESS_DENIED].Execute(w, nil)
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
		MTemplates[BLOCK_TEMP_ACCESS_DENIED].Execute(w, nil)
		return
	}

	slog.Warn("This service cannot be disabled.", slog.String("service", "WebServer"))

	nav_settings(w, r)
}
