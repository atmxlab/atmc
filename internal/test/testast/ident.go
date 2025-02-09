package testast

import (
	"github.com/atmxlab/atmcfg/internal/parser/ast"
	"github.com/atmxlab/atmcfg/internal/test/testutils"
	"github.com/atmxlab/atmcfg/internal/types"
)

func NewIdent(name string, hooks ...func(ident *ast.Ident)) ast.Ident {
	id := ast.NewIdent(name, types.Location{})
	testutils.ApplyHooks(&id, hooks)

	return id
}
