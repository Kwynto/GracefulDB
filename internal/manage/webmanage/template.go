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

	BLOCK_TEMP_DATABASES = "ui/html/databases.html"
	// BLOCK_TEMP_DATABASE_REQUEST_ANSWER = "ui/html/dbreqanswer.html"

	BLOCK_TEMP_CONSOLE                = "ui/html/console.html"
	BLOCK_TEMP_CONSOLE_REQUEST_ANSWER = "ui/html/termreqanswer.html"

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

var MTemplates = make(map[string]*template.Template)

var (
	emHtmlDir   *embed.FS
	emStaticDir *embed.FS
)

func SetUiDirs(emHtmlFS *embed.FS, emStaticFS *embed.FS) {
	emHtmlDir = emHtmlFS
	emStaticDir = emStaticFS
}

func LoadTemplateFromString(sName string, sTemp string) (err error) {
	sOperation := "internal -> WebManage -> isolated -> loadTemplateFromVar"
	defer func() { e.Wrapper(sOperation, err) }()

	stTemp, err := template.New(sName).Parse(sTemp)
	if err != nil {
		slog.Debug("Error reading the template", slog.String("err", err.Error()))
		return err
	}
	MTemplates[sName] = stTemp
	return nil
}

func LoadTemplateFromEmbed(sName string) (err error) {
	sOperation := "internal -> WebManage -> isolated -> loadTemplateFromEmbed"
	defer func() { e.Wrapper(sOperation, err) }()

	slBytes, err := emHtmlDir.ReadFile(sName)
	if err != nil {
		slog.Debug("Error reading the template from Embed", slog.String("err", err.Error()))
		return err
	}
	sText := string(slBytes)

	stTemp, err := template.New(sName).Parse(sText)
	if err != nil {
		slog.Debug("Error reading the template", slog.String("err", err.Error()))
		return err
	}
	MTemplates[sName] = stTemp
	return nil
}

func parseTemplates() {
	LoadTemplateFromEmbed(HOME_TEMP_NAME)
	LoadTemplateFromEmbed(AUTH_TEMP_NAME)

	LoadTemplateFromEmbed(BLOCK_TEMP_DEFAULT)
	LoadTemplateFromEmbed(BLOCK_TEMP_ACCESS_DENIED)

	LoadTemplateFromEmbed(BLOCK_TEMP_DASHBOARD)

	LoadTemplateFromEmbed(BLOCK_TEMP_DATABASES)
	// LoadTemplateFromEmbed(BLOCK_TEMP_DATABASE_REQUEST_ANSWER)

	LoadTemplateFromEmbed(BLOCK_TEMP_CONSOLE)
	LoadTemplateFromEmbed(BLOCK_TEMP_CONSOLE_REQUEST_ANSWER)

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
