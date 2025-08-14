package ast

import (
	"github.com/atmxlab/atmcfg/internal/types"
	"github.com/atmxlab/atmcfg/pkg/errors"
)

type Spread struct {
	expressionNode
	v Var
}

func (Spread) isEntry() {}

func (s Spread) Var() Var {
	return s.v
}

func NewSpread(v Var, loc types.Location) Spread {
	s := Spread{v: v}
	s.loc = loc

	return s
}

func (s Spread) inspect(handler func(node Node) error) error {
	if err := handler(s); err != nil {
		return errors.Wrap(err, `failed to inspect spread`)
	}

	if err := s.v.inspect(handler); err != nil {
		return errors.Wrap(err, `failed to inspect spread var`)
	}

	return nil
}
