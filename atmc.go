package atmcfg

import (
	"github.com/atmxlab/atmcfg/internal/adapter"
	"github.com/atmxlab/atmcfg/internal/analyzer"
	"github.com/atmxlab/atmcfg/internal/lexer"
	"github.com/atmxlab/atmcfg/internal/linker"
	"github.com/atmxlab/atmcfg/internal/parser"
	"github.com/atmxlab/atmcfg/internal/processor"
	"github.com/atmxlab/atmcfg/pkg/errors"
)

type ATMC struct {
	processor *processor.Processor
}

func (c *ATMC) Load(path string) (*Scanner, error) {
	a, err := c.processor.Process(path)
	if err != nil {
		return nil, errors.Wrap(err, "processor.Process")
	}

	return NewScanner(a), nil
}

func NewATMC() *ATMC {
	return &ATMC{
		processor: processor.New(
			lexer.New(),
			parser.New(),
			analyzer.New(),
			linker.New(),
			adapter.NewOS(),
		),
	}
}
