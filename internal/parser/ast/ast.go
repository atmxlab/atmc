package ast

import "github.com/atmxlab/atmcfg/pkg/errors"

type Ast struct {
	root File
}

func (a Ast) Root() File {
	return a.root
}

func NewAst(root File) Ast {
	return Ast{root: root}
}

func (a Ast) Inspect(handler func(node Node) error) error {
	if err := handler(a.root); err != nil {
		return errors.Wrap(err, "inspecting file node")
	}

	return a.root.inspect(handler)
}
