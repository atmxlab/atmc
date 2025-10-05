package testlinkedast

import linkedast "github.com/atmxlab/atmcfg/internal/linker/ast"

type Builder struct {
	obj linkedast.Object
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

func (b *Builder) Build() linkedast.Ast {
	return linkedast.NewAst(b.obj)
}

type ObjectBuilder struct {
	kv []linkedast.KV
}

func NewObjectBuilder() *ObjectBuilder {
	return &ObjectBuilder{}
}

func (b *ObjectBuilder) KV(hook func(kvb *KVBuilder)) *ObjectBuilder {
	kvb := NewKVBuilder()
	hook(kvb)
	b.kv = append(b.kv, kvb.Build())
	return b
}

func (b *ObjectBuilder) Build() linkedast.Object {
	return linkedast.NewObject(b.kv)
}

type KVBuilder struct {
	key   linkedast.Ident
	value linkedast.Expression
}

func NewKVBuilder() *KVBuilder {
	return &KVBuilder{}
}

func (b *KVBuilder) Key(key string) *KVBuilder {
	b.key = linkedast.NewIdent(key)
	return b
}

func (b *KVBuilder) Value(value linkedast.Expression) *KVBuilder {
	b.value = value
	return b
}

func (b *KVBuilder) Build() linkedast.KV {
	return linkedast.NewKV(b.key, b.value)
}

type ArrayBuilder struct {
	elements []linkedast.Expression
}

func NewArrayBuilder() *ArrayBuilder {
	return &ArrayBuilder{}
}

func (b *ArrayBuilder) Element(element linkedast.Expression) *ArrayBuilder {
	b.elements = append(b.elements, element)
	return b
}

func (b *ArrayBuilder) Build() linkedast.Array {
	return linkedast.NewArray(b.elements)
}
