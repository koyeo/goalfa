package service

import (
	"github.com/koyeo/buck"
)

func NewUserRouter(service UserService) *UserRouter {
	return &UserRouter{service: service}
}

type UserRouter struct {
	service UserService
}

func (p UserRouter) Routes() []buck.Route {
	return []buck.Route{
		{
			Prefix: "/api",
			Children: []buck.Route{
				{
					Prefix: "",
					Children: []buck.Route{
						{Method: buck.GET, Path: "/ping", Handler: p.service.Ping},
						{Method: buck.GET, Path: "/inner-error", Handler: p.service.InnerError},
						{Method: buck.GET, Handler: p.service.ValidateError},
						{Method: buck.GET, Handler: p.service.ForbiddenError},
					},
				},
				{
					Prefix: "/guard",
					Children: []buck.Route{
						{Handler: p.service.TestPost},
						{Handler: p.service.TestPostArray},
						{Method: buck.GET, Handler: p.service.TestGet},
						{Method: buck.GET, Handler: p.service.TestGetArray},
						{Method: buck.PUT, Handler: p.service.TestPut},
						{Method: buck.PUT, Handler: p.service.TestPutArray},
						{Method: buck.DELETE, Handler: p.service.TestDelete},
						{Method: buck.DELETE, Handler: p.service.TestDeleteArray},
					},
				},
			},
		},
	}
}
