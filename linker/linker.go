package linker

import (
	ast3 "github.com/atmxlab/atmc/linker/ast"
	ast2 "github.com/atmxlab/atmc/parser/ast"
	"github.com/atmxlab/atmc/pkg/errors"
	"github.com/atmxlab/atmc/pkg/orderedset"
	"github.com/samber/lo"
)

type Linker struct {
	// Необходим, чтобы линковать импортированные ast.
	astByPath map[string]ast2.WithPath
	// Необходим, чтобы обрабатывать повторные импорты.
	linkedByPath map[string]ast3.Ast
	// Необходим, чтобы резолвить переменные среды.
	env map[string]string
}

func New() *Linker {
	return &Linker{
		astByPath:    make(map[string]ast2.WithPath),
		linkedByPath: make(map[string]ast3.Ast),
		env:          make(map[string]string),
	}
}

type scope struct {
	// Необходим, чтобы добираться до внутренностей переменных по названию.
	linkedByName map[string]ast3.Ast
	ast          ast2.WithPath
}

func newScope(a ast2.WithPath) scope {
	return scope{
		linkedByName: make(map[string]ast3.Ast),
		ast:          a,
	}
}

type LinkParam struct {
	// AST основного конфигурационного файла.
	MainAst ast2.WithPath
	// AST по пути нахождения файла.
	ASTByPath map[string]ast2.WithPath
	// Переменные среды.
	Env map[string]string
}

func (l *Linker) Link(param LinkParam) (ast3.Ast, error) {
	l.astByPath = param.ASTByPath
	l.env = param.Env

	return l.link(newScope(param.MainAst))
}

func (l *Linker) link(scp scope) (ast3.Ast, error) {
	for _, imp := range scp.ast.Imports() {
		absPath, ok := scp.ast.ImportPath(imp.Path().String())
		if !ok {
			return ast3.Ast{}, errors.New("get import absolute path by relative path")
		}

		if alreadyLinked, ok := l.linkedByPath[absPath]; ok {
			scp.linkedByName[imp.Name().String()] = alreadyLinked
			return alreadyLinked, nil
		}

		astForLink, ok := l.astByPath[absPath]
		if !ok {
			return ast3.Ast{}, errors.New("ast for link not found")
		}

		linked, err := l.link(newScope(astForLink))
		if err != nil {
			return ast3.Ast{}, errors.Wrapf(err, "link ast, path: [%s]", imp.Path().String())
		}

		l.linkedByPath[absPath] = linked
		scp.linkedByName[imp.Name().String()] = linked
	}

	obj, err := l.linkObject(scp, scp.ast.Root().Object())
	if err != nil {
		return ast3.Ast{}, errors.Wrap(err, "link object")
	}

	return ast3.NewAst(obj), nil
}

func (l *Linker) linkObject(scp scope, obj ast2.Object) (ast3.Object, error) {
	entries, err := l.linkEntries(scp, obj.Entries())
	if err != nil {
		return ast3.Object{}, errors.Wrap(err, "link entries")
	}

	return ast3.NewObject(entries), nil
}

func (l *Linker) linkEntries(scp scope, entries []ast2.Entry) ([]ast3.KV, error) {
	kvMap := orderedset.New[ast3.Ident, ast3.KV](0)

	for _, entry := range entries {
		switch e := entry.(type) {
		case ast2.KV:
			ent, err := l.linkKV(scp, e)
			if err != nil {
				return nil, errors.Wrap(err, "link kv")
			}
			if existingEntry, exist := kvMap.Get(ent.Key()); exist {
				kvMap.Set(ent.Key(), l.mergeEntries(existingEntry, ent))
			} else {
				kvMap.Set(ent.Key(), ent)
			}
		case ast2.Spread:
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

func (l *Linker) linkKV(scp scope, kv ast2.KV) (ast3.KV, error) {
	var value ast3.Expression
	switch v := kv.Value().(type) {
	case ast2.Object:
		exp, err := l.linkObject(scp, v)
		if err != nil {
			return ast3.KV{}, errors.Wrap(err, "link object")
		}

		value = exp
	case ast2.Array:
		exp, err := l.linkArray(scp, v)
		if err != nil {
			return ast3.KV{}, errors.Wrap(err, "link array")
		}

		value = exp
	case ast2.Var:
		node, err := l.findVariableExp(scp, v)
		if err != nil {
			return ast3.KV{}, errors.Wrap(err, "find variable")
		}

		value = node
	case ast2.Env:
		value = ast3.NewString(l.getEnv(v.Name().String()))
	case ast2.Bool:
		value = ast3.NewBool(v.Value())
	case ast2.String:
		value = ast3.NewString(v.Value())
	case ast2.Int:
		value = ast3.NewInt(v.Value())
	case ast2.Float:
		value = ast3.NewFloat(v.Value())
	default:
		return ast3.KV{}, errors.New("unknown value type")
	}

	return ast3.NewKV(ast3.NewIdent(kv.Key().String()), value), nil
}

func (l *Linker) linkObjectSpread(scp scope, spread ast2.Spread) ([]ast3.KV, error) {
	node, err := l.findVariableExp(scp, spread.Var())
	if err != nil {
		return nil, errors.Wrap(err, "find variable node")
	}

	obj, ok := node.(ast3.Object)
	if !ok {
		return nil, errors.Wrap(ErrUnexpectedNodeType, "expected: Object")
	}

	return obj.KV(), nil
}

func (l *Linker) linkArraySpread(scp scope, spread ast2.Spread) ([]ast3.Expression, error) {
	node, err := l.findVariableExp(scp, spread.Var())
	if err != nil {
		return nil, errors.Wrap(err, "find variable node")
	}

	arr, ok := node.(ast3.Array)
	if !ok {
		return nil, errors.Wrap(ErrUnexpectedNodeType, "expected: Array")
	}

	return arr.Elements(), nil
}

func (l *Linker) linkArray(scp scope, array ast2.Array) (ast3.Array, error) {
	elems := make([]ast3.Expression, 0, len(array.Elements()))

	for _, elem := range array.Elements() {
		switch v := elem.(type) {
		case ast2.Object:
			exp, err := l.linkObject(scp, v)
			if err != nil {
				return ast3.Array{}, errors.Wrap(err, "link object")
			}
			elems = append(elems, exp)
		case ast2.Spread:
			exps, err := l.linkArraySpread(scp, v)
			if err != nil {
				return ast3.Array{}, errors.Wrap(err, "link spread")
			}

			elems = append(elems, exps...)
		case ast2.Array:
			exp, err := l.linkArray(scp, v)
			if err != nil {
				return ast3.Array{}, errors.Wrap(err, "link array")
			}
			elems = append(elems, exp)
		case ast2.Var:
			node, err := l.findVariableExp(scp, v)
			if err != nil {
				return ast3.Array{}, errors.Wrap(err, "find variable")
			}

			elems = append(elems, node)
		case ast2.Env:
			elems = append(elems, ast3.NewString(l.getEnv(v.Name().String())))
		case ast2.Bool:
			elems = append(elems, ast3.NewBool(v.Value()))
		case ast2.String:
			elems = append(elems, ast3.NewString(v.Value()))
		case ast2.Int:
			elems = append(elems, ast3.NewInt(v.Value()))
		case ast2.Float:
			elems = append(elems, ast3.NewFloat(v.Value()))
		default:
			return ast3.Array{}, errors.New("unknown value type")
		}
	}

	return ast3.NewArray(elems), nil
}

func (l *Linker) findVariableExp(scp scope, v ast2.Var) (ast3.Expression, error) {
	linkedAst, ok := scp.linkedByName[v.Path()[0].String()]
	if !ok {
		return nil, newErrNotFoundVariable(v.Path()[0].String())
	}

	node, err := linkedAst.FindExpByPath(lo.Map(v.Path()[1:], func(item ast2.Ident, _ int) ast3.Ident {
		return ast3.NewIdent(item.String())
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

func (l *Linker) mergeEntries(entry1, entry2 ast3.KV) ast3.KV {
	v1, ok1 := entry1.Value().(ast3.Object)
	v2, ok2 := entry2.Value().(ast3.Object)
	if !ok1 || !ok2 {
		return entry2
	}

	kvMap := orderedset.New[ast3.Ident, ast3.KV](0)
	for _, v := range v1.KV() {
		if existingEntry, exist := kvMap.Get(v.Key()); exist {
			kvMap.Set(v.Key(), l.mergeEntries(existingEntry, v))
		} else {
			kvMap.Set(v.Key(), v)
		}
	}
	for _, v := range v2.KV() {
		if existingEntry, exist := kvMap.Get(v.Key()); exist {
			kvMap.Set(v.Key(), l.mergeEntries(existingEntry, v))
		} else {
			kvMap.Set(v.Key(), v)
		}
	}

	return ast3.NewKV(entry1.Key(), ast3.NewObject(kvMap.Values()))
}
