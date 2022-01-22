package service

import (
	"context"
)

type UserService interface {
	Register(ctx context.Context, in *RegisterIn) (out *RegisterOut, err error)
	GetProfile(ctx context.Context) (out *GetProfileOut, err error)
}

type RegisterIn struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type RegisterOut struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Status   int    `json:"status"`
}

type GetProfileOut struct {
	Text string `json:"text"`
}
