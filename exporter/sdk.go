package exporter

import (
	"encoding/json"
	"fmt"
)

type Folder struct {
	Name      string  `json:"name"`
	Namespace string  `json:"namespace"`
	Files     []*File `json:"files"`
}

type File struct {
	Name    string `json:"name"`
	Content string `json:"content"`
}

func NewSDK(methods []*Method) *SDK {
	return &SDK{methods: methods}
}

type SDK struct {
	methods []*Method
}

func (p SDK) Make(lang, pkg string) ([]byte, error) {
	var maker Maker
	switch lang {
	case "go":
		maker = new(GoMaker)
	case "axios":
		maker = new(AxiosMaker)
	case "angular":
		maker = new(AngularMaker)
	default:
		return nil, fmt.Errorf("unsupport sdk lang: '%s'", lang)
	}
	files, err := maker.Make(pkg, p.methods)
	if err != nil {
		return nil, err
	}
	return json.Marshal(files)
}

func (p SDK) make() {

}
