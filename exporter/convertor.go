package exporter

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
		return o
	}
}
