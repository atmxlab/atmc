package ast

type Import struct {
	node
	name Name
	path Path
}

type Name struct {
	identNode
	string
}

type Path struct {
	identNode
	string
}

func NewName(string string) Name {
	return Name{string: string}
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
