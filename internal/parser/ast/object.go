package ast

type Object struct {
	entryNode
	params []KeyValue
}

func NewObject(params []KeyValue) Object {
	return Object{params: params}
}

type Key struct {
	identNode
	string
}

func NewKey(string string) Key {
	return Key{string: string}
}

type KeyValue struct {
	node
	Key   Key
	Value Entry // object | array | literal
}

func NewKeyValue(key Key, value Entry) KeyValue {
	return KeyValue{Key: key, Value: value}
}
