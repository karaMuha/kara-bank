package utils

import "time"

type TokenMaker interface {
	CreateToken(email string, duration time.Duration) (string, error)

	VerifyToken(token string) (*TokenPayload, error)
}
