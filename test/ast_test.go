package test

import (
	"testing"

	"github.com/atmxlab/atmcfg/internal/lexer"
	"github.com/atmxlab/atmcfg/internal/lexer/tokenmover"
	"github.com/atmxlab/atmcfg/internal/parser/ast"
	"github.com/atmxlab/atmcfg/internal/parser/parser"
	"github.com/atmxlab/atmcfg/internal/test/testast"
	"github.com/atmxlab/atmcfg/internal/types"
	"github.com/stretchr/testify/require"
)

func TestAstGenerate(t *testing.T) {
	t.Parallel()

	input := `
	common ./dir/common.atmc
	kafka ./dir/kafka.atmc
	abs /dir/abs.atmc

	{
		common...
		int: 123
		float: 123.123
		bool1: true
		bool2: false
		string: "test string \"escaped"
		kafka.nested...
		object1: {
			common.field.field...
			field1: 123
			field2: {
				field1: 321
			}
		}
		array1: [
			123, 321,
			444, 555
		]
		envVar: $POSTGRES_PASSWORD
		arrayWithComma: ["test text", 444, -321, 123.123, true, false, common.field.field..., kafka.field.field]
		arrayWithoutComma: ["test text" 444 321 123.123 true false common.field.field... kafka.field.field]
		arrayWithNestedObject: [{field1: 99999999999, kafka.field.field...}, $POSTGRES_HOST]
	}
`

	expectedAst := ast.NewAst(
		ast.NewFile(
			[]ast.Import{
				ast.NewImport(
					ast.NewIdent("common", types.NewLocation(
						types.NewPosition(2, 1, 2),
						types.NewPosition(2, 7, 8),
					)),
					ast.NewPath("./dir/common.atmc", types.NewLocation(
						types.NewPosition(2, 8, 9),
						types.NewPosition(2, 25, 26),
					)),
				),
				ast.NewImport(
					ast.NewIdent("kafka", types.NewLocation(
						types.NewPosition(3, 1, 28),
						types.NewPosition(3, 6, 33),
					)),
					ast.NewPath("./dir/kafka.atmc", types.NewLocation(
						types.NewPosition(3, 7, 34),
						types.NewPosition(3, 23, 50),
					)),
				),
				ast.NewImport(
					ast.NewIdent("abs", types.NewLocation(
						types.NewPosition(4, 1, 52),
						types.NewPosition(4, 4, 55),
					)),
					ast.NewPath("/dir/abs.atmc", types.NewLocation(
						types.NewPosition(4, 5, 56),
						types.NewPosition(4, 18, 69),
					)),
				),
			},
			ast.NewObject(
				[]ast.Entry{
					ast.NewSpread(
						ast.NewVar([]ast.Ident{
							ast.NewIdent("common", types.NewLocation(
								types.NewPosition(7, 2, 76),
								types.NewPosition(7, 8, 82),
							)),
						}),
						types.NewLocation(
							types.NewPosition(7, 2, 76),
							types.NewPosition(7, 11, 85),
						),
					),
					ast.NewKV(
						ast.NewIdent("int", types.NewLocation(
							types.NewPosition(8, 2, 88),
							types.NewPosition(8, 5, 91),
						)),
						testast.MustNewIntWithLocation(t, "123", types.NewLocation(
							types.NewPosition(8, 7, 93),
							types.NewPosition(8, 10, 96),
						)),
					),
					ast.NewKV(
						ast.NewIdent("float", types.NewLocation(
							types.NewPosition(0, 0, 0),
							types.NewPosition(0, 0, 0),
						)),
						testast.MustNewFloatWithLocation(t, "123.123", types.NewLocation(
							types.NewPosition(0, 0, 0),
							types.NewPosition(0, 0, 0),
						)),
					),
					ast.NewKV(
						ast.NewIdent("bool1", types.NewLocation(
							types.NewPosition(0, 0, 0),
							types.NewPosition(0, 0, 0),
						)),
						testast.MustNewBoolWithLocation(t, "true", types.NewLocation(
							types.NewPosition(0, 0, 0),
							types.NewPosition(0, 0, 0),
						)),
					),
					ast.NewKV(
						ast.NewIdent("bool2", types.NewLocation(
							types.NewPosition(0, 0, 0),
							types.NewPosition(0, 0, 0),
						)),
						testast.MustNewBoolWithLocation(t, "false", types.NewLocation(
							types.NewPosition(0, 0, 0),
							types.NewPosition(0, 0, 0),
						)),
					),
					ast.NewKV(
						ast.NewIdent("bool1", types.NewLocation(
							types.NewPosition(0, 0, 0),
							types.NewPosition(0, 0, 0),
						)),
						ast.NewString(`"test string \"escaped"`, types.NewLocation(
							types.NewPosition(0, 0, 0),
							types.NewPosition(0, 0, 0),
						)),
					),
					ast.NewSpread(
						ast.NewVar([]ast.Ident{
							ast.NewIdent("kafka", types.NewLocation(
								types.NewPosition(0, 0, 0),
								types.NewPosition(0, 0, 0),
							)),
							ast.NewIdent("nested", types.NewLocation(
								types.NewPosition(0, 0, 0),
								types.NewPosition(0, 0, 0),
							)),
						}),
						types.NewLocation(
							types.NewPosition(7, 2, 76),
							types.NewPosition(7, 11, 85),
						),
					),

					ast.NewKV(
						ast.NewIdent("object1", types.NewLocation(
							types.NewPosition(0, 0, 0),
							types.NewPosition(0, 0, 0),
						)),
						ast.NewObject(
							[]ast.Entry{
								ast.NewSpread(
									ast.NewVar([]ast.Ident{
										ast.NewIdent("common", types.NewLocation(
											types.NewPosition(0, 0, 0),
											types.NewPosition(0, 0, 0),
										)),
										ast.NewIdent("field", types.NewLocation(
											types.NewPosition(0, 0, 0),
											types.NewPosition(0, 0, 0),
										)),
										ast.NewIdent("field", types.NewLocation(
											types.NewPosition(0, 0, 0),
											types.NewPosition(0, 0, 0),
										)),
									}),
									types.NewLocation(
										types.NewPosition(7, 2, 76),
										types.NewPosition(7, 11, 85),
									),
								),
								ast.NewKV(
									ast.NewIdent("field1", types.NewLocation(
										types.NewPosition(0, 0, 0),
										types.NewPosition(0, 0, 0),
									)),
									testast.MustNewIntWithLocation(t, "123", types.NewLocation(
										types.NewPosition(0, 0, 0),
										types.NewPosition(0, 0, 0),
									)),
								),
								ast.NewKV(
									ast.NewIdent("field2", types.NewLocation(
										types.NewPosition(0, 0, 0),
										types.NewPosition(0, 0, 0),
									)),
									ast.NewObject(
										[]ast.Entry{
											ast.NewKV(
												ast.NewIdent("field1", types.NewLocation(
													types.NewPosition(0, 0, 0),
													types.NewPosition(0, 0, 0),
												)),
												testast.MustNewIntWithLocation(t, "321", types.NewLocation(
													types.NewPosition(0, 0, 0),
													types.NewPosition(0, 0, 0),
												)),
											),
										},
										types.NewLocation(
											types.NewPosition(0, 0, 0),
											types.NewPosition(0, 0, 0),
										),
									),
								),
							},
							types.NewLocation(
								types.NewPosition(0, 0, 0),
								types.NewPosition(0, 0, 0),
							),
						),
					),
					ast.NewKV(
						ast.NewIdent("array1", types.NewLocation(
							types.NewPosition(0, 0, 0),
							types.NewPosition(0, 0, 0),
						)),
						ast.NewArray(
							[]ast.Expression{
								testast.MustNewIntWithLocation(t, "123", types.NewLocation(
									types.NewPosition(0, 0, 0),
									types.NewPosition(0, 0, 0),
								)),
								testast.MustNewIntWithLocation(t, "321", types.NewLocation(
									types.NewPosition(0, 0, 0),
									types.NewPosition(0, 0, 0),
								)),
								testast.MustNewIntWithLocation(t, "444", types.NewLocation(
									types.NewPosition(0, 0, 0),
									types.NewPosition(0, 0, 0),
								)),
								testast.MustNewIntWithLocation(t, "555", types.NewLocation(
									types.NewPosition(0, 0, 0),
									types.NewPosition(0, 0, 0),
								)),
							},
							types.NewLocation(
								types.NewPosition(0, 0, 0),
								types.NewPosition(0, 0, 0),
							),
						),
					),
					ast.NewKV(
						ast.NewIdent("envVar", types.NewLocation(
							types.NewPosition(0, 0, 0),
							types.NewPosition(0, 0, 0),
						)),
						ast.NewEnv(
							ast.NewIdent("POSTGRES_PASSWORD", types.NewLocation(
								types.NewPosition(0, 0, 0),
								types.NewPosition(0, 0, 0),
							)),
							types.NewLocation(
								types.NewPosition(0, 0, 0),
								types.NewPosition(0, 0, 0),
							),
						),
					),
					ast.NewKV(
						ast.NewIdent("arrayWithComma", types.NewLocation(
							types.NewPosition(0, 0, 0),
							types.NewPosition(0, 0, 0),
						)),
						ast.NewArray(
							[]ast.Expression{
								ast.NewString("test text", types.NewLocation(
									types.NewPosition(0, 0, 0),
									types.NewPosition(0, 0, 0),
								)),
								testast.MustNewIntWithLocation(t, "444", types.NewLocation(
									types.NewPosition(0, 0, 0),
									types.NewPosition(0, 0, 0),
								)),
								testast.MustNewIntWithLocation(t, "-321", types.NewLocation(
									types.NewPosition(0, 0, 0),
									types.NewPosition(0, 0, 0),
								)),
								testast.MustNewFloatWithLocation(t, "123.123", types.NewLocation(
									types.NewPosition(0, 0, 0),
									types.NewPosition(0, 0, 0),
								)),
								testast.MustNewBoolWithLocation(t, "true", types.NewLocation(
									types.NewPosition(0, 0, 0),
									types.NewPosition(0, 0, 0),
								)),
								testast.MustNewBoolWithLocation(t, "false", types.NewLocation(
									types.NewPosition(0, 0, 0),
									types.NewPosition(0, 0, 0),
								)),
								ast.NewSpread(
									ast.NewVar([]ast.Ident{
										ast.NewIdent("common", types.NewLocation(
											types.NewPosition(0, 0, 0),
											types.NewPosition(0, 0, 0),
										)),
										ast.NewIdent("field", types.NewLocation(
											types.NewPosition(0, 0, 0),
											types.NewPosition(0, 0, 0),
										)),
										ast.NewIdent("field", types.NewLocation(
											types.NewPosition(0, 0, 0),
											types.NewPosition(0, 0, 0),
										)),
									}),
									types.NewLocation(
										types.NewPosition(0, 0, 0),
										types.NewPosition(0, 0, 0),
									),
								),
								ast.NewVar([]ast.Ident{
									ast.NewIdent("kafka", types.NewLocation(
										types.NewPosition(0, 0, 0),
										types.NewPosition(0, 0, 0),
									)),
									ast.NewIdent("field", types.NewLocation(
										types.NewPosition(0, 0, 0),
										types.NewPosition(0, 0, 0),
									)),
									ast.NewIdent("field", types.NewLocation(
										types.NewPosition(0, 0, 0),
										types.NewPosition(0, 0, 0),
									)),
								}),
							},
							types.NewLocation(
								types.NewPosition(0, 0, 0),
								types.NewPosition(0, 0, 0),
							),
						),
					),
					ast.NewKV(
						ast.NewIdent("arrayWithoutComma", types.NewLocation(
							types.NewPosition(0, 0, 0),
							types.NewPosition(0, 0, 0),
						)),
						ast.NewArray(
							[]ast.Expression{
								ast.NewString("test text", types.NewLocation(
									types.NewPosition(0, 0, 0),
									types.NewPosition(0, 0, 0),
								)),
								testast.MustNewIntWithLocation(t, "444", types.NewLocation(
									types.NewPosition(0, 0, 0),
									types.NewPosition(0, 0, 0),
								)),
								testast.MustNewIntWithLocation(t, "-321", types.NewLocation(
									types.NewPosition(0, 0, 0),
									types.NewPosition(0, 0, 0),
								)),
								testast.MustNewFloatWithLocation(t, "123.123", types.NewLocation(
									types.NewPosition(0, 0, 0),
									types.NewPosition(0, 0, 0),
								)),
								testast.MustNewBoolWithLocation(t, "true", types.NewLocation(
									types.NewPosition(0, 0, 0),
									types.NewPosition(0, 0, 0),
								)),
								testast.MustNewBoolWithLocation(t, "false", types.NewLocation(
									types.NewPosition(0, 0, 0),
									types.NewPosition(0, 0, 0),
								)),
								ast.NewSpread(
									ast.NewVar([]ast.Ident{
										ast.NewIdent("common", types.NewLocation(
											types.NewPosition(0, 0, 0),
											types.NewPosition(0, 0, 0),
										)),
										ast.NewIdent("field", types.NewLocation(
											types.NewPosition(0, 0, 0),
											types.NewPosition(0, 0, 0),
										)),
										ast.NewIdent("field", types.NewLocation(
											types.NewPosition(0, 0, 0),
											types.NewPosition(0, 0, 0),
										)),
									}),
									types.NewLocation(
										types.NewPosition(0, 0, 0),
										types.NewPosition(0, 0, 0),
									),
								),
								ast.NewVar([]ast.Ident{
									ast.NewIdent("kafka", types.NewLocation(
										types.NewPosition(0, 0, 0),
										types.NewPosition(0, 0, 0),
									)),
									ast.NewIdent("field", types.NewLocation(
										types.NewPosition(0, 0, 0),
										types.NewPosition(0, 0, 0),
									)),
									ast.NewIdent("field", types.NewLocation(
										types.NewPosition(0, 0, 0),
										types.NewPosition(0, 0, 0),
									)),
								}),
							},
							types.NewLocation(
								types.NewPosition(0, 0, 0),
								types.NewPosition(0, 0, 0),
							),
						),
					),
					ast.NewKV(
						ast.NewIdent("arrayWithNestedObject", types.NewLocation(
							types.NewPosition(0, 0, 0),
							types.NewPosition(0, 0, 0),
						)),
						ast.NewArray(
							[]ast.Expression{
								ast.NewObject(
									[]ast.Entry{
										ast.NewKV(
											ast.NewIdent("field1", types.NewLocation(
												types.NewPosition(0, 0, 0),
												types.NewPosition(0, 0, 0),
											)),
											testast.MustNewIntWithLocation(t, "99999999999", types.NewLocation(
												types.NewPosition(0, 0, 0),
												types.NewPosition(0, 0, 0),
											)),
										),
										ast.NewSpread(
											ast.NewVar([]ast.Ident{
												ast.NewIdent("kafka", types.NewLocation(
													types.NewPosition(0, 0, 0),
													types.NewPosition(0, 0, 0),
												)),
												ast.NewIdent("field", types.NewLocation(
													types.NewPosition(0, 0, 0),
													types.NewPosition(0, 0, 0),
												)),
												ast.NewIdent("field", types.NewLocation(
													types.NewPosition(0, 0, 0),
													types.NewPosition(0, 0, 0),
												)),
											}),
											types.NewLocation(
												types.NewPosition(0, 0, 0),
												types.NewPosition(0, 0, 0),
											),
										),
									},
									types.NewLocation(
										types.NewPosition(0, 0, 0),
										types.NewPosition(0, 0, 0),
									),
								),
								ast.NewEnv(
									ast.NewIdent("POSTGRES_PASSWORD", types.NewLocation(
										types.NewPosition(0, 0, 0),
										types.NewPosition(0, 0, 0),
									)),
									types.NewLocation(
										types.NewPosition(0, 0, 0),
										types.NewPosition(0, 0, 0),
									),
								),
							},
							types.NewLocation(
								types.NewPosition(0, 0, 0),
								types.NewPosition(0, 0, 0),
							),
						),
					),
				},
				types.NewLocation(
					types.NewPosition(6, 1, 72),
					types.NewPosition(28, 1, 609),
				),
			),
		),
	)

	lex := lexer.New(input)
	tokens, err := lex.Tokenize()
	require.NoError(t, err)

	mover := tokenmover.New(tokens)

	p := parser.New(mover)

	a, err := p.Parse()
	require.NoError(t, err)
	require.Equal(t, expectedAst.Root().Imports(), a.Root().Imports())
}
