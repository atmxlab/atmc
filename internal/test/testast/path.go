package testast

import (
	"github.com/atmxlab/atmcfg/internal/parser/ast"
	"github.com/atmxlab/atmcfg/internal/test/testutils"
	"github.com/atmxlab/atmcfg/internal/types"
)

func NewPath(name string, hooks ...func(ident *ast.Path)) ast.Path {
	path := ast.NewPath(name, types.Location{})
	testutils.ApplyHooks(&path, hooks)

	return path
}
