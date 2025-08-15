package processor

import (
	"github.com/atmxlab/atmcfg/internal/lexer/tokenmover"
	"github.com/atmxlab/atmcfg/internal/linker"
	linkedast "github.com/atmxlab/atmcfg/internal/linker/ast"
	"github.com/atmxlab/atmcfg/internal/parser"
	"github.com/atmxlab/atmcfg/internal/parser/ast"
)

type OS interface {
	ReadFile(string) ([]byte, error)
	AbsPath(baseDir, relPath string) (string, error)
}

type Parser interface {
	Parse(mover parser.TokenMover) (ast.Ast, error)
}

type Lexer interface {
	Tokenize(input string) (*tokenmover.TokenMover, error)
}

type Linker interface {
	Link(param linker.LinkParam) (linkedast.Ast, error)
}
