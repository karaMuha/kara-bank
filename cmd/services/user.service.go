package services

import (
	"context"
	"errors"
	db "kara-bank/db/repositories"
	"kara-bank/dto"
	"kara-bank/utils"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5"
	"golang.org/x/crypto/bcrypt"
)

type UserServiceImpl struct {
	store      db.Store
	tokenMaker utils.TokenMaker
}

func NewUserService(store db.Store, tokenMaker utils.TokenMaker) *UserServiceImpl {
	return &UserServiceImpl{
		store:      store,
		tokenMaker: tokenMaker,
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
		UserRole:       utils.CustomerRole,
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

	accessToken, _, err := u.tokenMaker.CreateToken(user.Email, user.UserRole, 30*time.Minute)

	if err != nil {
		return "", &dto.ResponseError{
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
	}

	refreshToken, refreshTokenPayload, err := u.tokenMaker.CreateToken(user.Email, user.UserRole, 168*time.Hour) // valid for a week

	if err != nil {
		return "", &dto.ResponseError{
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
	}

	sessionParams := &db.CreateSessionParams{
		ID:           refreshTokenPayload.ID,
		Email:        refreshTokenPayload.Email,
		RefreshToken: refreshToken,
		UserAgent:    arg.UserAgent,
		ClientIp:     arg.ClientIp,
		IsBlocked:    false,
		ExpiresAt:    refreshTokenPayload.ExpiredAt,
	}

	_, err = u.store.CreateSession(ctx, sessionParams)

	if err != nil {
		return "", &dto.ResponseError{
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
	}

	return accessToken, nil
}

var _ UserServiceInterface = (*UserServiceImpl)(nil)
