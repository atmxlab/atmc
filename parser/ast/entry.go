package ast

import (
	"github.com/atmxlab/atmc/pkg/errors"
	"github.com/atmxlab/atmc/types"
)

type KV struct {
	entryNode
	key   Ident
	value Expression
}

func (kv KV) Key() Ident {
	return kv.key
}

func (kv KV) Value() Expression {
	return kv.value
}

func NewKV(key Ident, value Expression) KV {
	e := KV{key: key, value: value}
	e.loc = types.NewLocation(
		key.Location().Start(),
		value.Location().End(),
	)

	return e
}

func (kv KV) inspect(handler func(node Node) error) error {
	if err := handler(kv); err != nil {
		return errors.Wrap(err, `failed to inspect key value`)
	}

	if err := kv.value.inspect(handler); err != nil {
		return errors.Wrap(err, `failed to inspect expression in key value`)
	}

	return nil
}
