package ast

import "github.com/atmxlab/atmcfg/internal/types"

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
