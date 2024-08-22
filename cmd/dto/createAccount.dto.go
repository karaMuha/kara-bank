package dto

type CreateAccountDto struct {
	Owner    string `validate:"required,email"`
	Currency string `validate:"required,oneof=EUR USD"`
}
