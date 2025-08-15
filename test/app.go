package test

import (
	"testing"

	"github.com/atmxlab/atmc/analyzer"
	"github.com/atmxlab/atmc/lexer"
	"github.com/atmxlab/atmc/linker"
	"github.com/atmxlab/atmc/parser"
	"github.com/atmxlab/atmc/processor"
	"github.com/atmxlab/atmc/test/testos"
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
