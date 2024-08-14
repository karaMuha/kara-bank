package services

import (
	db "kara-bank/db/repositories"
	"kara-bank/dto"
)

type AccountServiceInterface interface {
	CreateAccount(args *dto.CreateAccountDto) (*db.CreateAccountParams, *dto.ResponseError)

	GetAccount(id string) (*db.Account, *dto.ResponseError)

	ListAccounts(args *dto.ListAccountsDto) ([]*db.Account, *dto.ResponseError)
}
