package webmanage

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/Kwynto/gosession"

	"github.com/Kwynto/GracefulDB/internal/base/basicsystem/gauth"
)

// Handler after authorization
func homeDefault(w http.ResponseWriter, r *http.Request, login string) {
	// This function is complete
	err := templatesMap[HOME_TEMP_NAME].Execute(w, nil)
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
		err := templatesMap[AUTH_TEMP_NAME].Execute(w, nil)
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
	err := templatesMap[BLOCK_TEMP_DEFAULT].Execute(w, nil)
	if err != nil {
		slog.Debug("Internal Server Error", slog.String("err", err.Error()))
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func nav_logout(w http.ResponseWriter, r *http.Request) {
	// This function is complete
	sesID := gosession.Start(&w, r)
	sesID.Remove("auth")
	w.Header().Set("HX-Redirect", "/log.out")
	// http.Redirect(w, r, "/", http.StatusFound)
}

func nav_dashboard(w http.ResponseWriter, r *http.Request) {
	err := templatesMap[BLOCK_TEMP_DASHBOARD].Execute(w, nil)
	if err != nil {
		slog.Debug("Internal Server Error", slog.String("err", err.Error()))
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func nav_databases(w http.ResponseWriter, r *http.Request) {
	err := templatesMap[BLOCK_TEMP_DATABASES].Execute(w, nil)
	if err != nil {
		slog.Debug("Internal Server Error", slog.String("err", err.Error()))
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func nav_accounts(w http.ResponseWriter, r *http.Request) {
	templatesMap[BLOCK_TEMP_ACCOUNTS].Execute(w, nil)
	// err := templatesMap[BLOCK_TEMP_ACCOUNTS].Execute(w, nil)
	// if err != nil {
	// 	slog.Debug("Internal Server Error", slog.String("err", err.Error()))
	// 	http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	// }
}

func nav_settings(w http.ResponseWriter, r *http.Request) {
	templatesMap[BLOCK_TEMP_SETTINGS].Execute(w, nil)
	// err := templatesMap[BLOCK_TEMP_SETTINGS].Execute(w, nil)
	// if err != nil {
	// 	slog.Debug("Internal Server Error", slog.String("err", err.Error()))
	// 	http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	// }
}
