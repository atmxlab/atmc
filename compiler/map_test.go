package compiler_test

import (
	"testing"

	"github.com/atmxlab/atmc/compiler"
	linkedast "github.com/atmxlab/atmc/linker/ast"
	"github.com/atmxlab/atmc/test/testlinkedast"
	"github.com/atmxlab/atmc/test/testutils"
	"github.com/stretchr/testify/require"
)

func TestMapCompiler(t *testing.T) {
	t.Parallel()

	t.Run("happy_path", func(t *testing.T) {
		t.Parallel()

		a := testlinkedast.NewBuilder().
			Object(func(ob *testlinkedast.ObjectBuilder) {
				ob.KV2("a", linkedast.NewInt(1))
				ob.KV2("b", linkedast.NewInt(2))
				ob.KV2("c", testlinkedast.NewObjectBuilder().
					KV2("d", linkedast.NewInt(3)).
					KV2(
						"e",
						testlinkedast.NewArrayBuilder().
							Element(linkedast.NewInt(4)).
							Element(linkedast.NewInt(5)).
							Element(linkedast.NewInt(6)).
							Element(linkedast.NewString("123456")).
							Build()).
					Build(),
				)
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

		actualMap := make(map[string]any)
		err := mc.Compile(actualMap, a)
		require.NoError(t, err)
		testutils.AssertEmptyDiff(t, expectedMap, actualMap)
	})
}
