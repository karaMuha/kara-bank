package dto

type CreateTransferDto struct {
	FromUser      string `validate:"required,email"`
	FromAccountId int64  `json:"from_account_id" validate:"required,min=1"`
	ToAccountId   int64  `json:"to_account_id" validate:"required,min=1"`
	Amount        int64  `json:"amount" validate:"required,gt=0"`
}
