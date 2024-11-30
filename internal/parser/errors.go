package parser

import "github.com/atmxlab/atmcfg/pkg/errors"

var (
	ErrUnexpectedToken = errors.New("unexpected token")
	ErrTokenMismatch   = errors.New("token mismatch")
)
