package parser

import (
	"github.com/atmxlab/atmcfg/internal/types/token"
	"github.com/atmxlab/atmcfg/pkg/errors"
)

func (p *Parser) match(tps ...token.Type) bool {
	for _, t := range tps {
		if t == p.mover.Token().Type() {
			return true
		}
	}

	return false
}

func (p *Parser) require(tps ...token.Type) error {
	if p.match(tps...) {
		return nil
	}

	return errors.Wrapf(NewErrUnexpectedToken(tps...), "got: %v", p.mover.Token().Type().String())
}

func (p *Parser) check(tps ...token.Type) error {
	if p.match(tps...) {
		return nil
	}

	return errors.Wrapf(NewErrTokenMismatch(tps...), "got: %v", p.mover.Token().Type().String())
}
