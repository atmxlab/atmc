package parser

import (
	"github.com/atmxlab/atmc/pkg/errors"
	"github.com/atmxlab/atmc/types/token"
)

func (p *Parser) match(tps ...token.Type) bool {
	if p.mover.IsEmpty() {
		return false
	}

	for _, t := range tps {
		if t == p.mover.Token().Type() {
			return true
		}
	}

	return false
}

func (p *Parser) require(tps ...token.Type) error {
	if p.mover.IsEmpty() {
		return NewErrTokenNotExist(tps...)
	}

	if p.match(tps...) {
		return nil
	}

	return errors.Wrapf(NewErrUnexpectedToken(tps...), "got: %v", p.mover.Token().Type().String())
}

func (p *Parser) check(tps ...token.Type) error {
	if p.mover.IsEmpty() {
		return NewErrTokenNotExist(tps...)
	}

	if p.match(tps...) {
		return nil
	}

	return errors.Wrapf(NewErrTokenMismatch(tps...), "got: %v", p.mover.Token().Type().String())
}
