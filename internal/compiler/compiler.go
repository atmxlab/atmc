package compiler

import (
	"os"
	"path/filepath"

	"github.com/atmxlab/atmcfg/pkg/errors"
)

type Compiler struct {
	semanticAnalyzer SemanticAnalyzer
	lexer            Lexer
	parser           Parser
	scope            *scope
}

func NewCompiler(semanticAnalyzer SemanticAnalyzer, lexer Lexer, parser Parser) *Compiler {
	return &Compiler{
		semanticAnalyzer: semanticAnalyzer,
		lexer:            lexer,
		parser:           parser,
		scope:            newScope(),
	}
}

type scope struct {
	// map path by path
	imported map[string]Ast
	// map by path
	importedStack map[string]struct{}
}

func newScope() *scope {
	return &scope{
		imported:      make(map[string]Ast),
		importedStack: make(map[string]struct{}),
	}
}

func (c *Compiler) Compile(configPath string) (a Ast, err error) {
	configPath, err = filepath.Abs(configPath)
	if err != nil {
		return nil, errors.Wrap(err, "get abs path")
	}

	if a, ok := c.scope.imported[configPath]; ok {
		return a, nil
	}
	if _, ok := c.scope.importedStack[configPath]; ok {
		return nil, errors.New("import cycle detected")
	}

	c.scope.importedStack[configPath] = struct{}{}

	mainAst, err := c.makeAst(configPath)
	if err != nil {
		return nil, errors.Wrap(err, "make ast")
	}

	// map by import name
	importedAst := make(map[string]Ast, len(mainAst.Imports()))

	for _, imp := range mainAst.Imports() {

		a, err = c.Compile(imp.Path())
		if err != nil {
			return nil, errors.Wrap(err, "compile ast")
		}

		c.scope.imported[configPath] = a
		importedAst[imp.Name()] = a
	}

	delete(c.scope.importedStack, configPath)

	return mainAst, nil
}

func (c *Compiler) makeAst(configPath string) (Ast, error) {
	code, err := c.readFileContent(configPath)
	if err != nil {
		return nil, errors.Wrap(err, "read file content")
	}

	tokenMover, err := c.lexer.Tokenize(code)
	if err != nil {
		return nil, errors.Wrap(err, "tokenize")
	}

	a, err := c.parser.Parse(tokenMover)
	if err != nil {
		return nil, errors.Wrap(err, "parse")
	}

	if err = c.semanticAnalyzer.Analyze(a); err != nil {
		return nil, errors.Combine(err, ErrSemantic)
	}

	return a, nil
}

func (c *Compiler) readFileContent(filePath string) (string, error) {
	// TODO: тут можно избавиться от зависимости от файловой системы
	content, err := os.ReadFile(filePath)
	if err != nil {
		return "", errors.Wrap(err, "failed to read file")
	}

	return string(content), nil
}
