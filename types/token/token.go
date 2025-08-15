package token

import (
	"github.com/atmxlab/atmc/types"
)

type Token struct {
	t        Type
	value    Value
	location types.Location
}

func New(t Type, value Value, location types.Location) Token {
	return Token{t: t, value: value, location: location}
}

func (t Token) Type() Type {
	return t.t
}

func (t Token) Value() Value {
	return t.value
}

func (t Token) Location() types.Location {
	return t.location
}

type Value string

func (v Value) String() string {
	return string(v)
}
