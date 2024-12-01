package tokenmover

import (
	"github.com/atmxlab/atmcfg/internal/types/token"
)

type TokenMover struct {
	tokens []token.Token
	pos    int
}

func New(tokens []token.Token) *TokenMover {
	return &TokenMover{tokens: tokens}
}

func (l *TokenMover) IsEmpty() bool {
	return l.pos >= len(l.tokens)
}

func (l *TokenMover) Next() {
	l.pos++
}
func (l *TokenMover) Prev() {
	l.pos--
}

func (l *TokenMover) Token() token.Token {
	if l.pos >= len(l.tokens) {
		panic("out of range")
	}

	return l.tokens[l.pos]
}
