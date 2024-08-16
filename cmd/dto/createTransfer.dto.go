package dto

type CreateTransferDto struct {
	FromAccountId int64
	ToAccountId   int64
	Amount        int64
}
