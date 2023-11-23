package webmanage

import (
	"fmt"
	"html/template"
	"log/slog"
	"net/http"

	"github.com/Kwynto/gosession"
)

func homeDefault(w http.ResponseWriter, r *http.Request, login string) {
	ts, err := template.ParseFiles("./ui/html/home.html")
	if err != nil {
		slog.Debug("Internal Server Error", slog.String("err", err.Error()))
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	err = ts.Execute(w, nil)
	if err != nil {
		slog.Debug("Internal Server Error", slog.String("err", err.Error()))
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func homeAuth(w http.ResponseWriter, r *http.Request) {
	ts, err := template.ParseFiles("./ui/html/auth.html")
	if err != nil {
		slog.Debug("Internal Server Error", slog.String("err", err.Error()))
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	err = ts.Execute(w, nil)
	if err != nil {
		slog.Debug("Internal Server Error", slog.String("err", err.Error()))
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
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
