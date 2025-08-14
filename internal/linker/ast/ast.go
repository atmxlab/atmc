package ast

import (
	"github.com/atmxlab/atmcfg/pkg/errors"
)

type Ast struct {
	object Object
}

func (a Ast) Object() Object {
	return a.object
}

func NewAst(object Object) Ast {
	return Ast{object: object}
}

func (a Ast) FindExpByPath(path []Ident) (Expression, error) {
	if len(path) == 0 {
		return a.object, nil
	}

	foundNode, err := a.object.FindExpByPath(path)
	if err != nil {
		return nil, errors.Wrapf(err, "error finding node by path")
	}

	return foundNode, nil
}
