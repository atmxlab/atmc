package ast

import (
	"github.com/atmxlab/atmcfg/internal/types"
	"github.com/atmxlab/atmcfg/pkg/errors"
)

type Import struct {
	statementNode
	name Ident
	path Path
}

type Path struct {
	identNode
}

func NewPath(string string, loc types.Location) Path {
	p := Path{identNode{string: string}}
	p.loc = loc

	return p
}

func NewImport(
	name Ident,
	path Path,
) Import {
	i := Import{name: name, path: path}
	i.loc = types.NewLocation(
		name.Location().Start(),
		path.Location().End(),
	)

	return i
}

func (i Import) Path() Path {
	return i.path
}

func (i Import) Name() Ident {
	return i.name
}

func (i Import) inspect(handler func(node Node) error) error {
	if err := handler(i); err != nil {
		return errors.Wrap(err, "inspect Import node")
	}

	return nil
}
