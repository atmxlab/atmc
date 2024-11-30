package parser

import (
	"github.com/atmxlab/atmcfg/internal/types/token"
	"github.com/atmxlab/atmcfg/pkg/errors"
)

var (
	ErrUnexpectedNode  = errors.New("unexpected node")
	ErrUnexpectedToken = errors.New("unexpected token")
	ErrTokenMismatch   = errors.New("token mismatch")
)

func NewErrUnexpectedToken(expectedTokens ...token.Type) error {
	return errors.Wrapf(ErrUnexpectedToken, "expeted tokens: %v", expectedTokens)
}

// NewUnexpectedNodeErr
// TODO: rename!
func NewUnexpectedNodeErr(expectedNodes ...string) error {
	return errors.Wrapf(ErrUnexpectedToken, "expeted nodes: %v", expectedNodes)
}
