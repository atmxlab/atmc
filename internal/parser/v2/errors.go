package v2

import (
	"github.com/atmxlab/atmcfg/internal/types/token"
	"github.com/atmxlab/atmcfg/pkg/errors"
	"github.com/samber/lo"
)

var (
	ErrTokenMismatch   = errors.New("token mismatch")
	ErrUnexpectedToken = errors.New("unexpected token")
)

func NewErrTokenMismatch(expectedTokens ...token.Type) error {
	expectedTokensStr := lo.Map(expectedTokens, func(tokType token.Type, _ int) string {
		return tokType.String()
	})

	return errors.Wrapf(
		ErrTokenMismatch,
		"expeted tokens: %v",
		expectedTokensStr,
	)
}

func NewErrUnexpectedToken(expectedTokens ...token.Type) error {
	expectedTokensStr := lo.Map(expectedTokens, func(tokType token.Type, _ int) string {
		return tokType.String()
	})

	return errors.Wrapf(
		ErrUnexpectedToken,
		"expected tokens: %v",
		expectedTokensStr,
	)
}
