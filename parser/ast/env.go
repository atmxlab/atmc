package ast

import (
	"github.com/atmxlab/atmc/pkg/errors"
	"github.com/atmxlab/atmc/types"
)

type Env struct {
	expressionNode
	name Ident
}

func (e Env) Name() Ident {
	return e.name
}

func NewEnv(name Ident, loc types.Location) Env {
	e := Env{name: name}
	e.loc = loc

	return e
}

func (e Env) inspect(handler func(node Node) error) error {
	if err := handler(e); err != nil {
		return errors.Wrap(err, `failed to inspect env`)
	}

	return nil
}
