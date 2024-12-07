package ast

type Object struct {
	expressionNode
	spreads []Spread
	entries []EntryNode
}

func (o Object) Spreads() []Spread {
	return o.spreads
}

func (o Object) Entries() []EntryNode {
	return o.entries
}

func NewObject(spreads []Spread, entries []EntryNode) Object {
	return Object{
		spreads: spreads,
		entries: entries,
	}
}

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
