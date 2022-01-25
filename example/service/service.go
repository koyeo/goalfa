package service

import (
	"context"
)

type UserService interface {
	Register(ctx context.Context, in *RegisterIn) (out *RegisterOut, err error)
	GetProfile(ctx context.Context, in *GetProfileInput) (out *GetProfileOut, err error)
	Wait(ctx context.Context, in *GetProfileInput) (out *GetProfileOut, err error)
	Ping(ctx context.Context) (err error)
}

type RegisterIn struct {
	Username string  `json:"username" validator:"required" label:"用户名" description:"" form:"date"`
	Password string  `json:"password" validator:"required" label:"密码"`
	Age      int     `json:"age" validator:"required,max(10),min(1)" label:"年龄"`
	Price    float64 `json:"price"`
}

type RegisterOut struct {
	Username string `json:"username,omitempty"`
	Password string `json:"password"`
	Status   int    `json:"status"`
}

type GetProfileInput struct {
	Channel string `json:"channel"`
}

type GetProfileOut struct {
	Text string `json:"text"`
}
