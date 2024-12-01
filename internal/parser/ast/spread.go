package ast

type Spread struct {
	expressionNode
	v Var
}

func NewSpread(v Var) Spread {
	return Spread{v: v}
}
