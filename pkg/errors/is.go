package errors

import "errors"

func Is(err error, target error) bool {
	return errors.Is(err, target)
}
