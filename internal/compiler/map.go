package compiler

import (
	linkedast "github.com/atmxlab/atmcfg/internal/linker/ast"
	"github.com/atmxlab/atmcfg/pkg/errors"
)

type MapCompiler struct {
}

func NewMapCompiler() *MapCompiler {
	return &MapCompiler{}
}

func (c *MapCompiler) Compile(t map[string]any, a linkedast.Ast) error {
	compiled, err := c.compileObj(a.Object())
	if err != nil {
		return err
	}
	for k, v := range compiled {
		t[k] = v
	}
	return nil
}

func (c *MapCompiler) compileObj(obj linkedast.Object) (map[string]any, error) {
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

func (c *MapCompiler) compileArr(arr linkedast.Array) ([]any, error) {
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

func (c *MapCompiler) compileExpr(exp linkedast.Expression) (any, error) {
	switch v := exp.(type) {
	case linkedast.Object:
		return c.compileObj(v)
	case linkedast.Array:
		return c.compileArr(v)
	case linkedast.String:
		return v.Value(), nil
	case linkedast.Bool:
		return v.Value(), nil
	case linkedast.Int:
		return v.Value(), nil
	case linkedast.Float:
		return v.Value(), nil
	default:
		return nil, errors.New("unexpected expression type")
	}
}
