package lexer

import (
	"github.com/atmxlab/atmcfg/internal/types"
	"github.com/atmxlab/atmcfg/pkg/errors"
)

func unexpectedTokenError(pos *types.Position) error {
	return errors.Newf("unexpected token at %v:%v", pos.Line, pos.Column)
}
