package atmc

import (
	"encoding/json"

	"github.com/atmxlab/atmc/adapter"
	"github.com/atmxlab/atmc/analyzer"
	"github.com/atmxlab/atmc/compiler"
	"github.com/atmxlab/atmc/lexer"
	"github.com/atmxlab/atmc/linker"
	"github.com/atmxlab/atmc/linker/ast"
	"github.com/atmxlab/atmc/parser"
	"github.com/atmxlab/atmc/pkg/errors"
	"github.com/atmxlab/atmc/processor"
)

const (
	defaultFieldTag = "atmc"
)

type config struct {
	fieldTag string
}

type option func(*config)

func WithFieldTag(tag string) option {
	return func(c *config) {
		c.fieldTag = tag
	}
}

type ATMC struct {
	processor *processor.Processor
	config    config
}

func New(opts ...option) *ATMC {
	cfg := config{
		fieldTag: defaultFieldTag,
	}
	for _, opt := range opts {
		opt(&cfg)
	}
	return &ATMC{
		processor: processor.New(
			lexer.New(),
			parser.New(),
			analyzer.New(),
			linker.New(),
			adapter.NewOS(),
		),
		config: cfg,
	}
}

func (c *ATMC) Load(path string) (*Scanner, error) {
	a, err := c.processor.Process(path)
	if err != nil {
		return nil, errors.Wrap(err, "processor.Process")
	}

	return c.makeScanner(a), nil
}

func (c *ATMC) makeScanner(a ast.Ast) *Scanner {
	return NewScanner(
		a,
		compiler.NewMapCompiler(),
		compiler.NewStructCompiler(c.config.fieldTag),
	)
}

func (c *ATMC) JSON(path string) ([]byte, error) {
	scanner, err := c.Load(path)
	if err != nil {
		return nil, errors.Wrap(err, "load")
	}

	m := make(map[string]any)
	if err = scanner.Scan(m); err != nil {
		return nil, errors.Wrap(err, "scanner.Scan")
	}

	bytes, err := json.Marshal(m)
	if err != nil {
		return nil, errors.Wrap(err, "json.Marshal")
	}

	return bytes, nil
}
