package linker

import (
	linkedast "github.com/atmxlab/atmcfg/internal/linker/ast"
	"github.com/atmxlab/atmcfg/internal/parser/ast"
	"github.com/atmxlab/atmcfg/pkg/errors"
	"github.com/atmxlab/atmcfg/pkg/orderedmap"
	"github.com/samber/lo"
)

type Linker struct {
	// Необходим, чтобы линковать импортированные ast.
	astByPath map[string]ast.WithPath
	// Необходим, чтобы обрабатывать повторные импорты.
	linkedByPath map[string]linkedast.Ast
	// Необходим, чтобы резолвить переменные среды.
	env map[string]string
}

func New() *Linker {
	return &Linker{
		astByPath:    make(map[string]ast.WithPath),
		linkedByPath: make(map[string]linkedast.Ast),
		env:          make(map[string]string),
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

type LinkParam struct {
	// AST основного конфигурационного файла.
	MainAst ast.WithPath
	// AST по пути нахождения файла.
	ASTByPath map[string]ast.WithPath
	// Переменные среды.
	Env map[string]string
}

func (l *Linker) Link(param LinkParam) (linkedast.Ast, error) {
	l.astByPath = param.ASTByPath
	l.env = param.Env

	return l.link(newScope(param.MainAst))
}

func (l *Linker) link(scp scope) (linkedast.Ast, error) {
	for _, imp := range scp.ast.Imports() {
		absPath, ok := scp.ast.ImportPath(imp.Path().String())
		if !ok {
			return linkedast.Ast{}, errors.New("get import absolute path by relative path")
		}

		if alreadyLinked, ok := l.linkedByPath[absPath]; ok {
			scp.linkedByName[imp.Name().String()] = alreadyLinked
			return alreadyLinked, nil
		}

		astForLink, ok := l.astByPath[absPath]
		if !ok {
			return linkedast.Ast{}, errors.New("ast for link not found")
		}

		linked, err := l.link(newScope(astForLink))
		if err != nil {
			return linkedast.Ast{}, errors.Wrapf(err, "link ast, path: [%s]", imp.Path().String())
		}

		l.linkedByPath[absPath] = linked
		scp.linkedByName[imp.Name().String()] = linked
	}

	obj, err := l.linkObject(scp, scp.ast.Root().Object())
	if err != nil {
		return linkedast.Ast{}, errors.Wrap(err, "link object")
	}

	return linkedast.NewAst(obj), nil
}

func (l *Linker) linkObject(scp scope, obj ast.Object) (linkedast.Object, error) {
	entries, err := l.linkEntries(scp, obj.Entries())
	if err != nil {
		return linkedast.Object{}, errors.Wrap(err, "link entries")
	}

	return linkedast.NewObject(entries), nil
}

func (l *Linker) linkEntries(scp scope, entries []ast.Entry) ([]linkedast.KV, error) {
	kvMap := orderedmap.NewOrderedMap[linkedast.Ident, linkedast.KV]()

	for _, entry := range entries {
		switch e := entry.(type) {
		case ast.KV:
			ent, err := l.linkKV(scp, e)
			if err != nil {
				return nil, errors.Wrap(err, "link kv")
			}
			kvMap.Set(ent.Key(), ent)
		case ast.Spread:
			spreadEntries, err := l.linkObjectSpread(scp, e)
			if err != nil {
				return nil, errors.Wrap(err, "link spread")
			}

			for _, spreadEntry := range spreadEntries {
				kvMap.Set(spreadEntry.Key(), spreadEntry)
			}
		default:
			return nil, errors.New("unknown entry type")
		}
	}

	return kvMap.Values(), nil
}

func (l *Linker) linkKV(scp scope, kv ast.KV) (linkedast.KV, error) {
	var value linkedast.Expression
	switch v := kv.Value().(type) {
	case ast.Object:
		exp, err := l.linkObject(scp, v)
		if err != nil {
			return linkedast.KV{}, errors.Wrap(err, "link object")
		}

		value = exp
	case ast.Array:
		exp, err := l.linkArray(scp, v)
		if err != nil {
			return linkedast.KV{}, errors.Wrap(err, "link array")
		}

		value = exp
	case ast.Var:
		node, err := l.findVariableExp(scp, v)
		if err != nil {
			return linkedast.KV{}, errors.Wrap(err, "find variable")
		}

		value = node
	case ast.Env:
		value = linkedast.NewString(l.getEnv(v.Name().String()))
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

func (l *Linker) linkObjectSpread(scp scope, spread ast.Spread) ([]linkedast.KV, error) {
	node, err := l.findVariableExp(scp, spread.Var())
	if err != nil {
		return nil, errors.Wrap(err, "find variable node")
	}

	obj, ok := node.(linkedast.Object)
	if !ok {
		return nil, errors.Wrap(ErrUnexpectedNodeType, "expected: Object")
	}

	return obj.KV(), nil
}

func (l *Linker) linkArraySpread(scp scope, spread ast.Spread) ([]linkedast.Expression, error) {
	node, err := l.findVariableExp(scp, spread.Var())
	if err != nil {
		return nil, errors.Wrap(err, "find variable node")
	}

	arr, ok := node.(linkedast.Array)
	if !ok {
		return nil, errors.Wrap(ErrUnexpectedNodeType, "expected: Array")
	}

	return arr.Elements(), nil
}

func (l *Linker) linkArray(scp scope, array ast.Array) (linkedast.Array, error) {
	elems := make([]linkedast.Expression, 0, len(array.Elements()))

	for _, elem := range array.Elements() {
		switch v := elem.(type) {
		case ast.Object:
			exp, err := l.linkObject(scp, v)
			if err != nil {
				return linkedast.Array{}, errors.Wrap(err, "link object")
			}
			elems = append(elems, exp)
		case ast.Spread:
			exps, err := l.linkArraySpread(scp, v)
			if err != nil {
				return linkedast.Array{}, errors.Wrap(err, "link spread")
			}

			elems = append(elems, exps...)
		case ast.Array:
			exp, err := l.linkArray(scp, v)
			if err != nil {
				return linkedast.Array{}, errors.Wrap(err, "link array")
			}
			elems = append(elems, exp)
		case ast.Var:
			node, err := l.findVariableExp(scp, v)
			if err != nil {
				return linkedast.Array{}, errors.Wrap(err, "find variable")
			}

			elems = append(elems, node)
		case ast.Env:
			elems = append(elems, linkedast.NewString(l.getEnv(v.Name().String())))
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

func (l *Linker) findVariableExp(scp scope, v ast.Var) (linkedast.Expression, error) {
	linkedAst, ok := scp.linkedByName[v.Path()[0].String()]
	if !ok {
		return nil, newErrNotFoundVariable(v.Path()[0].String())
	}

	node, err := linkedAst.FindExpByPath(lo.Map(v.Path()[1:], func(item ast.Ident, _ int) linkedast.Ident {
		return linkedast.NewIdent(item.String())
	}))
	switch {
	case errors.Is(err, errors.ErrNotFound):
		return nil, newErrNotFoundVariable(v.StringPath()...)
	case err != nil:
		return nil, errors.Wrap(err, "find node by path")
	}

	return node, nil
}

func (l *Linker) getEnv(name string) string {
	return l.env[name]
}
