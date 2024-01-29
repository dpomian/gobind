package token

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

const minSecretKeySize = 32

type JWTTokenMaker struct {
	secretKey string
}

func NewJWTTokenMaker(secretKey string) (TokenMaker, error) {
	if len(secretKey) < minSecretKeySize {
		return nil, fmt.Errorf("invalid key size: must be at least %d characters long", minSecretKeySize)
	}

	return &JWTTokenMaker{secretKey: secretKey}, nil
}

func (tokenMaker *JWTTokenMaker) CreateToken(userId string, duration time.Duration) (string, *Payload, error) {
	userIdUUID, err := uuid.Parse(userId)
	if err != nil {
		return "", nil, fmt.Errorf("invalid userId")
	}

	payload, err := NewPayload(userIdUUID, duration)
	if err != nil {
		return "", nil, nil
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	token, err := jwtToken.SignedString([]byte(tokenMaker.secretKey))
	return token, payload, err
}

func (tokenMaker *JWTTokenMaker) VerifyToken(token string) (*Payload, error) {
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, ErrorInvalidToken
		}
		return []byte(tokenMaker.secretKey), nil
	}

	jwtToken, err := jwt.ParseWithClaims(token, &Payload{}, keyFunc)
	if err != nil {
		return nil, err
	}

	payload, ok := jwtToken.Claims.(*Payload)
	if !ok {
		return nil, ErrorInvalidToken
	}

	return payload, nil
}
