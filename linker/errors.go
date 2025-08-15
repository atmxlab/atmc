package linker

import (
	"strings"

	"github.com/atmxlab/atmc/pkg/errors"
)

var (
	ErrUnexpectedNodeType = errors.New("unexpected node type")
	ErrNotFoundVariable   = errors.New("not found variable")
)

func newErrNotFoundVariable(variable ...string) error {
	return errors.Wrapf(ErrNotFoundVariable, "expected: %s", strings.Join(variable, "."))
}
