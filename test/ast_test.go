package test

//
// func TestAstGenerate(t *testing.T) {
// 	t.Parallel()
//
// 	input := `
// 	common ./dir/common.atmc
// 	kafka ./dir/kafka.atmc
//
// 	{
// 		common...
// 		field1: 123
// 		field2: 123.123
// 		field3: true
// 		field4: false
// 		field5: "test string \"escaped"
// 		kafka.nested...
// 		field6: {
// 			common.field.field...
// 			field1: 123
// 			field2: {
// 				field1: 321
// 			}
// 		}
// 		field7: [
// 			123, 321,
// 			444, 555
// 		]
// 		field8: ["test text", 444, 321, 123.123, true, false, common.field.field..., kafka.field.field]
// 		field9: [{field1: 99999999999, kafka.field.field...}]
// 	}
// `
//
// 	expectedAst := ast.NewAst(ast.NewFile(
// 		[]ast.Import{
// 			ast.NewImport(ast.NewIdent("common"), ast.NewPath("./dir/common.atmc")),
// 			ast.NewImport(ast.NewIdent("kafka"), ast.NewPath("./dir/kafka.atmc")),
// 		},
// 		ast.NewObject(
// 			[]ast.Spread{
// 				ast.NewSpread(ast.NewVar([]ast.Ident{ast.NewIdent("common")})),
// 				ast.NewSpread(ast.NewVar([]ast.Ident{
// 					ast.NewIdent("kafka"),
// 					ast.NewIdent("nested"),
// 				})),
// 			},
// 			[]ast.KeyValue{
// 				ast.NewKeyValue(ast.NewKey("field1"), testast.MustNewInt(t, "123")),
// 				ast.NewKeyValue(ast.NewKey("field2"), testast.MustNewFloat(t, "123.123")),
// 				ast.NewKeyValue(ast.NewKey("field3"), testast.MustNewBool(t, "true")),
// 				ast.NewKeyValue(ast.NewKey("field4"), testast.MustNewBool(t, "false")),
// 				ast.NewKeyValue(ast.NewKey("field5"), ast.NewString("\"test string \\\"escaped\"")),
// 				ast.NewKeyValue(ast.NewKey("field6"), ast.NewObject(
// 					[]ast.Spread{
// 						ast.NewSpread(ast.NewVar([]ast.Ident{
// 							ast.NewIdent("common"),
// 							ast.NewIdent("field"),
// 							ast.NewIdent("field"),
// 						})),
// 					},
// 					[]ast.KeyValue{
// 						ast.NewKeyValue(ast.NewKey("field1"), testast.MustNewInt(t, "123")),
// 						ast.NewKeyValue(ast.NewKey("field2"), ast.NewObject(
// 							[]ast.Spread{},
// 							[]ast.KeyValue{
// 								ast.NewKeyValue(ast.NewKey("field1"), testast.MustNewInt(t, "321")),
// 							},
// 						)),
// 					},
// 				)),
// 				ast.NewKeyValue(ast.NewKey("field7"), ast.NewArray([]ast.Node{
// 					testast.MustNewInt(t, "123"),
// 					testast.MustNewInt(t, "321"),
// 					testast.MustNewInt(t, "444"),
// 					testast.MustNewInt(t, "555"),
// 				})),
// 				ast.NewKeyValue(ast.NewKey("field8"), ast.NewArray([]ast.Node{
// 					ast.NewString("\"test text\""),
// 					testast.MustNewInt(t, "444"),
// 					testast.MustNewInt(t, "321"),
// 					testast.MustNewFloat(t, "123.123"),
// 					testast.MustNewBool(t, "true"),
// 					testast.MustNewBool(t, "false"),
// 					ast.NewSpread(ast.NewVar([]ast.Ident{
// 						ast.NewIdent("common"),
// 						ast.NewIdent("field"),
// 						ast.NewIdent("field"),
// 					})),
// 					ast.NewVar([]ast.Ident{
// 						ast.NewIdent("kafka"),
// 						ast.NewIdent("field"),
// 						ast.NewIdent("field"),
// 					}),
// 				})),
// 				ast.NewKeyValue(ast.NewKey("field9"), ast.NewArray([]ast.Node{
// 					ast.NewObject(
// 						[]ast.Spread{
// 							ast.NewSpread(ast.NewVar([]ast.Ident{
// 								ast.NewIdent("kafka"),
// 								ast.NewIdent("field"),
// 								ast.NewIdent("field"),
// 							})),
// 						},
// 						[]ast.KeyValue{
// 							ast.NewKeyValue(ast.NewKey("field1"), testast.MustNewInt(t, "99999999999")),
// 						},
// 					),
// 				})),
// 			},
// 		),
// 	))
//
// 	lex := lexer.New(input)
// 	tokens, err := lex.Tokenize()
// 	require.NoError(t, err)
//
// 	mover := tokenmover.New(tokens)
//
// 	prsr := parser.New(mover)
//
// 	a, err := prsr.Parse()
// 	require.NoError(t, err)
// 	require.Equal(t, expectedAst, a)
// }
