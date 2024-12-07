package ast

type File struct {
	node
	imports []Import
	object  Object
}

func (f File) Imports() []Import {
	return f.imports
}

func (f File) Object() Object {
	return f.object
}

func NewFile(imports []Import, object Object) File {
	return File{imports: imports, object: object}
}
