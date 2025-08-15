package testast

import (
	ast2 "github.com/atmxlab/atmc/parser/ast"
	"github.com/atmxlab/atmc/types"
)

type AstBuilder struct {
	imports []ast2.Import
	object  ast2.Object
}

func NewAstBuilder() *AstBuilder {
	return &AstBuilder{}
}

func (b *AstBuilder) Import(hook func(ib *ImportBuilder)) {
	ib := NewImportBuilder()
	hook(ib)
	b.imports = append(b.imports, ib.Build())
}

func (b *AstBuilder) Object(hook func(ob *ObjectBuilder)) {
	ob := NewObjectBuilder()
	hook(ob)
	b.object = ob.Build()
}

func (b *AstBuilder) Build() ast2.Ast {
	return ast2.NewAst(ast2.NewFile(b.imports, b.object))
}

type ImportBuilder struct {
	name ast2.Ident
	path ast2.Path
}

func NewImportBuilder() *ImportBuilder {
	return &ImportBuilder{}
}

func (b *ImportBuilder) Name(name ast2.Ident) {
	b.name = name
}

func (b *ImportBuilder) Path(path ast2.Path) {
	b.path = path
}

func (b *ImportBuilder) Build() ast2.Import {
	return ast2.NewImport(b.name, b.path)
}

type ObjectBuilder struct {
	location types.Location
	entries  []ast2.Entry
}

func NewObjectBuilder() *ObjectBuilder {
	return &ObjectBuilder{}
}

func (b *ObjectBuilder) Location(location types.Location) {
	b.location = location
}

func (b *ObjectBuilder) Spread(hook func(sb *SpreadBuilder)) {
	sb := NewSpreadBuilder()
	hook(sb)
	b.entries = append(b.entries, sb.Build())
}

func (b *ObjectBuilder) KV(hook func(kb *KVBuilder)) {
	sb := NewKVBuilder()
	hook(sb)
	b.entries = append(b.entries, sb.Build())
}

func (b *ObjectBuilder) Build() ast2.Object {
	return ast2.NewObject(b.entries, b.location)
}

type SpreadBuilder struct {
	location types.Location
	v        ast2.Var
}

func NewSpreadBuilder() *SpreadBuilder {
	return &SpreadBuilder{}
}

func (b *SpreadBuilder) Location(location types.Location) {
	b.location = location
}

func (b *SpreadBuilder) Var(hook func(vb *VarBuilder)) {
	vb := NewVarBuilder()
	hook(vb)
	b.v = vb.Build()
}

func (b *SpreadBuilder) Build() ast2.Spread {
	return ast2.NewSpread(b.v, b.location)
}

type VarBuilder struct {
	path []ast2.Ident
}

func NewVarBuilder() *VarBuilder {
	return &VarBuilder{}
}

func (b *VarBuilder) Part(part ast2.Ident) {
	b.path = append(b.path, part)
}

func (b *VarBuilder) Build() ast2.Var {
	return ast2.NewVar(b.path)
}

type KVBuilder struct {
	key   ast2.Ident
	value ast2.Expression
}

func NewKVBuilder() *KVBuilder {
	return &KVBuilder{}
}

func (b *KVBuilder) Key(key ast2.Ident) *KVBuilder {
	b.key = key
	return b
}

func (b *KVBuilder) Value(value ast2.Expression) *KVBuilder {
	b.value = value
	return b
}

func (b *KVBuilder) Var(hook func(vb *VarBuilder)) *KVBuilder {
	vb := NewVarBuilder()
	hook(vb)
	return b.Value(vb.Build())
}

func (b *KVBuilder) Build() ast2.KV {
	return ast2.NewKV(b.key, b.value)
}
