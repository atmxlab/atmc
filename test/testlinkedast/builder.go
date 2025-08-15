package testlinkedast

import (
	"github.com/atmxlab/atmc/linker/ast"
)

type Builder struct {
	obj ast.Object
}

func NewBuilder() *Builder {
	return &Builder{}
}

func (b *Builder) Object(hook func(ob *ObjectBuilder)) *Builder {
	ob := NewObjectBuilder()
	hook(ob)
	b.obj = ob.Build()
	return b
}

func (b *Builder) Build() ast.Ast {
	return ast.NewAst(b.obj)
}

type ObjectBuilder struct {
	kv []ast.KV
}

func NewObjectBuilder() *ObjectBuilder {
	return &ObjectBuilder{}
}

func (b *ObjectBuilder) KV2(key string, value ast.Expression) *ObjectBuilder {
	b.kv = append(b.kv, ast.NewKV(ast.NewIdent(key), value))
	return b
}

func (b *ObjectBuilder) Build() ast.Object {
	return ast.NewObject(b.kv)
}

type ArrayBuilder struct {
	elements []ast.Expression
}

func NewArrayBuilder() *ArrayBuilder {
	return &ArrayBuilder{}
}

func (b *ArrayBuilder) Element(element ast.Expression) *ArrayBuilder {
	b.elements = append(b.elements, element)
	return b
}

func (b *ArrayBuilder) Build() ast.Array {
	return ast.NewArray(b.elements)
}
