package dto

type RegisterUserDto struct {
	Email     string `validate:"required,email"`
	Password  string `validate:"required,min=8"`
	FirstName string `validate:"required,alpha"`
	LastName  string `validate:"required,alpha"`
}
