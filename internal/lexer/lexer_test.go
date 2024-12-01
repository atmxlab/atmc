package lexer_test

import (
	"testing"

	"github.com/atmxlab/atmcfg/internal/lexer"
	"github.com/atmxlab/atmcfg/internal/types"
	"github.com/atmxlab/atmcfg/internal/types/token"
	"github.com/atmxlab/atmcfg/pkg/collect"
	"github.com/stretchr/testify/require"
)

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
				common ./common.atmx`,
			expectedTokens: []token.Token{
				token.New(
					token.Ident,
					"common",
					types.NewPosition(2, 10, 11),
				),
				token.New(
					token.Path,
					"./common.atmx",
					types.NewPosition(2, 24, 25),
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
			input: `common ./common.atmx`,
			expectedTypes: []token.Type{
				token.Ident,
				token.Path,
			},
		},
		{
			name:  "nested import import",
			input: `common /dir1/dir2/common.atmx`,
			expectedTypes: []token.Type{
				token.Ident,
				token.Path,
			},
		},
		{
			name:  "nested import with prev dir import",
			input: `common /dir1/dir2/../common.atmx`,
			expectedTypes: []token.Type{
				token.Ident,
				token.Path,
			},
		},
		{
			name:          "import with invalid path",
			input:         `common /dir1/dir2/..1#/common.atmx`,
			expectedTypes: []token.Type{},
			hasError:      true,
		},
		{
			name:  "ident",
			input: `ident1 ident2 ident_3 Ident4 IDENT5 IDENT_______6 IdEnT7`,
			expectedTypes: []token.Type{
				token.Ident,
				token.Ident,
				token.Ident,
				token.Ident,
				token.Ident,
				token.Ident,
				token.Ident,
			},
			hasError: false,
		},
		{
			name:  "simple object",
			input: `{ key: "str value" }`,
			expectedTypes: []token.Type{
				token.LBrace,
				token.Ident,
				token.Colon,
				token.String,
				token.RBrace,
			},
			hasError: false,
		},
		{
			name:  "object with key start with number",
			input: `{ 1key: "str value" }`,
			expectedTypes: []token.Type{
				token.LBrace,
				token.Int,
				token.Ident,
				token.Colon,
				token.String,
				token.RBrace,
			},
			hasError: false,
		},
		{
			name:  "simple array",
			input: `[123, 123, 124]`,
			expectedTypes: []token.Type{
				token.LBracket,
				token.Int,
				token.Comma,
				token.Int,
				token.Comma,
				token.Int,
				token.RBracket,
			},
			hasError: false,
		},
		{
			name:  "one line object",
			input: `{key1: 123 key2: 123.123 key3: "test string 12313 123.123" key_4: "test string 12313 \"test\" 123.123" kEy5: true Key6: false key7: {key8: 123} key9: ["123", "321"]}`,
			expectedTypes: []token.Type{
				token.LBrace,
				token.Ident,
				token.Colon,
				token.Int,
				token.Ident,
				token.Colon,
				token.Float,
				token.Ident,
				token.Colon,
				token.String,
				token.Ident,
				token.Colon,
				token.String,
				token.Ident,
				token.Colon,
				token.Bool,
				token.Ident,
				token.Colon,
				token.Bool,
				token.Ident,
				token.Colon,
				token.LBrace,
				token.Ident,
				token.Colon,
				token.Int,
				token.RBrace,
				token.Ident,
				token.Colon,
				token.LBracket,
				token.String,
				token.Comma,
				token.String,
				token.RBracket,
				token.RBrace,
			},
			hasError: false,
		},
		{
			name: "with import, spread, var, array, literals",
			input: `
common ./common.atmc

{
    common...
    common.nested...
    common.nested.nested2...

    nested: common.nested1.nested2
    spreadObj: {
        common.nested1...
        ident: 123
        common.nested1.nested2...
    }

    spreadArr: [
        common.nested1...
        123
		"test"
    ]
}
`,
			expectedTypes: []token.Type{
				token.Ident,
				token.Path,

				token.LBrace,

				token.Ident,
				token.Spread,

				token.Ident,
				token.Dot,
				token.Ident,
				token.Spread,

				token.Ident,
				token.Dot,
				token.Ident,
				token.Dot,
				token.Ident,
				token.Spread,

				token.Ident,
				token.Colon,
				token.Ident,
				token.Dot,
				token.Ident,
				token.Dot,
				token.Ident,

				token.Ident,
				token.Colon,

				token.LBrace,
				token.Ident,
				token.Dot,
				token.Ident,
				token.Spread,

				token.Ident,
				token.Colon,
				token.Int,

				token.Ident,
				token.Dot,
				token.Ident,
				token.Dot,
				token.Ident,
				token.Spread,
				token.RBrace,

				token.Ident,
				token.Colon,

				token.LBracket,

				token.Ident,
				token.Dot,
				token.Ident,
				token.Spread,

				token.Int,

				token.String,

				token.RBracket,

				token.RBrace,
			},
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
