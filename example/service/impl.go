package service

import (
	"context"
	"time"
)

type UserImplService struct {
}

func (u UserImplService) Ping(ctx context.Context) (err error) {
	return
}

func (u UserImplService) Wait(ctx context.Context, in *GetProfileInput) (out *GetProfileOut, err error) {
	time.Sleep(3 * time.Second)
	out = &GetProfileOut{}
	return
}

func (u UserImplService) Register(ctx context.Context, in *RegisterIn) (out *RegisterOut, err error) {
	out = new(RegisterOut)
	out.Username = in.Username
	out.Password = in.Password
	return
}

func (u UserImplService) GetProfile(ctx context.Context, in *GetProfileInput) (out *GetProfileOut, err error) {
	out = new(GetProfileOut)
	out.Text = "Hello world!"
	return
}
