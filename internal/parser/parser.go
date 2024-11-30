package parser

import (
	"github.com/atmxlab/atmcfg/internal/parser/ast"
	"github.com/atmxlab/atmcfg/internal/types/token"
)

type Lexer interface {
	NextToken() token.Token
}

type Parser struct {
	lexer Lexer
}

func New(lexer Lexer) Parser {
	return Parser{lexer: lexer}
}

func (p Parser) Parse() ast.Ast {
	return ast.Ast{}
}
