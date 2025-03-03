package semantic

import (
	"github.com/atmxlab/atmcfg/internal/compiler"
	"github.com/atmxlab/atmcfg/internal/parser/ast"
	"github.com/atmxlab/atmcfg/pkg/errors"
)

type Analyzer struct {
	scope *scope
}

func NewAnalyzer() *Analyzer {
	return &Analyzer{scope: newScope()}
}

func (ar *Analyzer) Analyze(a compiler.Ast) error {
	err := a.Inspect(ar.Visit)
	if err != nil {
		return errors.Wrap(err, "inspect")
	}

	err = ar.scope.checkVariableRefs()
	if err != nil {
		return errors.Wrap(err, "check variables refs")
	}

	return nil
}

func (ar *Analyzer) Visit(node ast.Node) error {
	switch n := node.(type) {
	case ast.File:
	case ast.Import:
		ar.scope.addVariable(n.Name().String())
	case ast.Object:
	case ast.Spread:
	case ast.KV:
	case ast.Array:
	case ast.Var:
		err := ar.checkVar(n)
		if err != nil {
			return errors.Wrap(err, "check variable")
		}
	case ast.Env:
	case ast.Int:
	case ast.Float:
	case ast.String:
	case ast.Bool:
	default:
		return errors.New("invalid node type")
	}

	return nil
}

func (ar *Analyzer) checkVar(v ast.Var) error {
	if len(v.Path()) == 0 {
		return errors.Newf("invalid variable. variable path is empty") // TODO: нормально возвращать ошибку (с доп информацией)
	}
	firstPartFromVarPath := v.Path()[0]

	if !ar.scope.hasVariable(firstPartFromVarPath.String()) {
		return errors.Wrapf(ErrUndefinedVariable, "undefined variable: %s", firstPartFromVarPath.String())
	}

	ar.scope.incrRef(firstPartFromVarPath.String())

	return nil
}
