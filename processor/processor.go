package processor

import (
	"path/filepath"

	"github.com/atmxlab/atmc/lexer/tokenmover"
	"github.com/atmxlab/atmc/linker"
	linkedast "github.com/atmxlab/atmc/linker/ast"
	ast2 "github.com/atmxlab/atmc/parser/ast"
	"github.com/atmxlab/atmc/pkg/errors"
)

type Processor struct {
	os        OS
	lexer     Lexer
	parser    Parser
	analyzer  Analyzer
	linker    Linker
	astByPath map[string]ast2.WithPath
}

func New(lexer Lexer, parser Parser, analyzer Analyzer, linker Linker, os OS) *Processor {
	return &Processor{
		os:        os,
		lexer:     lexer,
		parser:    parser,
		analyzer:  analyzer,
		linker:    linker,
		astByPath: make(map[string]ast2.WithPath),
	}
}

func (p *Processor) Process(path string) (linkedast.Ast, error) {
	absPath, err := p.os.AbsPath(path, ".")
	if err != nil {
		return linkedast.Ast{}, errors.Wrap(err, "get abs path")
	}

	if err = p.process(absPath, newEmptyImportStack()); err != nil {
		return linkedast.Ast{}, errors.Wrap(err, "process")
	}

	linkedAst, err := p.linker.Link(linker.LinkParam{
		MainAst:   p.astByPath[absPath],
		ASTByPath: p.astByPath,
		Env:       p.os.EnvVariables(),
	})
	if err != nil {
		return linkedast.Ast{}, errors.Wrap(err, "linker.Link")
	}

	return linkedAst, nil
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

	p.astByPath[path] = ast2.NewWithPath(madeAST, path, importPathByRelPath)

	return nil
}

func (p *Processor) makeAst(path string) (ast2.Ast, error) {
	code, err := p.readFileContent(path)
	if err != nil {
		return ast2.Ast{}, errors.Wrap(err, "read file content")
	}

	tokens, err := p.lexer.Tokenize(code)
	if err != nil {
		return ast2.Ast{}, errors.Wrap(err, "tokenize")
	}

	a, err := p.parser.Parse(tokenmover.New(tokens))
	if err != nil {
		return ast2.Ast{}, errors.Wrap(err, "parse")
	}

	if err = p.analyzer.Analyze(a); err != nil {
		return ast2.Ast{}, errors.Wrapf(err, "semantic analyzer error: file: %s", path)
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
