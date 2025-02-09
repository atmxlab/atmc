package ast

import (
	"github.com/atmxlab/atmcfg/internal/types"
	"github.com/atmxlab/atmcfg/pkg/errors"
)

type Object struct {
	expressionNode
	entries []Entry
}

func (o Object) Entries() []Entry {
	return o.entries
}

func NewObject(entries []Entry, loc types.Location) Object {
	o := Object{
		entries: entries,
	}
	o.loc = loc

	return o
}

func (o Object) inspect(handler func(node Node) error) error {
	for _, entry := range o.entries {
		if err := entry.inspect(handler); err != nil {
			return errors.Wrap(err, "inspect entry")
		}
	}

	return nil
}
