package webmanage

import (
	"embed"
	"html/template"
	"log/slog"

	"github.com/Kwynto/GracefulDB/pkg/lib/e"
	"github.com/Kwynto/GracefulDB/pkg/lib/helpers/masquerade/htmx_masq"
)

const (
	// The names of the templates in the cache
	HOME_TEMP_NAME = "ui/html/home.html"
	AUTH_TEMP_NAME = "ui/html/auth.html"

	BLOCK_TEMP_DEFAULT       = "ui/html/default.html"
	BLOCK_TEMP_ACCESS_DENIED = "ui/html/accessdenied.html"
	BLOCK_TEMP_DASHBOARD     = "ui/html/dashboard.html"

	BLOCK_TEMP_DATABASES               = "ui/html/databases.html"
	BLOCK_TEMP_DATABASE_REQUEST_ANSWER = "ui/html/dbreqanswer.html"

	BLOCK_TEMP_ACCOUNTS                  = "ui/html/accounts.html"
	BLOCK_TEMP_ACCOUNT_CREATE_FORM_OK    = "ui/html/account-create-form-ok.html"
	BLOCK_TEMP_ACCOUNT_CREATE_FORM_LOAD  = "ui/html/account-create-form-load.html"
	BLOCK_TEMP_ACCOUNT_CREATE_FORM_ERROR = "ui/html/account-create-form-error.html"
	BLOCK_TEMP_ACCOUNT_EDIT_FORM_OK      = "ui/html/account-edit-form-ok.html"
	BLOCK_TEMP_ACCOUNT_EDIT_FORM_LOAD    = "ui/html/account-edit-form-load.html"
	BLOCK_TEMP_ACCOUNT_EDIT_FORM_ERROR   = "ui/html/account-edit-form-error.html"
	BLOCK_TEMP_ACCOUNT_BAN_FORM_OK       = "ui/html/account-ban-form-ok.html"
	BLOCK_TEMP_ACCOUNT_BAN_FORM_LOAD     = "ui/html/account-ban-form-load.html"
	BLOCK_TEMP_ACCOUNT_BAN_FORM_ERROR    = "ui/html/account-ban-form-error.html"
	BLOCK_TEMP_ACCOUNT_UNBAN_FORM_OK     = "ui/html/account-unban-form-ok.html"
	BLOCK_TEMP_ACCOUNT_UNBAN_FORM_LOAD   = "ui/html/account-unban-form-load.html"
	BLOCK_TEMP_ACCOUNT_UNBAN_FORM_ERROR  = "ui/html/account-unban-form-error.html"
	BLOCK_TEMP_ACCOUNT_DEL_FORM_OK       = "ui/html/account-del-form-ok.html"
	BLOCK_TEMP_ACCOUNT_DEL_FORM_LOAD     = "ui/html/account-del-form-load.html"
	BLOCK_TEMP_ACCOUNT_DEL_FORM_ERROR    = "ui/html/account-del-form-error.html"
	BLOCK_TEMP_ACCOUNT_SELFEDIT_OK       = "AccountSelfeditFormOk"
	BLOCK_TEMP_ACCOUNT_SELFEDIT_LOAD     = "AccountSelfeditFormLoad"
	BLOCK_TEMP_ACCOUNT_SELFEDIT_ERROR    = "AccountSelfeditFormError"

	BLOCK_TEMP_SETTINGS = "ui/html/settings.html"
)

var TemplatesMap = make(map[string]*template.Template)

var (
	uiHtmlDir   *embed.FS
	uiStaticDir *embed.FS
)

func SetUiDirs(uiHtmlFS *embed.FS, uiStaticFS *embed.FS) {
	uiHtmlDir = uiHtmlFS
	uiStaticDir = uiStaticFS
}

func LoadTemplateFromString(name string, temp string) (err error) {
	op := "internal -> WebManage -> isolated -> loadTemplateFromVar"
	defer func() { e.Wrapper(op, err) }()

	ts, err := template.New(name).Parse(temp)
	if err != nil {
		slog.Debug("Error reading the template", slog.String("err", err.Error()))
		return err
	}
	TemplatesMap[name] = ts
	return nil
}

func LoadTemplateFromEmbed(name string) (err error) {
	op := "internal -> WebManage -> isolated -> loadTemplateFromEmbed"
	defer func() { e.Wrapper(op, err) }()

	bytes, err := uiHtmlDir.ReadFile(name)
	if err != nil {
		slog.Debug("Error reading the template from Embed", slog.String("err", err.Error()))
		return err
	}
	str := string(bytes)

	ts, err := template.New(name).Parse(str)
	if err != nil {
		slog.Debug("Error reading the template", slog.String("err", err.Error()))
		return err
	}
	TemplatesMap[name] = ts
	return nil
}

func parseTemplates() {
	// LoadTemplateFromString(HOME_TEMP_NAME, home_masq.HtmlHome)
	LoadTemplateFromEmbed(HOME_TEMP_NAME)
	// LoadTemplateFromString(AUTH_TEMP_NAME, auth_masq.HtmlAuth)
	LoadTemplateFromEmbed(AUTH_TEMP_NAME)

	// LoadTemplateFromString(BLOCK_TEMP_DEFAULT, htmx_masq.Default)
	LoadTemplateFromEmbed(BLOCK_TEMP_DEFAULT)
	// LoadTemplateFromString(BLOCK_TEMP_ACCESS_DENIED, htmx_masq.AccessDenied)
	LoadTemplateFromEmbed(BLOCK_TEMP_ACCESS_DENIED)

	// LoadTemplateFromString(BLOCK_TEMP_DASHBOARD, htmx_masq.Dashboard)
	LoadTemplateFromEmbed(BLOCK_TEMP_DASHBOARD)

	// LoadTemplateFromString(BLOCK_TEMP_DATABASES, htmx_masq.Databases)
	LoadTemplateFromEmbed(BLOCK_TEMP_DATABASES)
	// LoadTemplateFromString(BLOCK_TEMP_DATABASE_REQUEST_ANSWER, htmx_masq.DatabaseRequestAnswer)
	LoadTemplateFromEmbed(BLOCK_TEMP_DATABASE_REQUEST_ANSWER)

	// LoadTemplateFromString(BLOCK_TEMP_ACCOUNTS, htmx_masq.Accounts)
	LoadTemplateFromEmbed(BLOCK_TEMP_ACCOUNTS)

	// LoadTemplateFromString(BLOCK_TEMP_ACCOUNT_CREATE_FORM_OK, htmx_masq.AccountCreateFormOk)
	LoadTemplateFromEmbed(BLOCK_TEMP_ACCOUNT_CREATE_FORM_OK)
	// LoadTemplateFromString(BLOCK_TEMP_ACCOUNT_CREATE_FORM_LOAD, htmx_masq.AccountCreateFormLoad)
	LoadTemplateFromEmbed(BLOCK_TEMP_ACCOUNT_CREATE_FORM_LOAD)
	// LoadTemplateFromString(BLOCK_TEMP_ACCOUNT_CREATE_FORM_ERROR, htmx_masq.AccountCreateFormError)
	LoadTemplateFromEmbed(BLOCK_TEMP_ACCOUNT_CREATE_FORM_ERROR)

	// LoadTemplateFromString(BLOCK_TEMP_ACCOUNT_EDIT_FORM_OK, htmx_masq.AccountEditFormOk)
	LoadTemplateFromEmbed(BLOCK_TEMP_ACCOUNT_EDIT_FORM_OK)
	// LoadTemplateFromString(BLOCK_TEMP_ACCOUNT_EDIT_FORM_LOAD, htmx_masq.AccountEditFormLoad)
	LoadTemplateFromEmbed(BLOCK_TEMP_ACCOUNT_EDIT_FORM_LOAD)
	// LoadTemplateFromString(BLOCK_TEMP_ACCOUNT_EDIT_FORM_ERROR, htmx_masq.AccountEditFormError)
	LoadTemplateFromEmbed(BLOCK_TEMP_ACCOUNT_EDIT_FORM_ERROR)

	// LoadTemplateFromString(BLOCK_TEMP_ACCOUNT_BAN_FORM_OK, htmx_masq.AccountBanFormOk)
	LoadTemplateFromEmbed(BLOCK_TEMP_ACCOUNT_BAN_FORM_OK)
	// LoadTemplateFromString(BLOCK_TEMP_ACCOUNT_BAN_FORM_LOAD, htmx_masq.AccountBanFormLoad)
	LoadTemplateFromEmbed(BLOCK_TEMP_ACCOUNT_BAN_FORM_LOAD)
	// LoadTemplateFromString(BLOCK_TEMP_ACCOUNT_BAN_FORM_ERROR, htmx_masq.AccountBanFormError)
	LoadTemplateFromEmbed(BLOCK_TEMP_ACCOUNT_BAN_FORM_ERROR)

	// LoadTemplateFromString(BLOCK_TEMP_ACCOUNT_UNBAN_FORM_OK, htmx_masq.AccountUnBanFormOk)
	LoadTemplateFromEmbed(BLOCK_TEMP_ACCOUNT_UNBAN_FORM_OK)
	// LoadTemplateFromString(BLOCK_TEMP_ACCOUNT_UNBAN_FORM_LOAD, htmx_masq.AccountUnBanFormLoad)
	LoadTemplateFromEmbed(BLOCK_TEMP_ACCOUNT_UNBAN_FORM_LOAD)
	// LoadTemplateFromString(BLOCK_TEMP_ACCOUNT_UNBAN_FORM_ERROR, htmx_masq.AccountUnBanFormError)
	LoadTemplateFromEmbed(BLOCK_TEMP_ACCOUNT_UNBAN_FORM_ERROR)

	// LoadTemplateFromString(BLOCK_TEMP_ACCOUNT_DEL_FORM_OK, htmx_masq.AccountDelFormOk)
	LoadTemplateFromEmbed(BLOCK_TEMP_ACCOUNT_DEL_FORM_OK)
	// LoadTemplateFromString(BLOCK_TEMP_ACCOUNT_DEL_FORM_LOAD, htmx_masq.AccountDelFormLoad)
	LoadTemplateFromEmbed(BLOCK_TEMP_ACCOUNT_DEL_FORM_LOAD)
	// LoadTemplateFromString(BLOCK_TEMP_ACCOUNT_DEL_FORM_ERROR, htmx_masq.AccountDelFormError)
	LoadTemplateFromEmbed(BLOCK_TEMP_ACCOUNT_DEL_FORM_ERROR)

	LoadTemplateFromString(BLOCK_TEMP_ACCOUNT_SELFEDIT_OK, htmx_masq.SelfEditFormOk)
	LoadTemplateFromString(BLOCK_TEMP_ACCOUNT_SELFEDIT_LOAD, htmx_masq.SelfEditFormLoad)
	LoadTemplateFromString(BLOCK_TEMP_ACCOUNT_SELFEDIT_ERROR, htmx_masq.SelfEditFormError)

	// LoadTemplateFromString(BLOCK_TEMP_SETTINGS, htmx_masq.Settings)
	LoadTemplateFromEmbed(BLOCK_TEMP_SETTINGS)
}
