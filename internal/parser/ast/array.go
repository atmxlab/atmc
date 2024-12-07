package ast

import "github.com/atmxlab/atmcfg/internal/types"

type Array struct {
	expressionNode
	elements []Expression
}

func (a Array) Elements() []Expression {
	return a.elements
}

func NewArray(elements []Expression, loc types.Location) Array {
	a := Array{elements: elements}
	a.loc = loc
	return a
}
