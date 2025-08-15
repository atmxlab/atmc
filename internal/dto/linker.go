package dto

import "github.com/atmxlab/atmcfg/internal/parser/ast"

type LinkParam struct {
	// AST основного конфигурационного файла.
	MainAst ast.WithPath
	// AST по пути нахождения файла.
	ASTByPath map[string]ast.WithPath
	// Переменные среды.
	Env map[string]string
}
