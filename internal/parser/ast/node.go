package ast

import "github.com/atmxlab/atmcfg/internal/types"

type node struct {
	loc types.Location
}

func (n node) Location() types.Location {
	return n.loc
}

type entryNode struct {
	node
}

func (entryNode) entryNodeMarker() {}

type identNode struct {
	node
	string
}

func (i identNode) String() string {
	return i.string
}

func (identNode) identNodeMarker() {}

type statementNode struct {
	node
	test string
}

func (statementNode) statementNodeMarker() {}

type expressionNode struct {
	node
}

func (expressionNode) expressionNodeMarker() {}

type literalNode[T any] struct {
	expressionNode
	value T
}

func (l literalNode[T]) literalNodeMarker() {}

func (l literalNode[T]) Value() T {
	return l.value
}
