package internal

import (
	"danger-dodgers/pkg/errors"
	"fmt"
)

type ReportFieldTooLargeError struct {
	errors.BadRequest
	max int
	field string
}

func (err *ReportFieldTooLargeError) Error() string {
	return fmt.Sprintf("report field %s is too large, max size is %d characters", err.field, err.max)
}

type ReportFieldBlankError struct {
	errors.BadRequest
	field string
}

func (err *ReportFieldBlankError) Error() string {
	return fmt.Sprintf("report field %s cannot be blank", err.field)
}

type ReportAlreadyExistsError struct {
	errors.BadRequest
}

func (err *ReportAlreadyExistsError) Error() string {
	return "the report with this id already exists"
}