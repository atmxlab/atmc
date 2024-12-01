package parser

import (
	"github.com/atmxlab/atmcfg/internal/parser/ast"
	"github.com/atmxlab/atmcfg/internal/types/token"
	"github.com/atmxlab/atmcfg/pkg/errors"
)

type Lexer interface {
	Token() token.Token
	Next()
	Prev()
	IsEmpty() bool
}

type Parser struct {
	lexer Lexer
}

func New(lexer Lexer) Parser {
	return Parser{lexer: lexer}
}

func (p Parser) Parse() (ast.Ast, error) {
	file, err := p.parseFile()
	if err != nil {
		return ast.Ast{}, errors.Wrap(err, "failed to parse file")
	}

	return ast.NewAst(file), nil
}

func (p Parser) parseFile() (ast.File, error) {
	var object ast.Object
	var foundObject bool
	imports := make([]ast.Import, 0)

	for !p.lexer.IsEmpty() {
		var err error
		var node ast.Node

		switch p.lexer.Token().Type() {
		case token.Ident:
			node, err = p.parseIndent()
			if err != nil {
				return ast.File{}, errors.Wrap(err, "failed to parse indent")
			}
		case token.LBrace:
			node, err = p.parseObject()
			if err != nil {
				return ast.File{}, errors.Wrap(err, "failed to parse object")
			}
		default:
			return ast.File{}, errors.New("unexpected token")
		}

		switch v := node.(type) {
		case ast.Import:
			imports = append(imports, v)
		case ast.Object:
			if foundObject {
				return ast.File{}, errors.New("found multiple objects")
			}

			object = v
		default:
			return ast.File{}, errors.Wrap(ErrUnexpectedToken, "unexpected node in file")
		}

		p.lexer.Next()
	}

	return ast.NewFile(imports, object), nil
}

func (p Parser) parseIndent() (ast.Node, error) {
	err := p.require(token.Ident)
	if err != nil {
		return nil, err
	}

	prevToken := p.lexer.Token()
	p.lexer.Next()

	switch p.lexer.Token().Type() {
	case token.Path:
		return ast.NewImport(
			ast.NewName(prevToken.Value().String()),
			ast.NewPath(p.lexer.Token().Value().String()),
		), nil
	case token.Colon:
		p.lexer.Next()
		entry, err := p.parseEntry()
		if err != nil {
			return nil, errors.Wrap(err, "failed to parse entry")
		}

		return ast.NewKeyValue(ast.NewKey(prevToken.Value().String()), entry), nil
	case token.Dot:
		expr, err := p.parseVar()
		if err != nil {
			return nil, errors.Wrap(err, "failed to parse spread")
		}

		return expr, nil
	case token.Spread:
		spreadExp, err := p.parseSpread()
		if err != nil {
			return nil, errors.Wrap(err, "failed to parse spread")
		}

		return spreadExp, nil
	case token.Comma, token.RBrace, token.RBracket:
		v := ast.NewVar([]ast.Ident{ast.NewName(p.lexer.Token().Value().String())})
		p.lexer.Prev()
		return v, nil
	default:
		return nil, NewErrUnexpectedToken(token.Path, token.Colon)
	}
}

func (p Parser) parseSpread() (ast.Expression, error) {
	if err := p.require(token.Spread); err != nil {
		return nil, err
	}

	p.lexer.Prev()

	v := ast.NewVar([]ast.Ident{
		ast.NewName(p.lexer.Token().Value().String()),
	})

	p.lexer.Next()

	return ast.NewSpread(v), nil
}

func (p Parser) parseVar() (ast.Expression, error) {
	idents := make([]ast.Ident, 0)

	if err := p.require(token.Dot); err != nil {
		return nil, err
	}

	p.lexer.Prev()

	if err := p.require(token.Ident); err != nil {
		return nil, err
	}

	idents = append(idents, ast.NewName(p.lexer.Token().Value().String()))

	p.lexer.Next()
	p.lexer.Next()

	if err := p.require(token.Ident); err != nil {
		return nil, err
	}

	idents = append(idents, ast.NewName(p.lexer.Token().Value().String()))

	node, err := p.parseIndent()
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse indent")
	}

	switch expr := node.(type) {
	case ast.Var:
		idents = append(idents, expr.Path()[1:]...)
		return ast.NewVar(idents), nil
	case ast.Spread:
		idents = append(idents, expr.Var().Path()[1:]...)
		spreadExp := ast.NewSpread(ast.NewVar(idents))
		return spreadExp, err
	}

	return ast.NewVar(idents), nil
}

func (p Parser) parseEntry() (ast.Entry, error) {
	switch p.lexer.Token().Type() {
	case token.Ident:
		node, err := p.parseIndent()
		if err != nil {
			return nil, errors.Wrap(err, "failed to parse indent")
		}

		switch v := node.(type) {
		case ast.Entry:
			return v, nil
		default:
			return nil, NewUnexpectedNodeErr("var", "spread")
		}
	case token.LBrace:
		obj, err := p.parseObject()
		if err != nil {
			return nil, errors.Wrap(err, "failed to parse object")
		}

		return obj, nil
	case token.LBracket:
		arr, err := p.parseArray()
		if err != nil {
			return nil, errors.Wrap(err, "failed to parse array")
		}
		return arr, nil
	case token.Int:
		i, err := ast.NewInt(p.lexer.Token().Value().String())
		if err != nil {
			return nil, errors.Wrap(err, "ast new int")
		}

		return i, nil
	case token.Float:
		i, err := ast.NewFloat(p.lexer.Token().Value().String())
		if err != nil {
			return nil, errors.Wrap(err, "ast new float")
		}

		return i, nil
	case token.String:
		return ast.NewString(p.lexer.Token().Value().String()), nil
	case token.Bool:
		i, err := ast.NewBool(p.lexer.Token().Value().String())
		if err != nil {
			return nil, errors.Wrap(err, "ast new bool")
		}

		return i, nil
	default:
		return nil, NewErrUnexpectedToken(
			token.LBrace,
			token.LBracket,
			token.Int,
			token.Float,
			token.String,
			token.Bool,
		)
	}
}

func (p Parser) parseObject() (ast.Entry, error) {
	err := p.require(token.LBrace)
	if err != nil {
		return nil, err
	}

	p.lexer.Next()

	keyValues := make([]ast.KeyValue, 0)
	spreads := make([]ast.Spread, 0)

	var node ast.Node
	for p.match(token.Ident) {
		node, err = p.parseIndent()
		switch v := node.(type) {
		case ast.KeyValue:
			keyValues = append(keyValues, v)
		case ast.Spread:
			spreads = append(spreads, v)
		default:
			return nil, NewUnexpectedNodeErr("KeyValue")
		}

		p.lexer.Next()
	}

	err = p.require(token.RBrace)
	if err != nil {
		return nil, err
	}

	return ast.NewObject(spreads, keyValues), nil
}

func (p Parser) parseArray() (ast.Entry, error) {
	if err := p.require(token.LBracket); err != nil {
		return nil, err
	}

	p.lexer.Next()

	check := func() error {
		p.lexer.Next()
		if err := p.require(token.Comma, token.RBracket); err != nil {
			return errors.Wrap(err, "after array elem must be comma or right bracket")
		}
		p.lexer.Prev()

		return nil
	}

	elements := make([]ast.Node, 0)

	for {
		switch p.lexer.Token().Type() {
		case
			token.Ident,
			token.LBrace,
			token.LBracket,
			token.Int,
			token.Float,
			token.String,
			token.Bool:
			node, err := p.parseEntry()
			if err != nil {
				return nil, errors.Wrap(err, "failed to parse entry¬")
			}

			elements = append(elements, node)
			if err = check(); err != nil {
				return nil, err
			}
		case token.Comma:
		case token.RBracket:
			return ast.NewArray(elements), nil
		default:
			return nil, NewErrUnexpectedToken(
				token.LBrace,
				token.LBracket,
				token.Int,
				token.Float,
				token.String,
				token.Bool,
				token.Comma,
				token.RBracket,
			)
		}

		p.lexer.Next()
	}
}
