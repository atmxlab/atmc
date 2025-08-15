package analyzer

import "github.com/atmxlab/atmc/pkg/errors"

var (
	ErrUnusedVariable    = errors.New("unused variable")
	ErrUndefinedVariable = errors.New("undefined variable")
)
