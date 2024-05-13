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

func Test_fnHomeDefault(t *testing.T) {
	gauth.Start()
	parseTemplates()

	t.Run("fnHomeDefault() function testing - negative", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/", nil)

		fnHomeDefault(w, r) // calling the tested function
		status := w.Code
		if status != http.StatusFound {
			t.Errorf("fnHomeDefault() error: %v", status)
		}
	})

	t.Run("fnHomeDefault() function testing - positive", func(t *testing.T) {
		w := httptest.NewRecorder()
		form := url.Values{}
		form.Add("username", "root")
		form.Add("password", "toor")
		r := httptest.NewRequest("POST", "/", strings.NewReader(form.Encode()))
		r.PostForm = form

		fnHomeAuth(w, r)
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

		fnHomeDefault(w1, r1) // calling the tested function
		status := w1.Code
		if status != http.StatusOK {
			t.Errorf("fnHomeDefault() error: %v", status)
		}
	})

	t.Run("fnHomeDefault() function testing - logout", func(t *testing.T) {
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

		fnHomeAuth(w, r)
		wCooks := w.Result().Cookies()

		delete(gauth.MAccess, randStr)

		w1 := httptest.NewRecorder()
		r1 := httptest.NewRequest("GET", "/", nil)
		for _, v := range wCooks {
			r1.AddCookie(&http.Cookie{
				Name:   v.Name,
				Value:  v.Value,
				MaxAge: v.MaxAge,
			})
		}

		fnHomeDefault(w1, r1) // calling the tested function
		status := w1.Code
		if status != http.StatusFound {
			t.Errorf("fnHomeDefault() error: %v", status)
		}
	})

	t.Run("fnHomeDefault() function testing - Template error", func(t *testing.T) {
		w := httptest.NewRecorder()
		form := url.Values{}
		form.Add("username", "root")
		form.Add("password", "toor")
		r := httptest.NewRequest("POST", "/", strings.NewReader(form.Encode()))
		r.PostForm = form

		fnHomeAuth(w, r)
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
		LoadTemplateFromString(HOME_TEMP_NAME, wrongStr)

		fnHomeDefault(w1, r1) // calling the tested function
		status := w1.Code
		if status != http.StatusOK {
			t.Errorf("fnHomeDefault() error: %v", status)
		}
	})
}

func Test_fnHomeAuth(t *testing.T) {
	gauth.Start()
	parseTemplates()

	t.Run("fnHomeAuth() function testing - GET negative", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/", nil)

		fnHomeAuth(w, r) // calling the tested function
		status := w.Code
		if status != http.StatusOK {
			t.Errorf("fnHomeAuth() error: %v", status)
		}
	})

	t.Run("fnHomeAuth() function testing - POST negative", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest("POST", "/", nil)
		r.PostForm = nil
		r.Body = nil

		fnHomeAuth(w, r)

		status := w.Code
		if status != http.StatusBadRequest {
			t.Errorf("fnHomeAuth() error: %v", status)
		}
	})

}

func Test_fnHome(t *testing.T) {
	gauth.Start()
	parseTemplates()

	t.Run("fnHome() function testing - negative", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/a", nil)

		fnHome(w, r) // calling the tested function
		status := w.Code
		if status != http.StatusFound {
			t.Errorf("fnHome() error: %v", status)
		}
	})

	t.Run("fnHome() function testing", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/", nil)

		fnHome(w, r) // calling the tested function
		status := w.Code
		if status != http.StatusOK {
			t.Errorf("fnHome() error: %v", status)
		}
	})

	t.Run("fnHome() function testing", func(t *testing.T) {
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

		fnHomeAuth(w, r)
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

		fnHome(w1, r1) // calling the tested function
		status := w1.Code
		if status != http.StatusOK {
			t.Errorf("fnHome() error: %v", status)
		}
	})
}

func Test_fnNavDefault(t *testing.T) {
	gauth.Start()
	parseTemplates()

	t.Run("fnNavDefault() function testing", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/", nil)

		fnNavDefault(w, r) // calling the tested function
		status := w.Code
		if status != http.StatusOK {
			t.Errorf("fnNavDefault() error: %v", status)
		}
	})
}

func Test_fnNavLogout(t *testing.T) {
	gauth.Start()
	parseTemplates()

	t.Run("fnNavLogout() function testing", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/", nil)

		fnNavLogout(w, r) // calling the tested function
		status := w.Code
		if status != http.StatusOK {
			t.Errorf("fnNavLogout() error: %v", status)
		}
	})
}

func Test_fnSelfeditLoadForm(t *testing.T) {
	gauth.Start()
	parseTemplates()

	t.Run("fnSelfeditLoadForm() function testing - negative", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/", nil)

		fnSelfeditLoadForm(w, r) // calling the tested function
		status := w.Code
		if status != http.StatusOK {
			t.Errorf("fnSelfeditLoadForm() error: %v", status)
		}
	})

	t.Run("fnSelfeditLoadForm() function testing - positive", func(t *testing.T) {
		w := httptest.NewRecorder()
		form := url.Values{}
		form.Add("username", "root")
		form.Add("password", "toor")
		r := httptest.NewRequest("POST", "/", strings.NewReader(form.Encode()))
		r.PostForm = form

		fnHomeAuth(w, r)
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

		fnSelfeditLoadForm(w1, r1) // calling the tested function
		status := w1.Code
		if status != http.StatusOK {
			t.Errorf("fnSelfeditLoadForm() error: %v", status)
		}
	})

	t.Run("fnSelfeditLoadForm() function testing", func(t *testing.T) {
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

		fnHomeAuth(w, r)
		wCooks := w.Result().Cookies()

		delete(gauth.MAccess, randStr)

		w1 := httptest.NewRecorder()
		r1 := httptest.NewRequest("GET", "/", nil)
		for _, v := range wCooks {
			r1.AddCookie(&http.Cookie{
				Name:   v.Name,
				Value:  v.Value,
				MaxAge: v.MaxAge,
			})
		}

		fnSelfeditLoadForm(w1, r1) // calling the tested function
		status := w1.Code
		if status != http.StatusOK {
			t.Errorf("fnSelfeditLoadForm() error: %v", status)
		}
	})
}

func Test_fnSelfeditOk(t *testing.T) {
	gauth.Start()
	parseTemplates()

	t.Run("fnSelfeditOk() function testing - Isolate negative", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/", nil)

		fnSelfeditOk(w, r) // calling the tested function
		status := w.Code
		if status != http.StatusOK {
			t.Errorf("fnSelfeditOk() error: %v", status)
		}
	})

	t.Run("fnSelfeditOk() function testing - POST negative", func(t *testing.T) {
		w := httptest.NewRecorder()
		form := url.Values{}
		form.Add("username", "root")
		form.Add("password", "toor")
		r := httptest.NewRequest("POST", "/", strings.NewReader(form.Encode()))
		r.PostForm = form

		fnHomeAuth(w, r)
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

		fnSelfeditOk(w1, r1) // calling the tested function
		status := w1.Code
		if status != http.StatusOK {
			t.Errorf("fnSelfeditOk() error: %v", status)
		}
	})

	t.Run("fnSelfeditOk() function testing POST positive and don't work ParseForm", func(t *testing.T) {
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

		fnHomeAuth(w, r)
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

		fnSelfeditOk(w1, r1) // calling the tested function
		status := w1.Code
		if status != http.StatusOK {
			t.Errorf("fnSelfeditOk() error: %v", status)
		}
	})

	t.Run("fnSelfeditOk() function testing POST positive and not value of password", func(t *testing.T) {
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

		fnHomeAuth(w, r)
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

		fnSelfeditOk(w1, r1) // calling the tested function
		status := w1.Code
		if status != http.StatusOK {
			t.Errorf("fnSelfeditOk() error: %v", status)
		}
	})

	t.Run("fnSelfeditOk() function testing - all right", func(t *testing.T) {
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

		fnHomeAuth(w, r)
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

		fnSelfeditOk(w1, r1) // calling the tested function
		status := w1.Code
		if status != http.StatusOK {
			t.Errorf("fnSelfeditOk() error: %v", status)
		}
	})
}

func Test_fnNavDashboard(t *testing.T) {
	gauth.Start()
	parseTemplates()

	t.Run("fnNavDashboard() function testing - Isolate negative", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/", nil)

		fnNavDashboard(w, r) // calling the tested function
		status := w.Code
		if status != http.StatusOK {
			t.Errorf("fnNavDashboard() error: %v", status)
		}
	})

	t.Run("fnNavDashboard() function testing - Isolate positive", func(t *testing.T) {
		w := httptest.NewRecorder()
		form := url.Values{}
		form.Add("username", "root")
		form.Add("password", "toor")
		r := httptest.NewRequest("POST", "/", strings.NewReader(form.Encode()))
		r.PostForm = form

		fnHomeAuth(w, r)
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

		fnNavDashboard(w1, r1) // calling the tested function
		status := w1.Code
		if status != http.StatusOK {
			t.Errorf("fnNavDashboard() error: %v", status)
		}
	})
}

func Test_fnNavDatabases(t *testing.T) {
	gauth.Start()
	parseTemplates()

	t.Run("fnNavDatabases() function testing - Isolate negative", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/", nil)

		fnNavDatabases(w, r) // calling the tested function
		status := w.Code
		if status != http.StatusOK {
			t.Errorf("fnNavDatabases() error: %v", status)
		}
	})

	t.Run("fnNavDatabases() function testing - Isolate positive", func(t *testing.T) {
		w := httptest.NewRecorder()
		form := url.Values{}
		form.Add("username", "root")
		form.Add("password", "toor")
		r := httptest.NewRequest("POST", "/", strings.NewReader(form.Encode()))
		r.PostForm = form

		fnHomeAuth(w, r)
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

		fnNavDatabases(w1, r1) // calling the tested function
		status := w1.Code
		if status != http.StatusOK {
			t.Errorf("fnNavDatabases() error: %v", status)
		}
	})
}

func Test_fnConsoleRequest(t *testing.T) {
	gauth.Start()
	parseTemplates()

	t.Run("fnConsoleRequest() function testing - Isolate negative", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/", nil)

		fnConsoleRequest(w, r) // calling the tested function
		status := w.Code
		if status != http.StatusOK {
			t.Errorf("fnConsoleRequest() error: %v", status)
		}
	})

	t.Run("fnConsoleRequest() function testing - POST negative", func(t *testing.T) {
		w := httptest.NewRecorder()
		form := url.Values{}
		form.Add("username", "root")
		form.Add("password", "toor")
		r := httptest.NewRequest("POST", "/", strings.NewReader(form.Encode()))
		r.PostForm = form

		fnHomeAuth(w, r)
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

		fnConsoleRequest(w1, r1) // calling the tested function
		status := w1.Code
		if status != http.StatusOK {
			t.Errorf("fnConsoleRequest() error: %v", status)
		}
	})

	t.Run("fnConsoleRequest() function testing POST positive and don't work ParseForm", func(t *testing.T) {
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

		fnHomeAuth(w, r)
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

		fnConsoleRequest(w1, r1) // calling the tested function
		status := w1.Code
		if status != http.StatusOK {
			t.Errorf("fnConsoleRequest() error: %v", status)
		}
	})

	t.Run("fnConsoleRequest() function testing - all right", func(t *testing.T) {
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

		fnHomeAuth(w, r)
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

		fnConsoleRequest(w1, r1) // calling the tested function
		status := w1.Code
		if status != http.StatusOK {
			t.Errorf("fnConsoleRequest() error: %v", status)
		}
	})
}

func Test_fnNavAccounts(t *testing.T) {
	gauth.Start()
	parseTemplates()

	t.Run("fnNavAccounts() function testing - Isolate negative", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/", nil)

		fnNavAccounts(w, r) // calling the tested function
		status := w.Code
		if status != http.StatusOK {
			t.Errorf("fnNavAccounts() error: %v", status)
		}
	})

	t.Run("fnNavAccounts() function testing - Isolate positive", func(t *testing.T) {
		w := httptest.NewRecorder()
		form := url.Values{}
		form.Add("username", "root")
		form.Add("password", "toor")
		r := httptest.NewRequest("POST", "/", strings.NewReader(form.Encode()))
		r.PostForm = form

		fnHomeAuth(w, r)
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

		fnNavAccounts(w1, r1) // calling the tested function
		status := w1.Code
		if status != http.StatusOK {
			t.Errorf("fnNavAccounts() error: %v", status)
		}
	})

	t.Run("fnNavAccounts() function testing - all coverage", func(t *testing.T) {
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

		fnHomeAuth(w, r)
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

		fnNavAccounts(w1, r1) // calling the tested function
		status := w1.Code
		if status != http.StatusOK {
			t.Errorf("fnNavAccounts() error: %v", status)
		}
	})
}

func Test_fnAccountCreateLoadForm(t *testing.T) {
	gauth.Start()
	parseTemplates()

	t.Run("fnAccountCreateLoadForm() function testing - Isolate negative", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/", nil)

		fnAccountCreateLoadForm(w, r) // calling the tested function
		status := w.Code
		if status != http.StatusOK {
			t.Errorf("fnAccountCreateLoadForm() error: %v", status)
		}
	})

	t.Run("fnAccountCreateLoadForm() function testing - Isolate positive", func(t *testing.T) {
		w := httptest.NewRecorder()
		form := url.Values{}
		form.Add("username", "root")
		form.Add("password", "toor")
		r := httptest.NewRequest("POST", "/", strings.NewReader(form.Encode()))
		r.PostForm = form

		fnHomeAuth(w, r)
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

		fnAccountCreateLoadForm(w1, r1) // calling the tested function
		status := w1.Code
		if status != http.StatusOK {
			t.Errorf("fnAccountCreateLoadForm() error: %v", status)
		}
	})
}

func Test_fnAccountCreateOk(t *testing.T) {
	gauth.Start()
	parseTemplates()

	t.Run("fnAccountCreateOk() function testing - Isolate negative", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/", nil)

		fnAccountCreateOk(w, r) // calling the tested function
		status := w.Code
		if status != http.StatusOK {
			t.Errorf("fnAccountCreateOk() error: %v", status)
		}
	})

	t.Run("fnAccountCreateOk() function testing - POST negative", func(t *testing.T) {
		w := httptest.NewRecorder()
		form := url.Values{}
		form.Add("username", "root")
		form.Add("password", "toor")
		r := httptest.NewRequest("POST", "/", strings.NewReader(form.Encode()))
		r.PostForm = form

		fnHomeAuth(w, r)
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

		fnAccountCreateOk(w1, r1) // calling the tested function
		status := w1.Code
		if status != http.StatusOK {
			t.Errorf("fnAccountCreateOk() error: %v", status)
		}
	})

	t.Run("fnAccountCreateOk() function testing POST positive and don't work ParseForm", func(t *testing.T) {
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

		fnHomeAuth(w, r)
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

		fnAccountCreateOk(w1, r1) // calling the tested function
		status := w1.Code
		if status != http.StatusOK {
			t.Errorf("fnAccountCreateOk() error: %v", status)
		}
	})

	t.Run("fnAccountCreateOk() function testing POST positive and not value of password", func(t *testing.T) {
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

		fnHomeAuth(w, r)
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

		fnAccountCreateOk(w1, r1) // calling the tested function
		status := w1.Code
		if status != http.StatusOK {
			t.Errorf("fnAccountCreateOk() error: %v", status)
		}
	})

	t.Run("fnAccountCreateOk() function testing - create error", func(t *testing.T) {
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

		fnHomeAuth(w, r)
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

		fnAccountCreateOk(w1, r1) // calling the tested function
		status := w1.Code
		if status != http.StatusOK {
			t.Errorf("fnAccountCreateOk() error: %v", status)
		}
	})

	t.Run("fnAccountCreateOk() function testing - all right", func(t *testing.T) {
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

		fnHomeAuth(w, r)
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

		fnAccountCreateOk(w1, r1) // calling the tested function
		status := w1.Code
		if status != http.StatusOK {
			t.Errorf("fnAccountCreateOk() error: %v", status)
		}
	})
}

func Test_fnAccountEditLoadForm(t *testing.T) {
	gauth.Start()
	parseTemplates()

	t.Run("fnAccountEditLoadForm() function testing - Isolate negative", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/", nil)

		fnAccountEditLoadForm(w, r) // calling the tested function
		status := w.Code
		if status != http.StatusOK {
			t.Errorf("fnAccountEditLoadForm() error: %v", status)
		}
	})

	t.Run("fnAccountEditLoadForm() function testing - Isolate positive and don't GetProfile", func(t *testing.T) {
		w := httptest.NewRecorder()
		form := url.Values{}
		form.Add("username", "root")
		form.Add("password", "toor")
		r := httptest.NewRequest("POST", "/", strings.NewReader(form.Encode()))
		r.PostForm = form

		fnHomeAuth(w, r)
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

		fnAccountEditLoadForm(w1, r1) // calling the tested function
		status := w1.Code
		if status != http.StatusOK {
			t.Errorf("fnAccountEditLoadForm() error: %v", status)
		}
	})

	t.Run("fnAccountEditLoadForm() function testing - all right", func(t *testing.T) {
		w := httptest.NewRecorder()
		form := url.Values{}
		form.Add("username", "root")
		form.Add("password", "toor")
		r := httptest.NewRequest("POST", "/", strings.NewReader(form.Encode()))
		r.PostForm = form

		fnHomeAuth(w, r)
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

		fnAccountEditLoadForm(w1, r1) // calling the tested function
		status := w1.Code
		if status != http.StatusOK {
			t.Errorf("fnAccountEditLoadForm() error: %v", status)
		}
	})
}

func Test_fnAccountEditOk(t *testing.T) {
	gauth.Start()
	parseTemplates()

	t.Run("fnAccountEditOk() function testing - Isolate negative", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/", nil)

		fnAccountEditOk(w, r) // calling the tested function
		status := w.Code
		if status != http.StatusOK {
			t.Errorf("fnAccountEditOk() error: %v", status)
		}
	})

	t.Run("fnAccountEditOk() function testing - POST negative", func(t *testing.T) {
		w := httptest.NewRecorder()
		form := url.Values{}
		form.Add("username", "root")
		form.Add("password", "toor")
		r := httptest.NewRequest("POST", "/", strings.NewReader(form.Encode()))
		r.PostForm = form

		fnHomeAuth(w, r)
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

		fnAccountEditOk(w1, r1) // calling the tested function
		status := w1.Code
		if status != http.StatusOK {
			t.Errorf("fnAccountEditOk() error: %v", status)
		}
	})

	t.Run("fnAccountEditOk() function testing POST positive and don't work ParseForm", func(t *testing.T) {
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

		fnHomeAuth(w, r)
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

		fnAccountEditOk(w1, r1) // calling the tested function
		status := w1.Code
		if status != http.StatusOK {
			t.Errorf("fnAccountEditOk() error: %v", status)
		}
	})

	t.Run("fnAccountEditOk() function testing POST positive and not value of login", func(t *testing.T) {
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

		fnHomeAuth(w, r)
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

		fnAccountEditOk(w1, r1) // calling the tested function
		status := w1.Code
		if status != http.StatusOK {
			t.Errorf("fnAccountEditOk() error: %v", status)
		}
	})

	t.Run("fnAccountEditOk() function testing POST positive and not value of password", func(t *testing.T) {
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

		fnHomeAuth(w, r)
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

		fnAccountEditOk(w1, r1) // calling the tested function
		status := w1.Code
		if status != http.StatusOK {
			t.Errorf("fnAccountEditOk() error: %v", status)
		}
	})

	t.Run("fnAccountEditOk() function testing - status error", func(t *testing.T) {
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

		fnHomeAuth(w, r)
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

		fnAccountEditOk(w1, r1) // calling the tested function
		status := w1.Code
		if status != http.StatusOK {
			t.Errorf("fnAccountEditOk() error: %v", status)
		}
	})

	t.Run("fnAccountEditOk() function testing - Roles", func(t *testing.T) {
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

		fnHomeAuth(w, r)
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

		fnAccountEditOk(w1, r1) // calling the tested function
		status := w1.Code
		if status != http.StatusOK {
			t.Errorf("fnAccountEditOk() error: %v", status)
		}
	})

	t.Run("fnAccountEditOk() function testing - root user", func(t *testing.T) {
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

		fnHomeAuth(w, r)
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

		fnAccountEditOk(w1, r1) // calling the tested function
		status := w1.Code
		if status != http.StatusOK {
			t.Errorf("fnAccountEditOk() error: %v", status)
		}
	})
}

func Test_fnAccountBanLoadForm(t *testing.T) {
	gauth.Start()
	parseTemplates()

	t.Run("fnAccountBanLoadForm() function testing - Isolate negative", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/", nil)

		fnAccountBanLoadForm(w, r) // calling the tested function
		status := w.Code
		if status != http.StatusOK {
			t.Errorf("fnAccountBanLoadForm() error: %v", status)
		}
	})

	t.Run("fnAccountBanLoadForm() function testing - root user", func(t *testing.T) {
		w := httptest.NewRecorder()
		form := url.Values{}
		form.Add("username", "root")
		form.Add("password", "toor")
		r := httptest.NewRequest("POST", "/", strings.NewReader(form.Encode()))
		r.PostForm = form

		fnHomeAuth(w, r)
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

		fnAccountBanLoadForm(w1, r1) // calling the tested function
		status := w1.Code
		if status != http.StatusOK {
			t.Errorf("fnAccountBanLoadForm() error: %v", status)
		}
	})

	t.Run("fnAccountBanLoadForm() function testing - all right", func(t *testing.T) {
		w := httptest.NewRecorder()
		form := url.Values{}
		form.Add("username", "root")
		form.Add("password", "toor")
		r := httptest.NewRequest("POST", "/", strings.NewReader(form.Encode()))
		r.PostForm = form

		fnHomeAuth(w, r)
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

		fnAccountBanLoadForm(w1, r1) // calling the tested function
		status := w1.Code
		if status != http.StatusOK {
			t.Errorf("fnAccountBanLoadForm() error: %v", status)
		}
	})
}

func Test_fnAccountBanOk(t *testing.T) {
	gauth.Start()
	parseTemplates()

	t.Run("fnAccountBanOk() function testing - Isolate negative", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/", nil)

		fnAccountBanOk(w, r) // calling the tested function
		status := w.Code
		if status != http.StatusOK {
			t.Errorf("fnAccountBanOk() error: %v", status)
		}
	})

	t.Run("fnAccountBanOk() function testing - POST negative", func(t *testing.T) {
		w := httptest.NewRecorder()
		form := url.Values{}
		form.Add("username", "root")
		form.Add("password", "toor")
		r := httptest.NewRequest("POST", "/", strings.NewReader(form.Encode()))
		r.PostForm = form

		fnHomeAuth(w, r)
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

		fnAccountBanOk(w1, r1) // calling the tested function
		status := w1.Code
		if status != http.StatusOK {
			t.Errorf("fnAccountBanOk() error: %v", status)
		}
	})

	t.Run("fnAccountBanOk() function testing POST positive and don't work ParseForm", func(t *testing.T) {
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

		fnHomeAuth(w, r)
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

		fnAccountBanOk(w1, r1) // calling the tested function
		status := w1.Code
		if status != http.StatusOK {
			t.Errorf("fnAccountBanOk() error: %v", status)
		}
	})

	t.Run("fnAccountBanOk() function testing POST positive and not value of login", func(t *testing.T) {
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

		fnHomeAuth(w, r)
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

		fnAccountBanOk(w1, r1) // calling the tested function
		status := w1.Code
		if status != http.StatusOK {
			t.Errorf("fnAccountBanOk() error: %v", status)
		}
	})

	t.Run("fnAccountBanOk() function testing POST positive and don't work BlockUser", func(t *testing.T) {
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

		fnHomeAuth(w, r)
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

		fnAccountBanOk(w1, r1) // calling the tested function
		status := w1.Code
		if status != http.StatusOK {
			t.Errorf("fnAccountBanOk() error: %v", status)
		}
	})

	t.Run("fnAccountBanOk() function testing - all right", func(t *testing.T) {
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

		fnHomeAuth(w, r)
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

		fnAccountBanOk(w1, r1) // calling the tested function
		status := w1.Code
		if status != http.StatusOK {
			t.Errorf("fnAccountBanOk() error: %v", status)
		}
	})
}

func Test_fnAccountUnbanLoadForm(t *testing.T) {
	gauth.Start()
	parseTemplates()

	t.Run("fnAccountUnbanLoadForm() function testing - Isolate negative", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/", nil)

		fnAccountUnbanLoadForm(w, r) // calling the tested function
		status := w.Code
		if status != http.StatusOK {
			t.Errorf("fnAccountUnbanLoadForm() error: %v", status)
		}
	})

	t.Run("fnAccountUnbanLoadForm() function testing - root user", func(t *testing.T) {
		w := httptest.NewRecorder()
		form := url.Values{}
		form.Add("username", "root")
		form.Add("password", "toor")
		r := httptest.NewRequest("POST", "/", strings.NewReader(form.Encode()))
		r.PostForm = form

		fnHomeAuth(w, r)
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

		fnAccountUnbanLoadForm(w1, r1) // calling the tested function
		status := w1.Code
		if status != http.StatusOK {
			t.Errorf("fnAccountUnbanLoadForm() error: %v", status)
		}
	})

	t.Run("fnAccountUnbanLoadForm() function testing - all right", func(t *testing.T) {
		w := httptest.NewRecorder()
		form := url.Values{}
		form.Add("username", "root")
		form.Add("password", "toor")
		r := httptest.NewRequest("POST", "/", strings.NewReader(form.Encode()))
		r.PostForm = form

		fnHomeAuth(w, r)
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

		fnAccountUnbanLoadForm(w1, r1) // calling the tested function
		status := w1.Code
		if status != http.StatusOK {
			t.Errorf("fnAccountUnbanLoadForm() error: %v", status)
		}
	})
}

func Test_fnAccountUnbanOk(t *testing.T) {
	gauth.Start()
	parseTemplates()

	t.Run("fnAccountUnbanOk() function testing - Isolate negative", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/", nil)

		fnAccountUnbanOk(w, r) // calling the tested function
		status := w.Code
		if status != http.StatusOK {
			t.Errorf("fnAccountUnbanOk() error: %v", status)
		}
	})

	t.Run("fnAccountUnbanOk() function testing - POST negative", func(t *testing.T) {
		w := httptest.NewRecorder()
		form := url.Values{}
		form.Add("username", "root")
		form.Add("password", "toor")
		r := httptest.NewRequest("POST", "/", strings.NewReader(form.Encode()))
		r.PostForm = form

		fnHomeAuth(w, r)
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

		fnAccountUnbanOk(w1, r1) // calling the tested function
		status := w1.Code
		if status != http.StatusOK {
			t.Errorf("fnAccountUnbanOk() error: %v", status)
		}
	})

	t.Run("fnAccountUnbanOk() function testing POST positive and don't work ParseForm", func(t *testing.T) {
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

		fnHomeAuth(w, r)
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

		fnAccountUnbanOk(w1, r1) // calling the tested function
		status := w1.Code
		if status != http.StatusOK {
			t.Errorf("fnAccountUnbanOk() error: %v", status)
		}
	})

	t.Run("fnAccountUnbanOk() function testing POST positive and not value of login", func(t *testing.T) {
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

		fnHomeAuth(w, r)
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

		fnAccountUnbanOk(w1, r1) // calling the tested function
		status := w1.Code
		if status != http.StatusOK {
			t.Errorf("fnAccountUnbanOk() error: %v", status)
		}
	})

	t.Run("fnAccountUnbanOk() function testing POST positive and don't work UnblockUser", func(t *testing.T) {
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

		fnHomeAuth(w, r)
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

		fnAccountUnbanOk(w1, r1) // calling the tested function
		status := w1.Code
		if status != http.StatusOK {
			t.Errorf("fnAccountUnbanOk() error: %v", status)
		}
	})

	t.Run("fnAccountUnbanOk() function testing - all right", func(t *testing.T) {
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

		fnHomeAuth(w, r)
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

		fnAccountUnbanOk(w1, r1) // calling the tested function
		status := w1.Code
		if status != http.StatusOK {
			t.Errorf("fnAccountUnbanOk() error: %v", status)
		}
	})
}

func Test_fnAccountDelLoadForm(t *testing.T) {
	gauth.Start()
	parseTemplates()

	t.Run("fnAccountDelLoadForm() function testing - Isolate negative", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/", nil)

		fnAccountDelLoadForm(w, r) // calling the tested function
		status := w.Code
		if status != http.StatusOK {
			t.Errorf("fnAccountDelLoadForm() error: %v", status)
		}
	})

	t.Run("fnAccountDelLoadForm() function testing - root user", func(t *testing.T) {
		w := httptest.NewRecorder()
		form := url.Values{}
		form.Add("username", "root")
		form.Add("password", "toor")
		r := httptest.NewRequest("POST", "/", strings.NewReader(form.Encode()))
		r.PostForm = form

		fnHomeAuth(w, r)
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

		fnAccountDelLoadForm(w1, r1) // calling the tested function
		status := w1.Code
		if status != http.StatusOK {
			t.Errorf("fnAccountDelLoadForm() error: %v", status)
		}
	})

	t.Run("fnAccountDelLoadForm() function testing - all right", func(t *testing.T) {
		w := httptest.NewRecorder()
		form := url.Values{}
		form.Add("username", "root")
		form.Add("password", "toor")
		r := httptest.NewRequest("POST", "/", strings.NewReader(form.Encode()))
		r.PostForm = form

		fnHomeAuth(w, r)
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

		fnAccountDelLoadForm(w1, r1) // calling the tested function
		status := w1.Code
		if status != http.StatusOK {
			t.Errorf("fnAccountDelLoadForm() error: %v", status)
		}
	})
}

func Test_fnAccountDelOk(t *testing.T) {
	gauth.Start()
	parseTemplates()

	t.Run("fnAccountDelOk() function testing - Isolate negative", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/", nil)

		fnAccountDelOk(w, r) // calling the tested function
		status := w.Code
		if status != http.StatusOK {
			t.Errorf("fnAccountDelOk() error: %v", status)
		}
	})

	t.Run("fnAccountDelOk() function testing - POST negative", func(t *testing.T) {
		w := httptest.NewRecorder()
		form := url.Values{}
		form.Add("username", "root")
		form.Add("password", "toor")
		r := httptest.NewRequest("POST", "/", strings.NewReader(form.Encode()))
		r.PostForm = form

		fnHomeAuth(w, r)
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

		fnAccountDelOk(w1, r1) // calling the tested function
		status := w1.Code
		if status != http.StatusOK {
			t.Errorf("fnAccountDelOk() error: %v", status)
		}
	})

	t.Run("fnAccountDelOk() function testing POST positive and don't work ParseForm", func(t *testing.T) {
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

		fnHomeAuth(w, r)
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

		fnAccountDelOk(w1, r1) // calling the tested function
		status := w1.Code
		if status != http.StatusOK {
			t.Errorf("fnAccountDelOk() error: %v", status)
		}
	})

	t.Run("fnAccountDelOk() function testing POST positive and not value of login", func(t *testing.T) {
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

		fnHomeAuth(w, r)
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

		fnAccountDelOk(w1, r1) // calling the tested function
		status := w1.Code
		if status != http.StatusOK {
			t.Errorf("fnAccountDelOk() error: %v", status)
		}
	})

	t.Run("fnAccountDelOk() function testing POST positive and don't work DeleteUser", func(t *testing.T) {
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

		fnHomeAuth(w, r)
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

		fnAccountDelOk(w1, r1) // calling the tested function
		status := w1.Code
		if status != http.StatusOK {
			t.Errorf("fnAccountDelOk() error: %v", status)
		}
	})

	t.Run("fnAccountDelOk() function testing - all right", func(t *testing.T) {
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

		fnHomeAuth(w, r)
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

		fnAccountDelOk(w1, r1) // calling the tested function
		status := w1.Code
		if status != http.StatusOK {
			t.Errorf("fnAccountDelOk() error: %v", status)
		}
	})
}

func Test_fnNavSettings(t *testing.T) {
	gauth.Start()
	parseTemplates()

	t.Run("fnNavSettings() function testing - Isolate negative", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/", nil)

		fnNavSettings(w, r) // calling the tested function
		status := w.Code
		if status != http.StatusOK {
			t.Errorf("fnNavSettings() error: %v", status)
		}
	})

	t.Run("fnNavSettings() function testing - Isolate positive", func(t *testing.T) {
		w := httptest.NewRecorder()
		form := url.Values{}
		form.Add("username", "root")
		form.Add("password", "toor")
		r := httptest.NewRequest("POST", "/", strings.NewReader(form.Encode()))
		r.PostForm = form

		fnHomeAuth(w, r)
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

		fnNavSettings(w1, r1) // calling the tested function
		status := w1.Code
		if status != http.StatusOK {
			t.Errorf("fnNavSettings() error: %v", status)
		}
	})
}

func Test_fnSettingsCoreFriendlyChangeSw(t *testing.T) {
	gauth.Start()
	parseTemplates()

	t.Run("fnSettingsCoreFriendlyChangeSw() function testing - Isolate negative", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/", nil)

		fnSettingsCoreFriendlyChangeSw(w, r) // calling the tested function
		status := w.Code
		if status != http.StatusOK {
			t.Errorf("fnSettingsCoreFriendlyChangeSw() error: %v", status)
		}
	})

	t.Run("fnSettingsCoreFriendlyChangeSw() function testing - off", func(t *testing.T) {
		config.StDefaultConfig.CoreSettings.FriendlyMode = true

		w := httptest.NewRecorder()
		form := url.Values{}
		form.Add("username", "root")
		form.Add("password", "toor")
		r := httptest.NewRequest("POST", "/", strings.NewReader(form.Encode()))
		r.PostForm = form

		fnHomeAuth(w, r)
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

		fnSettingsCoreFriendlyChangeSw(w1, r1) // calling the tested function
		status := w1.Code
		if status != http.StatusOK || !core.StLocalCoreSettings.FriendlyMode {
			t.Errorf("fnSettingsCoreFriendlyChangeSw() error: %v", status)
		}
	})

	t.Run("fnSettingsCoreFriendlyChangeSw() function testing - on", func(t *testing.T) {
		config.StDefaultConfig.CoreSettings.FriendlyMode = false

		w := httptest.NewRecorder()
		form := url.Values{}
		form.Add("username", "root")
		form.Add("password", "toor")
		r := httptest.NewRequest("POST", "/", strings.NewReader(form.Encode()))
		r.PostForm = form

		fnHomeAuth(w, r)
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

		fnSettingsCoreFriendlyChangeSw(w1, r1) // calling the tested function
		status := w1.Code
		if status != http.StatusOK || core.StLocalCoreSettings.FriendlyMode {
			t.Errorf("fnSettingsCoreFriendlyChangeSw() error: %v", status)
		}
	})
}

func Test_fnSettingsWScChangeSw(t *testing.T) {
	gauth.Start()
	parseTemplates()

	t.Run("fnSettingsWScChangeSw() function testing - Isolate negative", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/", nil)

		fnSettingsWScChangeSw(w, r) // calling the tested function
		status := w.Code
		if status != http.StatusOK {
			t.Errorf("fnSettingsWScChangeSw() error: %v", status)
		}
	})

	t.Run("fnSettingsWScChangeSw() function testing - shutdown", func(t *testing.T) {
		config.StDefaultConfig.WebSocketConnector.Enable = true
		go websocketconn.Start(&config.StDefaultConfig)
		closer.AddHandler(websocketconn.Shutdown) // Register a shutdown handler.

		w := httptest.NewRecorder()
		form := url.Values{}
		form.Add("username", "root")
		form.Add("password", "toor")
		r := httptest.NewRequest("POST", "/", strings.NewReader(form.Encode()))
		r.PostForm = form

		fnHomeAuth(w, r)
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

		fnSettingsWScChangeSw(w1, r1) // calling the tested function
		status := w1.Code
		if status != http.StatusOK {
			t.Errorf("fnSettingsWScChangeSw() error: %v", status)
		}
	})

	t.Run("fnSettingsWScChangeSw() function testing - start", func(t *testing.T) {
		config.StDefaultConfig.WebSocketConnector.Enable = false
		closer.RunAndDelHandler(websocketconn.Shutdown)

		w := httptest.NewRecorder()
		form := url.Values{}
		form.Add("username", "root")
		form.Add("password", "toor")
		r := httptest.NewRequest("POST", "/", strings.NewReader(form.Encode()))
		r.PostForm = form

		fnHomeAuth(w, r)
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

		fnSettingsWScChangeSw(w1, r1) // calling the tested function
		status := w1.Code
		if status != http.StatusOK {
			t.Errorf("fnSettingsWScChangeSw() error: %v", status)
		}
	})
}

func Test_fnSettingsRestChangeSw(t *testing.T) {
	gauth.Start()
	parseTemplates()

	t.Run("fnSettingsRestChangeSw() function testing - Isolate negative", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/", nil)

		fnSettingsRestChangeSw(w, r) // calling the tested function
		status := w.Code
		if status != http.StatusOK {
			t.Errorf("fnSettingsRestChangeSw() error: %v", status)
		}
	})

	t.Run("fnSettingsRestChangeSw() function testing - shutdown", func(t *testing.T) {
		config.StDefaultConfig.RestConnector.Enable = true
		go rest.Start(&config.StDefaultConfig)
		closer.AddHandler(rest.Shutdown) // Register a shutdown handler.

		w := httptest.NewRecorder()
		form := url.Values{}
		form.Add("username", "root")
		form.Add("password", "toor")
		r := httptest.NewRequest("POST", "/", strings.NewReader(form.Encode()))
		r.PostForm = form

		fnHomeAuth(w, r)
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

		fnSettingsRestChangeSw(w1, r1) // calling the tested function
		status := w1.Code
		if status != http.StatusOK {
			t.Errorf("fnSettingsRestChangeSw() error: %v", status)
		}
	})

	t.Run("fnSettingsRestChangeSw() function testing - start", func(t *testing.T) {
		config.StDefaultConfig.RestConnector.Enable = false
		closer.RunAndDelHandler(rest.Shutdown)

		w := httptest.NewRecorder()
		form := url.Values{}
		form.Add("username", "root")
		form.Add("password", "toor")
		r := httptest.NewRequest("POST", "/", strings.NewReader(form.Encode()))
		r.PostForm = form

		fnHomeAuth(w, r)
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

		fnSettingsRestChangeSw(w1, r1) // calling the tested function
		status := w1.Code
		if status != http.StatusOK {
			t.Errorf("fnSettingsRestChangeSw() error: %v", status)
		}
	})
}

func Test_fnSettingsGrpcChangeSw(t *testing.T) {
	gauth.Start()
	parseTemplates()

	t.Run("fnSettingsGrpcChangeSw() function testing - Isolate negative", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/", nil)

		fnSettingsGrpcChangeSw(w, r) // calling the tested function
		status := w.Code
		if status != http.StatusOK {
			t.Errorf("fnSettingsGrpcChangeSw() error: %v", status)
		}
	})
}

func Test_fnSettingsWebChangeSw(t *testing.T) {
	gauth.Start()
	parseTemplates()

	t.Run("fnSettingsWebChangeSw() function testing - Isolate negative", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/", nil)

		fnSettingsWebChangeSw(w, r) // calling the tested function
		status := w.Code
		if status != http.StatusOK {
			t.Errorf("fnSettingsWebChangeSw() error: %v", status)
		}
	})

	t.Run("fnSettingsWebChangeSw() function testing - Isolate positive", func(t *testing.T) {
		w := httptest.NewRecorder()
		form := url.Values{}
		form.Add("username", "root")
		form.Add("password", "toor")
		r := httptest.NewRequest("POST", "/", strings.NewReader(form.Encode()))
		r.PostForm = form

		fnHomeAuth(w, r)
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

		fnSettingsWebChangeSw(w1, r1) // calling the tested function
		status := w1.Code
		if status != http.StatusOK {
			t.Errorf("fnSettingsWebChangeSw() error: %v", status)
		}
	})

	t.Run("fnSettingsWebChangeSw() function testing - Isolate undefined user", func(t *testing.T) {
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

		fnHomeAuth(w, r)
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

		fnSettingsWebChangeSw(w1, r1) // calling the tested function
		status := w1.Code
		if status != http.StatusOK {
			t.Errorf("fnSettingsWebChangeSw() error: %v", status)
		}
	})
}
