package internal

import (
	"danger-dodgers/pkg/errors"
	"fmt"
)

type UserAlreadyExistsError struct {
	errors.BadRequest
}

func (err *UserAlreadyExistsError) Error() string {
	return "user with the given identifier already exists"
}

type BlankEmailError struct {
	errors.BadRequest
}

func (err *BlankEmailError) Error() string {
	return "user email cannot be blank"
}

type InvalidEmailError struct{
	errors.BadRequest
}

func (err *InvalidEmailError) Error() string {
	return "the provided email is invalid"
}

type EmailTooLargeError struct {
	errors.BadRequest
	max int
}

func (err *EmailTooLargeError) Error() string {
	return fmt.Sprintf("the provided email is too large, must be %d characters maximum", err.max)
}

type UsernameTooLargeError struct {
	errors.BadRequest
	max int
}

func (err *UsernameTooLargeError) Error() string {
	return fmt.Sprintf("the provided username is too large, must be %d characters maximum", err.max)
}

type BlankUsernameError struct {
	errors.BadRequest
}

func (err *BlankUsernameError) Error() string {
	return "the username must not be blank"
}

type PasswordTooLargeError struct {
	errors.BadRequest
	max int
}

func (err *PasswordTooLargeError) Error() string {
	return fmt.Sprintf("the provided password is too large, must be %d characters maximum", err.max)
}

type PasswordTooSmallError struct {
	errors.BadRequest
	min int
}

func (err *PasswordTooSmallError) Error() string {
	return fmt.Sprintf("the provided password is too small, must be %d characters minimum", err.min)
}

type NameTooLargeError struct {
	errors.BadRequest
	max int
}

func (err *NameTooLargeError) Error() string {
	return fmt.Sprintf("the provided name is too large, must be %d characters maximum", err.max)
}

type BlankNameError struct {
	errors.BadRequest
}

func (err *BlankNameError) Error() string {
	return "the name must not be blank"
}

type InvalidPasswordError struct {
	errors.BadRequest
}

func (err *InvalidPasswordError) Error() string {
	return "the provided password is invalid"
}

