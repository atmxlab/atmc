package ast

import (
	"github.com/atmxlab/atmcfg/internal/types"
	"github.com/atmxlab/atmcfg/pkg/errors"
)

type Env struct {
	expressionNode
	name Ident
}

func NewEnv(name Ident, loc types.Location) Env {
	e := Env{name: name}
	e.loc = loc // TODO:  можно вычислять!

	return e
}

func (e Env) inspect(handler func(node Node) error) error {
	if err := handler(e); err != nil {
		return errors.Wrap(err, `failed to inspect env`)
	}

	return nil
}
