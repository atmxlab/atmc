package token_test

import (
	"testing"

	"github.com/atmxlab/atmcfg/internal/types/token"
	"github.com/stretchr/testify/require"
)

func TestType_EOL_Regexp(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name     string
		input    string
		expected []int
	}{
		{
			name: "start with \\n",
			input: `
				{
					test: "123"
				}
				
			`,
			expected: []int{0, 1},
		},
		{
			name: "not start with \\n",
			input: ` {
					test: "123"
				}
				
			`,
			expected: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			indexes := token.EOL.Regexp().FindStringIndex(tc.input)
			require.Equal(t, tc.expected, indexes)
		})
	}
}

func TestType_WS_Regexp(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name     string
		input    string
		expected []int
	}{
		{
			name: "start with spaces",
			input: `    {
					test: "123"
				}
				
			`,
			expected: []int{0, 1},
		},
		{
			name:     "start with tabs",
			input:    "\t \t \n test: 123",
			expected: []int{0, 1},
		},
		{
			name:     "start with \\r",
			input:    "\r \t \n test: 123",
			expected: []int{0, 1},
		},
		{
			name:     "not start with whitespace",
			input:    "{\r \t \n test: 123",
			expected: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			indexes := token.WS.Regexp().FindStringIndex(tc.input)
			require.Equal(t, tc.expected, indexes)
		})
	}
}

func TestType_Import_Regexp(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name     string
		input    string
		expected []int
	}{
		{
			name:     "start with",
			input:    `import test from ./test.atmx`,
			expected: []int{0, 6},
		},
		{
			name:     "not start with",
			input:    " import test from ./test.atmx",
			expected: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			indexes := token.Import.Regexp().FindStringIndex(tc.input)
			require.Equal(t, tc.expected, indexes)
		})
	}
}

func TestType_From_Regexp(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name     string
		input    string
		expected []int
	}{
		{
			name:     "start with",
			input:    `from ./test.atmx`,
			expected: []int{0, 4},
		},
		{
			name:     "not start with",
			input:    "test from ./test.atmx",
			expected: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			indexes := token.From.Regexp().FindStringIndex(tc.input)
			require.Equal(t, tc.expected, indexes)
		})
	}
}

func TestType_LBrace_Regexp(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name     string
		input    string
		expected []int
	}{
		{
			name:     "start with",
			input:    `{sdawdsd:dwads""`,
			expected: []int{0, 1},
		},
		{
			name:     "not start with",
			input:    "}{sdawd{}{}{{[]][]",
			expected: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			indexes := token.LBrace.Regexp().FindStringIndex(tc.input)
			require.Equal(t, tc.expected, indexes)
		})
	}
}

func TestType_RBrace_Regexp(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name     string
		input    string
		expected []int
	}{
		{
			name:     "start with",
			input:    `}sdawdsd:dwads""`,
			expected: []int{0, 1},
		},
		{
			name:     "not start with",
			input:    "{{sdawd{}{}{{[]][]",
			expected: nil,
		},
		{
			name: "not start with and start from new lind",
			input: `
}
`,
			expected: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			indexes := token.RBrace.Regexp().FindStringIndex(tc.input)
			require.Equal(t, tc.expected, indexes)
		})
	}
}

func TestType_LBracket_Regexp(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name     string
		input    string
		expected []int
	}{
		{
			name:     "start with",
			input:    `[test{}[][]231:sda2131from||||import`,
			expected: []int{0, 1},
		},
		{
			name:     "not start with",
			input:    `][][][test{}[][]231:sda2131from||||import`,
			expected: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			indexes := token.LBracket.Regexp().FindStringIndex(tc.input)
			require.Equal(t, tc.expected, indexes)
		})
	}
}

func TestType_RBracket_Regexp(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name     string
		input    string
		expected []int
	}{
		{
			name:     "start with",
			input:    `]test{}[][]231:sda2131from||||import`,
			expected: []int{0, 1},
		},
		{
			name:     "not start with",
			input:    `[[][][test{}[][]231:sda2131from||||import`,
			expected: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			indexes := token.RBracket.Regexp().FindStringIndex(tc.input)
			require.Equal(t, tc.expected, indexes)
		})
	}
}

func TestType_Spread_Regexp(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name     string
		input    string
		expected []int
	}{
		{
			name:     "start with",
			input:    `...]test{}[][]231...:sda2131from||||import`,
			expected: []int{0, 3},
		},
		{
			name:     "not start with",
			input:    `[...[][][test{}[][]231:sda2131from||||import...`,
			expected: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			indexes := token.Spread.Regexp().FindStringIndex(tc.input)
			require.Equal(t, tc.expected, indexes)
		})
	}
}
func TestType_Colon_Regexp(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name     string
		input    string
		expected []int
	}{
		{
			name:     "start with",
			input:    `::::...]test{}[][]231...:sda2131from||||import`,
			expected: []int{0, 1},
		},
		{
			name:     "not start with",
			input:    `[:...[][][test{}[][]231:::sda2131from||||import...`,
			expected: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			indexes := token.Colon.Regexp().FindStringIndex(tc.input)
			require.Equal(t, tc.expected, indexes)
		})
	}
}

func TestType_Int_Regexp(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name     string
		input    string
		expected []int
	}{
		{
			name:     "start with",
			input:    `123::::...]test{}[][]231...:sda2131from||||import`,
			expected: []int{0, 3},
		},

		{
			name:     "start with but negative",
			input:    `-123::::...]test{}[][]231...:sda2131from||||import {test: 123}`,
			expected: []int{0, 4},
		},
		{
			name:     "not start with",
			input:    `[:...[][][test{}[][]231:::sda2131from||||import...`,
			expected: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			indexes := token.Int.Regexp().FindStringIndex(tc.input)
			require.Equal(t, tc.expected, indexes)
		})
	}
}

func TestType_Float_Regexp(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name     string
		input    string
		expected []int
	}{
		{
			name:     "start with",
			input:    `123.123::::...]test{}[][]231...:sda2131from||||import`,
			expected: []int{0, 7},
		},
		{
			name:     "start with int",
			input:    `123::::...]test{}[][]231...:sda2131from||||import 123.123`,
			expected: nil,
		},
		{
			name:     "start with but negative",
			input:    `-123.123::::...]test{}[][]231...:sda2131from||||import {test: 123}`,
			expected: []int{0, 8},
		},
		{
			name:     "not start with",
			input:    `[:...[][][test{}[][]231:::sda2131from||||import... 123.123`,
			expected: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			indexes := token.Float.Regexp().FindStringIndex(tc.input)
			require.Equal(t, tc.expected, indexes)
		})
	}
}

func TestType_Bool_Regexp(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name     string
		input    string
		expected []int
	}{
		{
			name:     "start with true",
			input:    `true 123.123::::...]test{}[][]231...:sda2131from||||import`,
			expected: []int{0, 4},
		},
		{
			name:     "start with true without space",
			input:    `true123::::...]test{}[][]231...:sda2131from||||import 123.123`,
			expected: nil,
		},
		{
			name:     "not start with true",
			input:    `[:...[][]true[test{}[][]231:::sda2131from||||import... 123.123 true `,
			expected: nil,
		},
		{
			name:     "start with false",
			input:    `false 123.123::::...]test{}[][]231...:sda2131from||||import`,
			expected: []int{0, 5},
		},
		{
			name:     "start with false without space",
			input:    `false123::::...]test{}[][]231...:sda2131from||||import 123.123`,
			expected: nil,
		},
		{
			name:     "not start with false",
			input:    `[:...[][]true[test{}[][]231:::sda2131from||||import... 123.123 false `,
			expected: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			indexes := token.Bool.Regexp().FindStringIndex(tc.input)
			require.Equal(t, tc.expected, indexes)
		})
	}
}

func TestType_String_Regexp(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name     string
		input    string
		expected []int
	}{
		{
			name:     "start with",
			input:    `"test string" true 123.123::::...]test{}[][]231...:sda2131from||||import`,
			expected: []int{0, 13},
		},
		{
			name:     "start with true without space",
			input:    `"test string"::::...]test{}[][]231...:sda2131from||||import 123.123`,
			expected: []int{0, 13},
		},
		{
			name:     "start with and has escape chars",
			input:    `"test \"string\""::::...]test{}[][]231...:sda2131from||||import 123.123`,
			expected: []int{0, 17},
		},
		{
			name:     "start with empty string",
			input:    `""::::...]test{}[][]231...:sda2131from||||import 123.123`,
			expected: []int{0, 2},
		},
		{
			name:     "start not with",
			input:    `::::...]test{}[][]231...:sda2131from"test"||||import 123.123`,
			expected: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			indexes := token.String.Regexp().FindStringIndex(tc.input)
			require.Equal(t, tc.expected, indexes)
		})
	}
}

func TestType_Ident_Regexp(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name     string
		input    string
		expected []int
	}{
		{
			name:     "start with",
			input:    `key true 123.123::::...]test{}[][]231...:sda2131from||||import`,
			expected: []int{0, 3},
		},
		{
			name:     "start with true without space",
			input:    `key::::...]test{}[][]231...:sda2131from||||import 123.123`,
			expected: []int{0, 3},
		},
		{
			name:     "start with underscore",
			input:    `_key"test \"string\""::::...]test{}[][]231...:sda2131from||||import 123.123`,
			expected: []int{0, 4},
		},
		{
			name:     "start with digit",
			input:    `1key""::::...]test{}[][]231...:sda2131from||||import 123.123`,
			expected: nil,
		},
		{
			name:     "start with but illegal symbol",
			input:    `-key""::::...]test{}[][]231...:sda2131from||||import 123.123`,
			expected: nil,
		},
		{
			name:     "start with but illegal symbol in the middle",
			input:    `key-key""::::...]test{}[][]231...:sda2131from||||import 123.123`,
			expected: []int{0, 3},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			indexes := token.Ident.Regexp().FindStringIndex(tc.input)
			require.Equal(t, tc.expected, indexes)
		})
	}
}
