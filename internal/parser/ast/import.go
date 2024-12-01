package ast

type Import struct {
	node
	name Name
	path Path
}

type Path struct {
	identNode
	string
}

func NewPath(string string) Path {
	return Path{string: string}
}

func NewImport(name Name, path Path) Import {
	return Import{name: name, path: path}
}

func (i Import) Path() Path {
	return i.path
}

func (i Import) Name() Name {
	return i.name
}
