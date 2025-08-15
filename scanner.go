package atmc

import (
	"github.com/atmxlab/atmc/compiler"
	linkedast "github.com/atmxlab/atmc/linker/ast"
	"github.com/atmxlab/atmc/pkg/errors"
)

type Scanner struct {
	ast            linkedast.Ast
	mapCompiler    *compiler.MapCompiler
	structCompiler *compiler.StructCompiler
}

func NewScanner(
	ast linkedast.Ast,
	mapCompiler *compiler.MapCompiler,
	structCompiler *compiler.StructCompiler,
) *Scanner {
	return &Scanner{
		ast:            ast,
		mapCompiler:    mapCompiler,
		structCompiler: structCompiler,
	}
}

func (s *Scanner) Scan(t any) error {
	switch v := t.(type) {
	case map[string]any:
		if err := s.mapCompiler.Compile(v, s.ast); err != nil {
			return errors.Wrap(err, "mapCompiler.Compile")
		}
	default:
		if err := s.structCompiler.Compile(v, s.ast); err != nil {
			return errors.Wrap(err, "structCompiler.Compile")
		}
	}

	return nil
}
