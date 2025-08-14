package processor

import (
	"github.com/atmxlab/atmcfg/internal/parser/ast"
	"github.com/atmxlab/atmcfg/pkg/errors"
)

type Processor struct {
	lexer  Lexer
	parser Parser
	linker Linker
	os     OS
	// Необходим, чтобы обрабатывать повторные импорты.
	astByPath map[string]ast.WithPath
}

func (p *Processor) Process(path string) error {
	absPath, err := p.os.AbsPath(path, ".")
	if err != nil {
		return errors.Wrap(err, "get abs path")
	}

	if err = p.process(absPath, make(map[string]struct{})); err != nil {
		return errors.Wrap(err, "process")
	}

	_, err = p.linker.Link(p.astByPath[absPath], p.astByPath)
	if err != nil {
		return errors.Wrap(err, "linker.Link")
	}

	return nil
}

func (p *Processor) process(path string, importStack map[string]struct{}) error {
	if _, ok := p.astByPath[path]; ok {
		return nil
	}

	if _, ok := importStack[path]; ok {
		return errors.New("import cycle detected")
	}

	importStack[path] = struct{}{}

	madeAST, err := p.makeAst(path)
	if err != nil {
		return errors.Wrap(err, "make ast")
	}

	importPathByRelPath := make(map[string]string, len(madeAST.Imports()))

	for _, imp := range madeAST.Imports() {
		importPath, err := p.os.AbsPath(path, imp.Path().String())
		if err != nil {
			return errors.Wrap(err, "get abs path")
		}

		importPathByRelPath[imp.Path().String()] = importPath

		if err := p.process(importPath, p.copyMap(importStack)); err != nil {
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

	tokenMover, err := p.lexer.Tokenize(code)
	if err != nil {
		return ast.Ast{}, errors.Wrap(err, "tokenize")
	}

	a, err := p.parser.Parse(tokenMover)
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

func (p *Processor) copyMap(importStack map[string]struct{}) map[string]struct{} {
	copied := make(map[string]struct{}, len(importStack))
	for path := range importStack {
		copied[path] = struct{}{}
	}

	return copied
}
