package test

import (
	"testing"

	"github.com/atmxlab/atmcfg/internal/lexer"
	"github.com/atmxlab/atmcfg/internal/lexer/tokenmover"
	"github.com/atmxlab/atmcfg/internal/parser"
	"github.com/atmxlab/atmcfg/internal/parser/ast"
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
		arrayWithoutComma: ["test text" 444 -321 123.123]
		arrayWithEnv: [ $POSTGRES_HOST $POSTGRES_DATABASE ]
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
							types.NewPosition(9, 2, 99),
							types.NewPosition(9, 7, 104),
						)),
						testast.MustNewFloatWithLocation(t, "123.123", types.NewLocation(
							types.NewPosition(9, 9, 106),
							types.NewPosition(9, 16, 113),
						)),
					),
					ast.NewKV(
						ast.NewIdent("bool1", types.NewLocation(
							types.NewPosition(10, 2, 116),
							types.NewPosition(10, 7, 121),
						)),
						testast.MustNewBoolWithLocation(t, "true", types.NewLocation(
							types.NewPosition(10, 9, 123),
							types.NewPosition(10, 13, 127),
						)),
					),
					ast.NewKV(
						ast.NewIdent("bool2", types.NewLocation(
							types.NewPosition(11, 2, 130),
							types.NewPosition(11, 7, 135),
						)),
						testast.MustNewBoolWithLocation(t, "false", types.NewLocation(
							types.NewPosition(11, 9, 137),
							types.NewPosition(11, 14, 142),
						)),
					),
					ast.NewKV(
						ast.NewIdent("string", types.NewLocation(
							types.NewPosition(12, 2, 145),
							types.NewPosition(12, 8, 151),
						)),
						ast.NewString(`"test string \"escaped"`, types.NewLocation(
							types.NewPosition(12, 10, 153),
							types.NewPosition(12, 33, 176),
						)),
					),
					ast.NewSpread(
						ast.NewVar([]ast.Ident{
							ast.NewIdent("kafka", types.NewLocation(
								types.NewPosition(13, 2, 179),
								types.NewPosition(13, 7, 184),
							)),
							ast.NewIdent("nested", types.NewLocation(
								types.NewPosition(13, 8, 185),
								types.NewPosition(13, 14, 191),
							)),
						}),
						types.NewLocation(
							types.NewPosition(13, 2, 179),
							types.NewPosition(13, 17, 194),
						),
					),
					ast.NewKV(
						ast.NewIdent("object1", types.NewLocation(
							types.NewPosition(14, 2, 197),
							types.NewPosition(14, 9, 204),
						)),
						ast.NewObject(
							[]ast.Entry{
								ast.NewSpread(
									ast.NewVar([]ast.Ident{
										ast.NewIdent("common", types.NewLocation(
											types.NewPosition(15, 3, 211),
											types.NewPosition(15, 9, 217),
										)),
										ast.NewIdent("field", types.NewLocation(
											types.NewPosition(15, 10, 218),
											types.NewPosition(15, 15, 223),
										)),
										ast.NewIdent("field", types.NewLocation(
											types.NewPosition(15, 16, 224),
											types.NewPosition(15, 21, 229),
										)),
									}),
									types.NewLocation(
										types.NewPosition(15, 3, 211),
										types.NewPosition(15, 24, 232),
									),
								),
								ast.NewKV(
									ast.NewIdent("field1", types.NewLocation(
										types.NewPosition(16, 3, 236),
										types.NewPosition(16, 9, 242),
									)),
									testast.MustNewIntWithLocation(t, "123", types.NewLocation(
										types.NewPosition(16, 11, 244),
										types.NewPosition(16, 14, 247),
									)),
								),
								ast.NewKV(
									ast.NewIdent("field2", types.NewLocation(
										types.NewPosition(17, 3, 251),
										types.NewPosition(17, 9, 257),
									)),
									ast.NewObject(
										[]ast.Entry{
											ast.NewKV(
												ast.NewIdent("field1", types.NewLocation(
													types.NewPosition(18, 4, 265),
													types.NewPosition(18, 10, 271),
												)),
												testast.MustNewIntWithLocation(t, "321", types.NewLocation(
													types.NewPosition(18, 12, 273),
													types.NewPosition(18, 15, 276),
												)),
											),
										},
										types.NewLocation(
											types.NewPosition(17, 11, 259),
											types.NewPosition(19, 4, 281),
										),
									),
								),
							},
							types.NewLocation(
								types.NewPosition(14, 11, 206),
								types.NewPosition(20, 3, 285),
							),
						),
					),
					ast.NewKV(
						ast.NewIdent("array1", types.NewLocation(
							types.NewPosition(21, 2, 288),
							types.NewPosition(21, 8, 294),
						)),
						ast.NewArray(
							[]ast.Expression{
								testast.MustNewIntWithLocation(t, "123", types.NewLocation(
									types.NewPosition(22, 3, 301),
									types.NewPosition(22, 6, 304),
								)),
								testast.MustNewIntWithLocation(t, "321", types.NewLocation(
									types.NewPosition(22, 8, 306),
									types.NewPosition(22, 11, 309),
								)),
								testast.MustNewIntWithLocation(t, "444", types.NewLocation(
									types.NewPosition(23, 3, 314),
									types.NewPosition(23, 6, 317),
								)),
								testast.MustNewIntWithLocation(t, "555", types.NewLocation(
									types.NewPosition(23, 8, 319),
									types.NewPosition(23, 11, 322),
								)),
							},
							types.NewLocation(
								types.NewPosition(21, 10, 296),
								types.NewPosition(24, 3, 326),
							),
						),
					),
					ast.NewKV(
						ast.NewIdent("envVar", types.NewLocation(
							types.NewPosition(25, 2, 329),
							types.NewPosition(25, 8, 335),
						)),
						ast.NewEnv(
							ast.NewIdent("POSTGRES_PASSWORD", types.NewLocation(
								types.NewPosition(25, 11, 338),
								types.NewPosition(25, 28, 355),
							)),
							types.NewLocation(
								types.NewPosition(25, 10, 337),
								types.NewPosition(25, 28, 355),
							),
						),
					),
					ast.NewKV(
						ast.NewIdent("arrayWithComma", types.NewLocation(
							types.NewPosition(26, 2, 358),
							types.NewPosition(26, 16, 372),
						)),
						ast.NewArray(
							[]ast.Expression{
								ast.NewString(`"test text"`, types.NewLocation(
									types.NewPosition(26, 19, 375),
									types.NewPosition(26, 30, 386),
								)),
								testast.MustNewIntWithLocation(t, "444", types.NewLocation(
									types.NewPosition(26, 32, 388),
									types.NewPosition(26, 35, 391),
								)),
								testast.MustNewIntWithLocation(t, "-321", types.NewLocation(
									types.NewPosition(26, 37, 393),
									types.NewPosition(26, 41, 397),
								)),
								testast.MustNewFloatWithLocation(t, "123.123", types.NewLocation(
									types.NewPosition(26, 43, 399),
									types.NewPosition(26, 50, 406),
								)),
								testast.MustNewBoolWithLocation(t, "true", types.NewLocation(
									types.NewPosition(26, 52, 408),
									types.NewPosition(26, 56, 412),
								)),
								testast.MustNewBoolWithLocation(t, "false", types.NewLocation(
									types.NewPosition(26, 58, 414),
									types.NewPosition(26, 63, 419),
								)),
								ast.NewSpread(
									ast.NewVar([]ast.Ident{
										ast.NewIdent("common", types.NewLocation(
											types.NewPosition(26, 65, 421),
											types.NewPosition(26, 71, 427),
										)),
										ast.NewIdent("field", types.NewLocation(
											types.NewPosition(26, 72, 428),
											types.NewPosition(26, 77, 433),
										)),
										ast.NewIdent("field", types.NewLocation(
											types.NewPosition(26, 78, 434),
											types.NewPosition(26, 83, 439),
										)),
									}),
									types.NewLocation(
										types.NewPosition(26, 65, 421),
										types.NewPosition(26, 86, 442),
									),
								),
								ast.NewVar([]ast.Ident{
									ast.NewIdent("kafka", types.NewLocation(
										types.NewPosition(26, 88, 444),
										types.NewPosition(26, 93, 449),
									)),
									ast.NewIdent("field", types.NewLocation(
										types.NewPosition(26, 94, 450),
										types.NewPosition(26, 99, 455),
									)),
									ast.NewIdent("field", types.NewLocation(
										types.NewPosition(26, 100, 456),
										types.NewPosition(26, 105, 461),
									)),
								}),
							},
							types.NewLocation(
								types.NewPosition(26, 18, 374),
								types.NewPosition(26, 106, 462),
							),
						),
					),
					ast.NewKV(
						ast.NewIdent("arrayWithoutComma", types.NewLocation(
							types.NewPosition(27, 2, 465),
							types.NewPosition(27, 19, 482),
						)),
						ast.NewArray(
							[]ast.Expression{
								ast.NewString(`"test text"`, types.NewLocation(
									types.NewPosition(27, 22, 485),
									types.NewPosition(27, 33, 496),
								)),
								testast.MustNewIntWithLocation(t, "444", types.NewLocation(
									types.NewPosition(27, 34, 497),
									types.NewPosition(27, 37, 500),
								)),
								testast.MustNewIntWithLocation(t, "-321", types.NewLocation(
									types.NewPosition(27, 38, 501),
									types.NewPosition(27, 42, 505),
								)),
								testast.MustNewFloatWithLocation(t, "123.123", types.NewLocation(
									types.NewPosition(27, 43, 506),
									types.NewPosition(27, 50, 513),
								)),
							},
							types.NewLocation(
								types.NewPosition(27, 21, 484),
								types.NewPosition(27, 51, 514),
							),
						),
					),
					ast.NewKV(
						ast.NewIdent("arrayWithEnv", types.NewLocation(
							types.NewPosition(28, 2, 517),
							types.NewPosition(28, 14, 529),
						)),
						ast.NewArray(
							[]ast.Expression{
								ast.NewEnv(
									ast.NewIdent("POSTGRES_HOST", types.NewLocation(
										types.NewPosition(28, 19, 534),
										types.NewPosition(28, 32, 547),
									)),
									types.NewLocation(
										types.NewPosition(28, 18, 533),
										types.NewPosition(28, 32, 547),
									),
								),
								ast.NewEnv(
									ast.NewIdent("POSTGRES_DATABASE", types.NewLocation(
										types.NewPosition(28, 34, 549),
										types.NewPosition(28, 51, 566),
									)),
									types.NewLocation(
										types.NewPosition(28, 33, 548),
										types.NewPosition(28, 51, 566),
									),
								),
							},
							types.NewLocation(
								types.NewPosition(28, 16, 531),
								types.NewPosition(28, 53, 568),
							),
						),
					),
				},
				types.NewLocation(
					types.NewPosition(6, 1, 72),
					types.NewPosition(29, 2, 571),
				),
			),
		),
	)

	lex := lexer.New(input)
	tokens, err := lex.Tokenize()
	require.NoError(t, err)

	mover := tokenmover.New(tokens)

	p := parser.New(mover)

	gotAst, err := p.Parse()
	require.NoError(t, err)

	require.Equal(t, expectedAst, gotAst)
}
