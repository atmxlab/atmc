package ast

type Ast struct {
	root File
}

func (a Ast) Root() File {
	return a.root
}

func NewAst(root File) Ast {
	return Ast{root: root}
}
