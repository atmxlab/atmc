package parser

import (
	"github.com/atmxlab/atmcfg/internal/types/token"
	"github.com/atmxlab/atmcfg/pkg/errors"
)

func (p Parser) match(tps ...token.Type) bool {
	for _, t := range tps {
		if t == p.lexer.Token().Type() {
			return true
		}
	}

	return false
}

func (p Parser) require(tps ...token.Type) error {
	if p.match(tps...) {
		return nil
	}

	return errors.Wrap(ErrTokenMismatch, "match")
}
