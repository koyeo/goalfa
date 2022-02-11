package exporter

import (
	"encoding/json"
	"fmt"
)

type SDKFile struct {
	Name    string
	Content string
}

func NewSDK(methods []*Method) *SDK {
	return &SDK{methods: methods}
}

type SDK struct {
	methods []*Method
}

func (p SDK) Make(lang string) ([]byte, error) {
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
	files, err := maker.Make(p.methods)
	if err != nil {
		return nil, err
	}
	return json.Marshal(files)
}

func (p SDK) make() {

}
