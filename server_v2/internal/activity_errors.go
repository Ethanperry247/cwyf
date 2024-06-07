package internal

import (
	"danger-dodgers/pkg/errors"
	"fmt"
)

type ActivityFieldTooLargeError struct {
	errors.BadRequest
	max   int
	field string
}

func (err *ActivityFieldTooLargeError) Error() string {
	return fmt.Sprintf("activity field %s is too large, max size is %d characters", err.field, err.max)
}

type ActivityFieldBlankError struct {
	errors.BadRequest
	field string
}

func (err *ActivityFieldBlankError) Error() string {
	return fmt.Sprintf("activity field %s cannot be blank", err.field)
}

type ActivityAlreadyExistsError struct {
	errors.BadRequest
}

func (err *ActivityAlreadyExistsError) Error() string {
	return "the activity with this id already exists"
}
