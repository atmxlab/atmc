package ast

// TODO: придумать название

type EntryNode struct {
	entryNode
	key   Ident
	value Expression
}

func (e EntryNode) Key() Ident {
	return e.key
}

func (e EntryNode) Value() Expression {
	return e.value
}

func NewEntryNode(entryNode entryNode, key Ident, value Expression) EntryNode {
	return EntryNode{entryNode: entryNode, key: key, value: value}
}
