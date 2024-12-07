package ast

import "github.com/atmxlab/atmcfg/internal/types"

// TODO: придумать название

type KV struct {
	entryNode
	key   Ident
	value Expression
}

func (e KV) Key() Ident {
	return e.key
}

func (e KV) Value() Expression {
	return e.value
}

func NewKV(key Ident, value Expression) KV {
	e := KV{key: key, value: value}
	e.loc = types.NewLocation(
		key.Location().Start(),
		value.Location().End(),
	)

	return e
}
