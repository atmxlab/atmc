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

func (a Ast) Inspect(handler func(node Node) error) error {
	return a.root.inspect(handler)
}
