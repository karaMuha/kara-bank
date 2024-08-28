package db

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func registerTestUser(t *testing.T, arg *RegisterUserParams) *User {
	user, err := testStore.RegisterUser(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, user)
	require.NotEmpty(t, user.CreatedAt)
	require.Equal(t, arg.Email, user.Email)

	return user
}
