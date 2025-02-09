package compiler

type SemanticAnalyzer interface {
	Analyze(Ast) error
}
