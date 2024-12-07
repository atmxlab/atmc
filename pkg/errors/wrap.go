package errors

import "fmt"

func Wrap(err error, msg string) error {
	if err == nil {
		return nil
	}

	return fmt.Errorf("%w: %s", err, msg)
}

func Wrapf(err error, msg string, a ...any) error {
	if err == nil {
		return nil
	}

	return fmt.Errorf("%w: %s", err, fmt.Sprintf(msg, a...))
}
