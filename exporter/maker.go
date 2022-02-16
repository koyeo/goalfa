package exporter

type Maker interface {
	Lang() string
	Make(methods []*Method) (files []*File, err error)
}
