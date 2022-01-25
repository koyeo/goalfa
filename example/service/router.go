package service

import (
	"github.com/gozelle/_api"
	"net/http"
)

func NewUserRouter(service UserService) *UserRouter {
	return &UserRouter{service: service}
}

type UserRouter struct {
	service UserService
}

func (p UserRouter) Routes() []_api.Route {
	return []_api.Route{
		{
			Prefix: "/api",
			Children: []_api.Route{
				{
					Prefix: "/guard",
					Children: []_api.Route{
						{
							Path:    "/GetProfile",
							Handler: p.service.GetProfile,
						},
						{
							Path:    "/wait",
							Handler: p.service.Wait,
						},
					},
				},
				{
					Prefix: "/open",
					Children: []_api.Route{
						{
							Path:        "/Register",
							Handler:     p.service.Register,
							Description: "用户注册",
						},
						{
							Path:        "/Ping",
							Method:      http.MethodGet,
							Handler:     p.service.Ping,
							Description: "Ping",
						},
					},
				},
			},
		},
	}
}
