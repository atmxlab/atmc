package parser

import (
	"github.com/atmxlab/atmc/pkg/errors"
	"github.com/atmxlab/atmc/types/token"
	"github.com/samber/lo"
)

var (
	ErrTokenMismatch   = errors.New("token mismatch")
	ErrUnexpectedToken = errors.New("unexpected token")
	ErrExpectedNode    = errors.New("expected node")
	ErrTokenNotExist   = errors.New("token not exist")
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

func NewErrExpectedNode(expectedNodes ...string) error {
	return errors.Wrapf(
		ErrExpectedNode,
		"expected nodes: %v",
		expectedNodes,
	)
}

func NewErrTokenNotExist(expectedTokens ...token.Type) error {
	expectedTokensStr := lo.Map(expectedTokens, func(tokType token.Type, _ int) string {
		return tokType.String()
	})
	return errors.Wrapf(
		ErrTokenNotExist,
		"expected tokens: %v",
		expectedTokensStr,
	)
}
