package webmanage

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/Kwynto/GracefulDB/internal/config"
	"github.com/Kwynto/GracefulDB/internal/connectors/rest"
	"github.com/Kwynto/GracefulDB/internal/connectors/websocketconn"
	"github.com/Kwynto/GracefulDB/internal/engine/basicsystem/gauth"
	"github.com/Kwynto/GracefulDB/internal/engine/core"
	"github.com/Kwynto/GracefulDB/pkg/lib/closer"
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

func Test_selfedit_load_form(t *testing.T) {
	gauth.Start()
	parseTemplates()

	t.Run("selfedit_load_form() function testing - negative", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/", nil)

		selfedit_load_form(w, r) // calling the tested function
		status := w.Code
		if status != http.StatusOK {
			t.Errorf("selfedit_load_form() error: %v", status)
		}
	})

	t.Run("selfedit_load_form() function testing - positive", func(t *testing.T) {
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

		selfedit_load_form(w1, r1) // calling the tested function
		status := w1.Code
		if status != http.StatusOK {
			t.Errorf("selfedit_load_form() error: %v", status)
		}
	})

	t.Run("selfedit_load_form() function testing", func(t *testing.T) {
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

		selfedit_load_form(w1, r1) // calling the tested function
		status := w1.Code
		if status != http.StatusOK {
			t.Errorf("selfedit_load_form() error: %v", status)
		}
	})
}

func Test_selfedit_ok(t *testing.T) {
	gauth.Start()
	parseTemplates()

	t.Run("selfedit_ok() function testing - Isolate negative", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/", nil)

		selfedit_ok(w, r) // calling the tested function
		status := w.Code
		if status != http.StatusOK {
			t.Errorf("selfedit_ok() error: %v", status)
		}
	})

	t.Run("selfedit_ok() function testing - POST negative", func(t *testing.T) {
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

		selfedit_ok(w1, r1) // calling the tested function
		status := w1.Code
		if status != http.StatusOK {
			t.Errorf("selfedit_ok() error: %v", status)
		}
	})

	t.Run("selfedit_ok() function testing POST positive and don't work ParseForm", func(t *testing.T) {
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
		r1 := httptest.NewRequest("POST", "/", nil)
		r1.PostForm = nil
		r1.Body = nil
		for _, v := range wCooks {
			r1.AddCookie(&http.Cookie{
				Name:   v.Name,
				Value:  v.Value,
				MaxAge: v.MaxAge,
			})
		}

		selfedit_ok(w1, r1) // calling the tested function
		status := w1.Code
		if status != http.StatusOK {
			t.Errorf("selfedit_ok() error: %v", status)
		}
	})

	t.Run("selfedit_ok() function testing POST positive and not value of password", func(t *testing.T) {
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
		r1 := httptest.NewRequest("POST", "/", nil)
		for _, v := range wCooks {
			r1.AddCookie(&http.Cookie{
				Name:   v.Name,
				Value:  v.Value,
				MaxAge: v.MaxAge,
			})
		}

		selfedit_ok(w1, r1) // calling the tested function
		status := w1.Code
		if status != http.StatusOK {
			t.Errorf("selfedit_ok() error: %v", status)
		}
	})

	t.Run("selfedit_ok() function testing - all right", func(t *testing.T) {
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
		form1 := url.Values{}
		form1.Add("password", randStr)
		form1.Add("desc", randStr)
		r1 := httptest.NewRequest("POST", "/", strings.NewReader(form1.Encode()))
		r1.PostForm = form1
		for _, v := range wCooks {
			r1.AddCookie(&http.Cookie{
				Name:   v.Name,
				Value:  v.Value,
				MaxAge: v.MaxAge,
			})
		}

		selfedit_ok(w1, r1) // calling the tested function
		status := w1.Code
		if status != http.StatusOK {
			t.Errorf("selfedit_ok() error: %v", status)
		}
	})
}

func Test_nav_dashboard(t *testing.T) {
	gauth.Start()
	parseTemplates()

	t.Run("nav_dashboard() function testing - Isolate negative", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/", nil)

		nav_dashboard(w, r) // calling the tested function
		status := w.Code
		if status != http.StatusOK {
			t.Errorf("nav_dashboard() error: %v", status)
		}
	})

	t.Run("nav_dashboard() function testing - Isolate positive", func(t *testing.T) {
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

		nav_dashboard(w1, r1) // calling the tested function
		status := w1.Code
		if status != http.StatusOK {
			t.Errorf("nav_dashboard() error: %v", status)
		}
	})
}

func Test_nav_databases(t *testing.T) {
	gauth.Start()
	parseTemplates()

	t.Run("nav_databases() function testing - Isolate negative", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/", nil)

		nav_databases(w, r) // calling the tested function
		status := w.Code
		if status != http.StatusOK {
			t.Errorf("nav_databases() error: %v", status)
		}
	})

	t.Run("nav_databases() function testing - Isolate positive", func(t *testing.T) {
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

		nav_databases(w1, r1) // calling the tested function
		status := w1.Code
		if status != http.StatusOK {
			t.Errorf("nav_databases() error: %v", status)
		}
	})
}

func Test_database_request(t *testing.T) {
	gauth.Start()
	parseTemplates()

	t.Run("database_request() function testing - Isolate negative", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/", nil)

		database_request(w, r) // calling the tested function
		status := w.Code
		if status != http.StatusOK {
			t.Errorf("database_request() error: %v", status)
		}
	})

	t.Run("database_request() function testing - POST negative", func(t *testing.T) {
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

		database_request(w1, r1) // calling the tested function
		status := w1.Code
		if status != http.StatusOK {
			t.Errorf("database_request() error: %v", status)
		}
	})

	t.Run("database_request() function testing POST positive and don't work ParseForm", func(t *testing.T) {
		randStr := gauth.GenerateTicket()
		prof := gauth.TProfile{
			Description: "Testing description",
			Status:      gauth.ACTIVE,
			Roles:       []gauth.TRole{gauth.ENGINEER},
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
		r1 := httptest.NewRequest("POST", "/", nil)
		r1.PostForm = nil
		r1.Body = nil
		for _, v := range wCooks {
			r1.AddCookie(&http.Cookie{
				Name:   v.Name,
				Value:  v.Value,
				MaxAge: v.MaxAge,
			})
		}

		database_request(w1, r1) // calling the tested function
		status := w1.Code
		if status != http.StatusOK {
			t.Errorf("database_request() error: %v", status)
		}
	})

	t.Run("database_request() function testing - all right", func(t *testing.T) {
		randStr := gauth.GenerateTicket()
		prof := gauth.TProfile{
			Description: "Testing description",
			Status:      gauth.ACTIVE,
			Roles:       []gauth.TRole{gauth.ENGINEER},
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
		form1 := url.Values{}
		form1.Add("request", randStr)
		r1 := httptest.NewRequest("POST", "/", strings.NewReader(form1.Encode()))
		r1.PostForm = form1
		for _, v := range wCooks {
			r1.AddCookie(&http.Cookie{
				Name:   v.Name,
				Value:  v.Value,
				MaxAge: v.MaxAge,
			})
		}

		database_request(w1, r1) // calling the tested function
		status := w1.Code
		if status != http.StatusOK {
			t.Errorf("database_request() error: %v", status)
		}
	})
}

func Test_nav_accounts(t *testing.T) {
	gauth.Start()
	parseTemplates()

	t.Run("nav_accounts() function testing - Isolate negative", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/", nil)

		nav_accounts(w, r) // calling the tested function
		status := w.Code
		if status != http.StatusOK {
			t.Errorf("nav_accounts() error: %v", status)
		}
	})

	t.Run("nav_accounts() function testing - Isolate positive", func(t *testing.T) {
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

		nav_accounts(w1, r1) // calling the tested function
		status := w1.Code
		if status != http.StatusOK {
			t.Errorf("nav_accounts() error: %v", status)
		}
	})

	t.Run("nav_accounts() function testing - all coverage", func(t *testing.T) {
		randStr := gauth.GenerateTicket()
		prof := gauth.TProfile{
			Description: "Testing description",
			Status:      gauth.ACTIVE,
			Roles:       []gauth.TRole{gauth.SYSTEM},
		}
		gauth.AddUser(randStr, randStr, prof)

		randStr1 := gauth.GenerateTicket()
		prof1 := gauth.TProfile{
			Description: "Testing description",
			Status:      gauth.BANED,
			Roles:       []gauth.TRole{gauth.USER},
		}
		gauth.AddUser(randStr1, randStr1, prof1)

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

		nav_accounts(w1, r1) // calling the tested function
		status := w1.Code
		if status != http.StatusOK {
			t.Errorf("nav_accounts() error: %v", status)
		}
	})
}

func Test_account_create_load_form(t *testing.T) {
	gauth.Start()
	parseTemplates()

	t.Run("account_create_load_form() function testing - Isolate negative", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/", nil)

		account_create_load_form(w, r) // calling the tested function
		status := w.Code
		if status != http.StatusOK {
			t.Errorf("account_create_load_form() error: %v", status)
		}
	})

	t.Run("account_create_load_form() function testing - Isolate positive", func(t *testing.T) {
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

		account_create_load_form(w1, r1) // calling the tested function
		status := w1.Code
		if status != http.StatusOK {
			t.Errorf("account_create_load_form() error: %v", status)
		}
	})
}

func Test_account_create_ok(t *testing.T) {
	gauth.Start()
	parseTemplates()

	t.Run("account_create_ok() function testing - Isolate negative", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/", nil)

		account_create_ok(w, r) // calling the tested function
		status := w.Code
		if status != http.StatusOK {
			t.Errorf("account_create_ok() error: %v", status)
		}
	})

	t.Run("account_create_ok() function testing - POST negative", func(t *testing.T) {
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

		account_create_ok(w1, r1) // calling the tested function
		status := w1.Code
		if status != http.StatusOK {
			t.Errorf("account_create_ok() error: %v", status)
		}
	})

	t.Run("account_create_ok() function testing POST positive and don't work ParseForm", func(t *testing.T) {
		randStr := gauth.GenerateTicket()
		prof := gauth.TProfile{
			Description: "Testing description",
			Status:      gauth.ACTIVE,
			Roles:       []gauth.TRole{gauth.MANAGER},
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
		r1 := httptest.NewRequest("POST", "/", nil)
		r1.PostForm = nil
		r1.Body = nil
		for _, v := range wCooks {
			r1.AddCookie(&http.Cookie{
				Name:   v.Name,
				Value:  v.Value,
				MaxAge: v.MaxAge,
			})
		}

		account_create_ok(w1, r1) // calling the tested function
		status := w1.Code
		if status != http.StatusOK {
			t.Errorf("account_create_ok() error: %v", status)
		}
	})

	t.Run("account_create_ok() function testing POST positive and not value of password", func(t *testing.T) {
		randStr := gauth.GenerateTicket()
		prof := gauth.TProfile{
			Description: "Testing description",
			Status:      gauth.ACTIVE,
			Roles:       []gauth.TRole{gauth.MANAGER},
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

		randStr1 := gauth.GenerateTicket()
		w1 := httptest.NewRecorder()
		form1 := url.Values{}
		form1.Add("login", randStr1)
		form1.Add("password", "")
		form1.Add("desc", randStr1)
		r1 := httptest.NewRequest("POST", "/", strings.NewReader(form1.Encode()))
		r1.PostForm = form1
		for _, v := range wCooks {
			r1.AddCookie(&http.Cookie{
				Name:   v.Name,
				Value:  v.Value,
				MaxAge: v.MaxAge,
			})
		}

		account_create_ok(w1, r1) // calling the tested function
		status := w1.Code
		if status != http.StatusOK {
			t.Errorf("account_create_ok() error: %v", status)
		}
	})

	t.Run("account_create_ok() function testing - create error", func(t *testing.T) {
		randStr := gauth.GenerateTicket()
		prof := gauth.TProfile{
			Description: "Testing description",
			Status:      gauth.ACTIVE,
			Roles:       []gauth.TRole{gauth.MANAGER},
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
		form1 := url.Values{}
		form1.Add("login", "root")
		form1.Add("password", randStr)
		form1.Add("desc", randStr)
		r1 := httptest.NewRequest("POST", "/", strings.NewReader(form1.Encode()))
		r1.PostForm = form1
		for _, v := range wCooks {
			r1.AddCookie(&http.Cookie{
				Name:   v.Name,
				Value:  v.Value,
				MaxAge: v.MaxAge,
			})
		}

		account_create_ok(w1, r1) // calling the tested function
		status := w1.Code
		if status != http.StatusOK {
			t.Errorf("account_create_ok() error: %v", status)
		}
	})

	t.Run("account_create_ok() function testing - all right", func(t *testing.T) {
		randStr := gauth.GenerateTicket()
		prof := gauth.TProfile{
			Description: "Testing description",
			Status:      gauth.ACTIVE,
			Roles:       []gauth.TRole{gauth.MANAGER},
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

		randStr1 := gauth.GenerateTicket()
		w1 := httptest.NewRecorder()
		form1 := url.Values{}
		form1.Add("login", randStr1)
		form1.Add("password", randStr1)
		form1.Add("desc", randStr1)
		r1 := httptest.NewRequest("POST", "/", strings.NewReader(form1.Encode()))
		r1.PostForm = form1
		for _, v := range wCooks {
			r1.AddCookie(&http.Cookie{
				Name:   v.Name,
				Value:  v.Value,
				MaxAge: v.MaxAge,
			})
		}

		account_create_ok(w1, r1) // calling the tested function
		status := w1.Code
		if status != http.StatusOK {
			t.Errorf("account_create_ok() error: %v", status)
		}
	})
}

func Test_account_edit_load_form(t *testing.T) {
	gauth.Start()
	parseTemplates()

	t.Run("account_edit_load_form() function testing - Isolate negative", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/", nil)

		account_edit_load_form(w, r) // calling the tested function
		status := w.Code
		if status != http.StatusOK {
			t.Errorf("account_edit_load_form() error: %v", status)
		}
	})

	t.Run("account_edit_load_form() function testing - Isolate positive and don't GetProfile", func(t *testing.T) {
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

		account_edit_load_form(w1, r1) // calling the tested function
		status := w1.Code
		if status != http.StatusOK {
			t.Errorf("account_edit_load_form() error: %v", status)
		}
	})

	t.Run("account_edit_load_form() function testing - all right", func(t *testing.T) {
		w := httptest.NewRecorder()
		form := url.Values{}
		form.Add("username", "root")
		form.Add("password", "toor")
		r := httptest.NewRequest("POST", "/", strings.NewReader(form.Encode()))
		r.PostForm = form

		homeAuth(w, r)
		wCooks := w.Result().Cookies()

		randStr := gauth.GenerateTicket()
		prof := gauth.TProfile{
			Description: "Testing description",
			Status:      gauth.ACTIVE,
			Roles:       []gauth.TRole{gauth.USER, gauth.SYSTEM},
		}
		gauth.AddUser(randStr, randStr, prof)

		w1 := httptest.NewRecorder()
		testurl := fmt.Sprintf("/?user=%s", randStr)
		r1 := httptest.NewRequest("GET", testurl, nil)
		for _, v := range wCooks {
			r1.AddCookie(&http.Cookie{
				Name:   v.Name,
				Value:  v.Value,
				MaxAge: v.MaxAge,
			})
		}

		account_edit_load_form(w1, r1) // calling the tested function
		status := w1.Code
		if status != http.StatusOK {
			t.Errorf("account_edit_load_form() error: %v", status)
		}
	})
}

func Test_account_edit_ok(t *testing.T) {
	gauth.Start()
	parseTemplates()

	t.Run("account_edit_ok() function testing - Isolate negative", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/", nil)

		account_edit_ok(w, r) // calling the tested function
		status := w.Code
		if status != http.StatusOK {
			t.Errorf("account_edit_ok() error: %v", status)
		}
	})

	t.Run("account_edit_ok() function testing - POST negative", func(t *testing.T) {
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

		account_edit_ok(w1, r1) // calling the tested function
		status := w1.Code
		if status != http.StatusOK {
			t.Errorf("account_edit_ok() error: %v", status)
		}
	})

	t.Run("account_edit_ok() function testing POST positive and don't work ParseForm", func(t *testing.T) {
		randStr := gauth.GenerateTicket()
		prof := gauth.TProfile{
			Description: "Testing description",
			Status:      gauth.ACTIVE,
			Roles:       []gauth.TRole{gauth.MANAGER},
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
		r1 := httptest.NewRequest("POST", "/", nil)
		r1.PostForm = nil
		r1.Body = nil
		for _, v := range wCooks {
			r1.AddCookie(&http.Cookie{
				Name:   v.Name,
				Value:  v.Value,
				MaxAge: v.MaxAge,
			})
		}

		account_edit_ok(w1, r1) // calling the tested function
		status := w1.Code
		if status != http.StatusOK {
			t.Errorf("account_edit_ok() error: %v", status)
		}
	})

	t.Run("account_edit_ok() function testing POST positive and not value of login", func(t *testing.T) {
		randStr := gauth.GenerateTicket()
		prof := gauth.TProfile{
			Description: "Testing description",
			Status:      gauth.ACTIVE,
			Roles:       []gauth.TRole{gauth.MANAGER},
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

		randStr1 := gauth.GenerateTicket()
		w1 := httptest.NewRecorder()
		form1 := url.Values{}
		form1.Add("login", "")
		form1.Add("password", "")
		form1.Add("desc", randStr1)
		r1 := httptest.NewRequest("POST", "/", strings.NewReader(form1.Encode()))
		r1.PostForm = form1
		for _, v := range wCooks {
			r1.AddCookie(&http.Cookie{
				Name:   v.Name,
				Value:  v.Value,
				MaxAge: v.MaxAge,
			})
		}

		account_edit_ok(w1, r1) // calling the tested function
		status := w1.Code
		if status != http.StatusOK {
			t.Errorf("account_edit_ok() error: %v", status)
		}
	})

	t.Run("account_edit_ok() function testing POST positive and not value of password", func(t *testing.T) {
		randStr := gauth.GenerateTicket()
		prof := gauth.TProfile{
			Description: "Testing description",
			Status:      gauth.ACTIVE,
			Roles:       []gauth.TRole{gauth.MANAGER},
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

		randStr1 := gauth.GenerateTicket()
		w1 := httptest.NewRecorder()
		form1 := url.Values{}
		form1.Add("login", randStr1)
		form1.Add("password", "")
		form1.Add("desc", randStr1)
		r1 := httptest.NewRequest("POST", "/", strings.NewReader(form1.Encode()))
		r1.PostForm = form1
		for _, v := range wCooks {
			r1.AddCookie(&http.Cookie{
				Name:   v.Name,
				Value:  v.Value,
				MaxAge: v.MaxAge,
			})
		}

		account_edit_ok(w1, r1) // calling the tested function
		status := w1.Code
		if status != http.StatusOK {
			t.Errorf("account_edit_ok() error: %v", status)
		}
	})

	t.Run("account_edit_ok() function testing - status error", func(t *testing.T) {
		randStr := gauth.GenerateTicket()
		prof := gauth.TProfile{
			Description: "Testing description",
			Status:      gauth.ACTIVE,
			Roles:       []gauth.TRole{gauth.MANAGER},
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
		form1 := url.Values{}
		form1.Add("login", randStr)
		form1.Add("password", randStr)
		form1.Add("desc", randStr)
		form1.Add("status", "0")
		r1 := httptest.NewRequest("POST", "/", strings.NewReader(form1.Encode()))
		r1.PostForm = form1
		for _, v := range wCooks {
			r1.AddCookie(&http.Cookie{
				Name:   v.Name,
				Value:  v.Value,
				MaxAge: v.MaxAge,
			})
		}

		account_edit_ok(w1, r1) // calling the tested function
		status := w1.Code
		if status != http.StatusOK {
			t.Errorf("account_edit_ok() error: %v", status)
		}
	})

	t.Run("account_edit_ok() function testing - Roles", func(t *testing.T) {
		randStr := gauth.GenerateTicket()
		prof := gauth.TProfile{
			Description: "Testing description",
			Status:      gauth.ACTIVE,
			Roles:       []gauth.TRole{gauth.MANAGER},
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

		randStr1 := gauth.GenerateTicket()
		prof1 := gauth.TProfile{
			Description: "Testing description",
			Status:      gauth.ACTIVE,
			Roles:       []gauth.TRole{gauth.USER},
		}
		gauth.AddUser(randStr1, randStr1, prof1)

		w1 := httptest.NewRecorder()
		form1 := url.Values{}
		form1.Add("login", randStr1)
		form1.Add("password", randStr1)
		form1.Add("desc", randStr1)
		form1.Add("status", "2")
		chtext := `["SYSTEM", "ADMIN", "MANAGER", "ENGINEER", "USER"]`
		chbase := base64.StdEncoding.EncodeToString([]byte(chtext))
		form1.Add("role_names", chbase)
		r1 := httptest.NewRequest("POST", "/", strings.NewReader(form1.Encode()))
		r1.PostForm = form1
		for _, v := range wCooks {
			r1.AddCookie(&http.Cookie{
				Name:   v.Name,
				Value:  v.Value,
				MaxAge: v.MaxAge,
			})
		}

		account_edit_ok(w1, r1) // calling the tested function
		status := w1.Code
		if status != http.StatusOK {
			t.Errorf("account_edit_ok() error: %v", status)
		}
	})

	t.Run("account_edit_ok() function testing - root user", func(t *testing.T) {
		randStr := gauth.GenerateTicket()
		prof := gauth.TProfile{
			Description: "Testing description",
			Status:      gauth.ACTIVE,
			Roles:       []gauth.TRole{gauth.MANAGER},
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
		form1 := url.Values{}
		form1.Add("login", "root")
		form1.Add("password", "toor")
		form1.Add("desc", "root")
		form1.Add("status", "2")
		chtext := `["SYSTEM", "ADMIN", "MANAGER", "ENGINEER", "USER"]`
		chbase := base64.StdEncoding.EncodeToString([]byte(chtext))
		form1.Add("role_names", chbase)
		r1 := httptest.NewRequest("POST", "/", strings.NewReader(form1.Encode()))
		r1.PostForm = form1
		for _, v := range wCooks {
			r1.AddCookie(&http.Cookie{
				Name:   v.Name,
				Value:  v.Value,
				MaxAge: v.MaxAge,
			})
		}

		account_edit_ok(w1, r1) // calling the tested function
		status := w1.Code
		if status != http.StatusOK {
			t.Errorf("account_edit_ok() error: %v", status)
		}
	})
}

func Test_account_ban_load_form(t *testing.T) {
	gauth.Start()
	parseTemplates()

	t.Run("account_ban_load_form() function testing - Isolate negative", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/", nil)

		account_ban_load_form(w, r) // calling the tested function
		status := w.Code
		if status != http.StatusOK {
			t.Errorf("account_ban_load_form() error: %v", status)
		}
	})

	t.Run("account_ban_load_form() function testing - root user", func(t *testing.T) {
		w := httptest.NewRecorder()
		form := url.Values{}
		form.Add("username", "root")
		form.Add("password", "toor")
		r := httptest.NewRequest("POST", "/", strings.NewReader(form.Encode()))
		r.PostForm = form

		homeAuth(w, r)
		wCooks := w.Result().Cookies()

		randStr := gauth.GenerateTicket()
		prof := gauth.TProfile{
			Description: "Testing description",
			Status:      gauth.ACTIVE,
			Roles:       []gauth.TRole{gauth.USER},
		}
		gauth.AddUser(randStr, randStr, prof)

		w1 := httptest.NewRecorder()
		testurl := "/?user=root"
		r1 := httptest.NewRequest("GET", testurl, nil)
		for _, v := range wCooks {
			r1.AddCookie(&http.Cookie{
				Name:   v.Name,
				Value:  v.Value,
				MaxAge: v.MaxAge,
			})
		}

		account_ban_load_form(w1, r1) // calling the tested function
		status := w1.Code
		if status != http.StatusOK {
			t.Errorf("account_ban_load_form() error: %v", status)
		}
	})

	t.Run("account_ban_load_form() function testing - all right", func(t *testing.T) {
		w := httptest.NewRecorder()
		form := url.Values{}
		form.Add("username", "root")
		form.Add("password", "toor")
		r := httptest.NewRequest("POST", "/", strings.NewReader(form.Encode()))
		r.PostForm = form

		homeAuth(w, r)
		wCooks := w.Result().Cookies()

		randStr := gauth.GenerateTicket()
		prof := gauth.TProfile{
			Description: "Testing description",
			Status:      gauth.ACTIVE,
			Roles:       []gauth.TRole{gauth.USER},
		}
		gauth.AddUser(randStr, randStr, prof)

		w1 := httptest.NewRecorder()
		testurl := fmt.Sprintf("/?user=%s", randStr)
		r1 := httptest.NewRequest("GET", testurl, nil)
		for _, v := range wCooks {
			r1.AddCookie(&http.Cookie{
				Name:   v.Name,
				Value:  v.Value,
				MaxAge: v.MaxAge,
			})
		}

		account_ban_load_form(w1, r1) // calling the tested function
		status := w1.Code
		if status != http.StatusOK {
			t.Errorf("account_ban_load_form() error: %v", status)
		}
	})
}

func Test_account_ban_ok(t *testing.T) {
	gauth.Start()
	parseTemplates()

	t.Run("account_ban_ok() function testing - Isolate negative", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/", nil)

		account_ban_ok(w, r) // calling the tested function
		status := w.Code
		if status != http.StatusOK {
			t.Errorf("account_ban_ok() error: %v", status)
		}
	})

	t.Run("account_ban_ok() function testing - POST negative", func(t *testing.T) {
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

		account_ban_ok(w1, r1) // calling the tested function
		status := w1.Code
		if status != http.StatusOK {
			t.Errorf("account_ban_ok() error: %v", status)
		}
	})

	t.Run("account_ban_ok() function testing POST positive and don't work ParseForm", func(t *testing.T) {
		randStr := gauth.GenerateTicket()
		prof := gauth.TProfile{
			Description: "Testing description",
			Status:      gauth.ACTIVE,
			Roles:       []gauth.TRole{gauth.MANAGER},
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
		r1 := httptest.NewRequest("POST", "/", nil)
		r1.PostForm = nil
		r1.Body = nil
		for _, v := range wCooks {
			r1.AddCookie(&http.Cookie{
				Name:   v.Name,
				Value:  v.Value,
				MaxAge: v.MaxAge,
			})
		}

		account_ban_ok(w1, r1) // calling the tested function
		status := w1.Code
		if status != http.StatusOK {
			t.Errorf("account_ban_ok() error: %v", status)
		}
	})

	t.Run("account_ban_ok() function testing POST positive and not value of login", func(t *testing.T) {
		randStr := gauth.GenerateTicket()
		prof := gauth.TProfile{
			Description: "Testing description",
			Status:      gauth.ACTIVE,
			Roles:       []gauth.TRole{gauth.MANAGER},
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
		form1 := url.Values{}
		form1.Add("login", "")
		r1 := httptest.NewRequest("POST", "/", strings.NewReader(form1.Encode()))
		r1.PostForm = form1
		for _, v := range wCooks {
			r1.AddCookie(&http.Cookie{
				Name:   v.Name,
				Value:  v.Value,
				MaxAge: v.MaxAge,
			})
		}

		account_ban_ok(w1, r1) // calling the tested function
		status := w1.Code
		if status != http.StatusOK {
			t.Errorf("account_ban_ok() error: %v", status)
		}
	})

	t.Run("account_ban_ok() function testing POST positive and don't work BlockUser", func(t *testing.T) {
		randStr := gauth.GenerateTicket()
		prof := gauth.TProfile{
			Description: "Testing description",
			Status:      gauth.ACTIVE,
			Roles:       []gauth.TRole{gauth.MANAGER},
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
		form1 := url.Values{}
		form1.Add("login", "root")
		r1 := httptest.NewRequest("POST", "/", strings.NewReader(form1.Encode()))
		r1.PostForm = form1
		for _, v := range wCooks {
			r1.AddCookie(&http.Cookie{
				Name:   v.Name,
				Value:  v.Value,
				MaxAge: v.MaxAge,
			})
		}

		account_ban_ok(w1, r1) // calling the tested function
		status := w1.Code
		if status != http.StatusOK {
			t.Errorf("account_ban_ok() error: %v", status)
		}
	})

	t.Run("account_ban_ok() function testing - all right", func(t *testing.T) {
		randStr := gauth.GenerateTicket()
		prof := gauth.TProfile{
			Description: "Testing description",
			Status:      gauth.NEW,
			Roles:       []gauth.TRole{gauth.MANAGER},
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

		randStr1 := gauth.GenerateTicket()
		prof1 := gauth.TProfile{
			Description: "Testing description",
			Status:      gauth.ACTIVE,
			Roles:       []gauth.TRole{gauth.MANAGER},
		}
		gauth.AddUser(randStr1, randStr1, prof1)

		w1 := httptest.NewRecorder()
		form1 := url.Values{}
		form1.Add("login", randStr1)
		r1 := httptest.NewRequest("POST", "/", strings.NewReader(form1.Encode()))
		r1.PostForm = form1
		for _, v := range wCooks {
			r1.AddCookie(&http.Cookie{
				Name:   v.Name,
				Value:  v.Value,
				MaxAge: v.MaxAge,
			})
		}

		account_ban_ok(w1, r1) // calling the tested function
		status := w1.Code
		if status != http.StatusOK {
			t.Errorf("account_ban_ok() error: %v", status)
		}
	})
}

func Test_account_unban_load_form(t *testing.T) {
	gauth.Start()
	parseTemplates()

	t.Run("account_unban_load_form() function testing - Isolate negative", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/", nil)

		account_unban_load_form(w, r) // calling the tested function
		status := w.Code
		if status != http.StatusOK {
			t.Errorf("account_unban_load_form() error: %v", status)
		}
	})

	t.Run("account_unban_load_form() function testing - root user", func(t *testing.T) {
		w := httptest.NewRecorder()
		form := url.Values{}
		form.Add("username", "root")
		form.Add("password", "toor")
		r := httptest.NewRequest("POST", "/", strings.NewReader(form.Encode()))
		r.PostForm = form

		homeAuth(w, r)
		wCooks := w.Result().Cookies()

		randStr := gauth.GenerateTicket()
		prof := gauth.TProfile{
			Description: "Testing description",
			Status:      gauth.ACTIVE,
			Roles:       []gauth.TRole{gauth.USER},
		}
		gauth.AddUser(randStr, randStr, prof)

		w1 := httptest.NewRecorder()
		testurl := "/?user=root"
		r1 := httptest.NewRequest("GET", testurl, nil)
		for _, v := range wCooks {
			r1.AddCookie(&http.Cookie{
				Name:   v.Name,
				Value:  v.Value,
				MaxAge: v.MaxAge,
			})
		}

		account_unban_load_form(w1, r1) // calling the tested function
		status := w1.Code
		if status != http.StatusOK {
			t.Errorf("account_unban_load_form() error: %v", status)
		}
	})

	t.Run("account_unban_load_form() function testing - all right", func(t *testing.T) {
		w := httptest.NewRecorder()
		form := url.Values{}
		form.Add("username", "root")
		form.Add("password", "toor")
		r := httptest.NewRequest("POST", "/", strings.NewReader(form.Encode()))
		r.PostForm = form

		homeAuth(w, r)
		wCooks := w.Result().Cookies()

		randStr := gauth.GenerateTicket()
		prof := gauth.TProfile{
			Description: "Testing description",
			Status:      gauth.ACTIVE,
			Roles:       []gauth.TRole{gauth.USER},
		}
		gauth.AddUser(randStr, randStr, prof)

		w1 := httptest.NewRecorder()
		testurl := fmt.Sprintf("/?user=%s", randStr)
		r1 := httptest.NewRequest("GET", testurl, nil)
		for _, v := range wCooks {
			r1.AddCookie(&http.Cookie{
				Name:   v.Name,
				Value:  v.Value,
				MaxAge: v.MaxAge,
			})
		}

		account_unban_load_form(w1, r1) // calling the tested function
		status := w1.Code
		if status != http.StatusOK {
			t.Errorf("account_unban_load_form() error: %v", status)
		}
	})
}

func Test_account_unban_ok(t *testing.T) {
	gauth.Start()
	parseTemplates()

	t.Run("account_unban_ok() function testing - Isolate negative", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/", nil)

		account_unban_ok(w, r) // calling the tested function
		status := w.Code
		if status != http.StatusOK {
			t.Errorf("account_unban_ok() error: %v", status)
		}
	})

	t.Run("account_unban_ok() function testing - POST negative", func(t *testing.T) {
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

		account_unban_ok(w1, r1) // calling the tested function
		status := w1.Code
		if status != http.StatusOK {
			t.Errorf("account_unban_ok() error: %v", status)
		}
	})

	t.Run("account_unban_ok() function testing POST positive and don't work ParseForm", func(t *testing.T) {
		randStr := gauth.GenerateTicket()
		prof := gauth.TProfile{
			Description: "Testing description",
			Status:      gauth.ACTIVE,
			Roles:       []gauth.TRole{gauth.MANAGER},
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
		r1 := httptest.NewRequest("POST", "/", nil)
		r1.PostForm = nil
		r1.Body = nil
		for _, v := range wCooks {
			r1.AddCookie(&http.Cookie{
				Name:   v.Name,
				Value:  v.Value,
				MaxAge: v.MaxAge,
			})
		}

		account_unban_ok(w1, r1) // calling the tested function
		status := w1.Code
		if status != http.StatusOK {
			t.Errorf("account_unban_ok() error: %v", status)
		}
	})

	t.Run("account_unban_ok() function testing POST positive and not value of login", func(t *testing.T) {
		randStr := gauth.GenerateTicket()
		prof := gauth.TProfile{
			Description: "Testing description",
			Status:      gauth.ACTIVE,
			Roles:       []gauth.TRole{gauth.MANAGER},
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
		form1 := url.Values{}
		form1.Add("login", "")
		r1 := httptest.NewRequest("POST", "/", strings.NewReader(form1.Encode()))
		r1.PostForm = form1
		for _, v := range wCooks {
			r1.AddCookie(&http.Cookie{
				Name:   v.Name,
				Value:  v.Value,
				MaxAge: v.MaxAge,
			})
		}

		account_unban_ok(w1, r1) // calling the tested function
		status := w1.Code
		if status != http.StatusOK {
			t.Errorf("account_unban_ok() error: %v", status)
		}
	})

	t.Run("account_unban_ok() function testing POST positive and don't work UnblockUser", func(t *testing.T) {
		randStr := gauth.GenerateTicket()
		prof := gauth.TProfile{
			Description: "Testing description",
			Status:      gauth.ACTIVE,
			Roles:       []gauth.TRole{gauth.MANAGER},
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
		form1 := url.Values{}
		form1.Add("login", "root")
		r1 := httptest.NewRequest("POST", "/", strings.NewReader(form1.Encode()))
		r1.PostForm = form1
		for _, v := range wCooks {
			r1.AddCookie(&http.Cookie{
				Name:   v.Name,
				Value:  v.Value,
				MaxAge: v.MaxAge,
			})
		}

		account_unban_ok(w1, r1) // calling the tested function
		status := w1.Code
		if status != http.StatusOK {
			t.Errorf("account_unban_ok() error: %v", status)
		}
	})

	t.Run("account_unban_ok() function testing - all right", func(t *testing.T) {
		randStr := gauth.GenerateTicket()
		prof := gauth.TProfile{
			Description: "Testing description",
			Status:      gauth.ACTIVE,
			Roles:       []gauth.TRole{gauth.MANAGER},
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

		randStr1 := gauth.GenerateTicket()
		prof1 := gauth.TProfile{
			Description: "Testing description",
			Status:      gauth.ACTIVE,
			Roles:       []gauth.TRole{gauth.MANAGER},
		}
		gauth.AddUser(randStr1, randStr1, prof1)

		w1 := httptest.NewRecorder()
		form1 := url.Values{}
		form1.Add("login", randStr1)
		r1 := httptest.NewRequest("POST", "/", strings.NewReader(form1.Encode()))
		r1.PostForm = form1
		for _, v := range wCooks {
			r1.AddCookie(&http.Cookie{
				Name:   v.Name,
				Value:  v.Value,
				MaxAge: v.MaxAge,
			})
		}

		account_unban_ok(w1, r1) // calling the tested function
		status := w1.Code
		if status != http.StatusOK {
			t.Errorf("account_unban_ok() error: %v", status)
		}
	})
}

func Test_account_del_load_form(t *testing.T) {
	gauth.Start()
	parseTemplates()

	t.Run("account_del_load_form() function testing - Isolate negative", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/", nil)

		account_del_load_form(w, r) // calling the tested function
		status := w.Code
		if status != http.StatusOK {
			t.Errorf("account_del_load_form() error: %v", status)
		}
	})

	t.Run("account_del_load_form() function testing - root user", func(t *testing.T) {
		w := httptest.NewRecorder()
		form := url.Values{}
		form.Add("username", "root")
		form.Add("password", "toor")
		r := httptest.NewRequest("POST", "/", strings.NewReader(form.Encode()))
		r.PostForm = form

		homeAuth(w, r)
		wCooks := w.Result().Cookies()

		w1 := httptest.NewRecorder()
		testurl := "/?user=root"
		r1 := httptest.NewRequest("GET", testurl, nil)
		for _, v := range wCooks {
			r1.AddCookie(&http.Cookie{
				Name:   v.Name,
				Value:  v.Value,
				MaxAge: v.MaxAge,
			})
		}

		account_del_load_form(w1, r1) // calling the tested function
		status := w1.Code
		if status != http.StatusOK {
			t.Errorf("account_del_load_form() error: %v", status)
		}
	})

	t.Run("account_del_load_form() function testing - all right", func(t *testing.T) {
		w := httptest.NewRecorder()
		form := url.Values{}
		form.Add("username", "root")
		form.Add("password", "toor")
		r := httptest.NewRequest("POST", "/", strings.NewReader(form.Encode()))
		r.PostForm = form

		homeAuth(w, r)
		wCooks := w.Result().Cookies()

		randStr := gauth.GenerateTicket()
		prof := gauth.TProfile{
			Description: "Testing description",
			Status:      gauth.ACTIVE,
			Roles:       []gauth.TRole{gauth.USER},
		}
		gauth.AddUser(randStr, randStr, prof)

		w1 := httptest.NewRecorder()
		testurl := fmt.Sprintf("/?user=%s", randStr)
		r1 := httptest.NewRequest("GET", testurl, nil)
		for _, v := range wCooks {
			r1.AddCookie(&http.Cookie{
				Name:   v.Name,
				Value:  v.Value,
				MaxAge: v.MaxAge,
			})
		}

		account_del_load_form(w1, r1) // calling the tested function
		status := w1.Code
		if status != http.StatusOK {
			t.Errorf("account_del_load_form() error: %v", status)
		}
	})
}

func Test_account_del_ok(t *testing.T) {
	gauth.Start()
	parseTemplates()

	t.Run("account_del_ok() function testing - Isolate negative", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/", nil)

		account_del_ok(w, r) // calling the tested function
		status := w.Code
		if status != http.StatusOK {
			t.Errorf("account_del_ok() error: %v", status)
		}
	})

	t.Run("account_del_ok() function testing - POST negative", func(t *testing.T) {
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

		account_del_ok(w1, r1) // calling the tested function
		status := w1.Code
		if status != http.StatusOK {
			t.Errorf("account_del_ok() error: %v", status)
		}
	})

	t.Run("account_del_ok() function testing POST positive and don't work ParseForm", func(t *testing.T) {
		randStr := gauth.GenerateTicket()
		prof := gauth.TProfile{
			Description: "Testing description",
			Status:      gauth.ACTIVE,
			Roles:       []gauth.TRole{gauth.MANAGER},
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
		r1 := httptest.NewRequest("POST", "/", nil)
		r1.PostForm = nil
		r1.Body = nil
		for _, v := range wCooks {
			r1.AddCookie(&http.Cookie{
				Name:   v.Name,
				Value:  v.Value,
				MaxAge: v.MaxAge,
			})
		}

		account_del_ok(w1, r1) // calling the tested function
		status := w1.Code
		if status != http.StatusOK {
			t.Errorf("account_del_ok() error: %v", status)
		}
	})

	t.Run("account_del_ok() function testing POST positive and not value of login", func(t *testing.T) {
		randStr := gauth.GenerateTicket()
		prof := gauth.TProfile{
			Description: "Testing description",
			Status:      gauth.ACTIVE,
			Roles:       []gauth.TRole{gauth.MANAGER},
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
		form1 := url.Values{}
		form1.Add("login", "")
		r1 := httptest.NewRequest("POST", "/", strings.NewReader(form1.Encode()))
		r1.PostForm = form1
		for _, v := range wCooks {
			r1.AddCookie(&http.Cookie{
				Name:   v.Name,
				Value:  v.Value,
				MaxAge: v.MaxAge,
			})
		}

		account_del_ok(w1, r1) // calling the tested function
		status := w1.Code
		if status != http.StatusOK {
			t.Errorf("account_del_ok() error: %v", status)
		}
	})

	t.Run("account_del_ok() function testing POST positive and don't work DeleteUser", func(t *testing.T) {
		randStr := gauth.GenerateTicket()
		prof := gauth.TProfile{
			Description: "Testing description",
			Status:      gauth.ACTIVE,
			Roles:       []gauth.TRole{gauth.MANAGER},
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
		form1 := url.Values{}
		form1.Add("login", "root")
		r1 := httptest.NewRequest("POST", "/", strings.NewReader(form1.Encode()))
		r1.PostForm = form1
		for _, v := range wCooks {
			r1.AddCookie(&http.Cookie{
				Name:   v.Name,
				Value:  v.Value,
				MaxAge: v.MaxAge,
			})
		}

		account_del_ok(w1, r1) // calling the tested function
		status := w1.Code
		if status != http.StatusOK {
			t.Errorf("account_del_ok() error: %v", status)
		}
	})

	t.Run("account_del_ok() function testing - all right", func(t *testing.T) {
		randStr := gauth.GenerateTicket()
		prof := gauth.TProfile{
			Description: "Testing description",
			Status:      gauth.ACTIVE,
			Roles:       []gauth.TRole{gauth.MANAGER},
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

		randStr1 := gauth.GenerateTicket()
		prof1 := gauth.TProfile{
			Description: "Testing description",
			Status:      gauth.ACTIVE,
			Roles:       []gauth.TRole{gauth.MANAGER},
		}
		gauth.AddUser(randStr1, randStr1, prof1)

		w1 := httptest.NewRecorder()
		form1 := url.Values{}
		form1.Add("login", randStr1)
		r1 := httptest.NewRequest("POST", "/", strings.NewReader(form1.Encode()))
		r1.PostForm = form1
		for _, v := range wCooks {
			r1.AddCookie(&http.Cookie{
				Name:   v.Name,
				Value:  v.Value,
				MaxAge: v.MaxAge,
			})
		}

		account_del_ok(w1, r1) // calling the tested function
		status := w1.Code
		if status != http.StatusOK {
			t.Errorf("account_del_ok() error: %v", status)
		}
	})
}

func Test_nav_settings(t *testing.T) {
	gauth.Start()
	parseTemplates()

	t.Run("nav_settings() function testing - Isolate negative", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/", nil)

		nav_settings(w, r) // calling the tested function
		status := w.Code
		if status != http.StatusOK {
			t.Errorf("nav_settings() error: %v", status)
		}
	})

	t.Run("nav_settings() function testing - Isolate positive", func(t *testing.T) {
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

		nav_settings(w1, r1) // calling the tested function
		status := w1.Code
		if status != http.StatusOK {
			t.Errorf("nav_settings() error: %v", status)
		}
	})
}

func Test_settings_core_freeze_change_sw(t *testing.T) {
	gauth.Start()
	parseTemplates()

	t.Run("settings_core_freeze_change_sw() function testing - Isolate negative", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/", nil)

		settings_core_freeze_change_sw(w, r) // calling the tested function
		status := w.Code
		if status != http.StatusOK {
			t.Errorf("settings_core_freeze_change_sw() error: %v", status)
		}
	})

	t.Run("settings_core_freeze_change_sw() function testing - off", func(t *testing.T) {
		config.DefaultConfig.CoreSettings.FreezeMode = true

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

		settings_core_freeze_change_sw(w1, r1) // calling the tested function
		status := w1.Code
		if status != http.StatusOK || core.LocalCoreSettings.FreezeMode {
			t.Errorf("settings_core_freeze_change_sw() error: %v", status)
		}
	})

	t.Run("settings_core_freeze_change_sw() function testing - on", func(t *testing.T) {
		config.DefaultConfig.CoreSettings.FreezeMode = false

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

		settings_core_freeze_change_sw(w1, r1) // calling the tested function
		status := w1.Code
		if status != http.StatusOK || !core.LocalCoreSettings.FreezeMode {
			t.Errorf("settings_core_freeze_change_sw() error: %v", status)
		}
	})
}

func Test_settings_wsc_change_sw(t *testing.T) {
	gauth.Start()
	parseTemplates()

	t.Run("settings_wsc_change_sw() function testing - Isolate negative", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/", nil)

		settings_wsc_change_sw(w, r) // calling the tested function
		status := w.Code
		if status != http.StatusOK {
			t.Errorf("settings_wsc_change_sw() error: %v", status)
		}
	})

	t.Run("settings_wsc_change_sw() function testing - shutdown", func(t *testing.T) {
		config.DefaultConfig.WebSocketConnector.Enable = true
		go websocketconn.Start(&config.DefaultConfig)
		closer.AddHandler(websocketconn.Shutdown) // Register a shutdown handler.

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

		settings_wsc_change_sw(w1, r1) // calling the tested function
		status := w1.Code
		if status != http.StatusOK {
			t.Errorf("settings_wsc_change_sw() error: %v", status)
		}
	})

	t.Run("settings_wsc_change_sw() function testing - start", func(t *testing.T) {
		config.DefaultConfig.WebSocketConnector.Enable = false
		closer.RunAndDelHandler(websocketconn.Shutdown)

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

		settings_wsc_change_sw(w1, r1) // calling the tested function
		status := w1.Code
		if status != http.StatusOK {
			t.Errorf("settings_wsc_change_sw() error: %v", status)
		}
	})
}

func Test_settings_rest_change_sw(t *testing.T) {
	gauth.Start()
	parseTemplates()

	t.Run("settings_rest_change_sw() function testing - Isolate negative", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/", nil)

		settings_rest_change_sw(w, r) // calling the tested function
		status := w.Code
		if status != http.StatusOK {
			t.Errorf("settings_rest_change_sw() error: %v", status)
		}
	})

	t.Run("settings_rest_change_sw() function testing - shutdown", func(t *testing.T) {
		config.DefaultConfig.RestConnector.Enable = true
		go rest.Start(&config.DefaultConfig)
		closer.AddHandler(rest.Shutdown) // Register a shutdown handler.

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

		settings_rest_change_sw(w1, r1) // calling the tested function
		status := w1.Code
		if status != http.StatusOK {
			t.Errorf("settings_rest_change_sw() error: %v", status)
		}
	})

	t.Run("settings_rest_change_sw() function testing - start", func(t *testing.T) {
		config.DefaultConfig.RestConnector.Enable = false
		closer.RunAndDelHandler(rest.Shutdown)

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

		settings_rest_change_sw(w1, r1) // calling the tested function
		status := w1.Code
		if status != http.StatusOK {
			t.Errorf("settings_rest_change_sw() error: %v", status)
		}
	})
}

func Test_settings_grpc_change_sw(t *testing.T) {
	gauth.Start()
	parseTemplates()

	t.Run("settings_grpc_change_sw() function testing - Isolate negative", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/", nil)

		settings_grpc_change_sw(w, r) // calling the tested function
		status := w.Code
		if status != http.StatusOK {
			t.Errorf("settings_grpc_change_sw() error: %v", status)
		}
	})
}

func Test_settings_web_change_sw(t *testing.T) {
	gauth.Start()
	parseTemplates()

	t.Run("settings_web_change_sw() function testing - Isolate negative", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/", nil)

		settings_web_change_sw(w, r) // calling the tested function
		status := w.Code
		if status != http.StatusOK {
			t.Errorf("settings_web_change_sw() error: %v", status)
		}
	})

	t.Run("settings_web_change_sw() function testing - Isolate positive", func(t *testing.T) {
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

		settings_web_change_sw(w1, r1) // calling the tested function
		status := w1.Code
		if status != http.StatusOK {
			t.Errorf("settings_web_change_sw() error: %v", status)
		}
	})

	t.Run("settings_web_change_sw() function testing - Isolate undefined user", func(t *testing.T) {
		randStr := gauth.GenerateTicket()
		prof := gauth.TProfile{
			Description: "Testing description",
			Status:      gauth.UNDEFINED,
			Roles:       []gauth.TRole{gauth.ADMIN},
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

		settings_web_change_sw(w1, r1) // calling the tested function
		status := w1.Code
		if status != http.StatusOK {
			t.Errorf("settings_web_change_sw() error: %v", status)
		}
	})
}
