package db

import (
	"context"
	"database/sql"
	"github.com/stretchr/testify/require"
	"simplebank/util"
	"testing"
	"time"
)

func createRandomUser(t *testing.T) User {
	arg := CreateUserParams{
		Username:       util.RandomOwner(),
		HashedPassword: "secret",
		FullName:       util.RandomOwner(),
		Email:          util.RandomEmail(),
	}

	user, err := testQueries.CreateUser(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, user)

	require.Equal(t, arg.Username, user.Username)
	require.Equal(t, arg.HashedPassword, user.HashedPassword)
	require.Equal(t, arg.FullName, user.FullName)
	require.Equal(t, arg.Email, user.Email)

	require.True(t, user.PasswordChangedAt.IsZero())
	require.NotZero(t, user.CreatedAt)

	return user
}

func TestCreateUser(t *testing.T) {
	createRandomUser(t)
}

func TestGetUser(t *testing.T) {
	// Create account
	user1 := createRandomUser(t)
	user2, err := testQueries.GetUser(context.Background(), user1.Username)
	require.NoError(t, err)
	require.NotEmpty(t, user2)

	require.Equal(t, user1.Username, user2.Username)
	require.Equal(t, user1.Email, user2.Email)
	require.Equal(t, user1.FullName, user2.FullName)

	require.WithinDuration(t, user1.CreatedAt, user2.CreatedAt, time.Second)
	require.WithinDuration(t, user1.PasswordChangedAt, user2.PasswordChangedAt, time.Second)
}

func TestUpdateUser(t *testing.T) {
	oldUser := createRandomUser(t)

	// Update the created user with a new full name
	newFullname := util.RandomOwner()
	updateUser, err := testQueries.UpdateUser(context.Background(), UpdateUserParams{
		FullName: sql.NullString{
			String: newFullname,
			Valid:  true,
		},
		Username: oldUser.Username,
	})
	require.NoError(t, err)
	require.NotEqual(t, oldUser.FullName, updateUser.FullName)
	require.Equal(t, updateUser.FullName, newFullname)
	require.Equal(t, updateUser.Email, oldUser.Email)
	require.Equal(t, updateUser.HashedPassword, oldUser.HashedPassword)
}

func TestUpdateEmail(t *testing.T) {
	oldUser := createRandomUser(t)

	// Update the created user with a new full name
	newEmail := util.RandomEmail()
	updateUser, err := testQueries.UpdateUser(context.Background(), UpdateUserParams{
		Email: sql.NullString{
			String: newEmail,
			Valid:  true,
		},
		Username: oldUser.Username,
	})
	require.NoError(t, err)
	require.NotEqual(t, oldUser.Email, updateUser.Email)
	require.Equal(t, updateUser.Email, newEmail)
	require.Equal(t, updateUser.Username, oldUser.Username)
	require.Equal(t, updateUser.HashedPassword, oldUser.HashedPassword)
}

func TestUpdateHashedPassword(t *testing.T) {
	oldUser := createRandomUser(t)

	// Update the created user with a new full name
	password := util.RandomString(6)
	newHashedPassword, err := util.HashPassword(password)
	require.NoError(t, err)
	updateUser, err := testQueries.UpdateUser(context.Background(), UpdateUserParams{
		HashedPassword: sql.NullString{
			String: newHashedPassword,
			Valid:  true,
		},
		Username: oldUser.Username,
	})
	require.NoError(t, err)
	require.NotEqual(t, oldUser.HashedPassword, updateUser.HashedPassword)
	require.Equal(t, updateUser.HashedPassword, newHashedPassword)
	require.Equal(t, updateUser.Email, oldUser.Email)
	require.Equal(t, updateUser.Username, oldUser.Username)
}

func TestUpdateAllFields(t *testing.T) {
	oldUser := createRandomUser(t)

	// Update the created user with a new full name
	newFullName := util.RandomOwner()
	newEmail := util.RandomEmail()
	password := util.RandomString(6)
	newHashedPassword, err := util.HashPassword(password)
	require.NoError(t, err)
	updateUser, err := testQueries.UpdateUser(context.Background(), UpdateUserParams{
		HashedPassword: sql.NullString{
			String: newHashedPassword,
			Valid:  true,
		},
		Email: sql.NullString{
			String: newEmail,
			Valid:  true,
		},
		FullName: sql.NullString{
			String: newFullName,
			Valid:  true,
		},
		Username: oldUser.Username,
	})
	require.NoError(t, err)
	require.NotEqual(t, oldUser.HashedPassword, updateUser.HashedPassword)
	require.NotEqual(t, oldUser.HashedPassword, updateUser.HashedPassword)
	require.NotEqual(t, updateUser.Email, oldUser.Email)
	require.Equal(t, updateUser.Username, oldUser.Username)
	require.Equal(t, updateUser.Email, newEmail)
	require.Equal(t, updateUser.HashedPassword, newHashedPassword)
	require.Equal(t, updateUser.FullName, newFullName)
}
