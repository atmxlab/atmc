package errors

import "errors"

var (
	ErrInvalidArgument  = errors.New("invalid argument")
	ErrNotFound         = errors.New("not found")
	ErrAlreadyExists    = errors.New("already exists")
	ErrDeadlineExceeded = errors.New("deadline exceeded")
)

func NotFoundf(msg string, a ...any) error {
	return Wrapf(ErrNotFound, msg, a...)
}

func NotFound(msg string) error {
	return Wrap(ErrNotFound, msg)
}

func InvalidArgumentf(msg string, a ...any) error {
	return Wrapf(ErrInvalidArgument, msg, a...)
}

func AlreadyExistsf(msg string, a ...any) error {
	return Wrapf(ErrAlreadyExists, msg, a...)
}

func AlreadyExists(msg string) error {
	return Wrap(ErrAlreadyExists, msg)
}

func DeadlineExceeded(msg string) error {
	return Wrap(ErrDeadlineExceeded, msg)
}
