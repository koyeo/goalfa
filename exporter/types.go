package exporter

type Options struct {
	Envs []*Env
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
	Name     string   `json:"name"`
	Type     string   `json:"type"`
	Array    bool     `json:"array"`
	Required bool     `json:"required"`
	Elem     []*Field `json:"elem"`
	Children []*Field `json:"children"`
}
