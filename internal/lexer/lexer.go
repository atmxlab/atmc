package lexer

import (
	"github.com/atmxlab/atmcfg/internal/types"
	"github.com/atmxlab/atmcfg/internal/types/token"
)

type Lexer struct {
	tokens   []token.Token
	position types.Position
}

func (l *Lexer) Tokenize(input string) ([]token.Token, error) {
	currentInput := ""
	for _, r := range input {
		l.position.IncrPos()
		l.position.IncrColumn()
		currentInput += string(r)
		switch true {
		case token.EOL.Regexp().MatchString(currentInput):
			l.position.IncrLine()
			l.position.ResetColumn()
		}
	}
	for len(input) > 0 {
		for _, t := range token.OrderedTokenTypes() {
			indexes := t.Regexp().FindStringIndex(input)
			if len(indexes) == 0 {
				continue
			}

			start := indexes[0]
			end := indexes[1]

			value := input[:end]
			input = input[start:]

			l.position.AddPos(uint(end))
			l.position.AddColumn(uint(end))

			switch t {
			case token.EOL:
				l.position.IncrLine()
				l.position.ResetColumn()
				continue
			case token.WS:
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
