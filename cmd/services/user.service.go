package services

import (
	"context"
	"errors"
	db "kara-bank/db/repositories"
	"kara-bank/dto"
	"net/http"

	"github.com/jackc/pgx/v5"
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
	user, err := u.store.GetUser(ctx, arg.Email)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", &dto.ResponseError{
				Message: err.Error(),
				Status:  http.StatusNotFound,
			}
		}

		return "", &dto.ResponseError{
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(arg.Password))

	if err != nil {
		return "", &dto.ResponseError{
			Message: err.Error(),
			Status:  http.StatusUnauthorized,
		}
	}

	// TODO: generate and return jwt

	return "token", nil
}