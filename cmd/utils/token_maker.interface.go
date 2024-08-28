package utils

import "time"

type TokenMaker interface {
	CreateToken(email string, role string, duration time.Duration) (string, *TokenPayload, error)

	VerifyToken(token string) (*TokenPayload, error)
}
