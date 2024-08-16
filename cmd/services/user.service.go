package services

import (
	"context"
	db "kara-bank/db/repositories"
	"kara-bank/dto"
)

type UserServiceImpl struct {
	store db.Store
}

func NewUserService(store db.Store) UserServiceInterface {
	return &UserServiceImpl{
		store: store,
	}
}

func (u *UserServiceImpl) RegisterUser(ctx context.Context, arg *dto.RegisterUserDto) (*db.User, *dto.ResponseError) {
	return nil, nil
}

func (u *UserServiceImpl) GetUser(ctx context.Context, email string) (*db.User, *dto.ResponseError) {
	return nil, nil
}

func (u *UserServiceImpl) LoginUser(ctx context.Context, arg *dto.LoginUserDto) (string, *dto.ResponseError) {
	return "", nil
}
