package compiler

import (
	"github.com/atmxlab/atmc/linker/ast"
	"github.com/atmxlab/atmc/pkg/errors"
)

type MapCompiler struct {
}

func NewMapCompiler() *MapCompiler {
	return &MapCompiler{}
}

func (c *MapCompiler) Compile(t map[string]any, a ast.Ast) error {
	compiled, err := c.compileObj(a.Object())
	if err != nil {
		return err
	}
	for k, v := range compiled {
		t[k] = v
	}
	return nil
}

func (c *MapCompiler) compileObj(obj ast.Object) (map[string]any, error) {
	m := make(map[string]any)
	for _, kv := range obj.KV() {
		value, err := c.compileExpr(kv.Value())
		if err != nil {
			return nil, err
		}

		m[kv.Key().String()] = value
	}

	return m, nil
}

func (c *MapCompiler) compileArr(arr ast.Array) ([]any, error) {
	a := make([]any, 0)
	for _, elem := range arr.Elements() {
		expr, err := c.compileExpr(elem)
		if err != nil {
			return nil, err
		}
		a = append(a, expr)
	}

	return a, nil
}

func (c *MapCompiler) compileExpr(exp ast.Expression) (any, error) {
	switch v := exp.(type) {
	case ast.Object:
		return c.compileObj(v)
	case ast.Array:
		return c.compileArr(v)
	case ast.String:
		return v.Value(), nil
	case ast.Bool:
		return v.Value(), nil
	case ast.Int:
		return v.Value(), nil
	case ast.Float:
		return v.Value(), nil
	default:
		return nil, errors.New("unexpected expression type")
	}
}
