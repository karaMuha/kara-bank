package util

import (
	"strings"
	"time"

	"aidanwoods.dev/go-paseto"
	"github.com/google/uuid"
)

type PasetoMaker struct {
	symmeticKey paseto.V4SymmetricKey
	implicit    []byte
}

func NewPasetoMaker(symmetricKey string) TokenMaker {
	return &PasetoMaker{
		symmeticKey: paseto.NewV4SymmetricKey(),
		implicit:    []byte(symmetricKey),
	}
}

func (p *PasetoMaker) CreateToken(email string, duration time.Duration) (string, error) {
	token := paseto.NewToken()
	tokenId, err := uuid.NewRandom()

	if err != nil {
		return "", err
	}

	token.Set("id", tokenId.String())
	token.Set("email", email)
	token.SetIssuedAt(time.Now())
	token.SetExpiration(time.Now().Add(duration))

	return token.V4Encrypt(p.symmeticKey, p.implicit), nil
}

func (p *PasetoMaker) VerifyToken(token string) (*TokenPayload, error) {
	parser := paseto.NewParser()
	parser.AddRule(paseto.NotExpired())
	pasrsedToken, err := parser.ParseV4Local(p.symmeticKey, token, p.implicit)

	if err != nil {
		if strings.Contains(err.Error(), "expired") {
			return nil, ErrExpiredToken
		}
		return nil, ErrInvalidToken
	}

	payload, err := getPayloadFromToken(pasrsedToken)

	if err != nil {
		return nil, ErrInvalidToken
	}

	return payload, nil
}

func getPayloadFromToken(token *paseto.Token) (*TokenPayload, error) {
	id, err := token.GetString("id")

	if err != nil {
		return nil, ErrInvalidToken
	}

	email, err := token.GetString("email")

	if err != nil {
		return nil, ErrInvalidToken
	}

	issuedAt, err := token.GetIssuedAt()

	if err != nil {
		return nil, ErrInvalidToken
	}

	expiredAt, err := token.GetExpiration()

	if err != nil {
		return nil, ErrInvalidToken
	}

	return &TokenPayload{
		ID:        uuid.MustParse(id),
		Email:     email,
		IssuedAt:  issuedAt,
		ExpiredAt: expiredAt,
	}, nil
}
