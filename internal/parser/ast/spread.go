package ast

import "github.com/atmxlab/atmcfg/internal/types"

type Spread struct {
	expressionNode
	v Var
}

func (Spread) entryNodeMarker() {}

func (s Spread) Var() Var {
	return s.v
}

func NewSpread(v Var, loc types.Location) Spread {
	s := Spread{v: v}
	s.loc = loc

	return s
}
