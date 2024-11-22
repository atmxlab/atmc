package token

import (
	"regexp"
)

type Type uint

const (
	WS Type = iota
	EOL
	Import
	From
	LBrace
	RBrace
	LBracket
	RBracket
	Spread
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
	Import:   regexp.MustCompile("^import\\b"),
	From:     regexp.MustCompile("^from\\b"),
	LBrace:   regexp.MustCompile("^\\{"),
	RBrace:   regexp.MustCompile("^}"),
	LBracket: regexp.MustCompile("^\\["),
	RBracket: regexp.MustCompile("^]"),
	Spread:   regexp.MustCompile("^\\.\\.\\."),
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
		Import,
		From,
		Bool,
		Float,
		Int,
		LBrace,
		RBrace,
		LBracket,
		RBracket,
		Spread,
		Colon,
		Ident,
	}
}
