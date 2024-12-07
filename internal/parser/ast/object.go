package ast

import "github.com/atmxlab/atmcfg/internal/types"

type Object struct {
	expressionNode
	spreads []Spread
	kvs     []KV
}

func (o Object) Spreads() []Spread {
	return o.spreads
}

func (o Object) Kvs() []KV {
	return o.kvs
}

func NewObject(spreads []Spread, kvs []KV, loc types.Location) Object {
	o := Object{
		spreads: spreads,
		kvs:     kvs,
	}
	o.loc = loc

	return o
}
