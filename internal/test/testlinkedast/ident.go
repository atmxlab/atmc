package testlinkedast

import (
	linkedast "github.com/atmxlab/atmcfg/internal/linker/ast"
	"github.com/atmxlab/atmcfg/internal/test/testutils"
)

func NewIdent(name string, hooks ...func(ident *linkedast.Ident)) linkedast.Ident {
	id := linkedast.NewIdent(name)
	testutils.ApplyHooks(&id, hooks)

	return id
}
