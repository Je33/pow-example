package errs

import (
	"errors"
	"pow-example/pkg/logger"
)

type Error struct {
	err error
}

func New(msg error) *Error {
	return &Error{msg}
}

func (e *Error) Error() string {
	return e.err.Error()
}

func (e *Error) Log(log logger.Logger) error {
	log.Error(e.Error())
	return e
}

func Is(err error, target error) bool {
	var e *Error
	switch {
	case errors.As(err, &e):
		return errors.Is(err.(*Error).err, target)
	default:
		return errors.Is(err, target)
	}
}
