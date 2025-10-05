package test

import (
	"testing"

	"github.com/atmxlab/atmcfg/internal/lexer"
	"github.com/atmxlab/atmcfg/internal/linker"
	"github.com/atmxlab/atmcfg/internal/parser"
	"github.com/atmxlab/atmcfg/internal/processor"
	"github.com/atmxlab/atmcfg/internal/test/testos"
	"github.com/atmxlab/atmcfg/pkg/errors"
	"github.com/stretchr/testify/require"
)

func TestProcessor_SimplePath(t *testing.T) {
	t.Parallel()

	var (
		content = "{a: 1, b: 2}"
		absPath = "/home/user/config.atmc"
	)

	os := testos.
		NewOSBuilder().
		Content(absPath, []byte(content)).
		Build()
	lex := lexer.New()
	parse := parser.New()
	link := linker.New()

	proc := processor.New(
		lex,
		parse,
		link,
		os,
	)

	err := proc.Process(absPath)
	require.NoError(t, err)
}

func TestProcessor_WithImport_ImportFileNotFound(t *testing.T) {
	t.Parallel()

	var (
		content = "var1 ./import.atmc {a: 1, b: 2}"
		absPath = "/home/user/config.atmc"
	)

	os := testos.
		NewOSBuilder().
		Content(absPath, []byte(content)).
		Build()
	lex := lexer.New()
	parse := parser.New()
	link := linker.New()

	proc := processor.New(
		lex,
		parse,
		link,
		os,
	)

	err := proc.Process(absPath)
	require.ErrorIs(t, err, errors.ErrNotFound)
}

func TestProcessor_WithImport_EmptyContentInImportFile(t *testing.T) {
	t.Parallel()

	var (
		content       = "var1 ./import.atmc {a: 1, b: 2}"
		absPath       = "/home/user/config.atmc"
		importAbsPath = "/home/user/import.atmc"
	)

	os := testos.
		NewOSBuilder().
		Content(absPath, []byte(content)).
		Content(importAbsPath, []byte("")).
		Build()
	lex := lexer.New()
	parse := parser.New()
	link := linker.New()

	proc := processor.New(
		lex,
		parse,
		link,
		os,
	)

	err := proc.Process(absPath)
	require.ErrorIs(t, err, parser.ErrTokenNotExist)
}

func TestProcessor_WithImport(t *testing.T) {
	t.Parallel()

	var (
		absPath         = "/home/user/config.atmc"
		content         = "var1 ./import.atmc {a: 1, b: 2}"
		importAbsPath   = "/home/user/import.atmc"
		importedContent = "{c: 3, d: 4}"
	)

	os := testos.
		NewOSBuilder().
		Content(absPath, []byte(content)).
		Content(importAbsPath, []byte(importedContent)).
		Build()
	lex := lexer.New()
	parse := parser.New()
	link := linker.New()

	proc := processor.New(
		lex,
		parse,
		link,
		os,
	)

	err := proc.Process(absPath)
	require.NoError(t, err)
}
