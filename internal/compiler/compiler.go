package compiler

import "github.com/atmxlab/atmcfg/internal/parser/ast"

type Ast interface {
	Inspect(func(node ast.Node))
}

type Compiler struct {
}

func (c Compiler) Compile(a Ast) map[string]any {
	a.Inspect(func(node ast.Node) {
		// switch n := node.(type) {
		// case ast.Object:
		//
		// }
	})
	
	return nil
}
