package dto

type CreateTransferDto struct {
	FromUser      string `validate:"required,email"`
	FromAccountId int64  `validate:"required,min=1"`
	ToAccountId   int64  `validate:"required,min=1"`
	Amount        int64  `validate:"required,gt=0"`
}
