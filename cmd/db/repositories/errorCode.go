package db

import (
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
)

const (
	UniqueViolation = "23505"
)

func ErrorCode(err error) string {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return pgErr.Code
	}

	return ""
}
