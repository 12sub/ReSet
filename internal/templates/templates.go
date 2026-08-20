package templates

import (
	"embed"
	"html/template"
)

//go:embed *.html
var FS embed.FS

var T *template.Template

func init() {
	T = template.Must(template.ParseFS(FS, "*.html"))
}