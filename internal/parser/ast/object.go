package ast

type Object struct {
	entryNode
	params []KeyValue
}

type Key struct {
	identNode
	string
}

type KeyValue struct {
	Key   Key
	Value Entry // object | array | literal
}
