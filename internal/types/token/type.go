package token

import (
	"fmt"
	"regexp"
)

type Type uint

func (t Type) String() string {
	switch t {
	case WS:
		return "white space"
	case EOL:
		return "end of line"
	case LBrace:
		return "left brace"
	case RBrace:
		return "right brace"
	case LBracket:
		return "left bracket"
	case RBracket:
		return "right bracket"
	case Spread:
		return "spread"
	case Comma:
		return "comma"
	case Dot:
		return "comma"
	case Colon:
		return "colon"
	case Int:
		return "int"
	case Float:
		return "float"
	case String:
		return "string"
	case Bool:
		return "bool"
	case Ident:
		return "ident"
	case Path:
		return "path"
	case Dollar:
		return "dollar"
	default:
		return fmt.Sprintf("undefined token type: %d", t)
	}
}

const (
	WS Type = iota
	EOL
	Dollar
	LBrace
	RBrace
	LBracket
	RBracket
	Spread
	Comma
	Colon
	Int
	Float
	String
	Bool
	Ident
	Path
	Dot
)

var typeRegexps = map[Type]*regexp.Regexp{
	WS:       regexp.MustCompile("^[ \\t\\r]"),
	EOL:      regexp.MustCompile("^\\n"),
	LBrace:   regexp.MustCompile("^\\{"),
	RBrace:   regexp.MustCompile("^}"),
	LBracket: regexp.MustCompile("^\\["),
	RBracket: regexp.MustCompile("^]"),
	Spread:   regexp.MustCompile("^\\.\\.\\."),
	Comma:    regexp.MustCompile("^,"),
	Dot:      regexp.MustCompile("^\\."),
	Colon:    regexp.MustCompile("^:"),
	Int:      regexp.MustCompile("^[-+]?[0-9]+"),
	Float:    regexp.MustCompile("^[-+]?[0-9]+(\\.[0-9]+)"),
	Bool:     regexp.MustCompile("^(true|false)\\b"),
	String:   regexp.MustCompile(`^"(?:[^\\"]|\\.|\\\\)*"`),
	Ident:    regexp.MustCompile("^[a-zA-Z_][a-zA-Z0-9_]*"),
	Path:     regexp.MustCompile("^(?:/|\\./)[a-zA-Z0-9._/-]+"),
	Dollar:   regexp.MustCompile("^\\$"),
}

func (t Type) Regexp() *regexp.Regexp {
	rgxp, ok := typeRegexps[t]
	if !ok {
		panic("unknown token type")
	}

	return rgxp
}

func OrderedTokenTypes() []Type {
	return []Type{
		WS,
		EOL,
		String,
		Path,
		Bool,
		Float,
		Int,
		LBrace,
		RBrace,
		LBracket,
		RBracket,
		Spread,
		Comma,
		Dot,
		Dollar,
		Colon,
		Ident,
	}
}
