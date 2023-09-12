package auth

type InvalidSymmetricKeySizeError struct {
}

func (err *InvalidSymmetricKeySizeError) Error() string {
	return "key size is incorrect"
}

type InvalidTokenError struct {
}

func (err *InvalidTokenError) Error() string {
	return "invalid token supplied"
}

type TokenExpiryError struct {
}

func (err *TokenExpiryError) Error() string {
	return "the token has expired"
}