package exporter

import (
	"reflect"
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
		field.Name = t.Field(i).Name
		field.Type = t.Field(i).Type.String()
		root.Children = append(root.Children, field)
	}
	fields = append(fields, root)
	return fields
}
