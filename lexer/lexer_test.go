package lexer_test

import (
	"testing"

	"github.com/atmxlab/atmc/lexer"
	"github.com/atmxlab/atmc/pkg/collect"
	"github.com/atmxlab/atmc/types"
	token2 "github.com/atmxlab/atmc/types/token"
	"github.com/stretchr/testify/require"
)

func TestLexer_Tokenize_Tokens(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name           string
		input          string
		expectedTokens []token2.Token
	}{
		{
			name: "import",
			input: `
				common ./common.atmx`,
			expectedTokens: []token2.Token{
				token2.New(
					token2.Ident,
					"common",
					types.NewLocation(
						types.NewPosition(2, 4, 5),
						types.NewPosition(2, 10, 11),
					),
				),
				token2.New(
					token2.Path,
					"./common.atmx",
					types.NewLocation(
						types.NewPosition(2, 11, 12),
						types.NewPosition(2, 24, 25),
					),
				),
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			l := lexer.New()

			tokens, err := l.Tokenize(tc.input)
			require.NoError(t, err)
			require.Equal(t, tc.expectedTokens, tokens)
		})
	}
}

func TestLexer_Tokenize_LexerLocation(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name        string
		input       string
		expectedPos types.Location
	}{
		{
			name: "2, 36, 37",
			input: `
				import common from ./common.atmx`,
			expectedPos: types.NewLocation(
				types.NewPosition(2, 23, 24),
				types.NewPosition(2, 36, 37),
			),
		},
		{
			name: "3, 0, 38",
			input: `
				import common from ./common.atmx
`,
			expectedPos: types.NewLocation(
				types.NewPosition(3, 0, 38),
				types.NewPosition(3, 0, 38),
			),
		},
		{
			name: "7, 0, 42",
			input: `


				import common from ./common.atmx


`,
			expectedPos: types.NewLocation(
				types.NewPosition(7, 0, 42),
				types.NewPosition(7, 0, 42),
			),
		},
		{
			name: "7, 0, 42",
			input: `


				import common from ./common.atmx


													import

			`,
			expectedPos: types.NewLocation(
				types.NewPosition(9, 2, 65),
				types.NewPosition(9, 3, 66),
			),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			l := lexer.New()

			_, err := l.Tokenize(tc.input)
			require.NoError(t, err)
			require.Equal(t, tc.expectedPos, l.Location())
		})
	}
}

func TestLexer_Tokenize_TokenTypes(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name          string
		input         string
		expectedTypes []token2.Type
		hasError      bool
	}{
		{
			name:  "simple import",
			input: `common ./common.atmx`,
			expectedTypes: []token2.Type{
				token2.Ident,
				token2.Path,
			},
		},
		{
			name:          "comment",
			input:         `// comment`,
			expectedTypes: []token2.Type{},
		},
		{
			name:  "comment",
			input: `common... // comment`,
			expectedTypes: []token2.Type{
				token2.Ident,
				token2.Spread,
			},
		},
		{
			name:  "nested import import",
			input: `common /dir1/dir2/common.atmx`,
			expectedTypes: []token2.Type{
				token2.Ident,
				token2.Path,
			},
		},
		{
			name:  "nested import with prev dir import",
			input: `common /dir1/dir2/../common.atmx`,
			expectedTypes: []token2.Type{
				token2.Ident,
				token2.Path,
			},
		},
		{
			name:          "import with invalid path",
			input:         `common /dir1/dir2/..1#/common.atmx`,
			expectedTypes: []token2.Type{},
			hasError:      true,
		},
		{
			name:  "ident",
			input: `ident1 ident2 ident_3 Ident4 IDENT5 IDENT_______6 IdEnT7`,
			expectedTypes: []token2.Type{
				token2.Ident,
				token2.Ident,
				token2.Ident,
				token2.Ident,
				token2.Ident,
				token2.Ident,
				token2.Ident,
			},
			hasError: false,
		},
		{
			name:  "simple object",
			input: `{ key: "str value" }`,
			expectedTypes: []token2.Type{
				token2.LBrace,
				token2.Ident,
				token2.Colon,
				token2.String,
				token2.RBrace,
			},
			hasError: false,
		},
		{
			name:  "object with key start with number",
			input: `{ 1key: "str value" }`,
			expectedTypes: []token2.Type{
				token2.LBrace,
				token2.Int,
				token2.Ident,
				token2.Colon,
				token2.String,
				token2.RBrace,
			},
			hasError: false,
		},
		{
			name:  "simple array",
			input: `[123, 123, 124]`,
			expectedTypes: []token2.Type{
				token2.LBracket,
				token2.Int,
				token2.Int,
				token2.Int,
				token2.RBracket,
			},
			hasError: false,
		},
		{
			name:  "one line object",
			input: `{key1: 123 key2: 123.123 key3: "test string 12313 123.123" key_4: "test string 12313 \"test\" 123.123" kEy5: true Key6: false key7: {key8: 123} key9: ["123", "321"]}`,
			expectedTypes: []token2.Type{
				token2.LBrace,
				token2.Ident,
				token2.Colon,
				token2.Int,
				token2.Ident,
				token2.Colon,
				token2.Float,
				token2.Ident,
				token2.Colon,
				token2.String,
				token2.Ident,
				token2.Colon,
				token2.String,
				token2.Ident,
				token2.Colon,
				token2.Bool,
				token2.Ident,
				token2.Colon,
				token2.Bool,
				token2.Ident,
				token2.Colon,
				token2.LBrace,
				token2.Ident,
				token2.Colon,
				token2.Int,
				token2.RBrace,
				token2.Ident,
				token2.Colon,
				token2.LBracket,
				token2.String,
				token2.String,
				token2.RBracket,
				token2.RBrace,
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
			expectedTypes: []token2.Type{
				token2.Ident,
				token2.Path,

				token2.LBrace,

				token2.Ident,
				token2.Spread,

				token2.Ident,
				token2.Dot,
				token2.Ident,
				token2.Spread,

				token2.Ident,
				token2.Dot,
				token2.Ident,
				token2.Dot,
				token2.Ident,
				token2.Spread,

				token2.Ident,
				token2.Colon,
				token2.Ident,
				token2.Dot,
				token2.Ident,
				token2.Dot,
				token2.Ident,

				token2.Ident,
				token2.Colon,

				token2.LBrace,
				token2.Ident,
				token2.Dot,
				token2.Ident,
				token2.Spread,

				token2.Ident,
				token2.Colon,
				token2.Int,

				token2.Ident,
				token2.Dot,
				token2.Ident,
				token2.Dot,
				token2.Ident,
				token2.Spread,
				token2.RBrace,

				token2.Ident,
				token2.Colon,

				token2.LBracket,

				token2.Ident,
				token2.Dot,
				token2.Ident,
				token2.Spread,

				token2.Int,

				token2.String,

				token2.RBracket,

				token2.RBrace,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			l := lexer.New()

			tokens, err := l.Tokenize(tc.input)
			if tc.hasError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}

			tokenTypes := collect.Map(tokens, func(item token2.Token) token2.Type {
				return item.Type()
			})
			require.Equal(t, tc.expectedTypes, tokenTypes)
		})
	}
}
