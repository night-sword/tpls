package demo

import (
	"github.com/night-sword/tpls"
)

type TemplateName uint

//go:generate enumer --type=TemplateName --extramethod --linecomment --output=template_name_enum.go
const (
	TemplateDemo TemplateName = iota + 1 // tpl/demo.tmpl
)

func convertSlice(slice TemplateNameSlice) []tpls.TemplateName {
	ts := make([]tpls.TemplateName, len([]TemplateName(slice)))
	for i := range slice {
		ts[i] = slice[i]
	}

	return ts
}
