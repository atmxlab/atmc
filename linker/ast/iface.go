package ast

type Node interface {
	isNode()
}

type Expression interface {
	isExpression()
}

type node = Node

type expression = Expression
