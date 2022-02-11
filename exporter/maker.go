package exporter

type Maker interface {
	Make(methods []*Method) ([]*SDKFile, error)
}
