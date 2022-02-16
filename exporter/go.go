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

func (g GoMaker) Lang() string {
	return Go
}

func (g GoMaker) Make(methods []*Method) (files []*File, err error) {
	data := MakeRenderData(g.Lang(), methods, GoNamer, GoTyper)
	serviceFile := new(File)
	serviceFile.Name = "service.make.go"
	serviceFile.Content, err = Render(goServiceTpl, data, GoFormatter)
	if err != nil {
		return
	}
	files = append(files, serviceFile, serviceFile)
	return
}

const goServiceTpl = `
package sdk

import "context"
{% for package in Packages %}import "{{ package.From }}"
{% endfor %}

type sdk interface {
{% for method in Methods %}    {{ method.Name }}(ctx context.Context{% if method.InputType !='' %},in {{ method.InputType }}{% endif %})({% if method.OutputType !='' %}out {{ method.OutputType }},{% endif %} err error) {% if method.Description %}// {{ method.Description }}{% endif %}
{% endfor %}
}

var _ sdk = new(SDK)

type SDK struct {
}

func (s SDK) request(method string, path string) {

}

{% for method in Methods %}
{% if method.Description %}// {{ method.Name }} {{ method.Description }}{% endif %}
func (s SDK){{ method.Name }}(ctx context.Context{% if method.InputType !='' %},in {{ method.InputType }}{% endif %})({% if method.OutputType !='' %}out {{ method.OutputType }},{% endif %} err error){
    // {{ method.Method }} {{ method.Path }}
	return
}
{% endfor %}

{% for struct in Structs %}
type {{ struct.Name }} struct {
	{% for field in struct.Fields %} {{ field.Name }} {{ field.Type }} ` + "{% if field.Param != '' %}`" + `json:"{{ field.Param }}"` + "`{% endif %}" + `   {% if field.Description or field.Label %}// {{ field.Label }} {{ field.Description }}{% endif %}
    {% endfor %}}
{% endfor %}
`
