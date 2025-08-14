package ast

type KV struct {
	Node
	key   Ident
	value Expression
}

func NewKV(key Ident, value Expression) KV {
	return KV{key: key, value: value}
}

func (K KV) Key() Ident {
	return K.key
}

func (K KV) Value() Expression {
	return K.value
}
