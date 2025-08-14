package processor

import (
	"github.com/atmxlab/atmcfg/internal/lexer/tokenmover"
	ast2 "github.com/atmxlab/atmcfg/internal/linker/ast"
	"github.com/atmxlab/atmcfg/internal/parser/ast"
)

type OS interface {
	ReadFile(string) ([]byte, error)
	AbsPath(baseDir, relPath string) (string, error)
}

type Parser interface {
	Parse(mover tokenmover.TokenMover) (ast.Ast, error)
}

type Lexer interface {
	Tokenize(input string) (tokenmover.TokenMover, error)
}

type Linker interface {
	Link(mainAst ast.WithPath, astByPath map[string]ast.WithPath) (ast2.Ast, error)
}
