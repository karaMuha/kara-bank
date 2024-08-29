package services

import (
	"context"
	"errors"
	db "kara-bank/db/repositories"
	"kara-bank/dto"
	"kara-bank/utils"
	"net/http"

	"github.com/jackc/pgx/v5"
)

type AccountServiceImpl struct {
	store db.Store
}

func NewAccountService(store db.Store) *AccountServiceImpl {
	return &AccountServiceImpl{
		store: store,
	}
}

func (a *AccountServiceImpl) CreateAccount(ctx context.Context, args *dto.CreateAccountDto) (*db.Account, *dto.ResponseError) {
	createAccountParams := &db.CreateAccountParams{
		Owner:    args.Owner,
		Currency: args.Currency,
		Balance:  0,
	}

	createdAccount, err := a.store.CreateAccount(ctx, createAccountParams)

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

	return createdAccount, nil
}

func (a *AccountServiceImpl) GetAccount(ctx context.Context, id int64, email string, role string) (*db.Account, *dto.ResponseError) {
	account, err := a.store.GetAccount(ctx, id)

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

	if account.Owner == email || role == utils.BankerRole || role == utils.AdminRole {
		return account, nil
	}

	return nil, &dto.ResponseError{
		Message: "You have no permission for this account",
		Status:  http.StatusUnauthorized,
	}
}

func (a AccountServiceImpl) ListAccounts(ctx context.Context, arg *dto.ListAccountsDto, role string) ([]*db.Account, *dto.ResponseError) {
	if role != utils.AdminRole && role != utils.BankerRole {
		return nil, &dto.ResponseError{
			Message: "You have no permission for this action",
			Status:  http.StatusUnauthorized,
		}
	}

	params := &db.ListAccountsParams{
		Limit:  arg.Limit,
		Offset: arg.Offset,
	}

	accountList, err := a.store.ListAccounts(ctx, params)

	if err != nil {
		return nil, &dto.ResponseError{
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
	}

	return accountList, nil
}

var _ AccountServiceInterface = (*AccountServiceImpl)(nil)
