package testast

import (
	"github.com/atmxlab/atmc/parser/ast"
	"github.com/atmxlab/atmc/test/testutils"
	"github.com/atmxlab/atmc/types"
)

func NewPath(name string, hooks ...func(ident *ast.Path)) ast.Path {
	path := ast.NewPath(name, types.Location{})
	testutils.ApplyHooks(&path, hooks)

	return path
}
