package analyzer

import (
	ast2 "github.com/atmxlab/atmc/parser/ast"
	"github.com/atmxlab/atmc/pkg/errors"
)

type Analyzer struct {
	scope *scope
}

func New() *Analyzer {
	return &Analyzer{scope: newScope()}
}

func (ar *Analyzer) Analyze(a ast2.Ast) error {
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

func (ar *Analyzer) Visit(node ast2.Node) error {
	switch n := node.(type) {
	case ast2.File:
	case ast2.Import:
		ar.scope.addVariable(n.Name().String())
	case ast2.Object:
	case ast2.Spread:
	case ast2.KV:
	case ast2.Array:
	case ast2.Var:
		err := ar.checkVar(n)
		if err != nil {
			return errors.Wrap(err, "check variable")
		}
	case ast2.Env:
	case ast2.Int:
	case ast2.Float:
	case ast2.String:
	case ast2.Bool:
	default:
		return errors.New("invalid node type")
	}

	return nil
}

func (ar *Analyzer) checkVar(v ast2.Var) error {
	if len(v.Path()) == 0 {
		return errors.Newf("invalid variable. variable path is empty")
	}
	firstPartFromVarPath := v.Path()[0]

	if !ar.scope.hasVariable(firstPartFromVarPath.String()) {
		return errors.Wrapf(ErrUndefinedVariable, "undefined variable: %s", firstPartFromVarPath.String())
	}

	ar.scope.incrRef(firstPartFromVarPath.String())

	return nil
}
