package internal

import "danger-dodgers/pkg/db"

func VerifyField(min, max int, field, fieldName, model string) error {
	if len(field) < min {
		return &FieldTooSmallError{
			min:   min,
			field: fieldName,
			model: model,
		}
	}

	if len(field) > max {
		return &FieldTooLargeError{
			max:   max,
			field: fieldName,
			model: model,
		}
	}

	return nil
}

func Peek(id, model string, database db.Peek) error {
	err := database.Peek(id)
	if err == nil {
		return &AlreadyExistsError{
			model: model,
		}
	}

	_, ok := err.(*db.NotFoundError)
	if !ok {
		return err
	}

	return nil
}

type FieldVerifier struct {
	min   int
	max   int
	model string
	field string
}

func (verifier FieldVerifier) WithField(field string) FieldVerifier {
	verifier.field = field
	return verifier
}

func (verifier FieldVerifier) WithModel(model string) FieldVerifier {
	verifier.model = model
	return verifier
}

func (verifier FieldVerifier) WithMin(min int) FieldVerifier {
	verifier.min = min
	return verifier
}

func (verifier FieldVerifier) WithMax(max int) FieldVerifier {
	verifier.max = max
	return verifier
}

func (verifier FieldVerifier) Build() func(string) error {
	return func(s string) error {
		return VerifyField(verifier.min, verifier.max, s, verifier.field, verifier.model)
	}
}

func CompareUserIDs(id, actual string) error {
	if id != actual {
		return &InvalidUserIDError{}
	}

	return nil
}