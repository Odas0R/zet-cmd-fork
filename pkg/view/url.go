package view

import (
	"fmt"

	"github.com/a-h/templ"
)

func url(format string, a ...any) templ.SafeURL {
	return templ.URL(fmt.Sprintf(format, a...))
}
