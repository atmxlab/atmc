package errors

import "fmt"

func Wrap(err error, msg string) error {
	if err == nil {
		return nil
	}

	return fmt.Errorf("%s: %w", msg, err)
}

func Wrapf(err error, msg string, a ...any) error {
	if err == nil {
		return nil
	}

	return Wrap(err, fmt.Sprintf(msg, a...))
}
