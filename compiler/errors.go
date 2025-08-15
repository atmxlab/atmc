package compiler

import "github.com/atmxlab/atmc/pkg/errors"

var (
	ErrUnsupportedType = errors.New("unsupported type")
	ErrInvalidType     = errors.New("invalid type")
	ErrTypeOverflow    = errors.New("type overflow")
)
