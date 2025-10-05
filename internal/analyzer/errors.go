package semantic

import "github.com/atmxlab/atmcfg/pkg/errors"

var (
	ErrUnusedVariable    = errors.New("unused variable")
	ErrUndefinedVariable = errors.New("undefined variable")
)
