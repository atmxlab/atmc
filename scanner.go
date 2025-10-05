package atmcfg

import (
	"github.com/atmxlab/atmcfg/internal/compiler"
	linkedast "github.com/atmxlab/atmcfg/internal/linker/ast"
	"github.com/atmxlab/atmcfg/pkg/errors"
)

type Scanner struct {
	ast         linkedast.Ast
	mapCompiler *compiler.MapCompiler
}

func (s *Scanner) Scan(t any) error {
	switch v := t.(type) {
	case map[string]any:
		if err := s.mapCompiler.Compile(v, s.ast); err != nil {
			return errors.Wrap(err, "mapCompiler.Compile")
		}
		return nil
	default:
		return errors.New("unsupported type")
	}
}

func NewScanner(ast linkedast.Ast) *Scanner {
	return &Scanner{
		ast:         ast,
		mapCompiler: compiler.NewMapCompiler(),
	}
}
