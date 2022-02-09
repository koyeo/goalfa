package exporter

type Options struct {
	Envs []*Env `json:"envs"`
}

type Env struct {
	Name string `json:"name"`
	Host string `json:"host"`
}

type Method struct {
	Name        string `json:"name,omitempty"`
	Path        string `json:"path,omitempty"`
	Method      string `json:"method,omitempty"`
	Description string `json:"description,omitempty"`
	Middlewares string `json:"middlewares,omitempty"`
	Input       *Field `json:"input,omitempty"`
	Output      *Field `json:"output,omitempty"`
}

func (p Method) Fork() *Method {
	n := new(Method)
	n.Name = p.Name
	n.Path = p.Path
	n.Method = p.Method
	n.Description = p.Description
	n.Middlewares = p.Middlewares
	if p.Input != nil {
		n.Input = p.Input.Fork()
	}
	if p.Output != nil {
		n.Output = p.Output.Fork()
	}
	return n
}

type Field struct {
	Name        string     `json:"name,omitempty"`
	Label       string     `json:"label,omitempty"`
	Type        string     `json:"type,omitempty"`
	Description string     `json:"description,omitempty"`
	Array       bool       `json:"array,omitempty"`
	Struct      bool       `json:"struct,omitempty"`
	Origin      string     `json:"origin,omitempty"`    // 原始类型
	Fields      []*Field   `json:"fields,omitempty"`    // 描述 Struct 成员变量
	Elem        *Field     `json:"elem,omitempty"`      // 描述 Slice/Array 子元素
	Validator   *Validator `json:"validator,omitempty"` // 定义校验器
	Form        string     `json:"form,omitempty"`      // 定义表单组件
}

func (p Field) Fork() *Field {
	n := new(Field)
	n.Name = p.Name
	n.Label = p.Label
	n.Type = p.Type
	n.Description = p.Description
	n.Array = p.Array
	n.Struct = p.Struct
	n.Origin = p.Origin
	n.Validator = p.Validator
	n.Form = p.Form
	if p.Elem != nil {
		n.Elem = p.Elem.Fork()
	}
	for _, v := range p.Fields {
		n.Fields = append(n.Fields, v.Fork())
	}
	return n
}

type Validator struct {
	Required bool     `json:"required,omitempty"`
	Max      *uint64  `json:"max,omitempty"`
	Min      *int64   `json:"min,omitempty"`
	Enums    []string `json:"enums,omitempty"`
}

type Component struct {
	Name string
}
