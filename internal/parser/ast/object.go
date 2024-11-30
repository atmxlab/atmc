package ast

type Object struct {
	entryNode
	keyValues []KeyValue
}

func NewObject(keyValues []KeyValue) Object {
	return Object{keyValues: keyValues}
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
