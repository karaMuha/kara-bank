package dto

type ListAccountsDto struct {
	Limit  int32 `validate:"required,min=1"`
	Offset int32 `validate:"required,gte=0"`
}
