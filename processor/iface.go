package processor

import (
	"github.com/atmxlab/atmc/linker"
	linkedast "github.com/atmxlab/atmc/linker/ast"
	"github.com/atmxlab/atmc/parser"
	"github.com/atmxlab/atmc/parser/ast"
	"github.com/atmxlab/atmc/types/token"
)

//go:generate mock OS
type OS interface {
	ReadFile(string) ([]byte, error)
	AbsPath(baseDir, relPath string) (string, error)
	EnvVariables() map[string]string
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

//go:generate mock Analyzer
type Analyzer interface {
	Analyze(a ast.Ast) error
}
