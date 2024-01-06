package htmx_masq

import (
	"errors"
	"strings"
)

var Default string = `
<div class="container-fluid pt-4 px-4">
    <div class="bg-secondary text-center rounded p-4">
		<i class="fa fa-exclamation-triangle"></i> Error: Bad request
    </div>
</div>
`

var AccessDenied string = `
<div class="container-fluid pt-4 px-4">
    <div class="bg-secondary text-center rounded p-4">
		<i class="fa fa-exclamation-circle"></i> Warning: You are not allowed access to this section.
    </div>
</div>
`

var Dashboard string = `
<div class="container-fluid pt-4 px-4">
    <div class="bg-secondary text-center rounded p-4">
	    <h4>Dashboard</h4>
    </div>
</div>
`

func DefaultBlank() (string, error) {
	str := strings.TrimSpace(Default)
	if str == "" {
		return str, errors.New("an empty template")
	}
	return Default, nil
}
