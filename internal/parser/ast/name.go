package ast

type Name struct {
	identNode
	string
}

func NewName(string string) Name {
	return Name{string: string}
}
