package ast

import (
	"github.com/atmxlab/atmc/pkg/errors"
	"github.com/atmxlab/atmc/types"
)

type ident struct {
	identNode
}

func NewIdent(string string, loc types.Location) Ident {
	i := ident{identNode{string: string}}
	i.loc = loc

	return i
}

func (i ident) inspect(handler func(Node) error) error {
	if err := handler(i); err != nil {
		return errors.Wrap(err, "failed to inspect ident")
	}

	return nil
}
