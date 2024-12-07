package ast

import "github.com/atmxlab/atmcfg/internal/types"

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

func NewEntryNode(key Ident, value Expression, loc types.Location) EntryNode {
	e := EntryNode{key: key, value: value}
	e.loc = loc

	return e
}
