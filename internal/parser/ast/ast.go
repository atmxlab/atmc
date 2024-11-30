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
	entryNode()
}

type Ident interface {
	Node
	identNode()
}

type Literal interface {
	Entry
	literalNode()
}
