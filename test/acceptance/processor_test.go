package acceptance

import (
	"testing"

	"github.com/atmxlab/atmc/analyzer"
	"github.com/atmxlab/atmc/linker"
	linkedast "github.com/atmxlab/atmc/linker/ast"
	"github.com/atmxlab/atmc/parser"
	"github.com/atmxlab/atmc/pkg/errors"
	"github.com/atmxlab/atmc/test"
	"github.com/atmxlab/atmc/test/testlinkedast"
	"github.com/atmxlab/atmc/test/testos"
	"github.com/stretchr/testify/require"
)

func TestProcessor_WithErrors(t *testing.T) {
	t.Parallel()

	t.Run("unused_variable", func(t *testing.T) {
		t.Parallel()

		mainFilePath := "/home/user/config.atmc"

		os := testos.NewOSBuilder().
			File(func(fb *testos.FileBuilder) {
				fb.
					Path(mainFilePath).
					Content("var1 ./import.atmc {a: 1, b: 2}")
			}).
			Build()

		app := test.NewApp(t, test.WithOS(os))

		_, err := app.Processor().Process(mainFilePath)

		require.ErrorIs(t, err, analyzer.ErrUnusedVariable)
	})

	t.Run("import_not_found", func(t *testing.T) {
		t.Parallel()

		mainFilePath := "/home/user/config.atmc"

		os := testos.NewOSBuilder().
			File(func(fb *testos.FileBuilder) {
				fb.
					Path(mainFilePath).
					Content("var1 ./import.atmc {a: 1, b: var1}")
			}).
			Build()

		app := test.NewApp(t, test.WithOS(os))

		_, err := app.Processor().Process(mainFilePath)

		require.ErrorIs(t, err, errors.ErrNotFound)
		require.ErrorContains(t, err, "file not found")
	})

	t.Run("empty_content_in_import", func(t *testing.T) {
		t.Parallel()

		mainFilePath := "/home/user/config.atmc"

		os := testos.NewOSBuilder().
			File(func(fb *testos.FileBuilder) {
				fb.
					Path(mainFilePath).
					Content("var1 ./import.atmc {a: 1, b: var1.c}")
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

	t.Run("with_array_spread_to_object", func(t *testing.T) {
		t.Parallel()

		mainFilePath := "/home/user/config.atmc"

		os := testos.NewOSBuilder().
			File(func(fb *testos.FileBuilder) {
				fb.
					Path(mainFilePath).
					Content(`
var1 ./import1.atmc

{
	a: var1.c, 
	b: var1.d,
	var1.c...
}
`)
			}).
			File(func(fb *testos.FileBuilder) {
				fb.
					Path("/home/user/import1.atmc").
					Content("{c: [3, 4, 5], d: 6}")
			}).
			Env(func(eb *testos.EnvBuilder) {
				eb.
					Key("password").
					Value("qwerty")
			}).
			Build()

		app := test.NewApp(t, test.WithOS(os))

		_, err := app.Processor().Process(mainFilePath)
		require.ErrorIs(t, err, linker.ErrUnexpectedNodeType)
		require.ErrorContains(t, err, "expected: Object")
	})

	t.Run("with_undefined_variable", func(t *testing.T) {
		t.Parallel()

		mainFilePath := "/home/user/config.atmc"

		os := testos.NewOSBuilder().
			File(func(fb *testos.FileBuilder) {
				fb.
					Path(mainFilePath).
					Content(`
var1 ./import1.atmc

{
	a: var1.c, 
	b: var2.b,
}
`)
			}).
			File(func(fb *testos.FileBuilder) {
				fb.
					Path("/home/user/import1.atmc").
					Content("{c: [3, 4, 5], d: 6}")
			}).
			Env(func(eb *testos.EnvBuilder) {
				eb.
					Key("password").
					Value("qwerty")
			}).
			Build()

		app := test.NewApp(t, test.WithOS(os))

		_, err := app.Processor().Process(mainFilePath)
		require.ErrorIs(t, err, analyzer.ErrUndefinedVariable)
		require.ErrorContains(t, err, "undefined variable: var2")
	})

	t.Run("with_undefined_nested_variable", func(t *testing.T) {
		t.Parallel()

		mainFilePath := "/home/user/config.atmc"

		os := testos.NewOSBuilder().
			File(func(fb *testos.FileBuilder) {
				fb.
					Path(mainFilePath).
					Content(`
var1 ./import1.atmc

{
	a: var1.c, 
	b: var1.j,
}
`)
			}).
			File(func(fb *testos.FileBuilder) {
				fb.
					Path("/home/user/import1.atmc").
					Content("{c: [3, 4, 5], d: 6}")
			}).
			Env(func(eb *testos.EnvBuilder) {
				eb.
					Key("password").
					Value("qwerty")
			}).
			Build()

		app := test.NewApp(t, test.WithOS(os))

		_, err := app.Processor().Process(mainFilePath)
		require.ErrorIs(t, err, linker.ErrNotFoundVariable)
		require.ErrorContains(t, err, "expected: var1.j")
	})
}

func TestProcessor_HappyPath(t *testing.T) {
	t.Parallel()

	t.Run("simple_code_without_import", func(t *testing.T) {
		t.Parallel()

		mainFilePath := "/home/user/config.atmc"

		os := testos.NewOSBuilder().
			File(func(fb *testos.FileBuilder) {
				fb.
					Path(mainFilePath).
					Content("{a: 1, b: 2}")
			}).
			Build()

		app := test.NewApp(t, test.WithOS(os))

		a, err := app.Processor().Process(mainFilePath)
		require.NoError(t, err)

		expectedAst := testlinkedast.NewBuilder().
			Object(func(ob *testlinkedast.ObjectBuilder) {
				ob.
					KV2("a", linkedast.NewInt(1)).
					KV2("b", linkedast.NewInt(2))
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

		expectedAst := testlinkedast.NewBuilder().
			Object(func(ob *testlinkedast.ObjectBuilder) {
				ob.
					KV2("a", linkedast.NewInt(3)).
					KV2("b", linkedast.NewInt(4))
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

		expectedAst := testlinkedast.NewBuilder().
			Object(func(ob *testlinkedast.ObjectBuilder) {
				ob.
					KV2("a", linkedast.NewInt(55)).
					KV2("b", linkedast.NewInt(4))
			}).
			Build()

		require.Equal(t, expectedAst, a)
	})

	t.Run("array_variable", func(t *testing.T) {
		t.Parallel()

		mainFilePath := "/home/user/config.atmc"

		os := testos.NewOSBuilder().
			File(func(fb *testos.FileBuilder) {
				fb.
					Path(mainFilePath).
					Content(`
var1 ./import.atmc

{
	a: var1.c.nested1.nested2,
	b: var1.d
}
`)
			}).
			File(func(fb *testos.FileBuilder) {
				fb.
					Path("/home/user/import.atmc").
					Content(`
{
	c: {
		nested1: {
					nested2: [1, 2, 3, 4, 5]
				}
	}, 
	d: 4
}
`)
			}).
			Build()

		app := test.NewApp(t, test.WithOS(os))

		a, err := app.Processor().Process(mainFilePath)
		require.NoError(t, err)

		expectedAst := testlinkedast.NewBuilder().
			Object(func(ob *testlinkedast.ObjectBuilder) {
				ob.
					KV2(
						"a",
						testlinkedast.NewArrayBuilder().
							Element(linkedast.NewInt(1)).
							Element(linkedast.NewInt(2)).
							Element(linkedast.NewInt(3)).
							Element(linkedast.NewInt(4)).
							Element(linkedast.NewInt(5)).
							Build(),
					).
					KV2("b", linkedast.NewInt(4))
			}).
			Build()

		require.Equal(t, expectedAst, a)
	})

	t.Run("several_imports", func(t *testing.T) {
		t.Parallel()

		mainFilePath := "/home/user/config.atmc"

		os := testos.NewOSBuilder().
			File(func(fb *testos.FileBuilder) {
				fb.
					Path(mainFilePath).
					Content(`
var1 ./import1.atmc
var2 ./import2.atmc

{
	a: var1.c, 
	b: var2.i
}
`)
			}).
			File(func(fb *testos.FileBuilder) {
				fb.
					Path("/home/user/import1.atmc").
					Content("{c: 3, d: 4}")
			}).
			File(func(fb *testos.FileBuilder) {
				fb.
					Path("/home/user/import2.atmc").
					Content("{i: 5, f: 6}")
			}).
			Build()

		app := test.NewApp(t, test.WithOS(os))

		a, err := app.Processor().Process(mainFilePath)
		require.NoError(t, err)

		expectedAst := testlinkedast.NewBuilder().
			Object(func(ob *testlinkedast.ObjectBuilder) {
				ob.
					KV2("a", linkedast.NewInt(3)).
					KV2("b", linkedast.NewInt(5))
			}).
			Build()

		require.Equal(t, expectedAst, a)
	})

	t.Run("with_env_variables", func(t *testing.T) {
		t.Parallel()

		mainFilePath := "/home/user/config.atmc"

		os := testos.NewOSBuilder().
			File(func(fb *testos.FileBuilder) {
				fb.
					Path(mainFilePath).
					Content(`
var1 ./import1.atmc

{
	a: var1.c, 
	b: var1.d,
	object1: {
		password: $password
	}
}
`)
			}).
			File(func(fb *testos.FileBuilder) {
				fb.
					Path("/home/user/import1.atmc").
					Content("{c: 3, d: 4}")
			}).
			Env(func(eb *testos.EnvBuilder) {
				eb.
					Key("password").
					Value("qwerty")
			}).
			Build()

		app := test.NewApp(t, test.WithOS(os))

		a, err := app.Processor().Process(mainFilePath)
		require.NoError(t, err)

		expectedAst := testlinkedast.NewBuilder().
			Object(func(ob *testlinkedast.ObjectBuilder) {
				ob.
					KV2("a", linkedast.NewInt(3)).
					KV2("b", linkedast.NewInt(4)).
					KV2(
						"object1",
						testlinkedast.
							NewObjectBuilder().
							KV2("password", linkedast.NewString("qwerty")).
							Build(),
					)
			}).
			Build()

		require.Equal(t, expectedAst, a)
	})

	t.Run("with_object_spread", func(t *testing.T) {
		t.Parallel()

		mainFilePath := "/home/user/config.atmc"

		os := testos.NewOSBuilder().
			File(func(fb *testos.FileBuilder) {
				fb.
					Path(mainFilePath).
					Content(`
var1 ./import1.atmc

{
	a: var1.c, 
	b: var1.d,
	var1...
}
`)
			}).
			File(func(fb *testos.FileBuilder) {
				fb.
					Path("/home/user/import1.atmc").
					Content("{c: 3, d: 4}")
			}).
			Build()

		app := test.NewApp(t, test.WithOS(os))

		a, err := app.Processor().Process(mainFilePath)
		require.NoError(t, err)

		expectedAst := testlinkedast.NewBuilder().
			Object(func(ob *testlinkedast.ObjectBuilder) {
				ob.
					KV2("a", linkedast.NewInt(3)).
					KV2("b", linkedast.NewInt(4)).
					KV2("c", linkedast.NewInt(3)).
					KV2("d", linkedast.NewInt(4))
			}).
			Build()

		require.Equal(t, expectedAst, a)
	})

	t.Run("with_array_spread", func(t *testing.T) {
		t.Parallel()

		mainFilePath := "/home/user/config.atmc"

		os := testos.NewOSBuilder().
			File(func(fb *testos.FileBuilder) {
				fb.
					Path(mainFilePath).
					Content(`
var1 ./import1.atmc

{
	a: var1.c, 
	b: [1, 2, var1.d..., 6, 7],
}
`)
			}).
			File(func(fb *testos.FileBuilder) {
				fb.
					Path("/home/user/import1.atmc").
					Content("{c: 3, d: [3, 4, 5]}")
			}).
			Build()

		app := test.NewApp(t, test.WithOS(os))

		a, err := app.Processor().Process(mainFilePath)
		require.NoError(t, err)

		expectedAst := testlinkedast.NewBuilder().
			Object(func(ob *testlinkedast.ObjectBuilder) {
				ob.
					KV2("a", linkedast.NewInt(3)).
					KV2("b", testlinkedast.NewArrayBuilder().
						Element(linkedast.NewInt(1)).
						Element(linkedast.NewInt(2)).
						Element(linkedast.NewInt(3)).
						Element(linkedast.NewInt(4)).
						Element(linkedast.NewInt(5)).
						Element(linkedast.NewInt(6)).
						Element(linkedast.NewInt(7)).
						Build())
			}).
			Build()

		require.Equal(t, expectedAst, a)
	})

	t.Run("override_with_spread", func(t *testing.T) {
		t.Parallel()

		mainFilePath := "/home/user/config.atmc"

		os := testos.NewOSBuilder().
			File(func(fb *testos.FileBuilder) {
				fb.
					Path(mainFilePath).
					Content(`
var1 ./import1.atmc

{
	a: 1, 
	b: 2,
	var1...
}
`)
			}).
			File(func(fb *testos.FileBuilder) {
				fb.
					Path("/home/user/import1.atmc").
					Content("{a: 3, b: 4, c: 5}")
			}).
			Build()

		app := test.NewApp(t, test.WithOS(os))

		a, err := app.Processor().Process(mainFilePath)
		require.NoError(t, err)

		expectedAst := testlinkedast.NewBuilder().
			Object(func(ob *testlinkedast.ObjectBuilder) {
				ob.
					KV2("a", linkedast.NewInt(3)).
					KV2("b", linkedast.NewInt(4)).
					KV2("c", linkedast.NewInt(5))
			}).
			Build()

		require.Equal(t, expectedAst, a)
	})

	t.Run("override_after_spread", func(t *testing.T) {
		t.Parallel()

		mainFilePath := "/home/user/config.atmc"

		os := testos.NewOSBuilder().
			File(func(fb *testos.FileBuilder) {
				fb.
					Path(mainFilePath).
					Content(`
var1 ./import1.atmc

{
	var1...,
	a: 1, 
	b: 2,
}
`)
			}).
			File(func(fb *testos.FileBuilder) {
				fb.
					Path("/home/user/import1.atmc").
					Content("{a: 3, b: 4, c: 5}")
			}).
			Build()

		app := test.NewApp(t, test.WithOS(os))

		a, err := app.Processor().Process(mainFilePath)
		require.NoError(t, err)

		expectedAst := testlinkedast.NewBuilder().
			Object(func(ob *testlinkedast.ObjectBuilder) {
				ob.
					KV2("a", linkedast.NewInt(1)).
					KV2("b", linkedast.NewInt(2)).
					KV2("c", linkedast.NewInt(5))
			}).
			Build()

		require.Equal(t, expectedAst, a)
	})

	t.Run("override_nested_entry_only", func(t *testing.T) {
		t.Parallel()

		mainFilePath := "/home/user/config.atmc"

		os := testos.NewOSBuilder().
			File(func(fb *testos.FileBuilder) {
				fb.
					Path(mainFilePath).
					Content(`
common ./common.atmc

{
	common...,
	logging: {
		enabled: false
	}
}
`)
			}).
			File(func(fb *testos.FileBuilder) {
				fb.
					Path("/home/user/common.atmc").
					Content(`{
logging: {
	enabled: true
	level: ["warn", "error"]
}
}`)
			}).
			Build()

		app := test.NewApp(t, test.WithOS(os))

		a, err := app.Processor().Process(mainFilePath)
		require.NoError(t, err)

		expectedAst := testlinkedast.NewBuilder().
			Object(func(ob *testlinkedast.ObjectBuilder) {
				ob.
					KV2(
						"logging",
						testlinkedast.NewObjectBuilder().
							KV2("enabled", linkedast.NewBool(false)).
							KV2(
								"level",
								testlinkedast.NewArrayBuilder().
									Element(linkedast.NewString("warn")).
									Element(linkedast.NewString("error")).
									Build(),
							).
							Build(),
					)
			}).
			Build()

		require.Equal(t, expectedAst, a)
	})

	t.Run("override_deep_nested_entry_only", func(t *testing.T) {
		t.Parallel()

		mainFilePath := "/home/user/config.atmc"

		os := testos.NewOSBuilder().
			File(func(fb *testos.FileBuilder) {
				fb.
					Path(mainFilePath).
					Content(`
common ./common.atmc

{
	common...,
	logging: {
		settings: {
			enabled: false
		}
	}
}
`)
			}).
			File(func(fb *testos.FileBuilder) {
				fb.
					Path("/home/user/common.atmc").
					Content(`{
logging: {
	settings: {
		enabled: true
		enableTracing: true
	}
	level: ["warn", "error"]
}
}`)
			}).
			Build()

		app := test.NewApp(t, test.WithOS(os))

		a, err := app.Processor().Process(mainFilePath)
		require.NoError(t, err)

		expectedAst := testlinkedast.NewBuilder().
			Object(func(ob *testlinkedast.ObjectBuilder) {
				ob.
					KV2(
						"logging",
						testlinkedast.NewObjectBuilder().
							KV2(
								"settings",
								testlinkedast.NewObjectBuilder().
									KV2("enabled", linkedast.NewBool(false)).
									KV2("enableTracing", linkedast.NewBool(true)).
									Build(),
							).
							KV2(
								"level",
								testlinkedast.NewArrayBuilder().
									Element(linkedast.NewString("warn")).
									Element(linkedast.NewString("error")).
									Build(),
							).
							Build(),
					)
			}).
			Build()

		require.Equal(t, expectedAst, a)
	})
}
