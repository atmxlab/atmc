package ast

type Ident struct {
	node
	string
}

func NewIdent(string string) Ident {
	return Ident{string: string}
}

func (i Ident) String() string {
	return i.string
}
