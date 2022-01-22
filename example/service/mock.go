package service

import "context"

type UserMockService struct {
}

func (u UserMockService) Register(ctx context.Context, in *RegisterIn) (out *RegisterOut, err error) {
	panic("implement me")
}

func (u UserMockService) GetProfile(ctx context.Context) (out *GetProfileOut, err error) {
	panic("implement me")
}
