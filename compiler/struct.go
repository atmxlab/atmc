package compiler

import (
	"math"
	"reflect"

	"github.com/atmxlab/atmc/linker/ast"
	"github.com/atmxlab/atmc/pkg/errors"
)

type StructCompiler struct {
	tagName string
}

func NewStructCompiler(tagName string) *StructCompiler {
	return &StructCompiler{tagName: tagName}
}

func (c *StructCompiler) Compile(t any, ast ast.Ast) error {
	if err := c.processObject(ast.Object(), reflect.ValueOf(t)); err != nil {
		return errors.Wrap(err, "processObject")
	}

	return nil
}

func (c *StructCompiler) processObject(obj ast.Object, v reflect.Value) error {
	if v.Kind() != reflect.Ptr {
		return errors.Newf("invalid elem kind: expected: [%s], actual: [%s]", reflect.Ptr, v.Kind())
	}

	if v.Elem().Kind() != reflect.Struct {
		return errors.Newf("invalid elem kind: expected: [%s], actual: [%s]", reflect.Struct, v.Elem().Kind())
	}

	for _, kv := range obj.KV() {
		field, err := c.getObjectField(v, kv.Key().String())
		if err != nil && !errors.Is(err, errors.ErrNotFound) {
			return errors.Wrap(err, "c.getObjectField")
		}
		if errors.Is(err, errors.ErrNotFound) {
			continue
		}

		switch astV := kv.Value().(type) {
		case ast.Array:
			if err = c.processArray(astV, field); err != nil {
				return errors.Wrap(err, "c.processArray")
			}
		case ast.Object:
			object := c.makeValueRecursive(field)

			if err = c.processObject(astV, object.Addr()); err != nil {
				return errors.Wrap(err, "c.processObject")
			}
		case ast.Int, ast.Float, ast.String, ast.Bool:
			literal := c.makeValueRecursive(field)

			if err = c.processLiteral(astV, literal); err != nil {
				return errors.Wrap(err, "c.processLiteral")
			}
		}
	}

	return nil
}

func (c *StructCompiler) processArray(arr ast.Array, field reflect.Value) error {
	result := reflect.MakeSlice(field.Type().Elem(), len(arr.Elements()), len(arr.Elements()))

	for idx, exp := range arr.Elements() {
		elem := result.Index(idx)

		switch astV := exp.(type) {
		case ast.Object:
			elem = c.makeValueRecursive(elem.Addr())

			if err := c.processObject(astV, elem.Addr()); err != nil {
				return errors.Wrap(err, "processObject")
			}
		case ast.Array:
			if err := c.processArray(astV, elem); err != nil {
				return errors.Wrap(err, "processArray")
			}
		case ast.Int, ast.Float, ast.String, ast.Bool:
			if err := c.processLiteral(astV, elem); err != nil {
				return errors.Wrap(err, "c.processLiteral")
			}
		}
	}

	field.Elem().Set(result)

	return nil
}

func (c *StructCompiler) processLiteral(
	exp ast.Expression,
	field reflect.Value,
) error {
	switch astV := exp.(type) {
	case ast.Int:
		switch field.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			if err := c.checkInt(astV, field.Kind()); err != nil {
				return errors.Wrap(err, "c.checkInt")
			}

			field.SetInt(astV.Value())
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
			if err := c.checkUInt(astV, field.Kind()); err != nil {
				return errors.Wrap(err, "c.checkUInt")
			}

			field.SetUint(uint64(astV.Value()))
		default:
			return errors.Wrapf(ErrInvalidType, "expected int, got [%s]", field.Kind())
		}
	case ast.Float:
		if err := c.checkFloat(astV, field.Kind()); err != nil {
			return errors.Wrap(err, "c.checkFloat")
		}

		field.SetFloat(astV.Value())
	case ast.String:
		field.SetString(astV.Value())
	case ast.Bool:
		field.SetBool(astV.Value())
	}

	return nil
}

func (c *StructCompiler) checkInt(astV ast.Int, kind reflect.Kind) error {
	var minV, maxV int64

	switch kind {
	case reflect.Int8:
		minV = math.MinInt8
		maxV = math.MaxInt8
	case reflect.Int16:
		minV = math.MinInt16
		maxV = math.MaxInt16
	case reflect.Int32:
		minV = math.MinInt32
		maxV = math.MaxInt32
	case reflect.Int64:
		minV = math.MinInt64
		maxV = math.MaxInt64
	case reflect.Int:
		minV = math.MinInt
		maxV = math.MaxInt
	default:
		return errors.Wrapf(ErrInvalidType, "expected int, got [%s]", kind)
	}

	if astV.Value() > maxV {
		return errors.Wrapf(ErrTypeOverflow, "integer value overflow: actual value: [%d], max value: [%d], kind: [%s]", astV.Value(), maxV, kind)
	}

	if astV.Value() < minV {
		return errors.Wrapf(ErrTypeOverflow, "integer value overflow: actual value: [%d], min value: [%d], kind: [%s]", astV.Value(), minV, kind)
	}

	return nil
}

func (c *StructCompiler) checkUInt(astV ast.Int, kind reflect.Kind) error {
	var minV, maxV uint64

	switch kind {
	case reflect.Uint8:
		maxV = math.MaxUint8
	case reflect.Uint16:
		maxV = math.MaxUint16
	case reflect.Uint32:
		maxV = math.MaxUint32
	case reflect.Uint64:
		maxV = math.MaxInt64
	case reflect.Uint:
		maxV = math.MaxInt
	default:
		return errors.Wrapf(ErrInvalidType, "expected uint, got [%s]", kind)
	}

	if astV.Value() > int64(maxV) {
		return errors.Wrapf(ErrTypeOverflow, "unsigned integer value overflow: actual value: [%d], max value: [%d], kind: [%s]", astV.Value(), maxV, kind)
	}

	if astV.Value() < int64(minV) {
		return errors.Wrapf(ErrTypeOverflow, "unsigned integer value overflow: actual value: [%d], min value: [%d], kind: [%s]", astV.Value(), minV, kind)
	}

	return nil
}

func (c *StructCompiler) checkFloat(astV ast.Float, kind reflect.Kind) error {
	var minV, maxV float64
	switch kind {
	case reflect.Float32:
		minV = -math.MaxFloat32
		maxV = math.MaxFloat32
	case reflect.Float64:
		minV = -math.MaxFloat64
		maxV = math.MaxFloat64
	default:
		return errors.Wrapf(ErrInvalidType, "expected int, got [%s]", kind)
	}

	if astV.Value() > maxV {
		return errors.Wrapf(ErrTypeOverflow, "float value overflow: actual value: [%f], max value: [%f], kind: [%s]", astV.Value(), maxV, kind)
	}

	if astV.Value() < minV {
		return errors.Wrapf(ErrTypeOverflow, "float value overflow: actual value: [%f], min value: [%f], kind: [%s]", astV.Value(), minV, kind)
	}

	return nil
}

func (c *StructCompiler) makeValueRecursive(value reflect.Value) reflect.Value {
	current := value

	if current.Kind() == reflect.Ptr {
		if current.Elem().Kind() == reflect.Ptr {
			current = reflect.New(current.Elem().Type().Elem())
			value.Elem().Set(current)
			return c.makeValueRecursive(current)
		}
	}

	return current.Elem()
}

func (c *StructCompiler) getObjectField(v reflect.Value, name string) (reflect.Value, error) {
	if v.Kind() != reflect.Ptr {
		return reflect.Value{}, errors.Newf("invalid elem kind: expected: [%s], actual: [%s]", reflect.Ptr, v.Elem().Kind())
	}

	if v.Elem().Kind() != reflect.Struct {
		return reflect.Value{}, errors.Newf("invalid elem kind: expected: [%s], actual: [%s]", reflect.Struct, v.Elem().Kind())
	}

	val := v.Elem()
	typ := val.Type()

	for i := 0; i < typ.NumField(); i++ {
		field := val.Field(i)
		fieldTyp := typ.Field(i)

		fieldName := fieldTyp.Name

		tag := fieldTyp.Tag.Get(c.tagName)
		if tag != "" {
			fieldName = tag
		}

		if fieldName == name {
			return field.Addr(), nil
		}
	}

	return reflect.Value{}, errors.NotFound("field not found")
}
