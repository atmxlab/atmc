package acceptance

import (
	"testing"

	linkedast "github.com/atmxlab/atmcfg/internal/linker/ast"
	"github.com/atmxlab/atmcfg/internal/parser"
	"github.com/atmxlab/atmcfg/internal/test"
	"github.com/atmxlab/atmcfg/internal/test/testlinkedast"
	"github.com/atmxlab/atmcfg/internal/test/testos"
	"github.com/atmxlab/atmcfg/pkg/errors"
	"github.com/stretchr/testify/require"
)

func TestProcessor_WithErrors(t *testing.T) {
	t.Parallel()

	t.Run("import_not_found", func(t *testing.T) {
		t.Parallel()

		mainFilePath := "/home/user/config.atmc"

		os := testos.
			NewOSBuilder().
			File(func(fb *testos.FileBuilder) {
				fb.
					Path(mainFilePath).
					Content("var1 ./import.atmc {a: 1, b: 2}")
			}).
			Build()

		app := test.NewApp(t, test.WithOS(os))

		_, err := app.Processor().Process(mainFilePath)

		require.ErrorIs(t, err, errors.ErrNotFound)
	})

	t.Run("empty_content_in_import", func(t *testing.T) {
		t.Parallel()

		mainFilePath := "/home/user/config.atmc"

		os := testos.
			NewOSBuilder().
			File(func(fb *testos.FileBuilder) {
				fb.
					Path(mainFilePath).
					Content("var1 ./import.atmc {a: 1, b: 2}")
			}).
			File(func(fb *testos.FileBuilder) {
				fb.
					Path("/home/user/import.atmc").
					Content("")
			}).
			Build()

		app := test.NewApp(t, test.WithOS(os))

		_, err := app.Processor().Process(mainFilePath)
		require.ErrorIs(t, err, parser.ErrTokenNotExist)
	})
}

func TestProcessor_HappyPath(t *testing.T) {
	t.Parallel()

	t.Run("simple_code_without_import", func(t *testing.T) {
		t.Parallel()

		mainFilePath := "/home/user/config.atmc"

		os := testos.
			NewOSBuilder().
			File(func(fb *testos.FileBuilder) {
				fb.
					Path(mainFilePath).
					Content("{a: 1, b: 2}")
			}).
			Build()

		app := test.NewApp(t, test.WithOS(os))

		a, err := app.Processor().Process(mainFilePath)
		require.NoError(t, err)

		expectedAst := testlinkedast.
			NewBuilder().
			Object(func(ob *testlinkedast.ObjectBuilder) {
				ob.
					KV(func(kvb *testlinkedast.KVBuilder) {
						kvb.Key(linkedast.NewIdent("a")).
							Value(linkedast.NewInt(1))
					}).
					KV(func(kvb *testlinkedast.KVBuilder) {
						kvb.Key(linkedast.NewIdent("b")).
							Value(linkedast.NewInt(2))
					})
			}).
			Build()

		require.Equal(t, expectedAst, a)
	})

	t.Run("simple_code_with_import_but_not_used_imported_data", func(t *testing.T) {
		t.Parallel()

		mainFilePath := "/home/user/config.atmc"

		os := testos.NewOSBuilder().
			File(func(fb *testos.FileBuilder) {
				fb.
					Path(mainFilePath).
					Content("var1 ./import.atmc {a: 1, b: 2}")
			}).
			File(func(fb *testos.FileBuilder) {
				fb.
					Path("/home/user/import.atmc").
					Content("{c: 3, d: 4}")
			}).
			Build()

		app := test.NewApp(t, test.WithOS(os))

		a, err := app.Processor().Process(mainFilePath)
		require.NoError(t, err)

		expectedAst := testlinkedast.
			NewBuilder().
			Object(func(ob *testlinkedast.ObjectBuilder) {
				ob.
					KV(func(kvb *testlinkedast.KVBuilder) {
						kvb.Key(linkedast.NewIdent("a")).
							Value(linkedast.NewInt(1))
					}).
					KV(func(kvb *testlinkedast.KVBuilder) {
						kvb.Key(linkedast.NewIdent("b")).
							Value(linkedast.NewInt(2))
					})
			}).
			Build()

		require.Equal(t, expectedAst, a)
	})

	t.Run("simple_code_with_import_and_used_imported_data", func(t *testing.T) {
		t.Parallel()

		mainFilePath := "/home/user/config.atmc"

		os := testos.NewOSBuilder().
			File(func(fb *testos.FileBuilder) {
				fb.
					Path(mainFilePath).
					Content("var1 ./import.atmc {a: var1.c, b: var1.d}")
			}).
			File(func(fb *testos.FileBuilder) {
				fb.
					Path("/home/user/import.atmc").
					Content("{c: 3, d: 4}")
			}).
			Build()

		app := test.NewApp(t, test.WithOS(os))

		a, err := app.Processor().Process(mainFilePath)
		require.NoError(t, err)

		expectedAst := testlinkedast.
			NewBuilder().
			Object(func(ob *testlinkedast.ObjectBuilder) {
				ob.
					KV(func(kvb *testlinkedast.KVBuilder) {
						kvb.
							Key(linkedast.NewIdent("a")).
							Value(linkedast.NewInt(3))
					}).
					KV(func(kvb *testlinkedast.KVBuilder) {
						kvb.
							Key(linkedast.NewIdent("b")).
							Value(linkedast.NewInt(4))
					})
			}).
			Build()

		require.Equal(t, expectedAst, a)
	})

	t.Run("nested_import_variable", func(t *testing.T) {
		t.Parallel()

		mainFilePath := "/home/user/config.atmc"

		os := testos.NewOSBuilder().
			File(func(fb *testos.FileBuilder) {
				fb.
					Path(mainFilePath).
					Content("var1 ./import.atmc {a: var1.c.nested1.nested2, b: var1.d}")
			}).
			File(func(fb *testos.FileBuilder) {
				fb.
					Path("/home/user/import.atmc").
					Content("{c: {nested1: {nested2: 55}}, d: 4}")
			}).
			Build()

		app := test.NewApp(t, test.WithOS(os))

		a, err := app.Processor().Process(mainFilePath)
		require.NoError(t, err)

		expectedAst := testlinkedast.
			NewBuilder().
			Object(func(ob *testlinkedast.ObjectBuilder) {
				ob.
					KV(func(kvb *testlinkedast.KVBuilder) {
						kvb.
							Key(linkedast.NewIdent("a")).
							Value(linkedast.NewInt(55))
					}).
					KV(func(kvb *testlinkedast.KVBuilder) {
						kvb.
							Key(linkedast.NewIdent("b")).
							Value(linkedast.NewInt(4))
					})
			}).
			Build()

		require.Equal(t, expectedAst, a)
	})
}
