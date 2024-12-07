package ast

import "github.com/atmxlab/atmcfg/internal/types"

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

func NewFile(imports []Import, object Object, loc types.Location) File {
	f := File{imports: imports, object: object}
	f.loc = loc

	return f
}
