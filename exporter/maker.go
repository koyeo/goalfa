package exporter

type Maker interface {
	Make(lang string, exporter *Exporter, methods []*Method) (files []*File, err error)
}
