package generator

import (
	"strings"
	"text/template"
)

var funcMap = template.FuncMap{
	"join":    strings.Join,
	"lower":   strings.ToLower,
	"title":   strings.Title, //nolint:staticcheck
	"replace": strings.ReplaceAll,
	"indent": func(n int, s string) string {
		pad := strings.Repeat(" ", n)
		return pad + strings.ReplaceAll(s, "\n", "\n"+pad)
	},
}
