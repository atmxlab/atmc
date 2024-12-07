package ast

import "github.com/atmxlab/atmcfg/internal/types"

type Var struct {
	expressionNode
	path []Ident
}

func (v Var) Path() []Ident {
	return v.path
}

func NewVar(path []Ident) Var {
	v := Var{path: path}

	start := types.Position{}
	end := types.Position{}

	if len(path) > 0 {
		start = path[0].Location().Start()
		end = path[len(path)-1].Location().End()
	}

	v.loc = types.NewLocation(start, end)

	return v
}
