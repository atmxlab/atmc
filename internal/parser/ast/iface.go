package ast

import (
	"fmt"

	"github.com/atmxlab/atmcfg/internal/types"
)

type Node interface {
	Location() types.Location
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
