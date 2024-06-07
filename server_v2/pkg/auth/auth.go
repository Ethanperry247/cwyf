package auth

import (
	"time"

	"github.com/aead/chacha20poly1305"
	"github.com/google/uuid"
	"github.com/o1egl/paseto"
)

type Authenticator interface {
	Generate(string, time.Duration) (string, error)
	Validate(string) (*Payload, error)
}

type TokenType int

const (
	Refresh TokenType = iota
	Authentication
)

type Payload struct {
	ID             uuid.UUID `json:"id"`
	UserIdentifier string    `json:"userIdentifier"`
	Issued         time.Time `json:"issued"`
	Expired        time.Time `json:"expired"`
	Type           TokenType `json:"type"`
}

func (payload *Payload) Valid() error {
	if time.Now().After(payload.Expired) {
		return &TokenExpiryError{}
	}
	return nil
}

type PasetoAuthenticator struct {
	key    []byte
	paseto *paseto.V2
	Type   TokenType
}

func NewPasetoAuthenticator(key []byte, paseto *paseto.V2, options ...func(*PasetoAuthenticator)) (*PasetoAuthenticator, error) {

	if len(key) != chacha20poly1305.KeySize {
		return nil, &InvalidSymmetricKeySizeError{}
	}

	auth := &PasetoAuthenticator{
		key:    key,
		paseto: paseto,
	}

	for _, option := range options {
		option(auth)
	}

	return auth, nil
}

func (authenticator *PasetoAuthenticator) Generate(identifier string, expiry time.Duration) (string, error) {
	time := time.Now()
	payload := &Payload{
		ID:             uuid.New(),
		UserIdentifier: identifier,
		Issued:         time,
		Expired:        time.Add(expiry),
		Type:           authenticator.Type,
	}

	return authenticator.paseto.Encrypt(authenticator.key, payload, nil)
}

func (authenticator *PasetoAuthenticator) Validate(token string) (*Payload, error) {
	payload := new(Payload)

	err := authenticator.paseto.Decrypt(token, authenticator.key, payload, nil)
	if err != nil {
		return nil, &InvalidTokenError{}
	}

	err = payload.Valid()
	if err != nil {
		return nil, err
	}

	return payload, nil
}
