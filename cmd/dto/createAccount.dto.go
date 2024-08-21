package dto

type CreateAccountDto struct {
	Owner    string
	Currency string `validate:"required,oneof=EUR USD"`
}
