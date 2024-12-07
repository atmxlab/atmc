package lexer

import (
	"github.com/atmxlab/atmcfg/internal/types"
	"github.com/atmxlab/atmcfg/internal/types/token"
)

type Lexer struct {
	input    string
	tokens   []token.Token
	location types.Location
}

func (l *Lexer) Location() types.Location {
	return l.location
}

func New(input string) *Lexer {
	return &Lexer{
		input:    input,
		tokens:   make([]token.Token, 0),
		location: types.NewInitialLocation(),
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
			case token.WS, token.Comma:
				// Ничего не делаем. Просто игнорируем пробелы.
			case token.EOL:
				l.location = l.location.SetEnd(
					l.location.End().
						IncrLine().
						ResetColumn(),
				)
				l.location = l.location.SetStart(
					l.location.End(),
				)
			default:
				l.addToken(t, value)
			}

			break
		}

		if !matched {
			return nil, unexpectedTokenError(l.location.End())
		}
	}

	return l.tokens, nil
}

func (l *Lexer) addToken(t token.Type, value string) {
	tok := token.New(
		t,
		token.Value(value),
		l.location,
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

	l.location = l.location.SetStart(
		l.location.End(),
	)
	l.location = l.location.SetEnd(
		l.location.End().
			AddPos(uint(end)).
			AddColumn(uint(end)),
	)

	return value, true
}
