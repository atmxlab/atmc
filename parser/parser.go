package parser

import (
	ast2 "github.com/atmxlab/atmc/parser/ast"
	"github.com/atmxlab/atmc/pkg/errors"
	"github.com/atmxlab/atmc/types"
	token2 "github.com/atmxlab/atmc/types/token"
)

type TokenMover interface {
	Token() token2.Token
	Next()
	Prev()
	IsEmpty() bool
	SavePoint()
	RemoveSavePoint()
	ReturnToSavePoint()
}

type Parser struct {
	mover TokenMover
}

func New() *Parser {
	return &Parser{}
}

func (p *Parser) Parse(mover TokenMover) (ast2.Ast, error) {
	p.mover = mover

	file, err := p.parseFile()
	if err != nil {
		// Тут сразу отдаем ошибку, потому что без файла ast быть не может!
		return ast2.Ast{}, errors.Wrap(err, "parse file")
	}
	return ast2.NewAst(file), nil
}

func (p *Parser) parseFile() (ast2.File, error) {
	imports, err := p.parseImports()
	if err != nil {
		return ast2.File{}, errors.Wrap(err, "parse imports")
	}

	object, err := p.parseObject()
	if err != nil {
		return ast2.File{}, errors.Wrap(err, "parse object")
	}

	return ast2.NewFile(imports, object), nil
}

func (p *Parser) parseImports() ([]ast2.Import, error) {
	imports := make([]ast2.Import, 0)
	if p.mover.IsEmpty() {
		return imports, nil
	}

	for {
		imp, err := p.parseImport()

		switch {
		case err == nil:
			imports = append(imports, imp)
		case errors.Is(err, ErrTokenMismatch):
			return imports, nil
		default:
			return imports, errors.Wrap(err, "parse import")
		}
	}
}

func (p *Parser) parseImport() (ast2.Import, error) {
	if err := p.check(token2.Ident); err != nil {
		return ast2.Import{}, err
	}

	importName := p.mover.Token()

	p.mover.Next()

	if err := p.require(token2.Path); err != nil {
		return ast2.Import{}, errors.Wrapf(err, "import statement expected %s token", token2.Path.String())
	}

	importPath := p.mover.Token()

	p.mover.Next()

	return ast2.NewImport(
		ast2.NewIdent(importName.Value().String(), importName.Location()),
		ast2.NewPath(importPath.Value().String(), importPath.Location()),
	), nil
}

func (p *Parser) parseObject() (ast2.Object, error) {
	if err := p.check(token2.LBrace); err != nil {
		return ast2.Object{}, err
	}

	start := p.mover.Token().Location().Start()

	p.mover.Next()

	entries := make([]ast2.Entry, 0)

	for !p.match(token2.RBrace) {
		entry, err := p.parseEntry()
		if err != nil {
			return ast2.Object{}, err
		}
		entries = append(entries, entry)
	}

	end := p.mover.Token().Location().End()

	obj := ast2.NewObject(
		entries,
		types.NewLocation(start, end),
	)

	p.mover.Next()

	return obj, nil
}

func (p *Parser) parseSpread() (ast2.Spread, error) {
	p.mover.SavePoint()
	defer p.mover.RemoveSavePoint()

	v, err := p.parseVar()
	if err != nil {
		return ast2.Spread{}, err
	}

	if err = p.check(token2.Spread); err != nil {
		p.mover.ReturnToSavePoint()
		return ast2.Spread{}, err
	}

	s := ast2.NewSpread(
		v,
		types.NewLocation(
			v.Location().Start(),
			p.mover.Token().Location().End(),
		),
	)

	p.mover.Next()

	return s, nil
}

func (p *Parser) parseVar() (ast2.Var, error) {
	if err := p.check(token2.Ident); err != nil {
		return ast2.Var{}, err
	}

	idents := make([]ast2.Ident, 0)

	idents = append(
		idents,
		ast2.NewIdent(
			p.mover.Token().Value().String(),
			types.NewLocation(
				p.mover.Token().Location().Start(),
				p.mover.Token().Location().End()),
		),
	)

	p.mover.Next()

	if p.match(token2.Dot) {
		p.mover.Next()

		v, err := p.parseVar()
		switch {
		case err == nil:
			idents = append(idents, v.Path()...)
		case errors.Is(err, ErrTokenMismatch):
			break
		default:
			return ast2.Var{}, errors.Wrap(err, "parse var")
		}
	}

	return ast2.NewVar(
		idents,
	), nil
}

func (p *Parser) parseEntry() (ast2.Entry, error) {
	kv, err := p.parseKV()
	switch {
	case err == nil:
		return kv, nil
	case errors.Is(err, ErrTokenMismatch):
	default:
		return nil, errors.Wrap(err, "parse entry")
	}

	spread, err := p.parseSpread()
	if err != nil {
		return nil, errors.Wrap(err, "parse spread")
	}

	return spread, nil
}

func (p *Parser) parseKV() (ast2.KV, error) {
	p.mover.SavePoint()
	defer p.mover.RemoveSavePoint()

	if err := p.check(token2.Ident); err != nil {
		return ast2.KV{}, err
	}

	key := ast2.NewIdent(
		p.mover.Token().Value().String(),
		types.NewLocation(
			p.mover.Token().Location().Start(),
			p.mover.Token().Location().End(),
		),
	)

	p.mover.Next()

	if err := p.check(token2.Colon); err != nil {
		p.mover.ReturnToSavePoint()
		return ast2.KV{}, err
	}

	p.mover.Next()

	expr, err := p.parseExpression()
	switch {
	case err == nil:
	case errors.Is(err, ErrTokenMismatch):
		return ast2.KV{}, NewErrExpectedNode("expression")
	default:
		return ast2.KV{}, errors.Wrap(err, "parse expression")
	}

	return ast2.NewKV(
		key,
		expr,
	), nil
}

func (p *Parser) parseExpression() (expr ast2.Expression, err error) {
	if err = p.require(
		token2.Ident,
		token2.LBrace,
		token2.LBracket,
		token2.Dollar,
		token2.String,
		token2.Int,
		token2.Float,
		token2.Bool,
	); err != nil {
		return nil, err
	}

	switch p.mover.Token().Type() {
	case token2.Ident:
		expr, err = p.parseSpread()
		switch {
		case err == nil:
			return expr, nil
		case errors.Is(err, ErrTokenMismatch):
		default:
			return nil, err
		}

		expr, err = p.parseVar()
		if err != nil {
			return nil, err
		}

		return expr, nil
	case token2.Dollar:
		expr, err = p.parseEnv()
		if err != nil {
			return nil, err
		}

		return expr, nil
	case token2.LBrace:
		expr, err = p.parseObject()
		if err != nil {
			return nil, err
		}

		return expr, nil

	case token2.LBracket:
		expr, err = p.parseArray()
		if err != nil {
			return nil, err
		}

		return expr, nil

	case token2.String:
		expr, err = p.parseString()
		if err != nil {
			return nil, err
		}

		return expr, nil

	case token2.Float:
		expr, err = p.parseFloat()
		if err != nil {
			return nil, err
		}

		return expr, nil

	case token2.Int:
		expr, err = p.parseInt()
		if err != nil {
			return nil, err
		}

		return expr, nil

	case token2.Bool:
		expr, err = p.parseBool()
		if err != nil {
			return nil, err
		}

		return expr, nil
	default:
		return nil, NewErrUnexpectedToken()
	}
}

func (p *Parser) parseString() (ast2.String, error) {
	if err := p.require(token2.String); err != nil {
		return ast2.String{}, err
	}

	s := ast2.NewString(
		p.mover.Token().Value().String(),
		p.mover.Token().Location(),
	)

	p.mover.Next()

	return s, nil
}

func (p *Parser) parseBool() (ast2.Bool, error) {
	if err := p.require(token2.Bool); err != nil {
		return ast2.Bool{}, err
	}

	b, err := ast2.NewBool(
		p.mover.Token().Value().String(),
		p.mover.Token().Location(),
	)
	if err != nil {
		return ast2.Bool{}, errors.Wrap(err, "parse bool")
	}

	p.mover.Next()

	return b, nil
}

func (p *Parser) parseInt() (ast2.Int, error) {
	if err := p.require(token2.Int); err != nil {
		return ast2.Int{}, err
	}

	i, err := ast2.NewInt(
		p.mover.Token().Value().String(),
		p.mover.Token().Location(),
	)
	if err != nil {
		return ast2.Int{}, errors.Wrap(err, "parse int")
	}

	p.mover.Next()

	return i, nil
}

func (p *Parser) parseFloat() (ast2.Float, error) {
	if err := p.require(token2.Float); err != nil {
		return ast2.Float{}, err
	}

	f, err := ast2.NewFloat(
		p.mover.Token().Value().String(),
		p.mover.Token().Location(),
	)
	if err != nil {
		return ast2.Float{}, errors.Wrap(err, "parse float")
	}

	p.mover.Next()

	return f, nil
}

func (p *Parser) parseEnv() (ast2.Env, error) {
	if err := p.require(token2.Dollar); err != nil {
		return ast2.Env{}, err
	}

	dollarToken := p.mover.Token()

	p.mover.Next()

	if err := p.require(token2.Ident); err != nil {
		return ast2.Env{}, err
	}

	env := ast2.NewEnv(
		ast2.NewIdent(
			p.mover.Token().Value().String(),
			p.mover.Token().Location(),
		),
		types.NewLocation(
			dollarToken.Location().Start(),
			p.mover.Token().Location().End(),
		),
	)

	p.mover.Next()

	return env, nil
}

func (p *Parser) parseArray() (ast2.Array, error) {
	if err := p.require(token2.LBracket); err != nil {
		return ast2.Array{}, err
	}

	start := p.mover.Token().Location().Start()

	p.mover.Next()

	elements := make([]ast2.Expression, 0)

	for !p.match(token2.RBracket) {
		expr, err := p.parseExpression()
		switch {
		case err == nil:
		case errors.Is(err, ErrTokenMismatch):
			return ast2.Array{}, NewErrExpectedNode("expression")
		default:
			return ast2.Array{}, errors.Wrap(err, "parse expression")
		}

		elements = append(elements, expr)
	}

	array := ast2.NewArray(
		elements,
		types.NewLocation(
			start,
			p.mover.Token().Location().End(),
		),
	)

	p.mover.Next()

	return array, nil
}
