package services

import (
	"context"
	"errors"
	db "kara-bank/db/repositories"
	"kara-bank/dto"
	"net/http"

	"github.com/jackc/pgx/v5"
)

type AccountServiceImpl struct {
	db.Querier
}

func NewAccountService(querier db.Querier) AccountServiceInterface {
	return &AccountServiceImpl{
		Querier: querier,
	}
}

func (a *AccountServiceImpl) CreateAccount(ctx context.Context, arg *dto.CreateAccountDto) (*db.Account, *dto.ResponseError) {
	params := &db.CreateAccountParams{
		Owner:    arg.Owner,
		Currency: arg.Currency,
		Balance:  0,
	}

	createdAccount, err := a.Querier.CreateAccount(ctx, params)

	if err != nil {
		return nil, &dto.ResponseError{
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
	}

	return createdAccount, nil
}

func (a *AccountServiceImpl) GetAccount(ctx context.Context, id int64) (*db.Account, *dto.ResponseError) {
	account, err := a.Querier.GetAccount(ctx, id)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, &dto.ResponseError{
				Message: err.Error(),
				Status:  http.StatusNotFound,
			}
		}
		return nil, &dto.ResponseError{
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
	}

	return account, nil
}

func (a AccountServiceImpl) ListAccounts(ctx context.Context, arg *dto.ListAccountsDto) ([]*db.Account, *dto.ResponseError) {
	params := &db.ListAccountsParams{
		Limit:  arg.Limit,
		Offset: arg.Offset,
	}

	accountList, err := a.Querier.ListAccounts(ctx, params)

	if err != nil {
		return nil, &dto.ResponseError{
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
	}

	return accountList, nil
}
