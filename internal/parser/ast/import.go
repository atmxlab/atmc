package ast

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

func NewImport(name Ident, path Path) Import {
	return Import{name: name, path: path}
}

func (i Import) Path() Path {
	return i.path
}

func (i Import) Name() Ident {
	return i.name
}
