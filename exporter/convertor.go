package exporter

import (
	"github.com/koyeo/buck/utils"
	"reflect"
	"strings"
)

// ReflectFields 反射转换输入输出的字段信息
func ReflectFields(name, label string, validator *Validator, t reflect.Type) (field *Field) {
	t = utils.TypeElem(t)
	field = new(Field)
	field.Name = name
	field.Label = label
	field.Type = t.String()
	field.Validator = validator
	if t.Kind() == reflect.Struct && field.Type != "decimal.Decimal" {
		field.Struct = true
		for i := 0; i < t.NumField(); i++ {
			sf := t.Field(i)
			_name := getJsonField(sf)
			_label := getFieldLabel(sf)
			_validator := getFieldValidator(sf)
			field.Fields = append(field.Fields, ReflectFields(_name, _label, _validator, sf.Type))
		}
	} else if t.Kind() == reflect.Slice || t.Kind() == reflect.Array {
		field.Array = true
		field.Elem = ReflectFields(name, label, validator, t.Elem())
	}
	return
}

func getFieldLabel(field reflect.StructField) string {
	return field.Tag.Get("label")
}

func getFieldValidator(field reflect.StructField) (validator *Validator) {
	required := strings.Contains(field.Tag.Get("validator"), "required")
	if required {
		validator = newIfNoValidator(validator)
		validator.Required = true
	}
	return
}

func newIfNoValidator(validator *Validator) *Validator {
	if validator == nil {
		validator = new(Validator)
	}
	return validator
}

func getJsonField(field reflect.StructField) string {
	n := field.Tag.Get("json")
	n = strings.ReplaceAll(n, ",omitempty", "")
	if n == "" {
		n = field.Name
	}
	return n
}

type TypeConverter func(string) string

var _ TypeConverter = typescriptTypeConverter

func typescriptTypeConverter(o string) string {
	switch o {
	case "int", "int8", "int16", "int32", "int64",
		"uint", "uint8", "uint16", "uint32", "uint64",
		"float32", "float64",
		"decimal.Decimal":
		return "number"
	case "bool":
		return "boolean"
	case "string":
		return "string"
	default:
		return "any"
	}
}
