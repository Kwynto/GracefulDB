package webmanage

import (
	"html/template"
	"log/slog"

	"github.com/Kwynto/GracefulDB/pkg/lib/helpers/masquerade/auth_masq"
	"github.com/Kwynto/GracefulDB/pkg/lib/helpers/masquerade/home_masq"
	"github.com/Kwynto/GracefulDB/pkg/lib/helpers/masquerade/htmx_masq"
)

const (
	// The names of the templates in the cache
	HOME_TEMP_NAME = "Home"
	AUTH_TEMP_NAME = "Auth"

	BLOCK_TEMP_DEFAULT                   = "Default"
	BLOCK_TEMP_ACCESS_DENIED             = "AccessDenied"
	BLOCK_TEMP_DASHBOARD                 = "Dashboard"
	BLOCK_TEMP_DATABASES                 = "Databases"
	BLOCK_TEMP_ACCOUNTS                  = "Accounts"
	BLOCK_TEMP_ACCOUNT_CREATE_FORM_OK    = "AccountCreateFormOk"
	BLOCK_TEMP_ACCOUNT_CREATE_FORM_LOAD  = "AccountCreateFormLoad"
	BLOCK_TEMP_ACCOUNT_CREATE_FORM_ERROR = "AccountCreateFormError"
	BLOCK_TEMP_ACCOUNT_EDIT_FORM_OK      = "AccountEditFormOk"
	BLOCK_TEMP_ACCOUNT_EDIT_FORM_LOAD    = "AccountEditFormLoad"
	BLOCK_TEMP_ACCOUNT_EDIT_FORM_ERROR   = "AccountEditFormError"
	BLOCK_TEMP_ACCOUNT_BAN_FORM_OK       = "AccountBanFormOk"
	BLOCK_TEMP_ACCOUNT_BAN_FORM_LOAD     = "AccountBanFormLoad"
	BLOCK_TEMP_ACCOUNT_BAN_FORM_ERROR    = "AccountBanFormError"
	BLOCK_TEMP_ACCOUNT_UNBAN_FORM_OK     = "AccountUnBanFormOk"
	BLOCK_TEMP_ACCOUNT_UNBAN_FORM_LOAD   = "AccountUnBanFormLoad"
	BLOCK_TEMP_ACCOUNT_UNBAN_FORM_ERROR  = "AccountUnBanFormError"
	BLOCK_TEMP_ACCOUNT_DEL_FORM_OK       = "AccountDelFormOk"
	BLOCK_TEMP_ACCOUNT_DEL_FORM_LOAD     = "AccountDelFormLoad"
	BLOCK_TEMP_ACCOUNT_DEL_FORM_ERROR    = "AccountDelFormError"
	BLOCK_TEMP_ACCOUNT_SELFEDIT_OK       = "AccountSelfeditFormOk"
	BLOCK_TEMP_ACCOUNT_SELFEDIT_LOAD     = "AccountSelfeditFormLoad"
	BLOCK_TEMP_ACCOUNT_SELFEDIT_ERROR    = "AccountSelfeditFormError"
	BLOCK_TEMP_SETTINGS                  = "Settings"
)

var TemplatesMap = make(map[string]*template.Template)

func loadTemplateFromVar(name string, temp string) {
	ts, err := template.New(name).Parse(temp)
	if err != nil {
		slog.Debug("Error reading the template", slog.String("err", err.Error()))
		return
	}
	TemplatesMap[name] = ts
}

func parseTemplates() {
	loadTemplateFromVar(HOME_TEMP_NAME, home_masq.HtmlHome)
	loadTemplateFromVar(AUTH_TEMP_NAME, auth_masq.HtmlAuth)

	loadTemplateFromVar(BLOCK_TEMP_DEFAULT, htmx_masq.Default)
	loadTemplateFromVar(BLOCK_TEMP_ACCESS_DENIED, htmx_masq.AccessDenied)

	loadTemplateFromVar(BLOCK_TEMP_DASHBOARD, htmx_masq.Dashboard)

	loadTemplateFromVar(BLOCK_TEMP_DATABASES, htmx_masq.Databases)

	loadTemplateFromVar(BLOCK_TEMP_ACCOUNTS, htmx_masq.Accounts)

	loadTemplateFromVar(BLOCK_TEMP_ACCOUNT_CREATE_FORM_OK, htmx_masq.AccountCreateFormOk)
	loadTemplateFromVar(BLOCK_TEMP_ACCOUNT_CREATE_FORM_LOAD, htmx_masq.AccountCreateFormLoad)
	loadTemplateFromVar(BLOCK_TEMP_ACCOUNT_CREATE_FORM_ERROR, htmx_masq.AccountCreateFormError)

	loadTemplateFromVar(BLOCK_TEMP_ACCOUNT_EDIT_FORM_OK, htmx_masq.AccountEditFormOk)
	loadTemplateFromVar(BLOCK_TEMP_ACCOUNT_EDIT_FORM_LOAD, htmx_masq.AccountEditFormLoad)
	loadTemplateFromVar(BLOCK_TEMP_ACCOUNT_EDIT_FORM_ERROR, htmx_masq.AccountEditFormError)

	loadTemplateFromVar(BLOCK_TEMP_ACCOUNT_BAN_FORM_OK, htmx_masq.AccountBanFormOk)
	loadTemplateFromVar(BLOCK_TEMP_ACCOUNT_BAN_FORM_LOAD, htmx_masq.AccountBanFormLoad)
	loadTemplateFromVar(BLOCK_TEMP_ACCOUNT_BAN_FORM_ERROR, htmx_masq.AccountBanFormError)

	loadTemplateFromVar(BLOCK_TEMP_ACCOUNT_UNBAN_FORM_OK, htmx_masq.AccountUnBanFormOk)
	loadTemplateFromVar(BLOCK_TEMP_ACCOUNT_UNBAN_FORM_LOAD, htmx_masq.AccountUnBanFormLoad)
	loadTemplateFromVar(BLOCK_TEMP_ACCOUNT_UNBAN_FORM_ERROR, htmx_masq.AccountUnBanFormError)

	loadTemplateFromVar(BLOCK_TEMP_ACCOUNT_DEL_FORM_OK, htmx_masq.AccountDelFormOk)
	loadTemplateFromVar(BLOCK_TEMP_ACCOUNT_DEL_FORM_LOAD, htmx_masq.AccountDelFormLoad)
	loadTemplateFromVar(BLOCK_TEMP_ACCOUNT_DEL_FORM_ERROR, htmx_masq.AccountDelFormError)

	loadTemplateFromVar(BLOCK_TEMP_ACCOUNT_SELFEDIT_OK, htmx_masq.SelfEditFormOk)
	loadTemplateFromVar(BLOCK_TEMP_ACCOUNT_SELFEDIT_LOAD, htmx_masq.SelfEditFormLoad)
	loadTemplateFromVar(BLOCK_TEMP_ACCOUNT_SELFEDIT_ERROR, htmx_masq.SelfEditFormError)

	loadTemplateFromVar(BLOCK_TEMP_SETTINGS, htmx_masq.Settings)
}
