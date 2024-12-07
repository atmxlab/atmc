package ast

type Spread struct {
	statementNode
	v Var
}

func (s Spread) Var() Var {
	return s.v
}

func NewSpread(v Var) Spread {
	return Spread{v: v}
}
