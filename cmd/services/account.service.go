package services

import (
	db "kara-bank/db/repositories"
	"kara-bank/dto"
)

type AccountServiceImpl struct {
	db.Querier
}

func NewAccountService(querier db.Querier) AccountServiceInterface {
	return &AccountServiceImpl{
		Querier: querier,
	}
}

func (a *AccountServiceImpl) CreateAccount(args *dto.CreateAccountDto) (*db.CreateAccountParams, *dto.ResponseError) {
	return nil, nil
}

func (a *AccountServiceImpl) GetAccount(id string) (*db.Account, *dto.ResponseError) {
	return nil, nil
}

func (a AccountServiceImpl) ListAccounts(args *dto.ListAccountsDto) ([]*db.Account, *dto.ResponseError) {
	return nil, nil
}
