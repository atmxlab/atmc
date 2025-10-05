package processor

import (
	"github.com/atmxlab/atmcfg/internal/linker"
	linkedast "github.com/atmxlab/atmcfg/internal/linker/ast"
	"github.com/atmxlab/atmcfg/internal/parser"
	"github.com/atmxlab/atmcfg/internal/parser/ast"
	"github.com/atmxlab/atmcfg/internal/types/token"
)

//go:generate mock OS
type OS interface {
	ReadFile(string) ([]byte, error)
	AbsPath(baseDir, relPath string) (string, error)
}

//go:generate mock Parser
type Parser interface {
	Parse(mover parser.TokenMover) (ast.Ast, error)
}

//go:generate mock Lexer
type Lexer interface {
	Tokenize(input string) ([]token.Token, error)
}

//go:generate mock Linker
type Linker interface {
	Link(param linker.LinkParam) (linkedast.Ast, error)
}
