package ast

import (
	"fmt"

	"github.com/atmxlab/atmc/types"
)

type Node interface {
	Location() types.Location
	inspect(func(Node) error) error
}

type Entry interface {
	Node
	isEntry()
}

type Ident interface {
	Node
	fmt.Stringer
	isIdent()
}

type Literal interface {
	Node
	isLiteral()
}

type Expression interface {
	Node
	isExpression()
}

type Statement interface {
	Node
	isStatement()
}
