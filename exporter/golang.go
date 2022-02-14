package exporter

import (
	"fmt"
	"go/format"
	"strings"
)

var GoNamer Namer = func(s string) string {
	return strings.Title(s)
}

var GoTyper Typer = func(s string, isStruct, isArray bool) string {
	if isArray {
		if isStruct {
			s = fmt.Sprintf("[]*%s", s)
		} else {
			s = fmt.Sprintf("[]%s", s)
		}
	} else {
		if isStruct {
			s = fmt.Sprintf("*%s", s)
		}
	}
	return s
}

var GoFormatter Formatter = func(s string) (r string, err error) {
	bytes, err := format.Source([]byte(s))
	if err != nil {
		lines := strings.Split(s, "\n")
		for k, v := range lines {
			fmt.Printf("%d: %s\n", k+1, v)
		}
		return
	}
	r = string(bytes)
	return
}

type GoMaker struct {
}

func (g GoMaker) Make(pkg string, methods []*Method) (files []*File, err error) {
	data := MakeRenderData(methods, GoNamer, GoTyper)
	methodsFile := new(File)
	methodsFile.Name = "methods.make.go"
	methodsFile.Content, err = Render(goServiceTpl, data, GoFormatter)
	if err != nil {
		return
	}
	typesFile := new(File)
	typesFile.Name = "types.make.go"
	typesFile.Content, err = Render(goStructTpl, data, GoFormatter)
	if err != nil {
		return
	}
	files = append(files, methodsFile, typesFile)
	return
}

const goServiceTpl = `
package sdk

import "context"

func request(method string, path string) {

}

{% for method in Methods %}
{% if method.Description %}// {{ method.Name }} {{ method.Description }}{% endif %}
func {{ method.Name }}(ctx context.Context{% if method.InputType !='' %},in {{ method.InputType }}{% endif %})({% if method.OutputType !='' %}out {{ method.OutputType }},{% endif %} err error){
    // {{ method.Method }} {{ method.Path }}
	return
}
{% endfor %}
`

const goStructTpl = `
package sdk

{% for struct in Structs %}
type {{ struct.Name }} struct {
	{% for field in struct.Fields %} {{ field.Name }} {{ field.Type }} ` + "`" + `json:"{{ field.Json }}"` + "`" + `   {% if field.Description %}// {{ field.Description }}{% endif %}
    {% endfor %}}
{% endfor %}
`
