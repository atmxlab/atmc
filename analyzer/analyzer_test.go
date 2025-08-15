package analyzer_test

import (
	"testing"

	"github.com/atmxlab/atmc/analyzer"
	testast2 "github.com/atmxlab/atmc/test/testast"
	"github.com/stretchr/testify/require"
)

func TestAnalyzer_Analyze(t *testing.T) {
	t.Parallel()

	t.Run("unused_import_variable", func(t *testing.T) {
		t.Parallel()

		b := testast2.NewAstBuilder()

		b.Import(func(ib *testast2.ImportBuilder) {
			ib.Name(testast2.NewIdent("import1"))
			ib.Path(testast2.NewPath("/path1/import1.atmc"))
		})

		a := analyzer.New()

		err := a.Analyze(b.Build())
		require.ErrorIs(t, err, analyzer.ErrUnusedVariable)
	})

	t.Run("undefined_spread_variable", func(t *testing.T) {
		t.Parallel()

		b := testast2.NewAstBuilder()

		b.Object(func(ob *testast2.ObjectBuilder) {
			ob.Spread(func(sb *testast2.SpreadBuilder) {
				sb.Var(func(vb *testast2.VarBuilder) {
					vb.Part(testast2.NewIdent("var_part1"))
					vb.Part(testast2.NewIdent("var_part2"))
				})
			})
		})

		a := analyzer.New()

		err := a.Analyze(b.Build())
		require.ErrorIs(t, err, analyzer.ErrUndefinedVariable)
	})

	t.Run("undefined_variable", func(t *testing.T) {
		t.Parallel()

		b := testast2.NewAstBuilder()

		b.Object(func(ob *testast2.ObjectBuilder) {
			ob.KV(func(kb *testast2.KVBuilder) {
				kb.
					Key(testast2.NewIdent("key1")).
					Var(func(vb *testast2.VarBuilder) {
						vb.Part(testast2.NewIdent("var1"))
					})
			})
		})

		a := analyzer.New()

		err := a.Analyze(b.Build())
		require.ErrorIs(t, err, analyzer.ErrUndefinedVariable)
	})
}
