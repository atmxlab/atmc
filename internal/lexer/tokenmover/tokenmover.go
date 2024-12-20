package tokenmover

import (
	"fmt"

	"github.com/atmxlab/atmcfg/internal/types/token"
)

type TokenMover struct {
	tokens          []token.Token
	pos             int
	savePointsStack []int
}

func (l *TokenMover) SavePoint() {
	l.savePointsStack = append(l.savePointsStack, l.pos)
}

func (l *TokenMover) RemoveSavePoint() {
	if len(l.savePointsStack) == 0 {
		return
	}

	l.savePointsStack = l.savePointsStack[:len(l.savePointsStack)-1]
}

func (l *TokenMover) ReturnToSavePoint() {
	if len(l.savePointsStack) == 0 {
		return
	}

	l.pos = l.savePointsStack[len(l.savePointsStack)-1]
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
		panic(fmt.Sprintf("token pos out of range [pos: %d; len: %d]", l.pos, len(l.tokens)))
	}

	return l.tokens[l.pos]
}
