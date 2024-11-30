package ast

type Array struct {
	entryNode
	elements []Node
}

func NewArray(elements []Node) Array {
	return Array{elements: elements}
}
