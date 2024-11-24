package lexer_test

import (
	"testing"

	"github.com/atmxlab/atmcfg/internal/lexer"
	"github.com/atmxlab/atmcfg/internal/types"
	"github.com/atmxlab/atmcfg/internal/types/token"
	"github.com/atmxlab/atmcfg/pkg/collect"
	"github.com/stretchr/testify/require"
)

func TestLexer_Tokenize_TokenTypes(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name          string
		input         string
		expectedTypes []token.Type
		hasError      bool
	}{
		{
			name:  "simple import",
			input: `import common from ./common.atmx`,
			expectedTypes: []token.Type{
				token.Import,
				token.Ident,
				token.From,
				token.Path,
			},
		},
		{
			name:  "nested import import",
			input: `import common from /dir1/dir2/common.atmx`,
			expectedTypes: []token.Type{
				token.Import,
				token.Ident,
				token.From,
				token.Path,
			},
		},
		{
			name:  "nested import with prev dir import",
			input: `import common from /dir1/dir2/../common.atmx`,
			expectedTypes: []token.Type{
				token.Import,
				token.Ident,
				token.From,
				token.Path,
			},
		},
		{
			name:          "import with invalid path",
			input:         `import common from /dir1/dir2/..1#/common.atmx`,
			expectedTypes: []token.Type{},
			hasError:      true,
		},
		{
			name:          "import with invalid path",
			input:         `import common from /dir1/dir2/..1#/common.atmx`,
			expectedTypes: []token.Type{},
			hasError:      true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			l := lexer.New(tc.input)

			tokens, err := l.Tokenize()
			if tc.hasError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}

			tokenTypes := collect.Map(tokens, func(item token.Token) token.Type {
				return item.Type()
			})
			require.Equal(t, tc.expectedTypes, tokenTypes)
		})
	}
}

func TestLexer_Tokenize_Tokens(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name           string
		input          string
		expectedTokens []token.Token
	}{
		{
			name: "import",
			input: `
				import common from ./common.atmx`,
			expectedTokens: []token.Token{
				token.New(
					token.Import,
					"import",
					types.NewPosition(2, 10, 11),
				),
				token.New(
					token.Ident,
					"common",
					types.NewPosition(2, 17, 18),
				),
				token.New(
					token.From,
					"from",
					types.NewPosition(2, 22, 23),
				),
				token.New(
					token.Path,
					"./common.atmx",
					types.NewPosition(2, 36, 37),
				),
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			l := lexer.New(tc.input)

			tokens, err := l.Tokenize()
			require.NoError(t, err)
			require.Equal(t, tc.expectedTokens, tokens)
		})
	}
}

func TestLexer_Tokenize_LexerPosition(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name        string
		input       string
		expectedPos *types.Position
	}{
		{
			name: "2, 36, 37",
			input: `
				import common from ./common.atmx`,
			expectedPos: types.NewPosition(2, 36, 37),
		},
		{
			name: "3, 0, 38",
			input: `
				import common from ./common.atmx
`,
			expectedPos: types.NewPosition(3, 0, 38),
		},
		{
			name: "7, 0, 42",
			input: `


				import common from ./common.atmx


`,
			expectedPos: types.NewPosition(7, 0, 42),
		},
		{
			name: "7, 0, 42",
			input: `


				import common from ./common.atmx

			
													import
			
			`,
			expectedPos: types.NewPosition(9, 3, 72),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			l := lexer.New(tc.input)

			_, err := l.Tokenize()
			require.NoError(t, err)
			require.Equal(t, tc.expectedPos, l.Position())
		})
	}
}
