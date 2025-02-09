package ast

import (
	"fmt"

	"github.com/atmxlab/atmcfg/internal/types"
)

type Node interface {
	Location() types.Location
	inspect(func(Node) error) error
}

type Entry interface {
	Node
	entryNodeMarker()
}

type Ident interface {
	Node
	fmt.Stringer
	identNodeMarker()
}

type Literal interface {
	Node
	literalNodeMarker()
}

type Expression interface {
	Node
	expressionNodeMarker()
}

type Statement interface {
	Node
	statementNodeMarker()
}
