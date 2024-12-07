package ast

type Array struct {
	expressionNode
	elements []Expression
}

func (a Array) Elements() []Expression {
	return a.elements
}

func NewArray(elements []Expression) Array {
	return Array{elements: elements}
}
