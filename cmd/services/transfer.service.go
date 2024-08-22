package services

import (
	"context"
	"errors"
	db "kara-bank/db/repositories"
	"kara-bank/dto"
	"net/http"

	"github.com/jackc/pgx/v5"
)

type TransferServiceImpl struct {
	store db.Store
}

func NewTransferService(store db.Store) TransferServiceInterface {
	return &TransferServiceImpl{
		store: store,
	}
}

func (t *TransferServiceImpl) CreateTransfer(ctx context.Context, arg *dto.CreateTransferDto) (*db.TransferTxResult, *dto.ResponseError) {
	respErr := t.validAccounts(ctx, arg.FromUser, arg.FromAccountId, arg.ToAccountId)

	if respErr != nil {
		return nil, respErr
	}

	queryParam := db.TransferTxParams{
		FromAccountID: arg.FromAccountId,
		ToAccountID:   arg.ToAccountId,
		Amount:        arg.Amount,
	}

	transfer, err := t.store.TransferTx(ctx, queryParam)

	if err != nil {
		return nil, &dto.ResponseError{
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
	}

	return &transfer, nil
}

func (t *TransferServiceImpl) validAccounts(ctx context.Context, fromUser string, fromAccountId int64, toAccountId int64) *dto.ResponseError {
	fromAccount, err := t.store.GetAccount(ctx, fromAccountId)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return &dto.ResponseError{
				Message: "fromAccount not found",
				Status:  http.StatusNotFound,
			}
		}

		return &dto.ResponseError{
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
	}

	if fromAccount.Owner != fromUser {
		return &dto.ResponseError{
			Message: "You cannot send money from accoutns other than yours",
			Status:  http.StatusUnauthorized,
		}
	}

	_, err = t.store.GetAccount(ctx, toAccountId)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return &dto.ResponseError{
				Message: "toAccount not found",
				Status:  http.StatusNotFound,
			}
		}

		return &dto.ResponseError{
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
	}

	return nil
}
