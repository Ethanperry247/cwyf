package internal

import (
	"danger-dodgers/pkg/auth"
	"time"
)

const TOKEN_LIFESPAN = time.Minute * 15

type TokenService struct {
	auth auth.Authenticator
}

func NewTokenService(auth auth.Authenticator) *TokenService {
	return &TokenService{
		auth: auth,
	}
}

func (s *TokenService) Token(user *User) (*Token, error) {
	token, err := s.auth.Generate(user.Email, TOKEN_LIFESPAN)
	if err != nil {
		return nil, err
	}

	return &Token{
		Token: token,
	}, nil
}
