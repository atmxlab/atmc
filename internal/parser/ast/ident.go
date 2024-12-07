package ast

import "github.com/atmxlab/atmcfg/internal/types"

type ident struct {
	identNode
	string
}

func NewIdent(string string, loc types.Location) Ident {
	i := ident{string: string}
	i.loc = loc

	return i
}
