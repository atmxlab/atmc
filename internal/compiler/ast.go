package compiler

import (
	"github.com/atmxlab/atmcfg/internal/parser/ast"
)

type Ast interface {
	Inspect(func(node ast.Node) error) error
	Imports() []Import
}

type Import interface {
	Name() string
	Path() string
}
