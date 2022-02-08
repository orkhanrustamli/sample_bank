package token

import (
	"testing"
	"time"

	"github.com/orkhanrustamli/simplebank/util"
	"github.com/stretchr/testify/require"
)

func TestPasetoManager(t *testing.T) {
	manager, err := NewPasetoManager(util.RandomString(32))
	require.NoError(t, err)
	require.NotEmpty(t, manager)

	username := util.RandomName()
	duration := time.Minute

	issuedAt := time.Now()
	expiredAt := time.Now().Add(duration)

	token, err := manager.CreateToken(username, duration)
	require.NoError(t, err)
	require.NotZero(t, token)

	payload, err := manager.VerifyToken(token)
	require.NoError(t, err)
	require.NotEmpty(t, payload)

	require.NotZero(t, payload.ID)
	require.Equal(t, payload.Username, username)
	require.WithinDuration(t, payload.IssuedAt, issuedAt, time.Second)
	require.WithinDuration(t, payload.ExpiredAt, expiredAt, time.Second)
}

func TestExpiredPasetoToken(t *testing.T) {
	manager, err := NewPasetoManager(util.RandomString(32))
	require.NoError(t, err)

	username := util.RandomName()

	token, err := manager.CreateToken(username, -time.Minute)
	require.NoError(t, err)
	require.NotZero(t, token)

	payload, err := manager.VerifyToken(token)
	require.Error(t, err)
	require.EqualError(t, err, ErrExpiredToken.Error())
	require.Nil(t, payload)
}
