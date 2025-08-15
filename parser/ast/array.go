package ast

import (
	"github.com/atmxlab/atmc/pkg/errors"
	"github.com/atmxlab/atmc/types"
)

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

func (a Array) inspect(handler func(node Node) error) error {
	if err := handler(a); err != nil {
		return errors.Wrap(err, `failed to inspect array`)
	}

	for _, element := range a.elements {
		if err := element.inspect(handler); err != nil {
			return errors.Wrap(err, `failed to inspect element`)
		}
	}

	return nil
}
