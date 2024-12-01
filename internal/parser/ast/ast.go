package ast

type Ast struct {
	file File
}

func NewAst(file File) Ast {
	return Ast{file: file}
}

type Node interface {
	Pos() uint
}

type Entry interface {
	Node
	entryNodeMarker()
}

type Ident interface {
	Node
	identNodeMarker()
}

type Literal interface {
	Entry
	literalNodeMarker()
}

type Expression interface {
	Entry
	expressionNodeMarker()
}
