package ast

import (
	"github.com/atmxlab/atmcfg/internal/types"
	"github.com/atmxlab/atmcfg/pkg/errors"
)

type ident struct {
	identNode
	string
}

func NewIdent(string string, loc types.Location) Ident {
	i := ident{string: string}
	i.loc = loc

	return i
}

func (i ident) inspect(handler func(Node) error) error {
	if err := handler(i); err != nil {
		return errors.Wrap(err, "failed to inspect ident")
	}

	return nil
}
