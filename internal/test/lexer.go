package test

import (
	"testing"

	"github.com/atmxlab/atmcfg/internal/types/token"
)

type Lexer struct {
	t      *testing.T
	tokens []token.Token
	pos    int
}

func NewLexer(t *testing.T, tokens []token.Token) *Lexer {
	return &Lexer{t: t, tokens: tokens}
}

func (l *Lexer) IsEmpty() bool {
	return l.pos >= len(l.tokens)
}

func (l *Lexer) Next() {
	l.pos++
}
func (l *Lexer) Prev() {
	l.pos--
}

func (l *Lexer) Token() token.Token {
	if l.pos >= len(l.tokens) {
		l.t.Fatalf("token pos out of range [pos: %d; len: %d]", l.pos, len(l.tokens))
	}

	return l.tokens[l.pos]
}
