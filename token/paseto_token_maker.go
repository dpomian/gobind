package token

import (
	"fmt"
	"time"

	"github.com/aead/chacha20poly1305"
	"github.com/google/uuid"
	"github.com/o1egl/paseto"
)

type PasetoTokenMaker struct {
	paseto       *paseto.V2
	symmetricKey []byte
}

func NewPasetoTokenMaker(symmetricKey string) (TokenMaker, error) {
	if len(symmetricKey) != chacha20poly1305.KeySize {
		return nil, fmt.Errorf("invaid key size: must have exactly %d characters", chacha20poly1305.KeySize)
	}

	pasetoTokenMaker := &PasetoTokenMaker{
		paseto:       paseto.NewV2(),
		symmetricKey: []byte(symmetricKey),
	}

	return pasetoTokenMaker, nil
}

func (tokenMaker *PasetoTokenMaker) CreateToken(userId string, duration time.Duration) (string, error) {
	userIdUUID, err := uuid.Parse(userId)
	if err != nil {
		return "", fmt.Errorf("invalid userId")
	}

	payload, err := NewPayload(userIdUUID, duration)
	if err != nil {
		return "", err
	}

	return tokenMaker.paseto.Encrypt(tokenMaker.symmetricKey, payload, nil)
}

func (tokenMaker *PasetoTokenMaker) VerifyToken(token string) (*Payload, error) {
	payload := &Payload{}
	err := tokenMaker.paseto.Decrypt(token, tokenMaker.symmetricKey, payload, nil)
	if err != nil {
		return nil, ErrorInvalidToken
	}

	err = payload.Validate()
	if err != nil {
		return nil, err
	}

	return payload, nil
}
