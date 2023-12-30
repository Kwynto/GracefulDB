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
	// ts, err := template.New(HOME_TEMP_NAME).Parse(home_masq.HtmlHome)
	// if err != nil {
	// 	slog.Debug("Error reading the template", slog.String("err", err.Error()))
	// 	return
	// }
	// TemplatesMap[HOME_TEMP_NAME] = ts

	loadTemplateFromVar(AUTH_TEMP_NAME, auth_masq.HtmlAuth)
	// ts, err = template.New(AUTH_TEMP_NAME).Parse(auth_masq.HtmlAuth)
	// if err != nil {
	// 	slog.Debug("Error reading the template", slog.String("err", err.Error()))
	// 	return
	// }
	// TemplatesMap[AUTH_TEMP_NAME] = ts

	loadTemplateFromVar(BLOCK_TEMP_DEFAULT, htmx_masq.Default)
	// ts, err = template.New(BLOCK_TEMP_DEFAULT).Parse(htmx_masq.Default)
	// if err != nil {
	// 	slog.Debug("Error reading the template", slog.String("err", err.Error()))
	// 	return
	// }
	// TemplatesMap[BLOCK_TEMP_DEFAULT] = ts

	loadTemplateFromVar(BLOCK_TEMP_ACCESS_DENIED, htmx_masq.AccessDenied)
	// ts, err = template.New(BLOCK_TEMP_ACCESS_DENIED).Parse(htmx_masq.AccessDenied)
	// if err != nil {
	// 	slog.Debug("Error reading the template", slog.String("err", err.Error()))
	// 	return
	// }
	// TemplatesMap[BLOCK_TEMP_ACCESS_DENIED] = ts

	loadTemplateFromVar(BLOCK_TEMP_DASHBOARD, htmx_masq.Dashboard)
	// TemplatesMap[BLOCK_TEMP_DASHBOARD] = ts
	// ts, err = template.New(BLOCK_TEMP_DASHBOARD).Parse(htmx_masq.Dashboard)
	// if err != nil {
	// 	slog.Debug("Error reading the template", slog.String("err", err.Error()))
	// 	return
	// }
	// TemplatesMap[BLOCK_TEMP_DASHBOARD] = ts

	loadTemplateFromVar(BLOCK_TEMP_DATABASES, htmx_masq.Databases)
	// ts, err = template.New(BLOCK_TEMP_DATABASES).Parse(htmx_masq.Databases)
	// if err != nil {
	// 	slog.Debug("Error reading the template", slog.String("err", err.Error()))
	// 	return
	// }
	// TemplatesMap[BLOCK_TEMP_DATABASES] = ts

	loadTemplateFromVar(BLOCK_TEMP_ACCOUNTS, htmx_masq.Accounts)
	// ts, err = template.New(BLOCK_TEMP_ACCOUNTS).Parse(htmx_masq.Accounts)
	// if err != nil {
	// 	slog.Debug("Error reading the template", slog.String("err", err.Error()))
	// 	return
	// }
	// TemplatesMap[BLOCK_TEMP_ACCOUNTS] = ts

	loadTemplateFromVar(BLOCK_TEMP_ACCOUNT_CREATE_FORM_OK, htmx_masq.AccountCreateFormOk)
	// ts, err = template.New(BLOCK_TEMP_ACCOUNT_CREATE_FORM_OK).Parse(htmx_masq.AccountCreateFormOk)
	// if err != nil {
	// 	slog.Debug("Error reading the template", slog.String("err", err.Error()))
	// 	return
	// }
	// TemplatesMap[BLOCK_TEMP_ACCOUNT_CREATE_FORM_OK] = ts

	loadTemplateFromVar(BLOCK_TEMP_ACCOUNT_CREATE_FORM_LOAD, htmx_masq.AccountCreateFormLoad)
	// ts, err = template.New(BLOCK_TEMP_ACCOUNT_CREATE_FORM_LOAD).Parse(htmx_masq.AccountCreateFormLoad)
	// if err != nil {
	// 	slog.Debug("Error reading the template", slog.String("err", err.Error()))
	// 	return
	// }
	// TemplatesMap[BLOCK_TEMP_ACCOUNT_CREATE_FORM_LOAD] = ts

	loadTemplateFromVar(BLOCK_TEMP_ACCOUNT_CREATE_FORM_ERROR, htmx_masq.AccountCreateFormError)
	// ts, err = template.New(BLOCK_TEMP_ACCOUNT_CREATE_FORM_ERROR).Parse(htmx_masq.AccountCreateFormError)
	// if err != nil {
	// 	slog.Debug("Error reading the template", slog.String("err", err.Error()))
	// 	return
	// }
	// TemplatesMap[BLOCK_TEMP_ACCOUNT_CREATE_FORM_ERROR] = ts

	loadTemplateFromVar(BLOCK_TEMP_ACCOUNT_EDIT_FORM_OK, htmx_masq.AccountEditFormOk)
	// ts, err = template.New(BLOCK_TEMP_ACCOUNT_EDIT_FORM_OK).Parse(htmx_masq.AccountEditFormOk)
	// if err != nil {
	// 	slog.Debug("Error reading the template", slog.String("err", err.Error()))
	// 	return
	// }
	// TemplatesMap[BLOCK_TEMP_ACCOUNT_EDIT_FORM_OK] = ts

	loadTemplateFromVar(BLOCK_TEMP_ACCOUNT_EDIT_FORM_LOAD, htmx_masq.AccountEditFormLoad)
	// ts, err = template.New(BLOCK_TEMP_ACCOUNT_EDIT_FORM_LOAD).Parse(htmx_masq.AccountEditFormLoad)
	// if err != nil {
	// 	slog.Debug("Error reading the template", slog.String("err", err.Error()))
	// 	return
	// }
	// TemplatesMap[BLOCK_TEMP_ACCOUNT_EDIT_FORM_LOAD] = ts

	loadTemplateFromVar(BLOCK_TEMP_ACCOUNT_EDIT_FORM_ERROR, htmx_masq.AccountEditFormError)
	// ts, err = template.New(BLOCK_TEMP_ACCOUNT_EDIT_FORM_ERROR).Parse(htmx_masq.AccountEditFormError)
	// if err != nil {
	// 	slog.Debug("Error reading the template", slog.String("err", err.Error()))
	// 	return
	// }
	// TemplatesMap[BLOCK_TEMP_ACCOUNT_EDIT_FORM_ERROR] = ts

	loadTemplateFromVar(BLOCK_TEMP_ACCOUNT_BAN_FORM_OK, htmx_masq.AccountBanFormOk)
	// ts, err = template.New(BLOCK_TEMP_ACCOUNT_BAN_FORM_OK).Parse(htmx_masq.AccountBanFormOk)
	// if err != nil {
	// 	slog.Debug("Error reading the template", slog.String("err", err.Error()))
	// 	return
	// }
	// TemplatesMap[BLOCK_TEMP_ACCOUNT_BAN_FORM_OK] = ts

	loadTemplateFromVar(BLOCK_TEMP_ACCOUNT_BAN_FORM_LOAD, htmx_masq.AccountBanFormLoad)
	// ts, err = template.New(BLOCK_TEMP_ACCOUNT_BAN_FORM_LOAD).Parse(htmx_masq.AccountBanFormLoad)
	// if err != nil {
	// 	slog.Debug("Error reading the template", slog.String("err", err.Error()))
	// 	return
	// }
	// TemplatesMap[BLOCK_TEMP_ACCOUNT_BAN_FORM_LOAD] = ts

	loadTemplateFromVar(BLOCK_TEMP_ACCOUNT_BAN_FORM_ERROR, htmx_masq.AccountBanFormError)
	// ts, err = template.New(BLOCK_TEMP_ACCOUNT_BAN_FORM_ERROR).Parse(htmx_masq.AccountBanFormError)
	// if err != nil {
	// 	slog.Debug("Error reading the template", slog.String("err", err.Error()))
	// 	return
	// }
	// TemplatesMap[BLOCK_TEMP_ACCOUNT_BAN_FORM_ERROR] = ts

	loadTemplateFromVar(BLOCK_TEMP_ACCOUNT_UNBAN_FORM_OK, htmx_masq.AccountUnBanFormOk)
	// ts, err = template.New(BLOCK_TEMP_ACCOUNT_UNBAN_FORM_OK).Parse(htmx_masq.AccountUnBanFormOk)
	// if err != nil {
	// 	slog.Debug("Error reading the template", slog.String("err", err.Error()))
	// 	return
	// }
	// TemplatesMap[BLOCK_TEMP_ACCOUNT_UNBAN_FORM_OK] = ts

	loadTemplateFromVar(BLOCK_TEMP_ACCOUNT_UNBAN_FORM_LOAD, htmx_masq.AccountUnBanFormLoad)
	// ts, err = template.New(BLOCK_TEMP_ACCOUNT_UNBAN_FORM_LOAD).Parse(htmx_masq.AccountUnBanFormLoad)
	// if err != nil {
	// 	slog.Debug("Error reading the template", slog.String("err", err.Error()))
	// 	return
	// }
	// TemplatesMap[BLOCK_TEMP_ACCOUNT_UNBAN_FORM_LOAD] = ts

	loadTemplateFromVar(BLOCK_TEMP_ACCOUNT_UNBAN_FORM_ERROR, htmx_masq.AccountUnBanFormError)
	// ts, err = template.New(BLOCK_TEMP_ACCOUNT_UNBAN_FORM_ERROR).Parse(htmx_masq.AccountUnBanFormError)
	// if err != nil {
	// 	slog.Debug("Error reading the template", slog.String("err", err.Error()))
	// 	return
	// }
	// TemplatesMap[BLOCK_TEMP_ACCOUNT_UNBAN_FORM_ERROR] = ts

	loadTemplateFromVar(BLOCK_TEMP_ACCOUNT_DEL_FORM_OK, htmx_masq.AccountDelFormOk)
	// ts, err = template.New(BLOCK_TEMP_ACCOUNT_DEL_FORM_OK).Parse(htmx_masq.AccountDelFormOk)
	// if err != nil {
	// 	slog.Debug("Error reading the template", slog.String("err", err.Error()))
	// 	return
	// }
	// TemplatesMap[BLOCK_TEMP_ACCOUNT_DEL_FORM_OK] = ts

	loadTemplateFromVar(BLOCK_TEMP_ACCOUNT_DEL_FORM_LOAD, htmx_masq.AccountDelFormLoad)
	// ts, err = template.New(BLOCK_TEMP_ACCOUNT_DEL_FORM_LOAD).Parse(htmx_masq.AccountDelFormLoad)
	// if err != nil {
	// 	slog.Debug("Error reading the template", slog.String("err", err.Error()))
	// 	return
	// }
	// TemplatesMap[BLOCK_TEMP_ACCOUNT_DEL_FORM_LOAD] = ts

	loadTemplateFromVar(BLOCK_TEMP_ACCOUNT_DEL_FORM_ERROR, htmx_masq.AccountDelFormError)
	// ts, err = template.New(BLOCK_TEMP_ACCOUNT_DEL_FORM_ERROR).Parse(htmx_masq.AccountDelFormError)
	// if err != nil {
	// 	slog.Debug("Error reading the template", slog.String("err", err.Error()))
	// 	return
	// }
	// TemplatesMap[BLOCK_TEMP_ACCOUNT_DEL_FORM_ERROR] = ts

	loadTemplateFromVar(BLOCK_TEMP_ACCOUNT_SELFEDIT_OK, htmx_masq.SelfEditFormOk)
	// ts, err = template.New(BLOCK_TEMP_ACCOUNT_SELFEDIT_OK).Parse(htmx_masq.SelfEditFormOk)
	// if err != nil {
	// 	slog.Debug("Error reading the template", slog.String("err", err.Error()))
	// 	return
	// }
	// TemplatesMap[BLOCK_TEMP_ACCOUNT_SELFEDIT_OK] = ts

	loadTemplateFromVar(BLOCK_TEMP_ACCOUNT_SELFEDIT_LOAD, htmx_masq.SelfEditFormLoad)
	// ts, err = template.New(BLOCK_TEMP_ACCOUNT_SELFEDIT_LOAD).Parse(htmx_masq.SelfEditFormLoad)
	// if err != nil {
	// 	slog.Debug("Error reading the template", slog.String("err", err.Error()))
	// 	return
	// }
	// TemplatesMap[BLOCK_TEMP_ACCOUNT_SELFEDIT_LOAD] = ts

	loadTemplateFromVar(BLOCK_TEMP_ACCOUNT_SELFEDIT_ERROR, htmx_masq.SelfEditFormError)
	// ts, err = template.New(BLOCK_TEMP_ACCOUNT_SELFEDIT_ERROR).Parse(htmx_masq.SelfEditFormError)
	// if err != nil {
	// 	slog.Debug("Error reading the template", slog.String("err", err.Error()))
	// 	return
	// }
	// TemplatesMap[BLOCK_TEMP_ACCOUNT_SELFEDIT_ERROR] = ts

	loadTemplateFromVar(BLOCK_TEMP_SETTINGS, htmx_masq.Settings)
	// ts, err = template.New(BLOCK_TEMP_SETTINGS).Parse(htmx_masq.Settings)
	// if err != nil {
	// 	slog.Debug("Error reading the template", slog.String("err", err.Error()))
	// 	return
	// }
	// TemplatesMap[BLOCK_TEMP_SETTINGS] = ts
}
