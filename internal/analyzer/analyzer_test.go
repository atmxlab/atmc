package analyzer_test

import (
	"testing"

	semantic2 "github.com/atmxlab/atmcfg/internal/analyzer"
	"github.com/atmxlab/atmcfg/internal/test/testast"
	"github.com/stretchr/testify/require"
)

func TestAnalyzer_Analyze(t *testing.T) {
	t.Parallel()

	t.Run("unused_import_variable", func(t *testing.T) {
		t.Parallel()

		b := testast.NewAstBuilder()

		b.Import(func(ib *testast.ImportBuilder) {
			ib.Name(testast.NewIdent("import1"))
			ib.Path(testast.NewPath("/path1/import1.atmc"))
		})

		a := semantic2.New()

		err := a.Analyze(b.Build())
		require.ErrorIs(t, err, semantic2.ErrUnusedVariable)
	})

	t.Run("undefined_spread_variable", func(t *testing.T) {
		t.Parallel()

		b := testast.NewAstBuilder()

		b.Object(func(ob *testast.ObjectBuilder) {
			ob.Spread(func(sb *testast.SpreadBuilder) {
				sb.Var(func(vb *testast.VarBuilder) {
					vb.Part(testast.NewIdent("var_part1"))
					vb.Part(testast.NewIdent("var_part2"))
				})
			})
		})

		a := semantic2.New()

		err := a.Analyze(b.Build())
		require.ErrorIs(t, err, semantic2.ErrUndefinedVariable)
	})

	t.Run("undefined_variable", func(t *testing.T) {
		t.Parallel()

		b := testast.NewAstBuilder()

		b.Object(func(ob *testast.ObjectBuilder) {
			ob.KV(func(kb *testast.KVBuilder) {
				kb.
					Key(testast.NewIdent("key1")).
					Var(func(vb *testast.VarBuilder) {
						vb.Part(testast.NewIdent("var1"))
					})
			})
		})

		a := semantic2.New()

		err := a.Analyze(b.Build())
		require.ErrorIs(t, err, semantic2.ErrUndefinedVariable)
	})
}
