package ast

import (
	"strings"

	"github.com/atmxlab/atmc/pkg/errors"
	"github.com/samber/lo"
)

type Object struct {
	node
	expression
	kv []KV
}

func NewObject(kv []KV) Object {
	return Object{kv: kv}
}

func (o Object) KV() []KV {
	return o.kv
}

// func (o Object) Get(key Ident) Expression {
// 	for _, kv := range o.kv {
// 		if kv.Key() == key {
// 			return kv.Value()
// 		}
// 	}
// }

func (o Object) FindExpByPath(path []Ident) (Expression, error) {
	if len(path) == 0 {
		return o, nil
	}

	foundNode, err := o.findExpByPath(path)
	if err != nil {
		return nil, errors.Wrapf(err, "error finding expression by path")
	}

	return foundNode, nil
}

func (o Object) findExpByPath(path []Ident) (Expression, error) {
	for _, kv := range o.kv {
		if kv.Key().String() != path[0].String() {
			continue
		}

		if len(path) == 1 {
			return kv.Value(), nil
		}

		switch v := kv.Value().(type) {
		case Object:
			return v.findExpByPath(path[1:])
		default:
			return nil, errors.New("unexpected value")
		}
	}

	return nil, errors.NotFoundf("expression by path not found: path: %s", strings.Join(lo.Map(path, func(
		item Ident,
		index int,
	) string {
		return item.String()
	}), "."))
}
