package db

import (
	"context"
	"fmt"
	"testing"

	"github.com/dpomian/gobind/utils"
	"github.com/stretchr/testify/require"
)

func createRandomUser(t *testing.T) User {
	password, err := utils.HashPassword(utils.RandomString(6))
	require.NoError(t, err)

	arg := CreateUserParams{
		ID:       testData.userId1,
		Username: "Username1",
		Email:    "email1@test.com",
		Password: password,
	}

	fmt.Println(arg)

	user, err := testQueries.CreateUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, user)

	require.Equal(t, arg.ID, user.ID)
	require.Equal(t, arg.Username, user.Username)
	require.Equal(t, arg.Email, user.Email)
	require.False(t, user.Suspended)
	require.NotZero(t, user.CreatedAt)
	require.True(t, user.PasswordChangedAt.IsZero())

	return user
}

func TestGetUser(t *testing.T) {
	user, err := testQueries.GetUserById(context.Background(), testData.userId1)
	require.NoError(t, err)
	require.NotEmpty(t, user)
}
