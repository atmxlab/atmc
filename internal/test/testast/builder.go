package testast

import (
	"github.com/atmxlab/atmcfg/internal/parser/ast"
	"github.com/atmxlab/atmcfg/internal/types"
)

type AstBuilder struct {
	imports []ast.Import
	object  ast.Object
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

func (b *AstBuilder) Build() ast.Ast {
	return ast.NewAst(ast.NewFile(b.imports, b.object))
}

type ImportBuilder struct {
	name ast.Ident
	path ast.Path
}

func NewImportBuilder() *ImportBuilder {
	return &ImportBuilder{}
}

func (b *ImportBuilder) Name(name ast.Ident) {
	b.name = name
}

func (b *ImportBuilder) Path(path ast.Path) {
	b.path = path
}

func (b *ImportBuilder) Build() ast.Import {
	return ast.NewImport(b.name, b.path)
}

type ObjectBuilder struct {
	location types.Location
	entries  []ast.Entry
}

func NewObjectBuilder() *ObjectBuilder {
	return &ObjectBuilder{}
}

func (b *ObjectBuilder) Location(location types.Location) {
	b.location = location
}

func (b *ObjectBuilder) Spread(hook func(sb *SpreadBuilder)) {
	sb := newSpreadBuilder()
	hook(sb)
	b.entries = append(b.entries, sb.Build())
}

func (b *ObjectBuilder) Build() ast.Object {
	return ast.NewObject(b.entries, b.location)
}

type SpreadBuilder struct {
	location types.Location
	v        ast.Var
}

func newSpreadBuilder() *SpreadBuilder {
	return &SpreadBuilder{}
}

func (b *SpreadBuilder) Location(location types.Location) {
	b.location = location
}

func (b *SpreadBuilder) Var(hook func(vb *VarBuilder)) {
	vb := newVarBuilder()
	hook(vb)
	b.v = vb.Build()
}

func (b *SpreadBuilder) Build() ast.Spread {
	return ast.NewSpread(b.v, b.location)
}

type VarBuilder struct {
	path []ast.Ident
}

func newVarBuilder() *VarBuilder {
	return &VarBuilder{}
}

func (b *VarBuilder) Part(part ast.Ident) {
	b.path = append(b.path, part)
}

func (b *VarBuilder) Build() ast.Var {
	return ast.NewVar(b.path)
}
