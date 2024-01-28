package token

import (
	"testing"
	"time"

	"github.com/dpomian/gobind/utils"
	"github.com/stretchr/testify/require"
)

func TestPasetoTokenMaker(t *testing.T) {
	tokenMaker, err := NewPasetoTokenMaker(utils.RandomString(32))
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

func TestExpiredPasetoToken(t *testing.T) {
	tokenMaker, err := NewPasetoTokenMaker(utils.RandomString(32))
	require.NoError(t, err)

	expiredToken, err := tokenMaker.CreateToken("user1@email.com", -1*time.Minute)
	require.NoError(t, err)
	require.NotEmpty(t, expiredToken)

	payload, err := tokenMaker.VerifyToken(expiredToken)
	require.Error(t, err)
	require.EqualError(t, err, ErrorExpiredToken.Error())
	require.Nil(t, payload)
}
