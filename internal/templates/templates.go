package templates

import (
	"embed"
	"html/template"
)

//go:embed *.html components/*.html
var FS embed.FS

var T *template.Template

func init() {
	T = template.Must(template.ParseFS(FS, "*.html", "components/*.html"))
}