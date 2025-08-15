package testlinkedast

import (
	linkedast "github.com/atmxlab/atmcfg/internal/linker/ast"
)

type AstBuilder struct {
	object linkedast.Object
}

func NewAstBuilder() *AstBuilder {
	return &AstBuilder{}
}

func (b *AstBuilder) Object(hook func(ob *ObjectBuilder)) {
	ob := NewObjectBuilder()
	hook(ob)
	b.object = ob.Build()
}

func (b *AstBuilder) Build() linkedast.Ast {
	return linkedast.NewAst(b.object)
}

type ObjectBuilder struct {
	kv []linkedast.KV
}

func NewObjectBuilder() *ObjectBuilder {
	return &ObjectBuilder{}
}

func (b *ObjectBuilder) Build() linkedast.Object {
	return linkedast.NewObject(b.kv)
}
