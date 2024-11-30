package ast

type Ast struct {
	root Node
}

func (a Ast) Root() Node {
	return a.root
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
