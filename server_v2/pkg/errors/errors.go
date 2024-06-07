package errors

type Error interface {
	error
	Code() Group
}

type Group int

const (
	BAD_REQUEST Group = iota
	NOT_FOUND
	FORBIDDEN
)

type BadRequest struct{}

func (BadRequest) Code() Group {
	return BAD_REQUEST
}

type NotFound struct{}

func (NotFound) Code() Group {
	return NOT_FOUND
}

type Forbidden struct{}

func (Forbidden) Code() Group {
	return FORBIDDEN
}
