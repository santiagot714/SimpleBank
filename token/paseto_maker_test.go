package token

import (
	"testing"
	"time"

	"github.com/santiagot714/SimpleBank/util"
	"github.com/stretchr/testify/require"
)

func TestPasetoMaker(t *testing.T) {
	maker, err := NewPasetoMaker(util.RandomSecretKey())
	require.NoError(t, err)
	require.NotNil(t, maker)
}

func TestCreateToken(t *testing.T) {
	maker, err := NewPasetoMaker(util.RandomSecretKey())
	username := util.RandomOwner()
	require.NoError(t, err)
	require.NotNil(t, maker)

	token, payload, err := maker.CreateToken(username, time.Minute)
	require.NoError(t, err)
	require.NotEmpty(t, token)
	require.NotNil(t, payload)

	payload, err = maker.VerifyToken(token)
	require.NoError(t, err)
	require.NotNil(t, payload)
	require.Equal(t, username, payload.Username)
	require.WithinDuration(t, time.Now(), payload.IssuedAt, time.Second)
	require.WithinDuration(t, time.Now().Add(time.Minute), payload.ExpiredAt, time.Second)
}

func TestExpiredToken(t *testing.T) {
	maker, err := NewPasetoMaker(util.RandomSecretKey())
	username := util.RandomOwner()
	require.NoError(t, err)
	require.NotNil(t, maker)

	token, payload, err := maker.CreateToken(username, -time.Minute)
	require.NoError(t, err)
	require.NotEmpty(t, token)
	require.NotNil(t, payload)

	payload, err = maker.VerifyToken(token)
	require.Error(t, err)
	require.EqualError(t, err, ErrExpiredToken.Error())
	require.Nil(t, payload)
}

func TestInvalidToken(t *testing.T) {
	maker, err := NewPasetoMaker(util.RandomSecretKey())
	require.NoError(t, err)
	require.NotNil(t, maker)

	token := "invalid.token.payload"
	payload, err := maker.VerifyToken(token)
	require.Error(t, err)
	require.EqualError(t, err, ErrInvalidToken.Error())
	require.Nil(t, payload)
}
