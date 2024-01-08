package rest

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func Test_squery(t *testing.T) {
	t.Run("squery() function testing - GET error", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/", nil)

		squery(w, r) // calling the tested function
		status := w.Code
		if status != http.StatusMethodNotAllowed {
			t.Error("squery() error. GET error.")
		}
	})

	t.Run("squery() function testing - POST - an empty query", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest("POST", "/", nil)

		squery(w, r) // calling the tested function
		status := w.Code
		if status != http.StatusBadRequest {
			t.Error("squery() error. POST - an empty query.")
		}
	})

	t.Run("squery() function testing - POST - positive", func(t *testing.T) {
		w := httptest.NewRecorder()

		form := url.Values{}
		form.Add("instruction", "instruction")
		form.Add("placeholder", "[]")
		r := httptest.NewRequest("POST", "/", strings.NewReader(form.Encode()))
		// r.Form = form
		r.PostForm = form

		squery(w, r) // calling the tested function
		status := w.Code
		if status != http.StatusOK {
			t.Errorf("squery() error. POST - an empty query. %v", status)
		}
	})
}
