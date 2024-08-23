package dto

type LoginUserDto struct {
	Email     string `validate:"required,email"`
	Password  string `validate:"required"`
	UserAgent string `validate:"required"`
	ClientIp  string `validate:"required"`
}
