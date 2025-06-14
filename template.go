package tpls

import (
	"bytes"
	"embed"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"text/template"

	"github.com/Masterminds/sprig/v3"
	"github.com/night-sword/kratos-kit/errors"
)

type Template struct {
	tpls TemplateMap
}

func NewTemplate(fs embed.FS, tplNames []TemplateName) *Template {
	inst := &Template{}
	tpls, err := inst.init(fs, tplNames)
	if err != nil {
		panic(err)
	}

	inst.tpls = tpls
	return inst
}

func (inst *Template) init(fs embed.FS, tplNames []TemplateName) (templates TemplateMap, err error) {
	templates = make(TemplateMap, len(tplNames))

	for _, name := range tplNames {
		meta := map[string]string{"name": name.String()}
		cnt, e := fs.ReadFile(name.String())
		if err = e; err != nil {
			err = errors.InternalServer("Read tpl cnt fail", err.Error()).WithCause(errors.Unrecoverable).WithMetadata(meta)
			return
		}

		tpl, e := template.New(name.String()).
			Funcs(sprig.TxtFuncMap()).
			Parse(string(cnt))

		if err = e; err != nil {
			err = errors.InternalServer("Parse template fail", err.Error()).WithCause(errors.Unrecoverable).WithMetadata(meta)
			return
		}

		templates[name] = tpl
	}

	return
}

func (inst *Template) RenderTrim(name TemplateName, params any, maxBlankLine int) (cnt *string, err error) {
	cnt, err = inst.Render(name, params)

	re := regexp.MustCompile(fmt.Sprintf(`\n{%d,}`, maxBlankLine+1))
	replace := strings.Repeat("\n", maxBlankLine)
	*cnt = re.ReplaceAllString(*cnt, replace)
	return
}

func (inst *Template) Render(name TemplateName, params any) (cnt *string, err error) {
	meta := map[string]string{"name": name.String()}
	tpl, ok := inst.tpls[name]
	if !ok {
		err = errors.InternalServer(errors.RsnInternal, "tpl not found").WithCause(errors.Unrecoverable).WithMetadata(meta)
		return
	}

	buffer := &bytes.Buffer{}
	err = tpl.Execute(buffer, params)
	if err != nil {
		j, e := json.Marshal(params)
		if e == nil {
			meta["params"] = string(j)
		}
		err = errors.InternalServer("render template fail", err.Error()).WithCause(errors.Unrecoverable).WithMetadata(meta)
		return
	}

	_cnt := buffer.String()
	cnt = &_cnt
	return
}
