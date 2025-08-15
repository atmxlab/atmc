package lexer

import (
	"github.com/atmxlab/atmc/pkg/errors"
	"github.com/atmxlab/atmc/types"
)

func unexpectedTokenError(pos types.Position) error {
	return errors.Newf("unexpected token at %v:%v", pos.Line(), pos.Column())
}
