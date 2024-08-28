package dto

type LoginUserDto struct {
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password" validate:"required"`
	UserAgent string `validate:"required"`
	ClientIp  string `validate:"required"`
}
