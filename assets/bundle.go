package assets

import (
	"embed"

	"github.com/Kwynto/GracefulDB/internal/manage/webmanage"
)

var (
	//go:embed ui/html
	uiHtmlDir embed.FS

	//go:embed ui/static
	uiStaticDir embed.FS
)

func init() {
	// Set UI
	webmanage.SetUiDirs(&uiHtmlDir, &uiStaticDir)
}
