package compiler

import (
	"github.com/atmxlab/atmcfg/internal/lexer/tokenmover"
)

type Parser interface {
	Parse(mover tokenmover.TokenMover) (Ast, error)
}
