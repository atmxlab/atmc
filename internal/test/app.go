package test

import (
	"testing"

	"github.com/atmxlab/atmcfg/internal/analyzer"
	"github.com/atmxlab/atmcfg/internal/lexer"
	"github.com/atmxlab/atmcfg/internal/linker"
	"github.com/atmxlab/atmcfg/internal/parser"
	"github.com/atmxlab/atmcfg/internal/processor"
	"github.com/atmxlab/atmcfg/internal/test/testos"
)

type config struct {
	os testos.OS
}

func newConfig() *config {
	return &config{
		os: testos.NewOSBuilder().Build(),
	}
}

type ConfigOpt func(*config)

func WithOS(os testos.OS) ConfigOpt {
	return func(c *config) {
		c.os = os
	}
}

type App struct {
	t         *testing.T
	processor *processor.Processor
}

func NewApp(t *testing.T, opts ...ConfigOpt) *App {
	cfg := newConfig()
	for _, opt := range opts {
		opt(cfg)
	}

	p := processor.New(
		lexer.New(),
		parser.New(),
		analyzer.New(),
		linker.New(),
		cfg.os,
	)

	return &App{
		t:         t,
		processor: p,
	}
}

func (a App) Processor() *processor.Processor {
	return a.processor
}
