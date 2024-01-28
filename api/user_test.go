package api

import (
	"testing"

	db "github.com/dpomian/gobind/db/sqlc"
	"github.com/dpomian/gobind/utils"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func randomUser(t *testing.T) db.User {
	password, err := utils.HashPassword(utils.RandomString(10))
	require.NoError(t, err)

	return db.User{
		ID:        uuid.New(),
		Username:  utils.RandomString(10),
		Email:     utils.RandomEmail(),
		Password:  password,
		Suspended: false,
	}
}
