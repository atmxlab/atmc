package ast

type File struct {
	node
	imports []Import
	object  Object
}

func NewFile(imports []Import, object Object) File {
	return File{imports: imports, object: object}
}
