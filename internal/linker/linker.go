package linker

import (
	linkedast "github.com/atmxlab/atmcfg/internal/linker/ast"
	"github.com/atmxlab/atmcfg/internal/parser/ast"
	"github.com/atmxlab/atmcfg/pkg/errors"
	"github.com/samber/lo"
)

// TODO: это чисто линкер, а не компилятор
//  нужно переименовать.
//  А компилятор будет из уже полностью готового AST компилировать
//  либо в структуру, либо в map либо в json либо в yml

// TODO: нужно разделить создание AST и линкер
//  тут не должно быть всяких чтений с диска, лексеров, парсеров - это не ответственность линкера

type Object = ast.Object

// Path to config.
type Path string

func (p Path) String() string {
	return string(p)
}

// Name of import.
type Name string

type OS interface {
	AbsPath(baseDir, relPath string) (string, error)
}

type Linker struct {
	// Необходим, чтобы обрабатывать повторные импорты.
	astByPath    map[string]ast.WithPath
	linkedByPath map[string]linkedast.Ast
}

func New() *Linker {
	return &Linker{
		astByPath:    make(map[string]ast.WithPath),
		linkedByPath: make(map[string]linkedast.Ast),
	}
}

type scope struct {
	// Необходим, чтобы добираться до внутренностей переменных по названию.
	linkedByName map[string]linkedast.Ast
	ast          ast.WithPath
}

func newScope(a ast.WithPath) scope {
	return scope{
		linkedByName: make(map[string]linkedast.Ast),
		ast:          a,
	}
}

func (c *Linker) Link(mainAst ast.WithPath, astByPath map[string]ast.WithPath) (linkedast.Ast, error) {
	c.astByPath = astByPath
	return c.link(newScope(mainAst))
}

func (c *Linker) link(scp scope) (linkedast.Ast, error) {
	for _, imp := range scp.ast.Imports() {
		absPath, ok := scp.ast.ImportPath(imp.Path().String())
		if !ok {
			return linkedast.Ast{}, errors.New("get import absolute path by relative path")
		}

		if alreadyLinked, ok := c.linkedByPath[absPath]; ok {
			scp.linkedByName[imp.Name().String()] = alreadyLinked
			return alreadyLinked, nil
		}

		astForLink, ok := c.astByPath[absPath]
		if !ok {
			return linkedast.Ast{}, errors.New("ast for link not found")
		}

		linked, err := c.link(newScope(astForLink))
		if err != nil {
			return linkedast.Ast{}, errors.Wrapf(err, "link ast, path: [%s]", imp.Path().String())
		}

		c.linkedByPath[absPath] = linked
		scp.linkedByName[imp.Name().String()] = linked
	}

	obj, err := c.linkObject(scp, scp.ast.Root().Object())
	if err != nil {
		return linkedast.Ast{}, errors.Wrap(err, "link object")
	}

	return linkedast.NewAst(obj), nil
}

func (c *Linker) linkObject(scp scope, obj ast.Object) (linkedast.Object, error) {
	entries, err := c.linkEntries(scp, obj.Entries())
	if err != nil {
		return linkedast.Object{}, errors.Wrap(err, "link entries")
	}

	return linkedast.NewObject(entries), nil
}

func (c *Linker) linkEntries(scp scope, entries []ast.Entry) ([]linkedast.KV, error) {
	kv := make([]linkedast.KV, 0, len(entries))

	for _, entry := range entries {
		switch e := entry.(type) {
		case ast.KV:
			ent, err := c.linkKV(scp, e)
			if err != nil {
				return nil, errors.Wrap(err, "link kv")
			}
			kv = append(kv, ent)
		case ast.Spread:
			spreadEntries, err := c.linkObjectSpread(scp, e)
			if err != nil {
				return nil, errors.Wrap(err, "link spread")
			}

			kv = append(kv, spreadEntries...)
		default:
			return nil, errors.New("unknown entry type")
		}
	}

	return kv, nil
}

func (c *Linker) linkKV(scp scope, kv ast.KV) (linkedast.KV, error) {
	var value linkedast.Expression
	switch v := kv.Value().(type) {
	case ast.Object:
		exp, err := c.linkObject(scp, v)
		if err != nil {
			return linkedast.KV{}, errors.Wrap(err, "link object")
		}

		value = exp
	case ast.Array:
		exp, err := c.linkArray(scp, v)
		if err != nil {
			return linkedast.KV{}, errors.Wrap(err, "link array")
		}

		value = exp
	case ast.Var:
		node, err := c.findVariableExp(scp, v)
		if err != nil {
			return linkedast.KV{}, errors.Wrap(err, "find variable")
		}

		value = node
	case ast.Bool:
		value = linkedast.NewBool(v.Value())
	case ast.String:
		value = linkedast.NewString(v.Value())
	case ast.Int:
		value = linkedast.NewInt(v.Value())
	case ast.Float:
		value = linkedast.NewFloat(v.Value())
	default:
		return linkedast.KV{}, errors.New("unknown value type")
	}

	return linkedast.NewKV(linkedast.NewIdent(kv.Key().String()), value), nil
}

func (c *Linker) linkObjectSpread(scp scope, spread ast.Spread) ([]linkedast.KV, error) {
	node, err := c.findVariableExp(scp, spread.Var())
	if err != nil {
		return nil, errors.Wrap(err, "find variable node")
	}

	obj, ok := node.(linkedast.Object)
	if !ok {
		return []linkedast.KV{}, errors.New("unexpected node type")
	}

	return obj.KV(), nil
}

func (c *Linker) linkArraySpread(scp scope, spread ast.Spread) ([]linkedast.Expression, error) {
	node, err := c.findVariableExp(scp, spread.Var())
	if err != nil {
		return nil, errors.Wrap(err, "find variable node")
	}

	arr, ok := node.(linkedast.Array)
	if !ok {
		return nil, errors.New("unexpected node type")
	}

	return arr.Elements(), nil
}

func (c *Linker) linkArray(scp scope, array ast.Array) (linkedast.Array, error) {
	elems := make([]linkedast.Expression, 0, len(array.Elements()))

	for _, elem := range array.Elements() {
		switch v := elem.(type) {
		case ast.Object:
			exp, err := c.linkObject(scp, v)
			if err != nil {
				return linkedast.Array{}, errors.Wrap(err, "link object")
			}
			elems = append(elems, exp)
		case ast.Spread:
			exps, err := c.linkArraySpread(scp, v)
			if err != nil {
				return linkedast.Array{}, errors.Wrap(err, "link spread")
			}

			elems = append(elems, exps...)
		case ast.Array:
			exp, err := c.linkArray(scp, v)
			if err != nil {
				return linkedast.Array{}, errors.Wrap(err, "link array")
			}
			elems = append(elems, exp)
		case ast.Var:
			node, err := c.findVariableExp(scp, v)
			if err != nil {
				return linkedast.Array{}, errors.Wrap(err, "find variable")
			}

			elems = append(elems, node)
		case ast.Bool:
			elems = append(elems, linkedast.NewBool(v.Value()))
		case ast.String:
			elems = append(elems, linkedast.NewString(v.Value()))
		case ast.Int:
			elems = append(elems, linkedast.NewInt(v.Value()))
		case ast.Float:
			elems = append(elems, linkedast.NewFloat(v.Value()))
		default:
			return linkedast.Array{}, errors.New("unknown value type")
		}
	}

	return linkedast.NewArray(elems), nil
}

func (c *Linker) findVariableExp(scp scope, v ast.Var) (linkedast.Expression, error) {
	linkedAst, ok := scp.linkedByName[v.Path()[0].String()]
	if !ok {
		return nil, errors.New("import for variable not found")
	}

	node, err := linkedAst.FindExpByPath(lo.Map(v.Path(), func(item ast.Ident, _ int) linkedast.Ident {
		return linkedast.NewIdent(item.String())
	}))
	if err != nil {
		return nil, errors.Wrap(err, "find node by path")
	}

	return node, nil
}
