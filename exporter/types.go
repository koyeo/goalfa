package exporter

type Options struct {
	Envs []*Env `json:"envs"`
}

type Env struct {
	Name string `json:"name"`
	Host string `json:"host"`
}

type Method struct {
	Name        string   `json:"name"`
	Path        string   `json:"path"`
	Method      string   `json:"method"`
	Description string   `json:"description"`
	Middlewares string   `json:"middlewares"`
	Input       []*Field `json:"input"`
	Output      []*Field `json:"output"`
}

type Field struct {
	Name        string     `json:"name,omitempty"`
	Label       string     `json:"label,omitempty"`
	Type        string     `json:"type,omitempty"`
	Description string     `json:"description,omitempty"`
	Array       bool       `json:"array"`
	Required    bool       `json:"required"`
	Max         int64      `json:"max,omitempty"`
	Min         int64      `json:"min,omitempty"`
	Form        string     `json:"form,omitempty"` // 定义表单组件
	Elem        []*Field   `json:"elem,omitempty"`
	Children    []*Field   `json:"children,omitempty"`
	Validator   *Validator `json:"validator,omitempty"`
}

type Validator struct {
	Required bool
	Max      uint64
	Min      int64
	Enums    []string
}

type Component struct {
	Name string
}
