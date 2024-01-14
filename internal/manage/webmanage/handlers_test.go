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
