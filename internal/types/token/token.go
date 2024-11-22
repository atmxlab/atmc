package token

import (
	"regexp"

	"github.com/atmxlab/atmcfg/internal/types"
)

type Token struct {
	_type    Type
	value    Value
	position types.Position
}

func NewToken(_type Type, value Value, position types.Position) Token {
	return Token{_type: _type, value: value, position: position}
}

func (t Token) Type() Type {
	return t._type
}

func (t Token) Value() Value {
	return t.value
}

func (t Token) Position() types.Position {
	return t.position
}

type Value string

type Type uint

const (
	Import Type = iota
	From
	LBrace
	RBrace
	LParen
	RParen
	LBracket
	RBracket
	Spread
	Colon
	Int
	Float
	String
	Bool
	Ident
	WS
	EOL
	EOF
)

func (t Type) Regexp() *regexp.Regexp {
	switch t {
	case Import:
		return regexp.MustCompile("")
	default:
		panic("unknown token type")
	}
}

func OrderedTokenTypes() []Type {
	return []Type{
		Import,
		From,
	}
}
