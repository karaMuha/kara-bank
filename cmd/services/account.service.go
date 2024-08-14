package services

import (
	"context"
	db "kara-bank/db/repositories"
	"kara-bank/dto"
	"net/http"
)

type AccountServiceImpl struct {
	db.Querier
}

func NewAccountService(querier db.Querier) AccountServiceInterface {
	return &AccountServiceImpl{
		Querier: querier,
	}
}

func (a *AccountServiceImpl) CreateAccount(ctx context.Context, args *dto.CreateAccountDto) (*db.Account, *dto.ResponseError) {
	createAccountParams := &db.CreateAccountParams{
		Owner:    args.Owner,
		Currency: args.Currency,
		Balance:  0,
	}

	createdAccount, err := a.Querier.CreateAccount(ctx, createAccountParams)

	if err != nil {
		return nil, &dto.ResponseError{
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
	}

	return createdAccount, nil
}

func (a *AccountServiceImpl) GetAccount(ctx context.Context, id string) (*db.Account, *dto.ResponseError) {
	return nil, nil
}

func (a AccountServiceImpl) ListAccounts(ctx context.Context, args *dto.ListAccountsDto) ([]*db.Account, *dto.ResponseError) {
	return nil, nil
}
