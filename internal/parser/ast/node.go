package ast

import "github.com/atmxlab/atmcfg/internal/types"

type node struct {
	pos types.Position
}

func (n node) Pos() uint {
	return n.pos.Pos()
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

type expressionNode struct {
	statementNode
}

func (expressionNode) expressionNodeMarker() {}

type statementNode struct {
	node
}

func (statementNode) statementNodeMarker() {}

type literalNode[T any] struct {
	expressionNode
	value T
}

func (l literalNode[T]) literalNodeMarker() {}

func (l literalNode[T]) Value() T {
	return l.value
}
