package parser_test

import (
	"testing"

	"github.com/atmxlab/atmcfg/internal/parser"
	"github.com/atmxlab/atmcfg/internal/parser/ast"
	"github.com/atmxlab/atmcfg/internal/test"
	"github.com/atmxlab/atmcfg/internal/test/gen"
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

}
