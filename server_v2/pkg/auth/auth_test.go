package auth

import (
	"crypto/rand"
	"testing"
	"time"

	"github.com/o1egl/paseto"
	"github.com/stretchr/testify/require"
)

func TestMain(m *testing.M) {
	m.Run()
}

const (
	username = "username"
	invalid  = "invalid"
)

func TestPasetoAuthenticator(t *testing.T) {
	paseto := paseto.NewV2()

	key := make([]byte, 32)
	_, err := rand.Read(key)
	require.NoError(t, err)

	authenticator, err := NewPasetoAuthenticator(key, paseto)
	require.NoError(t, err)

	token, err := authenticator.Generate(username, time.Minute)
	require.NoError(t, err)

	payload, err := authenticator.Validate(token)
	require.NoError(t, err)

	require.Equal(t, username, payload.UserIdentifier)
}

func TestInvalidSymmetricKey(t *testing.T) {
	paseto := paseto.NewV2()

	// Generate wrong sized key.
	key := make([]byte, 16)
	_, err := rand.Read(key)
	require.NoError(t, err)

	_, err = NewPasetoAuthenticator(key, paseto)
	require.Equal(t, &InvalidSymmetricKeySizeError{}, err)
}

func TestInvalidToken(t *testing.T) {
	paseto := paseto.NewV2()

	key := make([]byte, 32)
	_, err := rand.Read(key)
	require.NoError(t, err)

	authenticator, err := NewPasetoAuthenticator(key, paseto)
	require.NoError(t, err)

	_, err = authenticator.Validate(invalid)
	require.Equal(t, &InvalidTokenError{}, err)
}

func TestExpiredToken(t *testing.T) {
	paseto := paseto.NewV2()

	key := make([]byte, 32)
	_, err := rand.Read(key)
	require.NoError(t, err)

	authenticator, err := NewPasetoAuthenticator(key, paseto)
	require.NoError(t, err)

	token, err := authenticator.Generate(username, 0)
	require.NoError(t, err)

	_, err = authenticator.Validate(token)
	require.Equal(t, &TokenExpiryError{}, err)
}
