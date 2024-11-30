package token

import (
	"fmt"
	"regexp"
)

type Type uint

func (t Type) String() string {
	switch t {

	case WS:
		return "WS"
	case EOL:
		return "EOL"
	case LBrace:
		return "LBrace"
	case RBrace:
		return "RBrace"
	case LBracket:
		return "LBracket"
	case RBracket:
		return "RBracket"
	case Spread:
		return "Spread"
	case Comma:
		return "Comma"
	case Colon:
		return "Colon"
	case Int:
		return "Int"
	case Float:
		return "Float"
	case String:
		return "String"
	case Bool:
		return "Bool"
	case Ident:
		return "Ident"
	case Path:
		return "Path"
	default:
		return fmt.Sprintf("undefined token type: %d", t)
	}
}

const (
	WS Type = iota
	EOL
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
	Colon:    regexp.MustCompile("^:"),
	Int:      regexp.MustCompile("^[-+]?[0-9]+"),
	Float:    regexp.MustCompile("^[-+]?[0-9]+(\\.[0-9]+)"),
	Bool:     regexp.MustCompile("^(true|false)\\b"),
	String:   regexp.MustCompile(`^"(?:[^\\"]|\\.|\\\\)*"`),
	Ident:    regexp.MustCompile("^[a-zA-Z_][a-zA-Z0-9_]*"),
	Path:     regexp.MustCompile("^(?:/|\\./)[a-zA-Z0-9._/-]+"),
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
		Colon,
		Ident,
	}
}
