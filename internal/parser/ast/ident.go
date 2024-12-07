package ast

type ident struct {
	identNode
	string
}

func NewName(string string) Ident {
	return ident{string: string}
}
