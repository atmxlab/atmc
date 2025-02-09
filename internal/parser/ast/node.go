package ast

import (
	"github.com/atmxlab/atmcfg/internal/types"
	"github.com/atmxlab/atmcfg/pkg/errors"
)

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

func (l literalNode[T]) inspect(handler func(node Node) error) error {
	if err := handler(l); err != nil {
		return errors.Wrap(err, `failed to inspect literal node`)
	}

	return nil
}
