package testast

import (
	ast2 "github.com/atmxlab/atmc/parser/ast"
	"github.com/atmxlab/atmc/test/testutils"
	"github.com/atmxlab/atmc/types"
)

func NewIdent(name string, hooks ...func(ident *ast2.Ident)) ast2.Ident {
	id := ast2.NewIdent(name, types.Location{})
	testutils.ApplyHooks(&id, hooks)

	return id
}
