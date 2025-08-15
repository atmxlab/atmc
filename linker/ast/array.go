package ast

type Array struct {
	node
	expression
	elements []Expression
}

func (a Array) Elements() []Expression {
	return a.elements
}

func NewArray(elements []Expression) Array {
	return Array{elements: elements}
}
