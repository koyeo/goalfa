package exporter

import (
	"fmt"
	"github.com/fatih/structs"
	"github.com/flosch/pongo2/v5"
	"strings"
)

type RenderData struct {
	Methods []*RenderMethod
	Structs []*RenderStruct
}

type RenderMethod struct {
	Name        string
	Input       string
	Output      string
	Description string
	Method      string
	Path        string
}

type RenderStruct struct {
	Name        string
	Description string
	Fields      []*RenderField
}

type RenderField struct {
	Name        string
	Type        string
	Description string
	Json        string
	Required    bool
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
	for _, v := range methods {
		data.Methods = append(data.Methods, makeRenderMethod(v, namer, typer))
		data.Structs = append(data.Structs, makeRenderStruct(v, namer, typer)...)
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
		renderMethod.Input = makeMethodIOName(method.Input, typer)
	}
	if method.Output != nil {
		renderMethod.Output = makeMethodIOName(method.Output, typer)
	}
	return
}

func makeMethodIOName(field *Field, typer Typer) string {
	if field.Struct {
		return typer(field.Name, true, false)
	} else if field.Array {
		if field.Elem.Struct {
			return typer(fmt.Sprintf("%s%s", field.Name, strings.Title(field.Elem.Name)), true, true)
		} else {
			return typer(field.Elem.Type, false, true)
		}
	} else {
		return typer(field.Type, false, false)
	}
}

func makeRenderStruct(method *Method, namer Namer, typer Typer) (structs []*RenderStruct) {
	if method.Input != nil {
		structs = append(structs, fieldToRenderStruct(method.Input.Name, method.Input, namer, typer)...)
	}
	if method.Output != nil {
		structs = append(structs, fieldToRenderStruct(method.Output.Name, method.Output, namer, typer)...)
	}
	return
}

func fieldToRenderStruct(name string, field *Field, namer Namer, typer Typer) (structs []*RenderStruct) {
	if !field.Struct {
		return
	}
	rs := new(RenderStruct)
	rs.Name = name
	structs = append(structs, rs)
	for _, v := range field.Fields {
		if v.Struct {
			tn := fmt.Sprintf("%s%s", name, strings.Title(v.Name))
			rs.Fields = append(rs.Fields, &RenderField{
				Name:        namer(v.Name),
				Type:        typer(tn, true, false),
				Description: v.Description,
				Required:    false,
				Json:        v.Name,
			})
			structs = append(structs, fieldToRenderStruct(tn, v, namer, typer)...)
		} else if v.Array {
			if v.Elem == nil {
				continue
			}
			v = v.Elem
			var _type string
			if v.Struct {
				t := fmt.Sprintf("%s%sItem", name, strings.Title(v.Name))
				_type = typer(t, true, true)
				structs = append(structs, fieldToRenderStruct(t, v, namer, typer)...)
			} else if v.Array {
				if v.Elem == nil {
					continue
				}
				v = v.Elem
				t := fmt.Sprintf("%s%sItemElem", name, strings.Title(v.Name))
				structs = append(structs, fieldToRenderStruct(t, v, namer, typer)...)
			} else {
				_type = typer(v.Type, false, true)
			}
			rs.Fields = append(rs.Fields, &RenderField{
				Name:        namer(v.Name),
				Type:        _type,
				Description: v.Description,
				Required:    false,
				Json:        v.Name,
			})
		} else {
			rs.Fields = append(rs.Fields, &RenderField{
				Name:        namer(v.Name),
				Type:        typer(v.Type, false, false),
				Description: v.Description,
				Required:    false,
				Json:        v.Name,
			})
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
