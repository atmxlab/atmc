package v2_test

import (
	"testing"

	"github.com/atmxlab/atmcfg/internal/parser/ast"
	v2 "github.com/atmxlab/atmcfg/internal/parser/v2"
	"github.com/atmxlab/atmcfg/internal/test"
	"github.com/atmxlab/atmcfg/internal/types"
	"github.com/atmxlab/atmcfg/internal/types/token"
	"github.com/stretchr/testify/require"
)

func TestParser_Parse(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name     string
		tokens   []token.Token
		expected ast.Ast
	}{
		{
			name: "with import and empty object",
			tokens: []token.Token{
				token.New(token.Ident, "importName", types.Location{}),
				token.New(token.Path, "./dir/dir/config.atmc", types.Location{}),
				token.New(token.LBrace, "", types.Location{}),
				token.New(token.RBrace, "", types.Location{}),
			},
			expected: ast.NewAst(
				ast.NewFile(
					[]ast.Import{
						ast.NewImport(
							ast.NewIdent("importName", types.Location{}),
							ast.NewPath("./dir/dir/config.atmc"),
							types.Location{},
						),
					},
					ast.NewObject(
						[]ast.Entry{},
						types.Location{},
					),
					types.Location{},
				),
			),
		},
		{
			name: "without import and empty object",
			tokens: []token.Token{
				token.New(token.LBrace, "", types.Location{}),
				token.New(token.RBrace, "", types.Location{}),
			},
			expected: ast.NewAst(
				ast.NewFile(
					[]ast.Import{},
					ast.NewObject(
						[]ast.Entry{},
						types.Location{},
					),
					types.Location{},
				),
			),
		},
		{
			name: "with spreads in object",
			tokens: []token.Token{
				token.New(token.LBrace, "", types.Location{}),
				token.New(token.Ident, "common1", types.Location{}),
				token.New(token.Spread, "", types.Location{}),
				token.New(token.Ident, "common2", types.Location{}),
				token.New(token.Spread, "", types.Location{}),
				token.New(token.RBrace, "", types.Location{}),
			},
			expected: ast.NewAst(
				ast.NewFile(
					[]ast.Import{},
					ast.NewObject(
						[]ast.Entry{
							ast.NewSpread(
								ast.NewVar(
									[]ast.Ident{
										ast.NewIdent("common1", types.Location{}),
									},
									types.Location{},
								),
								types.Location{},
							),
							ast.NewSpread(
								ast.NewVar(
									[]ast.Ident{
										ast.NewIdent("common2", types.Location{}),
									},
									types.Location{},
								),
								types.Location{},
							),
						},
						types.Location{},
					),
					types.Location{},
				),
			),
		},
		{
			name: "with nested var spreads in object",
			tokens: []token.Token{
				token.New(token.LBrace, "", types.Location{}),
				token.New(token.Ident, "common1", types.Location{}),
				token.New(token.Dot, "", types.Location{}),
				token.New(token.Ident, "nested1", types.Location{}),
				token.New(token.Spread, "", types.Location{}),
				token.New(token.Ident, "common2", types.Location{}),
				token.New(token.Dot, "", types.Location{}),
				token.New(token.Ident, "nested1", types.Location{}),
				token.New(token.Dot, "", types.Location{}),
				token.New(token.Ident, "nested2", types.Location{}),
				token.New(token.Spread, "", types.Location{}),
				token.New(token.RBrace, "", types.Location{}),
			},
			expected: ast.NewAst(
				ast.NewFile(
					[]ast.Import{},
					ast.NewObject(
						[]ast.Entry{
							ast.NewSpread(
								ast.NewVar(
									[]ast.Ident{
										ast.NewIdent("common1", types.Location{}),
										ast.NewIdent("nested1", types.Location{}),
									},
									types.Location{},
								),
								types.Location{},
							),
							ast.NewSpread(
								ast.NewVar(
									[]ast.Ident{
										ast.NewIdent("common2", types.Location{}),
										ast.NewIdent("nested1", types.Location{}),
										ast.NewIdent("nested2", types.Location{}),
									},
									types.Location{},
								),
								types.Location{},
							),
						},
						types.Location{},
					),
					types.Location{},
				),
			),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			mover := test.NewTokenMover(t, tc.tokens)

			p := v2.NewParser(mover)

			a, err := p.Parse()
			require.NoError(t, err)

			require.Equal(t, tc.expected, a)
		})
	}
}
