package home_masq

import (
	"errors"
	"fmt"
	"strings"
)

var HtmlHome string = fmt.Sprint(html1, bootstrap_min_css, html2, style, html3, htmx_min_js, html4_jsblank, main_js, html5)

func Default() (string, error) {
	str := strings.TrimSpace(HtmlHome)
	if str == "" {
		return str, errors.New("an empty template")
	}
	return HtmlHome, nil
}
