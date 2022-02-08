package token

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/orkhanrustamli/simplebank/util"
	"github.com/stretchr/testify/require"
)

func TestJWTManager(t *testing.T) {
	manager, err := NewJWTManager(util.RandomString(32))
	require.NoError(t, err)

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

func TestExpiredToken(t *testing.T) {
	manager, err := NewJWTManager(util.RandomString(32))
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

func TestInvalidTokenAlgNone(t *testing.T) {
	payload, err := NewPayload(util.RandomName(), time.Minute)
	require.NoError(t, err)
	require.NotEmpty(t, payload)

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodNone, payload)
	token, err := jwtToken.SignedString(jwt.UnsafeAllowNoneSignatureType)
	require.NoError(t, err)
	require.NotZero(t, token)

	manager, err := NewJWTManager(util.RandomString(32))
	require.NoError(t, err)

	payload, err = manager.VerifyToken(token)
	require.Error(t, err)
	require.EqualError(t, err, ErrInvalidToken.Error())
	require.Nil(t, payload)
}
