package compiler_test

import (
	"testing"

	"github.com/atmxlab/atmcfg/internal/compiler"
	linkedast "github.com/atmxlab/atmcfg/internal/linker/ast"
	"github.com/atmxlab/atmcfg/internal/test/testlinkedast"
	"github.com/atmxlab/atmcfg/internal/test/testutils"
	"github.com/stretchr/testify/require"
)

func TestMapCompiler(t *testing.T) {
	t.Parallel()

	t.Run("happy_path", func(t *testing.T) {
		t.Parallel()

		a := testlinkedast.NewBuilder().
			Object(func(ob *testlinkedast.ObjectBuilder) {
				ob.KV(func(kvb *testlinkedast.KVBuilder) {
					kvb.Key("a").Value(linkedast.NewInt(1))
				})
				ob.KV(func(kvb *testlinkedast.KVBuilder) {
					kvb.Key("b").Value(linkedast.NewInt(2))
				})
				ob.KV(func(kvb *testlinkedast.KVBuilder) {
					kvb.Key("c").Value(
						testlinkedast.
							NewObjectBuilder().
							KV(func(kvb *testlinkedast.KVBuilder) {
								kvb.Key("d").Value(linkedast.NewInt(3))
							}).
							KV(func(kvb *testlinkedast.KVBuilder) {
								kvb.Key("e").Value(
									testlinkedast.
										NewArrayBuilder().
										Element(linkedast.NewInt(4)).
										Element(linkedast.NewInt(5)).
										Element(linkedast.NewInt(6)).
										Element(linkedast.NewString("123456")).
										Build(),
								)
							}).
							Build(),
					)
				})
			}).
			Build()

		expectedMap := map[string]any{
			"a": int64(1),
			"b": int64(2),
			"c": map[string]any{
				"d": int64(3),
				"e": []any{
					int64(4), int64(5), int64(6), "123456",
				},
			},
		}

		mc := compiler.NewMapCompiler()

		m, err := mc.Compile(a)
		require.NoError(t, err)
		testutils.AssertEmptyDiff(t, expectedMap, m)
	})
}
