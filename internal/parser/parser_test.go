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
				ast.NewObject([]ast.KeyValue{}),
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
				ast.NewObject([]ast.KeyValue{}),
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
				ast.NewObject([]ast.KeyValue{
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
				ast.NewObject([]ast.KeyValue{
					ast.NewKeyValue(ast.NewKey("field1"), testast.MustNewInt(t, "123")),
					ast.NewKeyValue(ast.NewKey("field2"), testast.MustNewFloat(t, "123.123")),
				}),
			),
		)

		a, err := p.Parse()
		require.NoError(t, err)
		require.Equal(t, expectedAst, a)
	})

}
