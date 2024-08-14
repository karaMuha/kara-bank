package services

import (
	"context"
	db "kara-bank/db/repositories"
	"kara-bank/dto"
)

type AccountServiceInterface interface {
	CreateAccount(ctx context.Context, args *dto.CreateAccountDto) (*db.Account, *dto.ResponseError)

	GetAccount(ctx context.Context, id string) (*db.Account, *dto.ResponseError)

	ListAccounts(ctx context.Context, args *dto.ListAccountsDto) ([]*db.Account, *dto.ResponseError)
}
