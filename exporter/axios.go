package exporter

type AxiosMaker struct {
}

func (a AxiosMaker) Lang() string {
	return Ts
}

func (a AxiosMaker) Make( methods []*Method) (files []*File, err error) {
	return
}
