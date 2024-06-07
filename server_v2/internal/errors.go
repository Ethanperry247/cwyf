package internal

import (
	"danger-dodgers/pkg/errors"
	"fmt"
)

type FieldTooLargeError struct {
	errors.BadRequest
	max   int
	field string
	model string
}

func (err *FieldTooLargeError) Error() string {
	return fmt.Sprintf("%s field %s is too large, max size is %d characters", err.model, err.field, err.max)
}

type FieldTooSmallError struct {
	errors.BadRequest
	min   int
	field string
	model string
}

func (err *FieldTooSmallError) Error() string {
	return fmt.Sprintf("%s field %s is too small, min size is %d characters", err.model, err.field, err.min)
}

type AlreadyExistsError struct {
	errors.BadRequest
	model string
}

func (err *AlreadyExistsError) Error() string {
	return fmt.Sprintf("%s with this ID already exists", err.model)
}