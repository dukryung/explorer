package errors

import (
	"errors"
	"fmt"
)

const (
	ErrInvalidRequestCode = iota + 100
	ErrInvalidResponseCode
	ErrInvalidParamsCode
	ErrInvalidMethodCode
	ErrInternalCode = 200
	ErrUnknownCode = 900
)

var (
	ErrInvalidRequest  = errors.New("invalid request")
	ErrInvalidResponse = errors.New("invalid response")

	ErrInvalidParams = errors.New("invalid params")
	ErrInvalidMethod = errors.New("invalid method")

	ErrInternal = errors.New("internal error")
	ErrUnknown = errors.New("unknown error")
)

var _ error = &WrappedError{}

type WrappedError struct {
	code int
	err  error
	msg  string
}

func (e *WrappedError) Error() string {
	return fmt.Sprintf("%s: %s", e.msg, e.err.Error())
}

func (e *WrappedError) Code() int {
	return e.code
}

func Wrap(err error, code int, description string) error {
	return &WrappedError{
		err: err,
		code: code,
		msg: description,
	}
}
