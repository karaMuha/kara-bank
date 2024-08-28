package dto

type CreateAccountDto struct {
	Owner    string `validate:"required,email"`
	Currency string `json:"currency" validate:"required,oneof=EUR USD"`
}
