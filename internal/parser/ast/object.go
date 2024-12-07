package ast

import "github.com/atmxlab/atmcfg/internal/types"

type Object struct {
	expressionNode
	spreads []Spread
	entries []Entry
}

func (o Object) Spreads() []Spread {
	return o.spreads
}

func (o Object) Entries() []Entry {
	return o.entries
}

func NewObject(spreads []Spread, entries []Entry, loc types.Location) Object {
	o := Object{
		spreads: spreads,
		entries: entries,
	}
	o.loc = loc

	return o
}
