package ast

import (
	"github.com/atmxlab/atmc/pkg/errors"
	"github.com/atmxlab/atmc/types"
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

func (o Object) FindNodeByPath(path []Ident) (Node, error) {
	if len(path) == 0 {
		return o, nil
	}

	foundNode, err := o.findNodeByPath(path)
	if err != nil {
		return nil, errors.Wrapf(err, "error finding node by path")
	}

	return foundNode, nil
}

func (o Object) findNodeByPath(path []Ident) (Node, error) {
	i := 0
	for _, entry := range o.Entries() {
		kv, ok := entry.(KV)
		if !ok {
			return nil, errors.New("unexpected entry")
		}
		if kv.Key().String() != path[i].String() {
			continue
		}

		if i == len(path)-1 {
			return kv.Value(), nil
		}

		switch v := kv.Value().(type) {
		case Object:
			return v.findNodeByPath(path[1:])
		default:
			return nil, errors.New("unexpected value")
		}
	}

	return nil, errors.New("node not found")
}
