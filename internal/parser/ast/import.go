package ast

import "github.com/atmxlab/atmcfg/internal/types"

type Import struct {
	statementNode
	name Ident
	path Path
}

type Path struct {
	identNode
	string
}

func NewPath(string string, loc types.Location) Path {
	p := Path{string: string}
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
