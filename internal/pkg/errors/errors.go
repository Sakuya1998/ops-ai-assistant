package errors

import (
	"errors"
	"fmt"
)

var (
	ErrNotFound         = errors.New("resource not found")
	ErrPermissionDenied = errors.New("permission denied")
	ErrInvalidParameter = errors.New("invalid parameter")
	ErrTimeout          = errors.New("operation timeout")
	ErrInternal         = errors.New("internal error")
)

type AppError struct {
	Code    int
	Message string
	Err     error
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

func (e *AppError) Unwrap() error {
	return e.Err
}

func New(code int, message string) *AppError {
	return &AppError{Code: code, Message: message}
}

func Wrap(err error, code int, message string) *AppError {
	return &AppError{Code: code, Message: message, Err: err}
}
