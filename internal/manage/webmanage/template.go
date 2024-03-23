package webmanage

import (
	"embed"
	"html/template"
	"log/slog"

	"github.com/Kwynto/GracefulDB/pkg/lib/e"
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
	BLOCK_TEMP_ACCOUNT_SELFEDIT_OK       = "ui/html/selfedit-form-ok.html"
	BLOCK_TEMP_ACCOUNT_SELFEDIT_LOAD     = "ui/html/selfedit-form-load.html"
	BLOCK_TEMP_ACCOUNT_SELFEDIT_ERROR    = "ui/html/selfedit-form-error.html"

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
	LoadTemplateFromEmbed(HOME_TEMP_NAME)
	LoadTemplateFromEmbed(AUTH_TEMP_NAME)

	LoadTemplateFromEmbed(BLOCK_TEMP_DEFAULT)
	LoadTemplateFromEmbed(BLOCK_TEMP_ACCESS_DENIED)

	LoadTemplateFromEmbed(BLOCK_TEMP_DASHBOARD)

	LoadTemplateFromEmbed(BLOCK_TEMP_DATABASES)
	LoadTemplateFromEmbed(BLOCK_TEMP_DATABASE_REQUEST_ANSWER)

	LoadTemplateFromEmbed(BLOCK_TEMP_ACCOUNTS)

	LoadTemplateFromEmbed(BLOCK_TEMP_ACCOUNT_CREATE_FORM_OK)
	LoadTemplateFromEmbed(BLOCK_TEMP_ACCOUNT_CREATE_FORM_LOAD)
	LoadTemplateFromEmbed(BLOCK_TEMP_ACCOUNT_CREATE_FORM_ERROR)

	LoadTemplateFromEmbed(BLOCK_TEMP_ACCOUNT_EDIT_FORM_OK)
	LoadTemplateFromEmbed(BLOCK_TEMP_ACCOUNT_EDIT_FORM_LOAD)
	LoadTemplateFromEmbed(BLOCK_TEMP_ACCOUNT_EDIT_FORM_ERROR)

	LoadTemplateFromEmbed(BLOCK_TEMP_ACCOUNT_BAN_FORM_OK)
	LoadTemplateFromEmbed(BLOCK_TEMP_ACCOUNT_BAN_FORM_LOAD)
	LoadTemplateFromEmbed(BLOCK_TEMP_ACCOUNT_BAN_FORM_ERROR)

	LoadTemplateFromEmbed(BLOCK_TEMP_ACCOUNT_UNBAN_FORM_OK)
	LoadTemplateFromEmbed(BLOCK_TEMP_ACCOUNT_UNBAN_FORM_LOAD)
	LoadTemplateFromEmbed(BLOCK_TEMP_ACCOUNT_UNBAN_FORM_ERROR)

	LoadTemplateFromEmbed(BLOCK_TEMP_ACCOUNT_DEL_FORM_OK)
	LoadTemplateFromEmbed(BLOCK_TEMP_ACCOUNT_DEL_FORM_LOAD)
	LoadTemplateFromEmbed(BLOCK_TEMP_ACCOUNT_DEL_FORM_ERROR)

	LoadTemplateFromEmbed(BLOCK_TEMP_ACCOUNT_SELFEDIT_OK)
	LoadTemplateFromEmbed(BLOCK_TEMP_ACCOUNT_SELFEDIT_LOAD)
	LoadTemplateFromEmbed(BLOCK_TEMP_ACCOUNT_SELFEDIT_ERROR)

	LoadTemplateFromEmbed(BLOCK_TEMP_SETTINGS)
}
