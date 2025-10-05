package processor

import (
	"path/filepath"

	"github.com/atmxlab/atmcfg/internal/lexer/tokenmover"
	"github.com/atmxlab/atmcfg/internal/linker"
	"github.com/atmxlab/atmcfg/internal/parser/ast"
	"github.com/atmxlab/atmcfg/pkg/errors"
)

type Processor struct {
	lexer     Lexer
	parser    Parser
	linker    Linker
	os        OS
	astByPath map[string]ast.WithPath
}

func New(lexer Lexer, parser Parser, linker Linker, os OS) *Processor {
	return &Processor{
		lexer:     lexer,
		parser:    parser,
		linker:    linker,
		os:        os,
		astByPath: make(map[string]ast.WithPath),
	}
}

func (p *Processor) Process(path string) error {
	absPath, err := p.os.AbsPath(path, ".")
	if err != nil {
		return errors.Wrap(err, "get abs path")
	}

	if err = p.process(absPath, newEmptyImportStack()); err != nil {
		return errors.Wrap(err, "process")
	}

	_, err = p.linker.Link(linker.LinkParam{
		MainAst:   p.astByPath[absPath],
		ASTByPath: p.astByPath,
		Env:       nil,
	})
	if err != nil {
		return errors.Wrap(err, "linker.Link")
	}

	return nil
}

func (p *Processor) process(path string, iStack importStack) error {
	if _, ok := p.astByPath[path]; ok {
		return nil
	}

	if _, ok := iStack[path]; ok {
		return errors.New("import cycle detected")
	}

	iStack[path] = struct{}{}

	madeAST, err := p.makeAst(path)
	if err != nil {
		return errors.Wrap(err, "make ast")
	}

	importPathByRelPath := make(map[string]string, len(madeAST.Imports()))

	for _, imp := range madeAST.Imports() {
		importPath, err := p.os.AbsPath(filepath.Dir(path), imp.Path().String())
		if err != nil {
			return errors.Wrap(err, "get abs path")
		}

		importPathByRelPath[imp.Path().String()] = importPath

		if err := p.process(importPath, iStack.Clone()); err != nil {
			return errors.Wrap(err, "process import")
		}
	}

	p.astByPath[path] = ast.NewWithPath(madeAST, path, importPathByRelPath)

	return nil
}

func (p *Processor) makeAst(path string) (ast.Ast, error) {
	code, err := p.readFileContent(path)
	if err != nil {
		return ast.Ast{}, errors.Wrap(err, "read file content")
	}

	tokens, err := p.lexer.Tokenize(code)
	if err != nil {
		return ast.Ast{}, errors.Wrap(err, "tokenize")
	}

	a, err := p.parser.Parse(tokenmover.New(tokens))
	if err != nil {
		return ast.Ast{}, errors.Wrap(err, "parse")
	}

	return a, nil
}

func (p *Processor) readFileContent(filePath string) (string, error) {
	content, err := p.os.ReadFile(filePath)
	if err != nil {
		return "", errors.Wrap(err, "failed to read file")
	}

	return string(content), nil
}
