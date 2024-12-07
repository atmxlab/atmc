package parser_test

import (
	"testing"

	"github.com/atmxlab/atmcfg/internal/parser"
	"github.com/atmxlab/atmcfg/internal/parser/ast"
	"github.com/atmxlab/atmcfg/internal/test"
	"github.com/atmxlab/atmcfg/internal/test/gen"
	"github.com/atmxlab/atmcfg/internal/test/testast"
	"github.com/atmxlab/atmcfg/internal/types/token"
	"github.com/stretchr/testify/require"
)

func TestParse(t *testing.T) {
	t.Parallel()

	t.Run("with import and empty object", func(t *testing.T) {
		t.Parallel()

		tokens := []token.Token{
			token.New(token.Ident, "importVar", gen.RandPosition()),
			token.New(token.Path, "./path/to/file.atmc", gen.RandPosition()),
			token.New(token.LBrace, "{", gen.RandPosition()),
			token.New(token.RBrace, "{", gen.RandPosition()),
		}

		testLexer := test.NewLexer(t, tokens)

		p := parser.New(testLexer)

		expectedAst := ast.NewAst(
			ast.NewFile(
				[]ast.Import{
					ast.NewImport(ast.NewName("importVar"), ast.NewPath("./path/to/file.atmc")),
				},
				ast.NewObject([]ast.Spread{}, []ast.KeyValue{}),
			),
		)

		a, err := p.Parse()
		require.NoError(t, err)
		require.Equal(t, expectedAst, a)
	})

	t.Run("without import and empty object", func(t *testing.T) {
		t.Parallel()

		tokens := []token.Token{
			token.New(token.LBrace, "{", gen.RandPosition()),
			token.New(token.RBrace, "}", gen.RandPosition()),
		}

		testLexer := test.NewLexer(t, tokens)

		p := parser.New(testLexer)

		expectedAst := ast.NewAst(
			ast.NewFile(
				[]ast.Import{},
				ast.NewObject([]ast.Spread{}, []ast.KeyValue{}),
			),
		)

		a, err := p.Parse()
		require.NoError(t, err)
		require.Equal(t, expectedAst, a)
	})

	t.Run("without import and non empty object", func(t *testing.T) {
		t.Parallel()

		tokens := []token.Token{
			token.New(token.LBrace, "{", gen.RandPosition()),
			token.New(token.Ident, "field1", gen.RandPosition()),
			token.New(token.Colon, ":", gen.RandPosition()),
			token.New(token.Int, "123", gen.RandPosition()),
			token.New(token.Ident, "field2", gen.RandPosition()),
			token.New(token.Colon, ":", gen.RandPosition()),
			token.New(token.Float, "123.123", gen.RandPosition()),
			token.New(token.RBrace, "}", gen.RandPosition()),
		}

		testLexer := test.NewLexer(t, tokens)

		p := parser.New(testLexer)

		expectedAst := ast.NewAst(
			ast.NewFile(
				[]ast.Import{},
				ast.NewObject([]ast.Spread{}, []ast.KeyValue{
					ast.NewKeyValue(ast.NewKey("field1"), testast.MustNewInt(t, "123")),
					ast.NewKeyValue(ast.NewKey("field2"), testast.MustNewFloat(t, "123.123")),
				}),
			),
		)

		a, err := p.Parse()
		require.NoError(t, err)
		require.Equal(t, expectedAst, a)
	})

	t.Run("with many import and non empty object", func(t *testing.T) {
		t.Parallel()

		tokens := []token.Token{
			token.New(token.Ident, "importVar1", gen.RandPosition()),
			token.New(token.Path, "./path/to/file1.atmc", gen.RandPosition()),
			token.New(token.Ident, "importVar2", gen.RandPosition()),
			token.New(token.Path, "./path/to/file2.atmc", gen.RandPosition()),
			token.New(token.LBrace, "{", gen.RandPosition()),
			token.New(token.Ident, "field1", gen.RandPosition()),
			token.New(token.Colon, ":", gen.RandPosition()),
			token.New(token.Int, "123", gen.RandPosition()),
			token.New(token.Ident, "field2", gen.RandPosition()),
			token.New(token.Colon, ":", gen.RandPosition()),
			token.New(token.Float, "123.123", gen.RandPosition()),
			token.New(token.RBrace, "}", gen.RandPosition()),
		}

		testLexer := test.NewLexer(t, tokens)

		p := parser.New(testLexer)

		expectedAst := ast.NewAst(
			ast.NewFile(
				[]ast.Import{
					ast.NewImport(ast.NewName("importVar1"), ast.NewPath("./path/to/file1.atmc")),
					ast.NewImport(ast.NewName("importVar2"), ast.NewPath("./path/to/file2.atmc")),
				},
				ast.NewObject([]ast.Spread{}, []ast.KeyValue{
					ast.NewKeyValue(ast.NewKey("field1"), testast.MustNewInt(t, "123")),
					ast.NewKeyValue(ast.NewKey("field2"), testast.MustNewFloat(t, "123.123")),
				}),
			),
		)

		a, err := p.Parse()
		require.NoError(t, err)
		require.Equal(t, expectedAst, a)
	})

	t.Run("imports after empty object", func(t *testing.T) {
		t.Parallel()

		tokens := []token.Token{
			token.New(token.LBrace, "{", gen.RandPosition()),
			token.New(token.RBrace, "}", gen.RandPosition()),
			token.New(token.Ident, "importVar1", gen.RandPosition()),
			token.New(token.Path, "./path/to/file1.atmc", gen.RandPosition()),
			token.New(token.Ident, "importVar2", gen.RandPosition()),
			token.New(token.Path, "./path/to/file2.atmc", gen.RandPosition()),
		}

		testLexer := test.NewLexer(t, tokens)

		p := parser.New(testLexer)

		// In general, these are meaningless imports, but we will check their meaning at the compilation stage
		expectedAst := ast.NewAst(
			ast.NewFile(
				[]ast.Import{
					ast.NewImport(ast.NewName("importVar1"), ast.NewPath("./path/to/file1.atmc")),
					ast.NewImport(ast.NewName("importVar2"), ast.NewPath("./path/to/file2.atmc")),
				},
				ast.NewObject([]ast.Spread{}, []ast.KeyValue{}),
			),
		)

		a, err := p.Parse()
		require.NoError(t, err)
		require.Equal(t, expectedAst, a)
	})

	t.Run("with imports and without object", func(t *testing.T) {
		t.Parallel()

		tokens := []token.Token{
			token.New(token.Ident, "importVar1", gen.RandPosition()),
			token.New(token.Path, "./path/to/file1.atmc", gen.RandPosition()),
			token.New(token.Ident, "importVar2", gen.RandPosition()),
			token.New(token.Path, "./path/to/file2.atmc", gen.RandPosition()),
		}

		expectedAst := ast.NewAst(
			ast.NewFile(
				[]ast.Import{
					ast.NewImport(ast.NewName("importVar1"), ast.NewPath("./path/to/file1.atmc")),
					ast.NewImport(ast.NewName("importVar2"), ast.NewPath("./path/to/file2.atmc")),
				},
				ast.NewObject(nil, nil),
			),
		)

		testLexer := test.NewLexer(t, tokens)

		p := parser.New(testLexer)

		a, err := p.Parse()
		require.NoError(t, err)
		require.Equal(t, expectedAst, a)
	})

	t.Run("invalid token sequence", func(t *testing.T) {
		t.Parallel()

		tokens := []token.Token{
			token.New(token.Ident, "importVar1", gen.RandPosition()),
			token.New(token.Ident, "importVar2", gen.RandPosition()),
		}

		testLexer := test.NewLexer(t, tokens)

		p := parser.New(testLexer)

		a, err := p.Parse()
		require.Error(t, err)
		require.ErrorIs(t, err, parser.ErrUnexpectedToken)
		require.Empty(t, a)
	})

	t.Run("with empty array value", func(t *testing.T) {
		t.Parallel()

		tokens := []token.Token{
			token.New(token.LBrace, "{", gen.RandPosition()),
			token.New(token.Ident, "arrIdent1", gen.RandPosition()),
			token.New(token.Colon, ":", gen.RandPosition()),
			token.New(token.LBracket, "[", gen.RandPosition()),
			token.New(token.RBracket, "]", gen.RandPosition()),
			token.New(token.RBrace, "}", gen.RandPosition()),
		}

		expectedAst := ast.NewAst(
			ast.NewFile(
				[]ast.Import{},
				ast.NewObject([]ast.Spread{}, []ast.KeyValue{
					ast.NewKeyValue(ast.NewKey("arrIdent1"), ast.NewArray([]ast.Node{})),
				}),
			),
		)

		testLexer := test.NewLexer(t, tokens)

		p := parser.New(testLexer)

		a, err := p.Parse()
		require.NoError(t, err)
		require.Equal(t, expectedAst, a)
	})

	t.Run("with non empty array with int type values", func(t *testing.T) {
		t.Parallel()

		tokens := []token.Token{
			token.New(token.LBrace, "{", gen.RandPosition()),
			token.New(token.Ident, "arrIdent1", gen.RandPosition()),
			token.New(token.Colon, ":", gen.RandPosition()),
			token.New(token.LBracket, "[", gen.RandPosition()),
			token.New(token.Int, "123", gen.RandPosition()),
			token.New(token.Comma, ",", gen.RandPosition()),
			token.New(token.Int, "456", gen.RandPosition()),
			token.New(token.Comma, ",", gen.RandPosition()),
			token.New(token.Int, "789", gen.RandPosition()),
			token.New(token.Comma, ",", gen.RandPosition()),
			token.New(token.RBracket, "]", gen.RandPosition()),
			token.New(token.RBrace, "}", gen.RandPosition()),
		}

		expectedAst := ast.NewAst(
			ast.NewFile(
				[]ast.Import{},
				ast.NewObject([]ast.Spread{}, []ast.KeyValue{
					ast.NewKeyValue(ast.NewKey("arrIdent1"), ast.NewArray([]ast.Node{
						testast.MustNewInt(t, "123"),
						testast.MustNewInt(t, "456"),
						testast.MustNewInt(t, "789"),
					})),
				}),
			),
		)

		testLexer := test.NewLexer(t, tokens)

		p := parser.New(testLexer)

		a, err := p.Parse()
		require.NoError(t, err)
		require.Equal(t, expectedAst, a)
	})

	t.Run("with non empty array with difference type values", func(t *testing.T) {
		t.Parallel()

		tokens := []token.Token{
			token.New(token.LBrace, "{", gen.RandPosition()),
			token.New(token.Ident, "arrIdent1", gen.RandPosition()),
			token.New(token.Colon, ":", gen.RandPosition()),
			token.New(token.LBracket, "[", gen.RandPosition()),
			token.New(token.Int, "123", gen.RandPosition()),
			token.New(token.Comma, ",", gen.RandPosition()),
			token.New(token.Float, "456.123123", gen.RandPosition()),
			token.New(token.Comma, ",", gen.RandPosition()),
			token.New(token.String, `text text`, gen.RandPosition()),
			token.New(token.Comma, ",", gen.RandPosition()),
			token.New(token.Bool, "true", gen.RandPosition()),
			token.New(token.Comma, ",", gen.RandPosition()),
			token.New(token.Bool, "false", gen.RandPosition()),
			token.New(token.RBracket, "]", gen.RandPosition()),
			token.New(token.RBrace, "}", gen.RandPosition()),
		}

		expectedAst := ast.NewAst(
			ast.NewFile(
				[]ast.Import{},
				ast.NewObject([]ast.Spread{}, []ast.KeyValue{
					ast.NewKeyValue(ast.NewKey("arrIdent1"), ast.NewArray([]ast.Node{
						testast.MustNewInt(t, "123"),
						testast.MustNewFloat(t, "456.123123"),
						ast.NewString(`text text`),
						testast.MustNewBool(t, "true"),
						testast.MustNewBool(t, "false"),
					})),
				}),
			),
		)

		testLexer := test.NewLexer(t, tokens)

		p := parser.New(testLexer)

		a, err := p.Parse()
		require.NoError(t, err)
		require.Equal(t, expectedAst, a)
	})

	t.Run("with non empty array with object type values", func(t *testing.T) {
		t.Parallel()

		tokens := []token.Token{
			token.New(token.LBrace, "{", gen.RandPosition()),
			token.New(token.Ident, "arrIdent1", gen.RandPosition()),
			token.New(token.Colon, ":", gen.RandPosition()),
			token.New(token.LBracket, "[", gen.RandPosition()),
			token.New(token.LBrace, "{", gen.RandPosition()),
			token.New(token.Ident, "nestedIdent1", gen.RandPosition()),
			token.New(token.Colon, ":", gen.RandPosition()),
			token.New(token.Int, "123", gen.RandPosition()),
			token.New(token.RBrace, "}", gen.RandPosition()),
			token.New(token.Comma, ",", gen.RandPosition()), // TODO: тут если пропустить запятую, нужно выдавать точную позицию ошибки.
			token.New(token.LBrace, "{", gen.RandPosition()),
			// Objects are different. We will check this at the compilation stage
			token.New(token.Ident, "nestedIdent2", gen.RandPosition()),
			token.New(token.Colon, ":", gen.RandPosition()),
			token.New(token.Int, "321", gen.RandPosition()),
			token.New(token.RBrace, "}", gen.RandPosition()),
			token.New(token.RBracket, "]", gen.RandPosition()),
			token.New(token.RBrace, "}", gen.RandPosition()),
		}

		expectedAst := ast.NewAst(
			ast.NewFile(
				[]ast.Import{},
				ast.NewObject([]ast.Spread{}, []ast.KeyValue{
					ast.NewKeyValue(
						ast.NewKey("arrIdent1"),
						ast.NewArray([]ast.Node{
							ast.NewObject(
								[]ast.Spread{},
								[]ast.KeyValue{
									ast.NewKeyValue(ast.NewKey("nestedIdent1"), testast.MustNewInt(t, "123")),
								},
							),
							ast.NewObject(
								[]ast.Spread{},
								[]ast.KeyValue{
									ast.NewKeyValue(ast.NewKey("nestedIdent2"), testast.MustNewInt(t, "321")),
								},
							),
						}),
					),
				}),
			),
		)

		testLexer := test.NewLexer(t, tokens)

		p := parser.New(testLexer)

		a, err := p.Parse()
		require.NoError(t, err)
		require.Equal(t, expectedAst, a)
	})

	t.Run("with non empty array with array type values", func(t *testing.T) {
		t.Parallel()

		tokens := []token.Token{
			token.New(token.LBrace, "{", gen.RandPosition()),
			token.New(token.Ident, "arrIdent1", gen.RandPosition()),
			token.New(token.Colon, ":", gen.RandPosition()),
			token.New(token.LBracket, "[", gen.RandPosition()),
			token.New(token.LBracket, "[", gen.RandPosition()),
			token.New(token.Int, "123", gen.RandPosition()),
			token.New(token.RBracket, "]", gen.RandPosition()),
			token.New(token.Comma, ",", gen.RandPosition()),
			token.New(token.LBracket, "[", gen.RandPosition()),
			token.New(token.Int, "321", gen.RandPosition()),
			token.New(token.RBracket, "]", gen.RandPosition()),
			token.New(token.RBracket, "]", gen.RandPosition()),
			token.New(token.RBrace, "}", gen.RandPosition()),
		}

		expectedAst := ast.NewAst(
			ast.NewFile(
				[]ast.Import{},
				ast.NewObject([]ast.Spread{}, []ast.KeyValue{
					ast.NewKeyValue(
						ast.NewKey("arrIdent1"),
						ast.NewArray([]ast.Node{
							ast.NewArray([]ast.Node{
								testast.MustNewInt(t, "123"),
							}),
							ast.NewArray([]ast.Node{
								testast.MustNewInt(t, "321"),
							}),
						}),
					),
				}),
			),
		)

		testLexer := test.NewLexer(t, tokens)

		p := parser.New(testLexer)

		a, err := p.Parse()
		require.NoError(t, err)
		require.Equal(t, expectedAst, a)
	})

	t.Run("with spread in object", func(t *testing.T) {
		t.Parallel()

		tokens := []token.Token{
			token.New(token.LBrace, "{", gen.RandPosition()),
			token.New(token.Ident, "common", gen.RandPosition()),
			token.New(token.Spread, "...", gen.RandPosition()),
			token.New(token.Ident, "field1", gen.RandPosition()),
			token.New(token.Colon, ":", gen.RandPosition()),
			token.New(token.Int, "123", gen.RandPosition()),
			token.New(token.Ident, "field2", gen.RandPosition()),
			token.New(token.Colon, ":", gen.RandPosition()),
			token.New(token.Float, "123.123", gen.RandPosition()),
			token.New(token.RBrace, "}", gen.RandPosition()),
		}

		testLexer := test.NewLexer(t, tokens)

		p := parser.New(testLexer)

		expectedAst := ast.NewAst(
			ast.NewFile(
				[]ast.Import{},
				ast.NewObject(
					[]ast.Spread{
						ast.NewSpread(ast.NewVar([]ast.Ident{ast.NewName("common")})),
					},
					[]ast.KeyValue{
						ast.NewKeyValue(ast.NewKey("field1"), testast.MustNewInt(t, "123")),
						ast.NewKeyValue(ast.NewKey("field2"), testast.MustNewFloat(t, "123.123")),
					}),
			),
		)

		a, err := p.Parse()
		require.NoError(t, err)
		require.Equal(t, expectedAst, a)
	})

	t.Run("with nested var spread in object", func(t *testing.T) {
		t.Parallel()

		tokens := []token.Token{
			token.New(token.LBrace, "{", gen.RandPosition()),
			token.New(token.Ident, "common", gen.RandPosition()),
			token.New(token.Dot, ".", gen.RandPosition()),
			token.New(token.Ident, "nested1", gen.RandPosition()),
			token.New(token.Spread, "...", gen.RandPosition()),
			token.New(token.Ident, "field1", gen.RandPosition()),
			token.New(token.Colon, ":", gen.RandPosition()),
			token.New(token.Int, "123", gen.RandPosition()),
			token.New(token.Ident, "field2", gen.RandPosition()),
			token.New(token.Colon, ":", gen.RandPosition()),
			token.New(token.Float, "123.123", gen.RandPosition()),
			token.New(token.RBrace, "}", gen.RandPosition()),
		}

		testLexer := test.NewLexer(t, tokens)

		p := parser.New(testLexer)

		expectedAst := ast.NewAst(
			ast.NewFile(
				[]ast.Import{},
				ast.NewObject(
					[]ast.Spread{
						ast.NewSpread(ast.NewVar([]ast.Ident{
							ast.NewName("common"),
							ast.NewName("nested1"),
						})),
					},
					[]ast.KeyValue{
						ast.NewKeyValue(ast.NewKey("field1"), testast.MustNewInt(t, "123")),
						ast.NewKeyValue(ast.NewKey("field2"), testast.MustNewFloat(t, "123.123")),
					}),
			),
		)

		a, err := p.Parse()
		require.NoError(t, err)
		require.Equal(t, expectedAst, a)
	})

	t.Run("with nested var without spread in object", func(t *testing.T) {
		t.Parallel()

		tokens := []token.Token{
			token.New(token.LBrace, "{", gen.RandPosition()),
			token.New(token.Ident, "common", gen.RandPosition()),
			token.New(token.Dot, ".", gen.RandPosition()),
			token.New(token.Ident, "nested1", gen.RandPosition()),
			token.New(token.Ident, "field1", gen.RandPosition()),
			token.New(token.Colon, ":", gen.RandPosition()),
			token.New(token.Int, "123", gen.RandPosition()),
			token.New(token.Ident, "field2", gen.RandPosition()),
			token.New(token.Colon, ":", gen.RandPosition()),
			token.New(token.Float, "123.123", gen.RandPosition()),
			token.New(token.RBrace, "}", gen.RandPosition()),
		}

		testLexer := test.NewLexer(t, tokens)

		p := parser.New(testLexer)

		a, err := p.Parse()
		require.Error(t, err)
		require.ErrorIs(t, err, parser.ErrUnexpectedToken)
		require.Empty(t, a)
	})

	t.Run("with var in array", func(t *testing.T) {
		t.Parallel()

		tokens := []token.Token{
			token.New(token.LBrace, "{", gen.RandPosition()),
			token.New(token.Ident, "arrIdent1", gen.RandPosition()),
			token.New(token.Colon, ":", gen.RandPosition()),
			token.New(token.LBracket, "[", gen.RandPosition()),
			token.New(token.Ident, "common", gen.RandPosition()),
			token.New(token.Dot, ".", gen.RandPosition()),
			token.New(token.Ident, "nested", gen.RandPosition()),
			token.New(token.RBracket, "]", gen.RandPosition()),
			token.New(token.RBrace, "}", gen.RandPosition()),
		}

		expectedAst := ast.NewAst(
			ast.NewFile(
				[]ast.Import{},
				ast.NewObject([]ast.Spread{}, []ast.KeyValue{
					ast.NewKeyValue(
						ast.NewKey("arrIdent1"),
						ast.NewArray([]ast.Node{
							ast.NewVar([]ast.Ident{
								ast.NewName("common"),
								ast.NewName("nested"),
							}),
						}),
					),
				}),
			),
		)

		testLexer := test.NewLexer(t, tokens)

		p := parser.New(testLexer)

		a, err := p.Parse()
		require.NoError(t, err)
		require.Equal(t, expectedAst, a)
	})

	t.Run("with spread in array", func(t *testing.T) {
		t.Parallel()

		tokens := []token.Token{
			token.New(token.LBrace, "{", gen.RandPosition()),
			token.New(token.Ident, "arrIdent1", gen.RandPosition()),
			token.New(token.Colon, ":", gen.RandPosition()),
			token.New(token.LBracket, "[", gen.RandPosition()),
			token.New(token.Ident, "common", gen.RandPosition()),
			token.New(token.Dot, ".", gen.RandPosition()),
			token.New(token.Ident, "nested", gen.RandPosition()),
			token.New(token.Dot, ".", gen.RandPosition()),
			token.New(token.Ident, "nested2", gen.RandPosition()),
			token.New(token.Spread, "...", gen.RandPosition()),
			token.New(token.RBracket, "]", gen.RandPosition()),
			token.New(token.RBrace, "}", gen.RandPosition()),
		}

		expectedAst := ast.NewAst(
			ast.NewFile(
				[]ast.Import{},
				ast.NewObject([]ast.Spread{}, []ast.KeyValue{
					ast.NewKeyValue(
						ast.NewKey("arrIdent1"),
						ast.NewArray([]ast.Node{
							ast.NewSpread(ast.NewVar([]ast.Ident{
								ast.NewName("common"),
								ast.NewName("nested"),
								ast.NewName("nested2"),
							})),
						}),
					),
				}),
			),
		)

		testLexer := test.NewLexer(t, tokens)

		p := parser.New(testLexer)

		a, err := p.Parse()
		require.NoError(t, err)
		require.Equal(t, expectedAst, a)
	})

	t.Run("with many difference types in array", func(t *testing.T) {
		t.Parallel()

		tokens := []token.Token{
			token.New(token.LBrace, "{", gen.RandPosition()),
			token.New(token.Ident, "arrIdent1", gen.RandPosition()),
			token.New(token.Colon, ":", gen.RandPosition()),
			token.New(token.LBracket, "[", gen.RandPosition()),
			token.New(token.Ident, "common", gen.RandPosition()),
			token.New(token.Dot, ".", gen.RandPosition()),
			token.New(token.Ident, "nested", gen.RandPosition()),
			token.New(token.Dot, ".", gen.RandPosition()),
			token.New(token.Ident, "nested2", gen.RandPosition()),
			token.New(token.Spread, "...", gen.RandPosition()),
			token.New(token.Comma, ",", gen.RandPosition()),
			token.New(token.Int, "123", gen.RandPosition()),
			token.New(token.Comma, ",", gen.RandPosition()),
			token.New(token.String, `test test`, gen.RandPosition()),
			token.New(token.Comma, ",", gen.RandPosition()),
			token.New(token.Bool, "true", gen.RandPosition()),
			token.New(token.Comma, ",", gen.RandPosition()),
			token.New(token.Ident, "common2", gen.RandPosition()),
			token.New(token.Dot, ".", gen.RandPosition()),
			token.New(token.Ident, "nested2", gen.RandPosition()),
			token.New(token.RBracket, "]", gen.RandPosition()),
			token.New(token.RBrace, "}", gen.RandPosition()),
		}

		expectedAst := ast.NewAst(
			ast.NewFile(
				[]ast.Import{},
				ast.NewObject([]ast.Spread{}, []ast.KeyValue{
					ast.NewKeyValue(
						ast.NewKey("arrIdent1"),
						ast.NewArray([]ast.Node{
							ast.NewSpread(ast.NewVar([]ast.Ident{
								ast.NewName("common"),
								ast.NewName("nested"),
								ast.NewName("nested2"),
							})),
							testast.MustNewInt(t, "123"),
							ast.NewString("test test"),
							testast.MustNewBool(t, "true"),
							ast.NewVar([]ast.Ident{
								ast.NewName("common2"),
								ast.NewName("nested2"),
							}),
						}),
					),
				}),
			),
		)

		testLexer := test.NewLexer(t, tokens)

		p := parser.New(testLexer)

		a, err := p.Parse()
		require.NoError(t, err)
		require.Equal(t, expectedAst, a)
	})

	t.Run("with many difference types in array", func(t *testing.T) {
		t.Parallel()

		tokens := []token.Token{
			token.New(token.LBrace, "{", gen.RandPosition()),
			token.New(token.Ident, "field", gen.RandPosition()),
			token.New(token.Colon, ":", gen.RandPosition()),
			token.New(token.LBracket, "[", gen.RandPosition()),
			token.New(token.LBrace, "{", gen.RandPosition()),
			token.New(token.Ident, "field", gen.RandPosition()),
			token.New(token.Colon, ":", gen.RandPosition()),
			token.New(token.Int, "99999999999", gen.RandPosition()),
			token.New(token.Comma, ",", gen.RandPosition()),
			token.New(token.Ident, "kafka", gen.RandPosition()),
			token.New(token.Dot, ".", gen.RandPosition()),
			token.New(token.Ident, "field", gen.RandPosition()),
			token.New(token.Dot, ".", gen.RandPosition()),
			token.New(token.Ident, "field", gen.RandPosition()),
			token.New(token.Spread, "...", gen.RandPosition()),
			token.New(token.RBrace, "}", gen.RandPosition()),
			token.New(token.RBracket, "]", gen.RandPosition()),
			token.New(token.RBrace, "}", gen.RandPosition()),
		}

		expectedAst := ast.NewAst(
			ast.NewFile(
				[]ast.Import{},
				ast.NewObject([]ast.Spread{}, []ast.KeyValue{
					ast.NewKeyValue(ast.NewKey("field"), ast.NewArray([]ast.Node{
						ast.NewObject(
							[]ast.Spread{
								ast.NewSpread(ast.NewVar([]ast.Ident{
									ast.NewName("kafka"),
									ast.NewName("field"),
									ast.NewName("field"),
								})),
							},
							[]ast.KeyValue{
								ast.NewKeyValue(ast.NewKey("field"), testast.MustNewInt(t, "99999999999")),
							},
						),
					})),
				}),
			),
		)

		testLexer := test.NewLexer(t, tokens)

		p := parser.New(testLexer)

		a, err := p.Parse()
		require.NoError(t, err)
		require.Equal(t, expectedAst, a)
	})
}
