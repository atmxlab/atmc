package semantic_test

import (
	"testing"

	"github.com/atmxlab/atmcfg/internal/compiler/semantic"
	"github.com/atmxlab/atmcfg/internal/test/testast"
	"github.com/stretchr/testify/require"
)

func TestAnalyzer_Analyze(t *testing.T) {
	t.Parallel()

	t.Run("unused import variable", func(t *testing.T) {
		t.Parallel()

		b := testast.NewAstBuilder()

		b.Import(func(ib *testast.ImportBuilder) {
			ib.Name(testast.NewIdent("import1"))
			ib.Path(testast.NewPath("/path1/import1.atmc"))
		})

		a := semantic.NewAnalyzer()

		err := a.Analyze(b.Build())
		require.ErrorIs(t, err, semantic.ErrUnusedVariable)
	})

	t.Run("undefined spread variable", func(t *testing.T) {
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

		a := semantic.NewAnalyzer()

		err := a.Analyze(b.Build())
		require.ErrorIs(t, err, semantic.ErrUndefinedVariable)
	})
}
