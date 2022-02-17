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
		{Method: buck.Get, Handler: f.service.Ping, Description: "测试"},
		{Method: buck.Get, Handler: f.service.GetHtml},
		{Method: buck.Get, Handler: f.service.GetText},
		{Method: buck.Get, Handler: f.service.GetInt},
		{Method: buck.Get, Handler: f.service.GetInt32},
		{Method: buck.Get, Handler: f.service.QueryPost},
		{Method: buck.Get, Handler: f.service.GetDecimal},
		{Method: buck.Get, Handler: f.service.GetBool},
		{Handler: f.service.AddPost},
		{Handler: f.service.TestGetArray},
		{Handler: f.service.TestPostArray},
		{Handler: f.service.Ping2},
		{Handler: f.service.PostShop},
	}
}
