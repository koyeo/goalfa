package exporter

import (
	"github.com/fatih/structs"
	"github.com/flosch/pongo2/v5"
)

type RenderData struct {
	Methods []*RenderMethod
	Structs []*RenderField
}

type RenderMethod struct {
	Name        string
	InputType   string
	OutputType  string
	Description string
	Method      string
	Path        string
}

type RenderField struct {
	Name        string
	Type        string
	Description string
	Json        string
	Required    bool
	Fields      []*RenderField
}

type Namer func(string) string
type Formatter func(string) (string, error)
type Typer func(s string, isStruct, isArray bool) string

var EmptyNamer Namer = func(s string) string {
	return s
}

var EmptyTyper Typer = func(s string, isStruct, isArray bool) string {
	return s
}

var EmptyFormatter Formatter = func(s string) (string, error) {
	return s, nil
}

func MakeRenderData(methods []*Method, namer Namer, typer Typer) (data *RenderData) {
	data = new(RenderData)
	checker := newRenderFieldChecker()
	for _, v := range methods {
		data.Methods = append(data.Methods, makeRenderMethod(v, namer, typer))
		data.Structs = append(data.Structs, makeRenderFields(v, namer, typer, checker)...)
	}
	return
}

func makeRenderMethod(method *Method, namer Namer, typer Typer) (renderMethod *RenderMethod) {
	renderMethod = new(RenderMethod)
	renderMethod.Name = namer(method.Name)
	renderMethod.Description = method.Description
	renderMethod.Method = method.Method
	renderMethod.Path = method.Path
	if method.Input != nil {
		renderMethod.InputType = makeMethodIOName(method.Input, typer)
	}
	if method.Output != nil {
		renderMethod.OutputType = makeMethodIOName(method.Output, typer)
	}
	return
}

func makeMethodIOName(field *Field, typer Typer) string {
	if field.Array {
		return getNestedType(field, typer)
	} else if field.Struct {
		return typer(field.Type, true, false)
	} else {
		return typer(field.Type, false, false)
	}
}

func getNestedType(field *Field, typer Typer) string {
	if field.Array {
		return typer(getNestedType(field.Elem, typer), field.Struct, field.Array)
	} else {
		return typer(field.Type, field.Struct, field.Array)
	}
}

func makeRenderFields(method *Method, namer Namer, typer Typer, checker *renderFieldChecker) (renderFields []*RenderField) {
	if method.Input != nil && (method.Input.Struct || (method.Input.Array && method.Input.Nested)) {
		toRenderFields(method.Input, namer, typer, &renderFields, nil, checker)
	}
	if method.Output != nil && (method.Output.Struct || (method.Output.Array && method.Output.Nested)) {
		toRenderFields(method.Output, namer, typer, &renderFields, nil, checker)
	}
	return
}

func toRenderFields(field *Field, namer Namer, typer Typer, renderFields *[]*RenderField, parent *RenderField, checker *renderFieldChecker) {
	if field.Array && field.Nested { // 处理嵌套数组对象
		toRenderFields(field.Elem, namer, typer, renderFields, nil, checker)
	} else {
		renderField := new(RenderField)
		name := field.Name
		if field.Struct { // 处理对象
			name = field.Type
			for _, v := range field.Fields {
				if v.Struct {
					toRenderFields(v, namer, typer, renderFields, renderField, checker)
				} else if v.Array && v.Nested {
					toRenderFields(v.Elem, namer, typer, renderFields, renderField, checker)
				} else {
					toRenderFields(v, namer, typer, renderFields, renderField, checker)
				}
			}
		}
		if field.Array { // 处理嵌套基础类型
			renderField.Type = getNestedType(field, typer)
		} else {
			renderField.Type = typer(field.Type, field.Struct, field.Array)
		}
		renderField.Name = namer(name)
		renderField.Description = field.Description
		renderField.Json = field.Name
		if parent != nil {
			parent.Fields = append(parent.Fields, renderField)
		}
		if field.Struct && !checker.Has(renderField.Name) {
			if checker != nil {
				checker.Add(renderField.Name)
			}
			*renderFields = append(*renderFields, renderField)
		}
	}
	return
}

func Render(tpl string, data interface{}, formatter Formatter) (result string, err error) {
	_tpl, err := pongo2.FromString(tpl)
	if err != nil {
		return
	}
	ctx := structs.Map(data)
	result, err = _tpl.Execute(ctx)
	if err != nil {
		return
	}
	result, err = formatter(result)
	if err != nil {
		return
	}
	return
}

func newRenderFieldChecker() *renderFieldChecker {
	return &renderFieldChecker{
		cache: map[string]bool{},
	}
}

type renderFieldChecker struct {
	cache map[string]bool
}

func (p *renderFieldChecker) Has(name string) bool {
	_, ok := p.cache[name]
	return ok
}

func (p *renderFieldChecker) Add(name string) {
	p.cache[name] = true
}
