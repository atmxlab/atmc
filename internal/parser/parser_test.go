package parser_test

import (
	"testing"

	"github.com/atmxlab/atmcfg/internal/lexer/tokenmover"
	"github.com/atmxlab/atmcfg/internal/parser"
	"github.com/atmxlab/atmcfg/internal/parser/ast"
	"github.com/atmxlab/atmcfg/internal/test/testast"
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
							ast.NewPath("./dir/dir/config.atmc", types.Location{}),
						),
					},
					ast.NewObject(
						[]ast.Entry{},
						types.Location{},
					),
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
								),
								types.Location{},
							),
							ast.NewSpread(
								ast.NewVar(
									[]ast.Ident{
										ast.NewIdent("common2", types.Location{}),
									},
								),
								types.Location{},
							),
						},
						types.Location{},
					),
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
								),
								types.Location{},
							),
						},
						types.Location{},
					),
				),
			),
		},
		{
			name: "object with entries with literal value",
			tokens: []token.Token{
				token.New(token.LBrace, "", types.Location{}),

				token.New(token.Ident, "key1", types.Location{}),
				token.New(token.Colon, "", types.Location{}),
				token.New(token.Int, "123", types.Location{}),

				token.New(token.Ident, "key2", types.Location{}),
				token.New(token.Colon, "", types.Location{}),
				token.New(token.Float, "123.321", types.Location{}),

				token.New(token.Ident, "key3", types.Location{}),
				token.New(token.Colon, "", types.Location{}),
				token.New(token.String, `"test string"`, types.Location{}),

				token.New(token.Ident, "key4", types.Location{}),
				token.New(token.Colon, "", types.Location{}),
				token.New(token.Bool, "false", types.Location{}),

				token.New(token.Ident, "key5", types.Location{}),
				token.New(token.Colon, "", types.Location{}),
				token.New(token.Bool, "true", types.Location{}),

				token.New(token.Ident, "key6", types.Location{}),
				token.New(token.Colon, "", types.Location{}),
				token.New(token.Int, "-123", types.Location{}),

				token.New(token.Ident, "key7", types.Location{}),
				token.New(token.Colon, "", types.Location{}),
				token.New(token.Float, "-123.321", types.Location{}),

				token.New(token.RBrace, "", types.Location{}),
			},
			expected: ast.NewAst(
				ast.NewFile(
					[]ast.Import{},
					ast.NewObject(
						[]ast.Entry{
							ast.NewKV(
								ast.NewIdent("key1", types.Location{}),
								testast.MustNewInt(t, "123"),
							),
							ast.NewKV(
								ast.NewIdent("key2", types.Location{}),
								testast.MustNewFloat(t, "123.321"),
							),
							ast.NewKV(
								ast.NewIdent("key3", types.Location{}),
								ast.NewString(`"test string"`, types.Location{}),
							),
							ast.NewKV(
								ast.NewIdent("key4", types.Location{}),
								testast.MustNewBool(t, "false"),
							),
							ast.NewKV(
								ast.NewIdent("key5", types.Location{}),
								testast.MustNewBool(t, "true"),
							),
							ast.NewKV(
								ast.NewIdent("key6", types.Location{}),
								testast.MustNewInt(t, "-123"),
							),
							ast.NewKV(
								ast.NewIdent("key7", types.Location{}),
								testast.MustNewFloat(t, "-123.321"),
							),
						},
						types.Location{},
					),
				),
			),
		},
		{
			name: "with nested var spreads and entries in object",
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

				token.New(token.Ident, "key1", types.Location{}),
				token.New(token.Colon, "", types.Location{}),
				token.New(token.Int, "123", types.Location{}),

				token.New(token.Ident, "key2", types.Location{}),
				token.New(token.Colon, "", types.Location{}),
				token.New(token.Float, "123.321", types.Location{}),

				token.New(token.Ident, "common3", types.Location{}),
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
								),
								types.Location{},
							),
							ast.NewKV(
								ast.NewIdent("key1", types.Location{}),
								testast.MustNewInt(t, "123"),
							),
							ast.NewKV(
								ast.NewIdent("key2", types.Location{}),
								testast.MustNewFloat(t, "123.321"),
							),
							ast.NewSpread(
								ast.NewVar(
									[]ast.Ident{
										ast.NewIdent("common3", types.Location{}),
									},
								),
								types.Location{},
							),
						},
						types.Location{},
					),
				),
			),
		},
		{
			name: "with object in value with mixed values",
			tokens: []token.Token{
				token.New(token.LBrace, "", types.Location{}),

				token.New(token.Ident, "common1", types.Location{}),
				token.New(token.Colon, "", types.Location{}),
				token.New(token.LBrace, "", types.Location{}),

				token.New(token.Ident, "nested1", types.Location{}),
				token.New(token.Colon, "", types.Location{}),
				token.New(token.Int, "123", types.Location{}),

				token.New(token.Ident, "nested2", types.Location{}),
				token.New(token.Colon, "", types.Location{}),
				token.New(token.Float, "123.321", types.Location{}),

				token.New(token.Ident, "common", types.Location{}),
				token.New(token.Spread, "", types.Location{}),

				token.New(token.RBrace, "", types.Location{}),

				token.New(token.RBrace, "", types.Location{}),
			},
			expected: ast.NewAst(
				ast.NewFile(
					[]ast.Import{},
					ast.NewObject(
						[]ast.Entry{
							ast.NewKV(
								ast.NewIdent("common1", types.Location{}),
								ast.NewObject(
									[]ast.Entry{
										ast.NewKV(
											ast.NewIdent("nested1", types.Location{}),
											testast.MustNewInt(t, "123"),
										),
										ast.NewKV(
											ast.NewIdent("nested2", types.Location{}),
											testast.MustNewFloat(t, "123.321"),
										),
										ast.NewSpread(
											ast.NewVar(
												[]ast.Ident{
													ast.NewIdent("common", types.Location{}),
												},
											),
											types.Location{},
										),
									},
									types.Location{},
								),
							),
						},
						types.Location{},
					),
				),
			),
		},
		{
			name: "with empty array",
			tokens: []token.Token{
				token.New(token.LBrace, "", types.Location{}),

				token.New(token.Ident, "common1", types.Location{}),
				token.New(token.Colon, "", types.Location{}),

				token.New(token.LBracket, "", types.Location{}),
				token.New(token.RBracket, "", types.Location{}),

				token.New(token.RBrace, "", types.Location{}),
			},
			expected: ast.NewAst(
				ast.NewFile(
					[]ast.Import{},
					ast.NewObject(
						[]ast.Entry{
							ast.NewKV(
								ast.NewIdent("common1", types.Location{}),
								ast.NewArray(
									[]ast.Expression{},
									types.Location{},
								),
							),
						},
						types.Location{},
					),
				),
			),
		},
		{
			name: "with array with mixed values",
			tokens: []token.Token{
				token.New(token.LBrace, "", types.Location{}),

				token.New(token.Ident, "common1", types.Location{}),
				token.New(token.Colon, "", types.Location{}),
				token.New(token.LBracket, "", types.Location{}),

				token.New(token.Ident, "nested1", types.Location{}),

				token.New(token.LBrace, "", types.Location{}),
				token.New(token.Ident, "nested2", types.Location{}),
				token.New(token.Colon, "", types.Location{}),
				token.New(token.Int, "123", types.Location{}),
				token.New(token.RBrace, "", types.Location{}),

				token.New(token.Bool, "false", types.Location{}),

				token.New(token.Dollar, "", types.Location{}),
				token.New(token.Ident, "ENV_VAR", types.Location{}),

				token.New(token.LBracket, "", types.Location{}),

				token.New(token.Ident, "nested1", types.Location{}),

				token.New(token.LBrace, "", types.Location{}),
				token.New(token.Ident, "nested2", types.Location{}),
				token.New(token.Colon, "", types.Location{}),
				token.New(token.Int, "123", types.Location{}),
				token.New(token.RBrace, "", types.Location{}),

				token.New(token.Bool, "false", types.Location{}),

				token.New(token.Dollar, "", types.Location{}),
				token.New(token.Ident, "ENV_VAR", types.Location{}),

				token.New(token.RBracket, "", types.Location{}),

				token.New(token.RBracket, "", types.Location{}),

				token.New(token.RBrace, "", types.Location{}),
			},
			expected: ast.NewAst(
				ast.NewFile(
					[]ast.Import{},
					ast.NewObject(
						[]ast.Entry{
							ast.NewKV(
								ast.NewIdent("common1", types.Location{}),
								ast.NewArray(
									[]ast.Expression{
										ast.NewVar(
											[]ast.Ident{
												ast.NewIdent("nested1", types.Location{}),
											},
										),
										ast.NewObject(
											[]ast.Entry{
												ast.NewKV(
													ast.NewIdent("nested2", types.Location{}),
													testast.MustNewInt(t, "123"),
												),
											},
											types.Location{},
										),
										testast.MustNewBool(t, "false"),
										ast.NewEnv(
											ast.NewIdent("ENV_VAR", types.Location{}),
											types.Location{},
										),

										ast.NewArray(
											[]ast.Expression{
												ast.NewVar(
													[]ast.Ident{
														ast.NewIdent("nested1", types.Location{}),
													},
												),
												ast.NewObject(
													[]ast.Entry{
														ast.NewKV(
															ast.NewIdent("nested2", types.Location{}),
															testast.MustNewInt(t, "123"),
														),
													},
													types.Location{},
												),
												testast.MustNewBool(t, "false"),
												ast.NewEnv(
													ast.NewIdent("ENV_VAR", types.Location{}),
													types.Location{},
												),
											},
											types.Location{},
										),
									},
									types.Location{},
								),
							),
						},
						types.Location{},
					),
				),
			),
		},
		{
			name: "array with spreads",
			tokens: []token.Token{
				token.New(token.LBrace, "", types.Location{}),

				token.New(token.Ident, "common1", types.Location{}),
				token.New(token.Colon, "", types.Location{}),

				token.New(token.LBracket, "", types.Location{}),

				token.New(token.Ident, "common1", types.Location{}),
				token.New(token.Spread, "", types.Location{}),

				token.New(token.Ident, "common2", types.Location{}),
				token.New(token.Dot, "", types.Location{}),
				token.New(token.Ident, "nested1", types.Location{}),
				token.New(token.Spread, "", types.Location{}),

				token.New(token.String, `"test test"`, types.Location{}),

				token.New(token.RBracket, "", types.Location{}),

				token.New(token.RBrace, "", types.Location{}),
			},
			expected: ast.NewAst(
				ast.NewFile(
					[]ast.Import{},
					ast.NewObject(
						[]ast.Entry{
							ast.NewKV(
								ast.NewIdent("common1", types.Location{}),
								ast.NewArray(
									[]ast.Expression{
										ast.NewSpread(
											ast.NewVar(
												[]ast.Ident{
													ast.NewIdent("common1", types.Location{}),
												},
											),
											types.Location{},
										),
										ast.NewSpread(
											ast.NewVar(
												[]ast.Ident{
													ast.NewIdent("common2", types.Location{}),
													ast.NewIdent("nested1", types.Location{}),
												},
											),
											types.Location{},
										),
										ast.NewString(`"test test"`, types.Location{}),
									},
									types.Location{},
								),
							),
						},
						types.Location{},
					),
				),
			),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			p := parser.New()

			a, err := p.Parse(tokenmover.New(tc.tokens))
			require.NoError(t, err)

			require.Equal(t, tc.expected, a)
		})
	}
}
