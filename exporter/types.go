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
	Middlewares string   `json:"middlewares"`
	Input       []*Field `json:"input"`
	Output      []*Field `json:"output"`
}

type Field struct {
	Name        string   `json:"name"`
	Label       string   `json:"label"`
	Type        string   `json:"type"`
	Description string   `json:"description"`
	Array       bool     `json:"array"`
	Required    bool     `json:"required"`
	Max         int64    `json:"max"`
	Min         int64    `json:"min"`
	Form        string   `json:"form"` // 表征表单的空间
	Elem        []*Field `json:"elem"`
	Children    []*Field `json:"children"`
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
