package services

import (
	"context"
	db "kara-bank/db/repositories"
	"kara-bank/dto"
	"net/http"

	"golang.org/x/crypto/bcrypt"
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
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(arg.Password), bcrypt.DefaultCost)

	if err != nil {
		return nil, &dto.ResponseError{
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
	}

	params := &db.RegisterUserParams{
		Email:          arg.Email,
		HashedPassword: string(hashedPassword),
		FirstName:      arg.FirstName,
		LastName:       arg.LastName,
	}

	user, err := u.store.RegisterUser(ctx, params)

	if err != nil {
		if db.ErrorCode(err) == db.UniqueViolation {
			return nil, &dto.ResponseError{
				Message: err.Error(),
				Status:  http.StatusConflict,
			}
		}

		return nil, &dto.ResponseError{
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
	}

	return user, nil
}

func (u *UserServiceImpl) GetUser(ctx context.Context, email string) (*db.User, *dto.ResponseError) {
	return nil, nil
}

func (u *UserServiceImpl) LoginUser(ctx context.Context, arg *dto.LoginUserDto) (string, *dto.ResponseError) {
	return "", nil
}
