package ast

type WithPath struct {
	ast                    Ast
	absPath                string
	importAbsPathByRelPath map[string]string
}

func NewWithPath(ast Ast, absPath string, importAbsPathByRelPath map[string]string) WithPath {
	return WithPath{ast: ast, absPath: absPath, importAbsPathByRelPath: importAbsPathByRelPath}
}

func (a WithPath) AST() Ast {
	return a.ast
}

func (a WithPath) Path() string {
	return a.absPath
}

func (a WithPath) ImportPath(relPath string) (string, bool) {
	path, ok := a.importAbsPathByRelPath[relPath]
	return path, ok
}

func (a WithPath) Root() File {
	return a.ast.Root()
}

func (a WithPath) Imports() []Import {
	return a.ast.Root().Imports()
}
