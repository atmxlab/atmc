package compiler

import (
	"github.com/atmxlab/atmcfg/internal/lexer/tokenmover"
)

type Lexer interface {
	Tokenize(input string) (tokenmover.TokenMover, error)
}
