package buck

type Status struct {
	Detail interface{}
}

func (p Status) Error() string {
	panic("implement me")
}
