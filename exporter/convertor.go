package exporter

import (
	"reflect"
	"strings"
)

// ReflectFields 反射转换输入输出的字段信息
func ReflectFields(name string, t reflect.Type) []*Field {
	for {
		if t.Kind() != reflect.Ptr {
			break
		}
		t = t.Elem()
	}
	var fields []*Field
	root := new(Field)
	root.Name = name
	for i := 0; i < t.NumField(); i++ {
		field := new(Field)
		field.Name = getJsonField(t.Field(i))
		field.Label = getFieldLabel(t.Field(i))
		field.Required = getFieldRequired(t.Field(i))
		field.Type = t.Field(i).Type.String()
		root.Children = append(root.Children, field)
	}
	fields = append(fields, root)
	return fields
}

func getFieldLabel(field reflect.StructField) string {
	return field.Tag.Get("label")
}

func getFieldRequired(field reflect.StructField) bool {
	return strings.Contains(field.Tag.Get("validator"), "required")
}

func getJsonField(field reflect.StructField) string {
	n := field.Tag.Get("json")
	n = strings.ReplaceAll(n, ",omitempty", "")
	if n == "" {
		n = field.Name
	}
	return n
}
