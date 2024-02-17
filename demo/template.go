package demo

import (
	"embed"
	_ "embed"

	"github.com/night-sword/tpls"
)

//go:embed tpl/*.tmpl
var templates embed.FS

type Template struct {
	*tpls.Template
}

func NewTemplate() *Template {
	return &Template{
		Template: tpls.NewTemplate(templates, convertSlice(TemplateNameValues())),
	}
}
