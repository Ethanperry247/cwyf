package internal

import "danger-dodgers/pkg/errors"

type MissingBearerTokenError struct {
	errors.BadRequest
}

func (err *MissingBearerTokenError) Error() string {
	return "bearer token missing"
}

type ImproperlyFormattedBearerTokenError struct {
	errors.BadRequest
}

func (err *ImproperlyFormattedBearerTokenError) Error() string {
	return "bearer token is improperly formatted"
}

type InvalidTokenTypeError struct {
	errors.BadRequest
}

func (err *InvalidTokenTypeError) Error() string {
	return "token type is invalid"
}

type InvalidUserIDError struct {
	errors.Forbidden
}

func (err *InvalidUserIDError) Error() string {
	return "user id does not match id of token"
}
