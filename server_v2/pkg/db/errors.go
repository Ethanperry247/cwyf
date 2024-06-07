package db

import (
	"fmt"

	"danger-dodgers/pkg/errors"
)

type AlreadyExistsError struct {
	errors.BadRequest
	id string
}

func (err *AlreadyExistsError) Error() string {
	return fmt.Sprintf("object with id %s already exists", err.id)
}

type NotFoundError struct {
	errors.NotFound
	id string
}

func (err *NotFoundError) Error() string {
	return fmt.Sprintf("object with the key %s was not found", err.id)
}

type InvalidKeyError struct {
	errors.BadRequest
	field string
}

func (err *InvalidKeyError) Error() string {
	return fmt.Sprintf("the %s must only contain letters, numbers, or the following special characters: _@!%+$.-", err.field)
}