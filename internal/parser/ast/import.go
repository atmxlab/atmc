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

func NewPath(string string) Path {
	return Path{string: string}
}

func NewImport(
	name Ident,
	path Path,
	loc types.Location,
) Import {
	i := Import{name: name, path: path}
	i.loc = loc

	return i
}

func (i Import) Path() Path {
	return i.path
}

func (i Import) Name() Ident {
	return i.name
}
