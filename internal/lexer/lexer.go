package lexer

import (
	"github.com/atmxlab/atmcfg/internal/types"
	"github.com/atmxlab/atmcfg/internal/types/token"
)

type Lexer struct {
	input    string
	tokens   []token.Token
	position types.Position
}

func (l *Lexer) Tokenize() ([]token.Token, error) {
	currentInput := ""
	for _, r := range l.input {
		l.position.IncrPos()
		l.position.IncrColumn()
		currentInput += string(r)
		switch true {
		case token.EOL.Regexp().MatchString(currentInput):
			l.position.IncrLine()
			l.position.ResetColumn()
		}
	}
	for len(l.input) > 0 {
		for _, t := range token.OrderedTokenTypes() {
			value, exists := l.Find(t)
			if !exists {
				continue
			}

			switch t {
			case token.WS:
				continue
			case token.EOL:
				l.position.IncrLine()
				l.position.ResetColumn()
				continue
			default:
				l.AddToken(t, value)
			}
		}
	}

	return l.tokens, nil
}

func (l *Lexer) AddToken(t token.Type, value string) {
	tok := token.NewToken(
		t,
		token.Value(value),
		l.position.Clone(),
	)

	l.tokens = append(
		l.tokens,
		tok,
	)
}

func (l *Lexer) Find(t token.Type) (value string, exists bool) {
	indexes := t.Regexp().FindStringIndex(l.input)
	if len(indexes) == 0 {
		return "", false
	}

	start := indexes[0]
	end := indexes[1]

	value = l.input[:end]
	l.input = l.input[start:]

	l.position.AddPos(uint(end))
	l.position.AddColumn(uint(end))

	return value, true
}
