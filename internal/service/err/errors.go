package err

import "errors"

type ValidationError struct {
	Err error
}

func NewValidationError(err error) *ValidationError {
	return &ValidationError{Err: err}
}

func (e *ValidationError) Error() string {
	return e.Err.Error()
}

func (e *ValidationError) Unwrap() error {
	return e.Err
}

var ErrIncorrectData = errors.New("incorrect input data")

var ErrUserAlreadyExists = errors.New("user already exists")
var ErrUserNotFound = errors.New("user not found")

var ErrOrderAlreadyExists = errors.New("order already exists")
var ErrOrderWrongOwner = errors.New("order has been already registered by another user")
var ErrOrderWrongNum = errors.New("invalid order number")
