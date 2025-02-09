package ast

import (
	"github.com/atmxlab/atmcfg/internal/types"
	"github.com/atmxlab/atmcfg/pkg/errors"
)

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
	f := File{imports: imports, object: object}

	start := object.Location().Start()

	if len(imports) > 0 {
		start = imports[0].Location().Start()
	}

	end := object.Location().End()

	f.loc = types.NewLocation(start, end)

	return f
}

func (f File) inspect(handler func(node Node) error) error {
	for _, imp := range f.imports {
		if err := imp.inspect(handler); err != nil {
			return errors.Wrapf(err, "inspecting import node")
		}
	}

	if err := f.object.inspect(handler); err != nil {
		return errors.Wrapf(err, "inspecting object node")
	}

	return nil
}
