package token

import (
	"fmt"
	"testing"
	"time"

	"github.com/dpomian/gobind/utils"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/require"
)

func TestJWTTokenMaker(t *testing.T) {
	tokenMaker, err := NewJWTTokenMaker(utils.RandomString(32))
	require.NoError(t, err)

	email := "user1@email.com"
	duration := 1 * time.Minute

	issuedAt := time.Now()
	expireAt := issuedAt.Add(duration)

	token, err := tokenMaker.CreateToken(email, duration)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	payload, err := tokenMaker.VerifyToken(token)
	require.NoError(t, err)
	require.NotEmpty(t, payload)

	require.NotZero(t, payload.ID)
	require.Equal(t, email, payload.Email)
	require.WithinDuration(t, issuedAt, payload.IssuedAt, 1*time.Second)
	require.WithinDuration(t, expireAt, payload.ExpiredAt, 1*time.Second)
}

func TestExpiredJWTToken(t *testing.T) {
	tokenMaker, err := NewJWTTokenMaker(utils.RandomString(32))
	require.NoError(t, err)

	expiredToken, err := tokenMaker.CreateToken("user1@email.com", -1*time.Minute)
	require.NoError(t, err)
	require.NotEmpty(t, expiredToken)

	expectedErr := fmt.Errorf("token has invalid claims: token has expired")
	payload, err := tokenMaker.VerifyToken(expiredToken)
	require.Error(t, err)
	require.EqualError(t, err, expectedErr.Error())
	require.Nil(t, payload)
}

func TestInvalidJWTToken(t *testing.T) {
	payload, err := NewPayload("user1@email.com", 1*time.Minute)
	require.NoError(t, err)

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodNone, payload)
	token, err := jwtToken.SignedString(jwt.UnsafeAllowNoneSignatureType)
	require.NoError(t, err)

	tokenMaker, err := NewJWTTokenMaker(utils.RandomString(32))
	require.NoError(t, err)

	expectedErr := fmt.Errorf("token is unverifiable: error while executing keyfunc: invalid token")
	payload, err = tokenMaker.VerifyToken(token)
	require.Error(t, err)
	require.EqualError(t, err, expectedErr.Error())
	require.Nil(t, payload)
}
