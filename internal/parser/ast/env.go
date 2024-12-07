package ast

import "github.com/atmxlab/atmcfg/internal/types"

type Env struct {
	expressionNode
	name Ident
}

func NewEnv(name Ident, loc types.Location) Env {
	e := Env{name: name}
	e.loc = loc

	return e
}
