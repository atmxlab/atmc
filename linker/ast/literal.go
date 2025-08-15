package ast

type literal[T comparable] struct {
	node
	expression
	value T
}

func (l literal[T]) Value() T {
	return l.value
}

type Int = literal[int64]

func NewInt(i int64) Int {
	return Int{value: i}
}

type Float = literal[float64]

func NewFloat(i float64) Float {
	return Float{value: i}
}

type String = literal[string]

func NewString(i string) String {
	return String{value: i}
}

type Bool = literal[bool]

func NewBool(i bool) Bool {
	return Bool{value: i}
}
