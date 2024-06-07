package passwords

import (
	"golang.org/x/crypto/bcrypt"
)

type Hasher interface {
	Hash(password string) (string, error)
	Compare(password string, hash string) error
}

type BCryptHasher struct {
	Cost int
}

func New(options ...func(*BCryptHasher)) *BCryptHasher {
	hasher := &BCryptHasher{
		Cost: 16,
	}

	for _, option := range options {
		option(hasher)
	}

	return hasher
}

func (hasher *BCryptHasher) Hash(password string) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(password), hasher.Cost)
	if err != nil {
		return "", err
	}

	return string(b), nil
}

func (hasher *BCryptHasher) Compare(password string, hash string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		return err
	}

	return nil
}
