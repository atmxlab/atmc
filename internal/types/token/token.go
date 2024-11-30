package token

import (
	"github.com/atmxlab/atmcfg/internal/types"
)

type Token struct {
	t        Type
	value    Value
	position *types.Position
}

func New(t Type, value Value, position *types.Position) Token {
	return Token{t: t, value: value, position: position}
}

func (t Token) Type() Type {
	return t.t
}

func (t Token) Value() Value {
	return t.value
}

func (t Token) Position() *types.Position {
	return t.position
}

type Value string

func (v Value) String() string {
	return string(v)
}
