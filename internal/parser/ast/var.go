package ast

import "github.com/atmxlab/atmcfg/internal/types"

type Var struct {
	expressionNode
	path []Ident
}

func (v Var) Path() []Ident {
	return v.path
}

func NewVar(path []Ident, loc types.Location) Var {
	v := Var{path: path}
	v.loc = loc

	return v
}
