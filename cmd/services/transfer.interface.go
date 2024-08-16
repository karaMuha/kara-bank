package services

import (
	"context"
	db "kara-bank/db/repositories"
	"kara-bank/dto"
)

type TransferServiceInterface interface {
	CreateTransfer(ctx context.Context, arg *dto.CreateTransferDto) (*db.TransferTxResult, *dto.ResponseError)
}
