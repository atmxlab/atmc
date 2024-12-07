package ast

import "github.com/atmxlab/atmcfg/internal/types"

type Object struct {
	expressionNode
	spreads []Spread
	entries []EntryNode
}

func (o Object) Spreads() []Spread {
	return o.spreads
}

func (o Object) Entries() []EntryNode {
	return o.entries
}

func NewObject(spreads []Spread, entries []EntryNode, loc types.Location) Object {
	o := Object{
		spreads: spreads,
		entries: entries,
	}
	o.loc = loc

	return o
}
