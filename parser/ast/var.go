package ast

import (
	"github.com/atmxlab/atmc/pkg/errors"
	"github.com/atmxlab/atmc/types"
	"github.com/samber/lo"
)

type Var struct {
	expressionNode
	path []Ident
}

func (v Var) Path() []Ident {
	return v.path
}

func (v Var) StringPath() []string {
	return lo.Map(v.path, func(item Ident, _ int) string {
		return item.String()
	})
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

func (v Var) inspect(handler func(node Node) error) error {
	if err := handler(v); err != nil {
		return errors.Wrap(err, `failed to inspect var`)
	}

	return nil
}
