package ast

import (
	"github.com/atmxlab/atmcfg/pkg/errors"
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

func (o Object) FindExpByPath(path []Ident) (Expression, error) {
	if len(path) == 0 {
		return o, nil
	}

	foundNode, err := o.findExpByPath(path)
	if err != nil {
		return nil, errors.Wrapf(err, "error finding node by path")
	}

	return foundNode, nil
}

func (o Object) findExpByPath(path []Ident) (Expression, error) {
	i := 0
	for _, kv := range o.kv {
		if kv.Key().String() != path[i].String() {
			continue
		}

		if i == len(path)-1 {
			return kv.Value(), nil
		}

		switch v := kv.Value().(type) {
		case Object:
			return v.findExpByPath(path[1:])
		default:
			return nil, errors.New("unexpected value")
		}
	}

	return nil, errors.New("node not found")
}
