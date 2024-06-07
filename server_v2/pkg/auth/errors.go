package auth

import "danger-dodgers/pkg/errors"

type InvalidSymmetricKeySizeError struct {
}

func (err *InvalidSymmetricKeySizeError) Error() string {
	return "key size is incorrect"
}

type InvalidTokenError struct {
	errors.Forbidden
}

func (err *InvalidTokenError) Error() string {
	return "invalid token supplied"
}

type TokenExpiryError struct {
	errors.Forbidden
}

func (err *TokenExpiryError) Error() string {
	return "the token has expired"
}
