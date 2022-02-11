package exporter

type Maker interface {
	Make(pkg string, methods []*Method) (files []*File, err error)
}
