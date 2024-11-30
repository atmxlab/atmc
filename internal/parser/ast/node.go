package ast

import "github.com/atmxlab/atmcfg/internal/types"

type node struct {
	pos types.Position
}

func (n node) Pos() uint {
	return n.pos.Pos()
}

type entryNode struct {
	node
}

func (entryNode) entryNode() {}

type identNode struct {
	node
}

func (identNode) identNode() {}
