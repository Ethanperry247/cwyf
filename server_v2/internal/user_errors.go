package internal

import (
	"danger-dodgers/pkg/errors"
	"fmt"
)

type InvalidPasswordError struct {
	errors.BadRequest
}

func (err *InvalidPasswordError) Error() string {
	return "the provided password is invalid"
}

type EmailAlreadyExistsError struct {
	errors.BadRequest
	email string
}

func (err *EmailAlreadyExistsError) Error() string {
	return fmt.Sprintf("a user with the email %s already exists", err.email)
}

type InvalidEmailError struct {
	errors.BadRequest
}

func (err *InvalidEmailError) Error() string {
	return "the provided email formatting is invalid"
}

type CannotUpdateEmailError struct {}

func (err *CannotUpdateEmailError) Error() string {
	return "email cannot be updated"
}