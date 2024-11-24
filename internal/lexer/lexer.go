package lexer

import (
	"github.com/atmxlab/atmcfg/internal/types"
	"github.com/atmxlab/atmcfg/internal/types/token"
)

type Lexer struct {
	input    string
	tokens   []token.Token
	position *types.Position
}

func (l *Lexer) Position() *types.Position {
	return l.position
}

func New(input string) *Lexer {
	return &Lexer{
		input:    input,
		tokens:   make([]token.Token, 0),
		position: types.NewInitialPosition(),
	}
}

func (l *Lexer) Tokenize() ([]token.Token, error) {
	orderedTokenTypes := token.OrderedTokenTypes()

	for len(l.input) > 0 {
		matched := false

		for _, t := range orderedTokenTypes {
			value, exists := l.find(t)
			if !exists {
				continue
			}

			matched = true

			switch t {
			case token.WS:
				// Ничего не делаем. Просто игнорируем пробелы.
			case token.EOL:
				l.position.IncrLine()
				l.position.ResetColumn()
			default:
				l.addToken(t, value)
			}

			break
		}

		if !matched {
			return nil, unexpectedTokenError(l.position)
		}
	}

	return l.tokens, nil
}

func (l *Lexer) addToken(t token.Type, value string) {
	tok := token.New(
		t,
		token.Value(value),
		l.position.Clone(),
	)

	l.tokens = append(l.tokens, tok)
}

func (l *Lexer) find(t token.Type) (value string, exists bool) {
	indexes := t.Regexp().FindStringIndex(l.input)
	if len(indexes) < 2 {
		return "", false
	}

	end := indexes[1]

	value = l.input[:end]
	l.input = l.input[end:]

	l.position.AddPos(uint(end))
	l.position.AddColumn(uint(end))

	return value, true
}
