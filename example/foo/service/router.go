package service

import (
	"github.com/koyeo/buck"
)

func NewFooRouter(service FooService) *FooRouter {
	return &FooRouter{service: service}
}

type FooRouter struct {
	service FooService
}

func (f FooRouter) Routes() []buck.Route {
	return []buck.Route{
		{Method: buck.Get, Handler: f.service.Ping},
		{Method: buck.Get, Handler: f.service.QueryPost},
		{Handler: f.service.AddPost},
	}
}
