package errors

import "errors"

type Joiner struct {
	err error
}

func NewJoiner() Joiner {
	return Joiner{err: nil}
}

func (j Joiner) Error() string {
	return j.err.Error()
}

func (j Joiner) Err() error {
	return j.err
}

func (j Joiner) Join(errs ...error) {
	errs = append(errs, j.err)
	j.err = Join(errs...)
}

func Join(errs ...error) error {
	return errors.Join(errs...)
}
