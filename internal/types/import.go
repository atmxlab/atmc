package types

type Import struct {
	name Name
	path Path
}

type Name string
type Path string

func NewImport(name Name, path Path) *Import {
	return &Import{name: name, path: path}
}

func (i Import) Path() Path {
	return i.path
}

func (i Import) Name() Name {
	return i.name
}
