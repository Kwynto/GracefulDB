package webmanage

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/Kwynto/gosession"

	"github.com/Kwynto/GracefulDB/internal/base/basicsystem/gauth"
	"github.com/Kwynto/GracefulDB/internal/config"
	"github.com/Kwynto/GracefulDB/internal/connectors/grpc"
	"github.com/Kwynto/GracefulDB/internal/connectors/rest"
	"github.com/Kwynto/GracefulDB/internal/connectors/websocketconn"

	"github.com/Kwynto/GracefulDB/pkg/lib/closer"
)

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

// Nav Menu Handlers
func nav_default(w http.ResponseWriter, r *http.Request) {
	// This function is complete
	TemplatesMap[BLOCK_TEMP_DEFAULT].Execute(w, nil)
}

func nav_logout(w http.ResponseWriter, r *http.Request) {
	// This function is complete
	w.Header().Set("HX-Redirect", "/log.out")
}

func nav_dashboard(w http.ResponseWriter, r *http.Request) {
	TemplatesMap[BLOCK_TEMP_DASHBOARD].Execute(w, nil)
}

func nav_databases(w http.ResponseWriter, r *http.Request) {
	TemplatesMap[BLOCK_TEMP_DATABASES].Execute(w, nil)
}

func nav_accounts(w http.ResponseWriter, r *http.Request) {
	TemplatesMap[BLOCK_TEMP_ACCOUNTS].Execute(w, nil)
}

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
