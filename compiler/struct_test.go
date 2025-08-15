package compiler_test

import (
	"math"
	"testing"

	"github.com/atmxlab/atmc/compiler"
	linkedast "github.com/atmxlab/atmc/linker/ast"
	"github.com/atmxlab/atmc/test/testlinkedast"
	"github.com/samber/lo"
	"github.com/stretchr/testify/require"
)

type TestType struct {
	FieldArray  []string `atmc:"field_array"`
	FieldStruct struct {
		FieldInt int `atmc:"field_int"`
	} `atmc:"field_struct"`
	FieldNestedTestType    *TestType   `atmc:"field_nested_test_type"`
	FieldPtrToPtr          **TestType  `atmc:"field_ptr_to_ptr"`
	FieldTestTypeArray     []TestType  `atmc:"field_test_type_array"`
	FieldTestTypePtrArray  []*TestType `atmc:"field_test_type_ptr_array"`
	FieldStr               string      `atmc:"field_str"`
	FieldInt               int         `atmc:"field_int"`
	FieldInt8              int8        `atmc:"field_int8"`
	FieldInt16             int16       `atmc:"field_int16"`
	FieldInt32             int32       `atmc:"field_int32"`
	FieldInt64             int64       `atmc:"field_int64"`
	FieldUint              uint        `atmc:"field_uint"`
	FieldUint8             uint8       `atmc:"field_uint8"`
	FieldUint16            uint16      `atmc:"field_uint16"`
	FieldUint32            uint32      `atmc:"field_uint32"`
	FieldUint64            uint64      `atmc:"field_uint64"`
	FieldFloat32           float32     `atmc:"field_float32"`
	FieldFloat64           float64     `atmc:"field_float64"`
	FieldBool              bool        `atmc:"field_bool"`
	FieldBoolPtr           *bool       `atmc:"field_bool_ptr"`
	FieldIntPtr            *int        `atmc:"field_int_ptr"`
	FieldUintPtr           *uint       `atmc:"field_uint_ptr"`
	FieldStringPtr         *string     `atmc:"field_string_ptr"`
	FieldBoolPtrToPtrToPtr ***bool     `atmc:"field_bool_ptr_to_ptr_to_ptr"`
	FieldWithoutTag        string
}

func TestStruct(t *testing.T) {
	t.Parallel()

	t.Run("with_tag_and_without_tag", func(t *testing.T) {
		t.Parallel()

		a := testlinkedast.
			NewBuilder().
			Object(func(ob *testlinkedast.ObjectBuilder) {
				ob.KV2("field_str", linkedast.NewString("field_str_value"))
				ob.KV2("FieldWithoutTag", linkedast.NewString("field without tag"))
			}).
			Build()

		c := compiler.NewStructCompiler("atmc")

		var v TestType
		err := c.Compile(&v, a)
		require.NoError(t, err)

		expected := TestType{
			FieldStr:        "field_str_value",
			FieldWithoutTag: "field without tag",
		}

		require.Equal(t, expected, v)
	})

	t.Run("with_not_exists_field", func(t *testing.T) {
		t.Parallel()

		a := testlinkedast.
			NewBuilder().
			Object(func(ob *testlinkedast.ObjectBuilder) {
				ob.KV2("field_str", linkedast.NewString("field_str_value"))
				ob.KV2("not_exists_field", linkedast.NewString("no_exists_field"))
			}).
			Build()

		c := compiler.NewStructCompiler("atmc")

		var v TestType
		err := c.Compile(&v, a)
		require.NoError(t, err)

		expected := TestType{
			FieldStr: "field_str_value",
		}

		require.Equal(t, expected, v)
	})

	t.Run("with_ptr", func(t *testing.T) {
		t.Parallel()

		a := testlinkedast.
			NewBuilder().
			Object(func(ob *testlinkedast.ObjectBuilder) {
				ob.KV2("field_bool_ptr", linkedast.NewBool(true))
				ob.KV2("field_int_ptr", linkedast.NewInt(-100))
				ob.KV2("field_uint_ptr", linkedast.NewInt(1000))
				ob.KV2("field_string_ptr", linkedast.NewString("field_str_ptr_value"))
			}).
			Build()

		c := compiler.NewStructCompiler("atmc")

		var v TestType
		err := c.Compile(&v, a)
		require.NoError(t, err)

		expected := TestType{
			FieldBoolPtr:   lo.ToPtr(true),
			FieldIntPtr:    lo.ToPtr(-100),
			FieldUintPtr:   lo.ToPtr[uint](1000),
			FieldStringPtr: lo.ToPtr("field_str_ptr_value"),
		}

		require.Equal(t, expected, v)
	})

	t.Run("with_specific_int_and_float", func(t *testing.T) {
		t.Parallel()

		a := testlinkedast.
			NewBuilder().
			Object(func(ob *testlinkedast.ObjectBuilder) {
				ob.KV2("field_int8", linkedast.NewInt(8))
				ob.KV2("field_int16", linkedast.NewInt(16))
				ob.KV2("field_int32", linkedast.NewInt(32))
				ob.KV2("field_int64", linkedast.NewInt(64))
				ob.KV2("field_uint8", linkedast.NewInt(8))
				ob.KV2("field_uint16", linkedast.NewInt(16))
				ob.KV2("field_uint32", linkedast.NewInt(32))
				ob.KV2("field_uint64", linkedast.NewInt(64))
				ob.KV2("field_float32", linkedast.NewFloat(32.32))
				ob.KV2("field_float64", linkedast.NewFloat(64.64))
			}).
			Build()

		c := compiler.NewStructCompiler("atmc")

		var v TestType
		err := c.Compile(&v, a)
		require.NoError(t, err)

		expected := TestType{
			FieldInt8:    8,
			FieldInt16:   16,
			FieldInt32:   32,
			FieldInt64:   64,
			FieldUint8:   8,
			FieldUint16:  16,
			FieldUint32:  32,
			FieldUint64:  64,
			FieldFloat32: 32.32,
			FieldFloat64: 64.64,
		}

		require.Equal(t, expected, v)
	})

	t.Run("with_overflow_int", func(t *testing.T) {
		t.Parallel()

		a := testlinkedast.
			NewBuilder().
			Object(func(ob *testlinkedast.ObjectBuilder) {
				ob.KV2("field_int8", linkedast.NewInt(500))
			}).
			Build()

		c := compiler.NewStructCompiler("atmc")

		var v TestType
		err := c.Compile(&v, a)
		require.ErrorIs(t, err, compiler.ErrTypeOverflow)
	})

	t.Run("with_float_overflow", func(t *testing.T) {
		t.Parallel()

		a := testlinkedast.
			NewBuilder().
			Object(func(ob *testlinkedast.ObjectBuilder) {
				ob.KV2("field_float32", linkedast.NewFloat(math.MaxFloat64))
			}).
			Build()

		c := compiler.NewStructCompiler("atmc")

		var v TestType
		err := c.Compile(&v, a)
		require.ErrorIs(t, err, compiler.ErrTypeOverflow)
	})

	t.Run("with_structs_in_array", func(t *testing.T) {
		t.Parallel()

		a := testlinkedast.
			NewBuilder().
			Object(func(ob *testlinkedast.ObjectBuilder) {
				ob.KV2("field_test_type_array",
					testlinkedast.
						NewArrayBuilder().
						Element(
							testlinkedast.
								NewObjectBuilder().
								KV2("field_int", linkedast.NewInt(10)).
								Build(),
						).
						Element(
							testlinkedast.
								NewObjectBuilder().
								KV2("field_int", linkedast.NewInt(20)).
								Build(),
						).
						Build(),
				)
			}).
			Build()

		c := compiler.NewStructCompiler("atmc")

		var v TestType
		err := c.Compile(&v, a)
		require.NoError(t, err)

		expected := TestType{
			FieldTestTypeArray: []TestType{
				{
					FieldInt: 10,
				},
				{
					FieldInt: 20,
				},
			},
		}

		require.Equal(t, expected, v)
	})

	t.Run("with_ptr_structs_in_array", func(t *testing.T) {
		t.Parallel()

		a := testlinkedast.
			NewBuilder().
			Object(func(ob *testlinkedast.ObjectBuilder) {
				ob.KV2("field_test_type_ptr_array",
					testlinkedast.
						NewArrayBuilder().
						Element(
							testlinkedast.
								NewObjectBuilder().
								KV2("field_int", linkedast.NewInt(30)).
								Build(),
						).
						Element(
							testlinkedast.
								NewObjectBuilder().
								KV2("field_int", linkedast.NewInt(40)).
								Build(),
						).
						Build(),
				)
			}).
			Build()

		c := compiler.NewStructCompiler("atmc")

		var v TestType
		err := c.Compile(&v, a)
		require.NoError(t, err)

		expected := TestType{
			FieldTestTypePtrArray: []*TestType{
				{
					FieldInt: 30,
				},
				{
					FieldInt: 40,
				},
			},
		}

		require.Equal(t, expected, v)
	})

	t.Run("with_ptr_structs_in_object", func(t *testing.T) {
		t.Parallel()

		a := testlinkedast.
			NewBuilder().
			Object(func(ob *testlinkedast.ObjectBuilder) {
				ob.KV2("field_nested_test_type", testlinkedast.
					NewObjectBuilder().
					KV2("field_int", linkedast.NewInt(10)).
					KV2("field_nested_test_type", testlinkedast.
						NewObjectBuilder().
						KV2("field_float32", linkedast.NewFloat(32.32)).
						Build(),
					).
					Build(),
				)
			}).
			Build()

		c := compiler.NewStructCompiler("atmc")

		var v TestType
		err := c.Compile(&v, a)
		require.NoError(t, err)

		expected := TestType{
			FieldNestedTestType: &TestType{
				FieldInt: 10,
				FieldNestedTestType: &TestType{
					FieldFloat32: 32.32,
				},
			},
		}

		require.Equal(t, expected, v)
	})

	t.Run("with_ptr_to_ptr_object", func(t *testing.T) {
		t.Parallel()

		a := testlinkedast.
			NewBuilder().
			Object(func(ob *testlinkedast.ObjectBuilder) {
				ob.KV2("field_ptr_to_ptr", testlinkedast.
					NewObjectBuilder().
					KV2("field_float32", linkedast.NewFloat(32.32)).
					Build(),
				)
			}).
			Build()

		c := compiler.NewStructCompiler("atmc")

		var v TestType
		err := c.Compile(&v, a)
		require.NoError(t, err)

		expected := TestType{
			FieldPtrToPtr: lo.ToPtr(lo.ToPtr(TestType{
				FieldFloat32: 32.32,
			})),
		}

		require.Equal(t, expected, v)
	})

	t.Run("with_ptr_to_ptr_to_ptr_bool", func(t *testing.T) {
		t.Parallel()

		a := testlinkedast.
			NewBuilder().
			Object(func(ob *testlinkedast.ObjectBuilder) {
				ob.KV2("field_bool_ptr_to_ptr_to_ptr", linkedast.NewBool(true))
			}).
			Build()

		c := compiler.NewStructCompiler("atmc")

		var v TestType
		err := c.Compile(&v, a)
		require.NoError(t, err)

		expected := TestType{
			FieldBoolPtrToPtrToPtr: lo.ToPtr(lo.ToPtr(lo.ToPtr(true))),
		}

		require.Equal(t, expected, v)
	})

	t.Run("with_embedded_struct", func(t *testing.T) {
		t.Parallel()

		a := testlinkedast.
			NewBuilder().
			Object(func(ob *testlinkedast.ObjectBuilder) {
				ob.KV2("field_str", linkedast.NewString("field_str_value"))
				ob.KV2("field_int", linkedast.NewInt(42))
				ob.KV2(
					"field_array",
					testlinkedast.
						NewArrayBuilder().
						Element(linkedast.NewString("field_array_value_1")).
						Element(linkedast.NewString("field_array_value_2")).
						Element(linkedast.NewString("field_array_value_3")).
						Build(),
				)
				ob.KV2(
					"field_struct",
					testlinkedast.
						NewObjectBuilder().
						KV2("field_int", linkedast.NewInt(44)).
						Build(),
				)
			}).
			Build()

		c := compiler.NewStructCompiler("atmc")

		var v TestType
		err := c.Compile(&v, a)
		require.NoError(t, err)

		expected := TestType{
			FieldArray: []string{"field_array_value_1", "field_array_value_2", "field_array_value_3"},
			FieldStruct: struct {
				FieldInt int `atmc:"field_int"`
			}{
				FieldInt: 44,
			},
			FieldStr: "field_str_value",
			FieldInt: 42,
		}

		require.Equal(t, expected, v)
	})
}
