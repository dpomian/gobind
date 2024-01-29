package token

import (
	"testing"
	"time"

	"github.com/dpomian/gobind/utils"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestPasetoTokenMaker(t *testing.T) {
	tokenMaker, err := NewPasetoTokenMaker(utils.RandomString(32))
	require.NoError(t, err)

	userId := uuid.New()
	duration := 1 * time.Minute

	issuedAt := time.Now()
	expireAt := issuedAt.Add(duration)

	token, _, err := tokenMaker.CreateToken(userId.String(), duration)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	payload, err := tokenMaker.VerifyToken(token)
	require.NoError(t, err)
	require.NotEmpty(t, payload)

	require.NotZero(t, payload.ID)
	require.Equal(t, userId, payload.UserId)
	require.WithinDuration(t, issuedAt, payload.IssuedAt, 1*time.Second)
	require.WithinDuration(t, expireAt, payload.ExpiredAt, 1*time.Second)
}

func TestExpiredPasetoToken(t *testing.T) {
	tokenMaker, err := NewPasetoTokenMaker(utils.RandomString(32))
	require.NoError(t, err)

	expiredToken, _, err := tokenMaker.CreateToken(uuid.NewString(), -1*time.Minute)
	require.NoError(t, err)
	require.NotEmpty(t, expiredToken)

	payload, err := tokenMaker.VerifyToken(expiredToken)
	require.Error(t, err)
	require.EqualError(t, err, ErrorExpiredToken.Error())
	require.Nil(t, payload)
}
