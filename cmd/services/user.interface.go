package services

import (
	"context"
	db "kara-bank/db/repositories"
	"kara-bank/dto"
)

type UserServiceInterface interface {
	RegisterUser(ctx context.Context, arg *dto.RegisterUserDto) (*db.User, *dto.ResponseError)

	GetUser(ctx context.Context, email string) (*db.User, *dto.ResponseError)

	LoginUser(ctx context.Context, arg *dto.LoginUserDto) (string, *dto.ResponseError)
}
