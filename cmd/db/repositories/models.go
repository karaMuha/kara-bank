// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0

package db

import (
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

type Account struct {
	ID        int64     `json:"id"`
	Owner     string    `json:"owner"`
	Balance   int64     `json:"balance"`
	Currency  string    `json:"currency"`
	CreatedAt time.Time `json:"created_at"`
}

type Entry struct {
	ID        int64       `json:"id"`
	AccountID pgtype.Int8 `json:"account_id"`
	// can be negative or positive
	Amount    int64     `json:"amount"`
	CreatedAt time.Time `json:"created_at"`
}

type Transfer struct {
	ID            int64       `json:"id"`
	FromAccountID pgtype.Int8 `json:"from_account_id"`
	ToAccountID   pgtype.Int8 `json:"to_account_id"`
	// must be positive
	Amount    int64     `json:"amount"`
	CreatedAt time.Time `json:"created_at"`
}
