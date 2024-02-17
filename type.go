package tpls

import (
	"text/template"
)

type TemplateName interface {
	String() string
}

type TemplateMap map[TemplateName]*template.Template
