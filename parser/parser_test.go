package parser_test

import (
	"testing"

	"github.com/atmxlab/atmc/lexer/tokenmover"
	"github.com/atmxlab/atmc/parser"
	ast2 "github.com/atmxlab/atmc/parser/ast"
	"github.com/atmxlab/atmc/test/testast"
	"github.com/atmxlab/atmc/types"
	token2 "github.com/atmxlab/atmc/types/token"
	"github.com/stretchr/testify/require"
)

func TestParser_Parse(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name     string
		tokens   []token2.Token
		expected ast2.Ast
	}{
		{
			name: "with import and empty object",
			tokens: []token2.Token{
				token2.New(token2.Ident, "importName", types.Location{}),
				token2.New(token2.Path, "./dir/dir/config.atmc", types.Location{}),
				token2.New(token2.LBrace, "", types.Location{}),
				token2.New(token2.RBrace, "", types.Location{}),
			},
			expected: ast2.NewAst(
				ast2.NewFile(
					[]ast2.Import{
						ast2.NewImport(
							ast2.NewIdent("importName", types.Location{}),
							ast2.NewPath("./dir/dir/config.atmc", types.Location{}),
						),
					},
					ast2.NewObject(
						[]ast2.Entry{},
						types.Location{},
					),
				),
			),
		},
		{
			name: "without import and empty object",
			tokens: []token2.Token{
				token2.New(token2.LBrace, "", types.Location{}),
				token2.New(token2.RBrace, "", types.Location{}),
			},
			expected: ast2.NewAst(
				ast2.NewFile(
					[]ast2.Import{},
					ast2.NewObject(
						[]ast2.Entry{},
						types.Location{},
					),
				),
			),
		},
		{
			name: "with spreads in object",
			tokens: []token2.Token{
				token2.New(token2.LBrace, "", types.Location{}),
				token2.New(token2.Ident, "common1", types.Location{}),
				token2.New(token2.Spread, "", types.Location{}),
				token2.New(token2.Ident, "common2", types.Location{}),
				token2.New(token2.Spread, "", types.Location{}),
				token2.New(token2.RBrace, "", types.Location{}),
			},
			expected: ast2.NewAst(
				ast2.NewFile(
					[]ast2.Import{},
					ast2.NewObject(
						[]ast2.Entry{
							ast2.NewSpread(
								ast2.NewVar(
									[]ast2.Ident{
										ast2.NewIdent("common1", types.Location{}),
									},
								),
								types.Location{},
							),
							ast2.NewSpread(
								ast2.NewVar(
									[]ast2.Ident{
										ast2.NewIdent("common2", types.Location{}),
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
			tokens: []token2.Token{
				token2.New(token2.LBrace, "", types.Location{}),
				token2.New(token2.Ident, "common1", types.Location{}),
				token2.New(token2.Dot, "", types.Location{}),
				token2.New(token2.Ident, "nested1", types.Location{}),
				token2.New(token2.Spread, "", types.Location{}),
				token2.New(token2.Ident, "common2", types.Location{}),
				token2.New(token2.Dot, "", types.Location{}),
				token2.New(token2.Ident, "nested1", types.Location{}),
				token2.New(token2.Dot, "", types.Location{}),
				token2.New(token2.Ident, "nested2", types.Location{}),
				token2.New(token2.Spread, "", types.Location{}),
				token2.New(token2.RBrace, "", types.Location{}),
			},
			expected: ast2.NewAst(
				ast2.NewFile(
					[]ast2.Import{},
					ast2.NewObject(
						[]ast2.Entry{
							ast2.NewSpread(
								ast2.NewVar(
									[]ast2.Ident{
										ast2.NewIdent("common1", types.Location{}),
										ast2.NewIdent("nested1", types.Location{}),
									},
								),
								types.Location{},
							),
							ast2.NewSpread(
								ast2.NewVar(
									[]ast2.Ident{
										ast2.NewIdent("common2", types.Location{}),
										ast2.NewIdent("nested1", types.Location{}),
										ast2.NewIdent("nested2", types.Location{}),
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
			tokens: []token2.Token{
				token2.New(token2.LBrace, "", types.Location{}),

				token2.New(token2.Ident, "key1", types.Location{}),
				token2.New(token2.Colon, "", types.Location{}),
				token2.New(token2.Int, "123", types.Location{}),

				token2.New(token2.Ident, "key2", types.Location{}),
				token2.New(token2.Colon, "", types.Location{}),
				token2.New(token2.Float, "123.321", types.Location{}),

				token2.New(token2.Ident, "key3", types.Location{}),
				token2.New(token2.Colon, "", types.Location{}),
				token2.New(token2.String, `"test string"`, types.Location{}),

				token2.New(token2.Ident, "key4", types.Location{}),
				token2.New(token2.Colon, "", types.Location{}),
				token2.New(token2.Bool, "false", types.Location{}),

				token2.New(token2.Ident, "key5", types.Location{}),
				token2.New(token2.Colon, "", types.Location{}),
				token2.New(token2.Bool, "true", types.Location{}),

				token2.New(token2.Ident, "key6", types.Location{}),
				token2.New(token2.Colon, "", types.Location{}),
				token2.New(token2.Int, "-123", types.Location{}),

				token2.New(token2.Ident, "key7", types.Location{}),
				token2.New(token2.Colon, "", types.Location{}),
				token2.New(token2.Float, "-123.321", types.Location{}),

				token2.New(token2.RBrace, "", types.Location{}),
			},
			expected: ast2.NewAst(
				ast2.NewFile(
					[]ast2.Import{},
					ast2.NewObject(
						[]ast2.Entry{
							ast2.NewKV(
								ast2.NewIdent("key1", types.Location{}),
								testast.MustNewInt(t, "123"),
							),
							ast2.NewKV(
								ast2.NewIdent("key2", types.Location{}),
								testast.MustNewFloat(t, "123.321"),
							),
							ast2.NewKV(
								ast2.NewIdent("key3", types.Location{}),
								ast2.NewString(`"test string"`, types.Location{}),
							),
							ast2.NewKV(
								ast2.NewIdent("key4", types.Location{}),
								testast.MustNewBool(t, "false"),
							),
							ast2.NewKV(
								ast2.NewIdent("key5", types.Location{}),
								testast.MustNewBool(t, "true"),
							),
							ast2.NewKV(
								ast2.NewIdent("key6", types.Location{}),
								testast.MustNewInt(t, "-123"),
							),
							ast2.NewKV(
								ast2.NewIdent("key7", types.Location{}),
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
			tokens: []token2.Token{
				token2.New(token2.LBrace, "", types.Location{}),

				token2.New(token2.Ident, "common1", types.Location{}),
				token2.New(token2.Dot, "", types.Location{}),
				token2.New(token2.Ident, "nested1", types.Location{}),
				token2.New(token2.Spread, "", types.Location{}),

				token2.New(token2.Ident, "common2", types.Location{}),
				token2.New(token2.Dot, "", types.Location{}),
				token2.New(token2.Ident, "nested1", types.Location{}),
				token2.New(token2.Dot, "", types.Location{}),
				token2.New(token2.Ident, "nested2", types.Location{}),
				token2.New(token2.Spread, "", types.Location{}),

				token2.New(token2.Ident, "key1", types.Location{}),
				token2.New(token2.Colon, "", types.Location{}),
				token2.New(token2.Int, "123", types.Location{}),

				token2.New(token2.Ident, "key2", types.Location{}),
				token2.New(token2.Colon, "", types.Location{}),
				token2.New(token2.Float, "123.321", types.Location{}),

				token2.New(token2.Ident, "common3", types.Location{}),
				token2.New(token2.Spread, "", types.Location{}),

				token2.New(token2.RBrace, "", types.Location{}),
			},
			expected: ast2.NewAst(
				ast2.NewFile(
					[]ast2.Import{},
					ast2.NewObject(
						[]ast2.Entry{
							ast2.NewSpread(
								ast2.NewVar(
									[]ast2.Ident{
										ast2.NewIdent("common1", types.Location{}),
										ast2.NewIdent("nested1", types.Location{}),
									},
								),
								types.Location{},
							),
							ast2.NewSpread(
								ast2.NewVar(
									[]ast2.Ident{
										ast2.NewIdent("common2", types.Location{}),
										ast2.NewIdent("nested1", types.Location{}),
										ast2.NewIdent("nested2", types.Location{}),
									},
								),
								types.Location{},
							),
							ast2.NewKV(
								ast2.NewIdent("key1", types.Location{}),
								testast.MustNewInt(t, "123"),
							),
							ast2.NewKV(
								ast2.NewIdent("key2", types.Location{}),
								testast.MustNewFloat(t, "123.321"),
							),
							ast2.NewSpread(
								ast2.NewVar(
									[]ast2.Ident{
										ast2.NewIdent("common3", types.Location{}),
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
			tokens: []token2.Token{
				token2.New(token2.LBrace, "", types.Location{}),

				token2.New(token2.Ident, "common1", types.Location{}),
				token2.New(token2.Colon, "", types.Location{}),
				token2.New(token2.LBrace, "", types.Location{}),

				token2.New(token2.Ident, "nested1", types.Location{}),
				token2.New(token2.Colon, "", types.Location{}),
				token2.New(token2.Int, "123", types.Location{}),

				token2.New(token2.Ident, "nested2", types.Location{}),
				token2.New(token2.Colon, "", types.Location{}),
				token2.New(token2.Float, "123.321", types.Location{}),

				token2.New(token2.Ident, "common", types.Location{}),
				token2.New(token2.Spread, "", types.Location{}),

				token2.New(token2.RBrace, "", types.Location{}),

				token2.New(token2.RBrace, "", types.Location{}),
			},
			expected: ast2.NewAst(
				ast2.NewFile(
					[]ast2.Import{},
					ast2.NewObject(
						[]ast2.Entry{
							ast2.NewKV(
								ast2.NewIdent("common1", types.Location{}),
								ast2.NewObject(
									[]ast2.Entry{
										ast2.NewKV(
											ast2.NewIdent("nested1", types.Location{}),
											testast.MustNewInt(t, "123"),
										),
										ast2.NewKV(
											ast2.NewIdent("nested2", types.Location{}),
											testast.MustNewFloat(t, "123.321"),
										),
										ast2.NewSpread(
											ast2.NewVar(
												[]ast2.Ident{
													ast2.NewIdent("common", types.Location{}),
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
			tokens: []token2.Token{
				token2.New(token2.LBrace, "", types.Location{}),

				token2.New(token2.Ident, "common1", types.Location{}),
				token2.New(token2.Colon, "", types.Location{}),

				token2.New(token2.LBracket, "", types.Location{}),
				token2.New(token2.RBracket, "", types.Location{}),

				token2.New(token2.RBrace, "", types.Location{}),
			},
			expected: ast2.NewAst(
				ast2.NewFile(
					[]ast2.Import{},
					ast2.NewObject(
						[]ast2.Entry{
							ast2.NewKV(
								ast2.NewIdent("common1", types.Location{}),
								ast2.NewArray(
									[]ast2.Expression{},
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
			tokens: []token2.Token{
				token2.New(token2.LBrace, "", types.Location{}),

				token2.New(token2.Ident, "common1", types.Location{}),
				token2.New(token2.Colon, "", types.Location{}),
				token2.New(token2.LBracket, "", types.Location{}),

				token2.New(token2.Ident, "nested1", types.Location{}),

				token2.New(token2.LBrace, "", types.Location{}),
				token2.New(token2.Ident, "nested2", types.Location{}),
				token2.New(token2.Colon, "", types.Location{}),
				token2.New(token2.Int, "123", types.Location{}),
				token2.New(token2.RBrace, "", types.Location{}),

				token2.New(token2.Bool, "false", types.Location{}),

				token2.New(token2.Dollar, "", types.Location{}),
				token2.New(token2.Ident, "ENV_VAR", types.Location{}),

				token2.New(token2.LBracket, "", types.Location{}),

				token2.New(token2.Ident, "nested1", types.Location{}),

				token2.New(token2.LBrace, "", types.Location{}),
				token2.New(token2.Ident, "nested2", types.Location{}),
				token2.New(token2.Colon, "", types.Location{}),
				token2.New(token2.Int, "123", types.Location{}),
				token2.New(token2.RBrace, "", types.Location{}),

				token2.New(token2.Bool, "false", types.Location{}),

				token2.New(token2.Dollar, "", types.Location{}),
				token2.New(token2.Ident, "ENV_VAR", types.Location{}),

				token2.New(token2.RBracket, "", types.Location{}),

				token2.New(token2.RBracket, "", types.Location{}),

				token2.New(token2.RBrace, "", types.Location{}),
			},
			expected: ast2.NewAst(
				ast2.NewFile(
					[]ast2.Import{},
					ast2.NewObject(
						[]ast2.Entry{
							ast2.NewKV(
								ast2.NewIdent("common1", types.Location{}),
								ast2.NewArray(
									[]ast2.Expression{
										ast2.NewVar(
											[]ast2.Ident{
												ast2.NewIdent("nested1", types.Location{}),
											},
										),
										ast2.NewObject(
											[]ast2.Entry{
												ast2.NewKV(
													ast2.NewIdent("nested2", types.Location{}),
													testast.MustNewInt(t, "123"),
												),
											},
											types.Location{},
										),
										testast.MustNewBool(t, "false"),
										ast2.NewEnv(
											ast2.NewIdent("ENV_VAR", types.Location{}),
											types.Location{},
										),

										ast2.NewArray(
											[]ast2.Expression{
												ast2.NewVar(
													[]ast2.Ident{
														ast2.NewIdent("nested1", types.Location{}),
													},
												),
												ast2.NewObject(
													[]ast2.Entry{
														ast2.NewKV(
															ast2.NewIdent("nested2", types.Location{}),
															testast.MustNewInt(t, "123"),
														),
													},
													types.Location{},
												),
												testast.MustNewBool(t, "false"),
												ast2.NewEnv(
													ast2.NewIdent("ENV_VAR", types.Location{}),
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
			tokens: []token2.Token{
				token2.New(token2.LBrace, "", types.Location{}),

				token2.New(token2.Ident, "common1", types.Location{}),
				token2.New(token2.Colon, "", types.Location{}),

				token2.New(token2.LBracket, "", types.Location{}),

				token2.New(token2.Ident, "common1", types.Location{}),
				token2.New(token2.Spread, "", types.Location{}),

				token2.New(token2.Ident, "common2", types.Location{}),
				token2.New(token2.Dot, "", types.Location{}),
				token2.New(token2.Ident, "nested1", types.Location{}),
				token2.New(token2.Spread, "", types.Location{}),

				token2.New(token2.String, `"test test"`, types.Location{}),

				token2.New(token2.RBracket, "", types.Location{}),

				token2.New(token2.RBrace, "", types.Location{}),
			},
			expected: ast2.NewAst(
				ast2.NewFile(
					[]ast2.Import{},
					ast2.NewObject(
						[]ast2.Entry{
							ast2.NewKV(
								ast2.NewIdent("common1", types.Location{}),
								ast2.NewArray(
									[]ast2.Expression{
										ast2.NewSpread(
											ast2.NewVar(
												[]ast2.Ident{
													ast2.NewIdent("common1", types.Location{}),
												},
											),
											types.Location{},
										),
										ast2.NewSpread(
											ast2.NewVar(
												[]ast2.Ident{
													ast2.NewIdent("common2", types.Location{}),
													ast2.NewIdent("nested1", types.Location{}),
												},
											),
											types.Location{},
										),
										ast2.NewString(`"test test"`, types.Location{}),
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
