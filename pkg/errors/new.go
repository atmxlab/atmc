package errors

import (
	"errors"
	"fmt"
)

func New(msg string) error {
	return errors.New(msg)
}

func Newf(format string, arg ...any) error {
	return fmt.Errorf(format, arg...)
}
