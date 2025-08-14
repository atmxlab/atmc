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

func (entryNode) isEntry() {}

type identNode struct {
	node
	string
}

func (i identNode) String() string {
	return i.string
}

func (identNode) isIdent() {}

type statementNode struct {
	node
	test string
}

func (statementNode) isStatement() {}

type expressionNode struct {
	node
}

func (expressionNode) isExpression() {}

type literalNode[T any] struct {
	expressionNode
	value T
}

func (l literalNode[T]) isLiteral() {}

func (l literalNode[T]) Value() T {
	return l.value
}

func (l literalNode[T]) inspect(handler func(node Node) error) error {
	if err := handler(l); err != nil {
		return errors.Wrap(err, `failed to inspect literal node`)
	}

	return nil
}
