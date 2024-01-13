package webmanage

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/Kwynto/GracefulDB/internal/engine/basicsystem/gauth"
)

func Test_homeDefault(t *testing.T) {
	gauth.Start()
	parseTemplates()

	t.Run("homeDefault() function testing - negative", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/", nil)

		homeDefault(w, r) // calling the tested function
		status := w.Code
		if status != http.StatusFound {
			t.Errorf("homeDefault() error: %v", status)
		}
	})

	t.Run("homeDefault() function testing - positive", func(t *testing.T) {
		w := httptest.NewRecorder()
		form := url.Values{}
		form.Add("username", "root")
		form.Add("password", "toor")
		r := httptest.NewRequest("POST", "/", strings.NewReader(form.Encode()))
		r.PostForm = form

		homeAuth(w, r)
		wCooks := w.Result().Cookies()

		w1 := httptest.NewRecorder()
		r1 := httptest.NewRequest("GET", "/", nil)
		for _, v := range wCooks {
			r1.AddCookie(&http.Cookie{
				Name:   v.Name,
				Value:  v.Value,
				MaxAge: v.MaxAge,
			})
		}

		homeDefault(w1, r1) // calling the tested function
		status := w1.Code
		if status != http.StatusOK {
			t.Errorf("homeDefault() error: %v", status)
		}
	})

	t.Run("homeDefault() function testing - logout", func(t *testing.T) {
		randStr := gauth.GenerateTicket()
		prof := gauth.TProfile{
			Description: "Testing description",
			Status:      gauth.ACTIVE,
			Roles:       []gauth.TRole{gauth.USER},
		}
		gauth.AddUser(randStr, randStr, prof)

		w := httptest.NewRecorder()
		form := url.Values{}
		form.Add("username", randStr)
		form.Add("password", randStr)
		r := httptest.NewRequest("POST", "/", strings.NewReader(form.Encode()))
		r.PostForm = form

		homeAuth(w, r)
		wCooks := w.Result().Cookies()

		delete(gauth.AccessMap, randStr)

		w1 := httptest.NewRecorder()
		r1 := httptest.NewRequest("GET", "/", nil)
		for _, v := range wCooks {
			r1.AddCookie(&http.Cookie{
				Name:   v.Name,
				Value:  v.Value,
				MaxAge: v.MaxAge,
			})
		}

		homeDefault(w1, r1) // calling the tested function
		status := w1.Code
		if status != http.StatusFound {
			t.Errorf("homeDefault() error: %v", status)
		}
	})

	t.Run("homeDefault() function testing - Template error", func(t *testing.T) {
		w := httptest.NewRecorder()
		form := url.Values{}
		form.Add("username", "root")
		form.Add("password", "toor")
		r := httptest.NewRequest("POST", "/", strings.NewReader(form.Encode()))
		r.PostForm = form

		homeAuth(w, r)
		wCooks := w.Result().Cookies()

		w1 := httptest.NewRecorder()
		r1 := httptest.NewRequest("GET", "/", nil)
		for _, v := range wCooks {
			r1.AddCookie(&http.Cookie{
				Name:   v.Name,
				Value:  v.Value,
				MaxAge: v.MaxAge,
			})
		}

		wrongStr := `
		<html>
			{{ .errorexp }}
		</html>
		`
		loadTemplateFromVar(HOME_TEMP_NAME, wrongStr)

		homeDefault(w1, r1) // calling the tested function
		status := w1.Code
		if status != http.StatusOK {
			t.Errorf("homeDefault() error: %v", status)
		}
	})
}

func Test_homeAuth(t *testing.T) {
	gauth.Start()
	parseTemplates()

	t.Run("homeAuth() function testing - GET negative", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/", nil)

		homeAuth(w, r) // calling the tested function
		status := w.Code
		if status != http.StatusOK {
			t.Errorf("homeAuth() error: %v", status)
		}
	})

	t.Run("homeAuth() function testing - POST negative", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest("POST", "/", nil)
		r.PostForm = nil
		r.Body = nil

		homeAuth(w, r)

		status := w.Code
		if status != http.StatusBadRequest {
			t.Errorf("homeAuth() error: %v", status)
		}
	})

}

func Test_home(t *testing.T) {
	gauth.Start()
	parseTemplates()

	t.Run("home() function testing - negative", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/a", nil)

		home(w, r) // calling the tested function
		status := w.Code
		if status != http.StatusFound {
			t.Errorf("home() error: %v", status)
		}
	})

	t.Run("home() function testing", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/", nil)

		home(w, r) // calling the tested function
		status := w.Code
		if status != http.StatusOK {
			t.Errorf("home() error: %v", status)
		}
	})

	t.Run("home() function testing", func(t *testing.T) {
		randStr := gauth.GenerateTicket()
		prof := gauth.TProfile{
			Description: "Testing description",
			Status:      gauth.ACTIVE,
			Roles:       []gauth.TRole{gauth.USER},
		}
		gauth.AddUser(randStr, randStr, prof)

		w := httptest.NewRecorder()
		form := url.Values{}
		form.Add("username", randStr)
		form.Add("password", randStr)
		r := httptest.NewRequest("POST", "/", strings.NewReader(form.Encode()))
		r.PostForm = form

		homeAuth(w, r)
		wCooks := w.Result().Cookies()

		w1 := httptest.NewRecorder()
		r1 := httptest.NewRequest("GET", "/", nil)
		for _, v := range wCooks {
			r1.AddCookie(&http.Cookie{
				Name:   v.Name,
				Value:  v.Value,
				MaxAge: v.MaxAge,
			})
		}

		home(w1, r1) // calling the tested function
		status := w1.Code
		if status != http.StatusOK {
			t.Errorf("home() error: %v", status)
		}
	})
}

func Test_nav_default(t *testing.T) {
	gauth.Start()
	parseTemplates()

	t.Run("nav_default() function testing", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/", nil)

		nav_default(w, r) // calling the tested function
		status := w.Code
		if status != http.StatusOK {
			t.Errorf("nav_default() error: %v", status)
		}
	})
}

func Test_nav_logout(t *testing.T) {
	gauth.Start()
	parseTemplates()

	t.Run("nav_logout() function testing", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/", nil)

		nav_logout(w, r) // calling the tested function
		status := w.Code
		if status != http.StatusOK {
			t.Errorf("nav_logout() error: %v", status)
		}
	})
}
