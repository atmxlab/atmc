package ast

type Var struct {
	expressionNode
	path []Ident
}

func (v Var) Path() []Ident {
	return v.path
}

func NewVar(path []Ident) Var {
	return Var{path: path}
}
