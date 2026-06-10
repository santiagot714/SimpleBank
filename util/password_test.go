package util

import (
	"testing"

	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

// TestPassword tests the HashPassword and CheckPassword functions.
func TestPassword(t *testing.T) {
	password := RandomString(10)

	hashedPassword1, err := HashPassword(password)
	require.NoError(t, err)
	require.NotEmpty(t, hashedPassword1)

	err = CheckPassword(password, hashedPassword1)
	require.NoError(t, err)

	wrongPassword := RandomString(10)
	err = CheckPassword(wrongPassword, hashedPassword1)
	require.Error(t, err)
	require.EqualError(t, err, bcrypt.ErrMismatchedHashAndPassword.Error())

	hashedPassword2, err := HashPassword(password)
	require.NoError(t, err)
	require.NotEmpty(t, hashedPassword2)
	// ensure that the hashed passwords are different
	require.NotEqual(t, hashedPassword1, hashedPassword2)
}
