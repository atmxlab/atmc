package ast

import "fmt"

type Node interface {
	Pos() uint
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
