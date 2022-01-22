package service

import (
	"context"
)

type UserImplService struct {
}

func (u UserImplService) Register(ctx context.Context, in *RegisterIn) (out *RegisterOut, err error) {
	out = new(RegisterOut)
	out.Username = in.Username
	out.Password = in.Password
	return
}

func (u UserImplService) GetProfile(ctx context.Context) (out *GetProfileOut, err error) {
	out = new(GetProfileOut)
	out.Text = "Hello world!"
	return
}
