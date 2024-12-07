package v2

import (
	"github.com/atmxlab/atmcfg/internal/parser/ast"
	"github.com/atmxlab/atmcfg/internal/types"
	"github.com/atmxlab/atmcfg/internal/types/token"
	"github.com/atmxlab/atmcfg/pkg/errors"
)

type TokenMover interface {
	Token() token.Token
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

func NewParser(mover TokenMover) *Parser {
	return &Parser{mover: mover}
}

func (p *Parser) Parse() (ast.Ast, error) {
	file, err := p.parseFile()
	if err != nil {
		// Тут сразу отдаем ошибку, потому что без файла ast быть не может!
		return ast.Ast{}, errors.Wrap(err, "parse file")
	}
	return ast.NewAst(file), nil
}

func (p *Parser) parseFile() (ast.File, error) {
	imports, err := p.parseImports()
	if err != nil {
		return ast.File{}, errors.Wrap(err, "parse imports")
	}

	object, err := p.parseObject()
	if err != nil {
		return ast.File{}, errors.Wrap(err, "file expected object node")
	}

	start := object.Location().Start()

	if len(imports) > 0 {
		start = imports[0].Location().Start()
	}

	end := object.Location().End()

	return ast.NewFile(imports, object, types.NewLocation(start, end)), nil
}

func (p *Parser) parseImports() ([]ast.Import, error) {
	imports := make([]ast.Import, 0)

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

func (p *Parser) parseImport() (ast.Import, error) {
	if err := p.check(token.Ident); err != nil {
		return ast.Import{}, err
	}

	importName := p.mover.Token()

	p.mover.Next()

	if err := p.require(token.Path); err != nil {
		return ast.Import{}, errors.Wrapf(err, "import statement expected %s token", token.Path.String())
	}

	importPath := p.mover.Token()

	p.mover.Next()

	return ast.NewImport(
		ast.NewIdent(importName.Value().String(), importName.Location()),
		ast.NewPath(importPath.Value().String()),
		types.NewLocation(importName.Location().Start(), importPath.Location().End()),
	), nil
}

func (p *Parser) parseObject() (ast.Object, error) {
	if err := p.check(token.LBrace); err != nil {
		return ast.Object{}, err
	}

	start := p.mover.Token().Location().Start()

	p.mover.Next()

	entries := make([]ast.Entry, 0)

	for !p.match(token.RBrace) {
		entry, err := p.parseEntry()
		if err != nil {
			return ast.Object{}, err
		}
		entries = append(entries, entry)

		p.mover.Next()
	}

	end := p.mover.Token().Location().End()

	return ast.NewObject(
		entries,
		types.NewLocation(start, end),
	), nil
}

func (p *Parser) parseSpread() (ast.Spread, error) {
	p.mover.SavePoint()
	defer p.mover.RemoveSavePoint()

	v, err := p.parseVar()
	if err != nil {
		return ast.Spread{}, err
	}

	p.mover.Next()

	if err = p.check(token.Spread); err != nil {
		p.mover.ReturnToSavePoint()
		return ast.Spread{}, err
	}

	return ast.NewSpread(
		v,
		types.NewLocation(
			v.Location().Start(),
			p.mover.Token().Location().End(),
		),
	), nil
}

func (p *Parser) parseVar() (ast.Var, error) {
	if err := p.check(token.Ident); err != nil {
		return ast.Var{}, err
	}

	p.mover.SavePoint()
	defer p.mover.RemoveSavePoint()

	idents := make([]ast.Ident, 0)

	idents = append(
		idents,
		ast.NewIdent(
			p.mover.Token().Value().String(),
			types.NewLocation(
				p.mover.Token().Location().Start(),
				p.mover.Token().Location().End()),
		),
	)

	p.mover.Next()

	if p.match(token.Dot) {
		p.mover.Next()

		v, err := p.parseVar()
		switch {
		case err == nil:
			idents = append(idents, v.Path()...)
		case errors.Is(err, ErrTokenMismatch):
			break
		default:
			return ast.Var{}, errors.Wrap(err, "parse var")
		}
	} else {
		p.mover.ReturnToSavePoint()
	}

	return ast.NewVar(
		idents,
		types.NewLocation(
			idents[0].Location().Start(),
			idents[len(idents)-1].Location().End(),
		),
	), nil
}

func (p *Parser) parseEntry() (ast.Entry, error) {
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

func (p *Parser) parseKV() (ast.KV, error) {
	p.mover.SavePoint()
	defer p.mover.RemoveSavePoint()

	if err := p.check(token.Ident); err != nil {
		return ast.KV{}, err
	}

	key := ast.NewIdent(
		p.mover.Token().Value().String(),
		types.NewLocation(
			p.mover.Token().Location().Start(),
			p.mover.Token().Location().End(),
		),
	)

	p.mover.Next()

	if err := p.check(token.Colon); err != nil {
		p.mover.ReturnToSavePoint()
		return ast.KV{}, err
	}

	p.mover.Next()

	expr, err := p.parseExpression()
	switch {
	case err == nil:
	case errors.Is(err, ErrTokenMismatch):
		return ast.KV{}, NewErrExpectedNode("expression")
	default:
		return ast.KV{}, errors.Wrap(err, "parse expression")
	}

	return ast.NewKV(
		key,
		expr,
		types.NewLocation(
			key.Location().Start(),
			expr.Location().End(),
		),
	), nil
}

func (p *Parser) parseExpression() (expr ast.Expression, err error) {
	if err = p.require(
		token.Ident,
		token.LBrace,
		token.RBracket,
		token.Dollar,
		token.String,
		token.Int,
		token.Float,
		token.Bool,
	); err != nil {
		return nil, err
	}

	defer p.mover.Next()

	switch p.mover.Token().Type() {
	case token.Ident:
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
	case token.Dollar:
		expr, err = p.parseEnv()
		if err != nil {
			return nil, err
		}

		return expr, nil
	case token.LBrace:
		expr, err = p.parseObject()
		if err != nil {
			return nil, err
		}

		return expr, nil

	case token.LBracket:
		expr, err = p.parseArray()
		if err != nil {
			return nil, err
		}

		return expr, nil

	case token.String:
		expr, err = p.parseString()
		if err != nil {
			return nil, err
		}

		return expr, nil

	case token.Float:
		expr, err = p.parseFloat()
		if err != nil {
			return nil, err
		}

		return expr, nil

	case token.Int:
		expr, err = p.parseInt()
		if err != nil {
			return nil, err
		}

		return expr, nil

	case token.Bool:
		expr, err = p.parseBool()
		if err != nil {
			return nil, err
		}

		return expr, nil
	default:
		return nil, NewErrUnexpectedToken()
	}
}

func (p *Parser) parseString() (ast.String, error) {
	if err := p.require(token.String); err != nil {
		return ast.String{}, err
	}

	return ast.NewString(
		p.mover.Token().Value().String(),
		p.mover.Token().Location(),
	), nil
}

func (p *Parser) parseBool() (ast.Bool, error) {
	if err := p.require(token.Bool); err != nil {
		return ast.Bool{}, err
	}

	b, err := ast.NewBool(
		p.mover.Token().Value().String(),
		p.mover.Token().Location(),
	)
	if err != nil {
		return ast.Bool{}, errors.Wrap(err, "parse bool")
	}

	return b, nil
}

func (p *Parser) parseInt() (ast.Int, error) {
	if err := p.require(token.Int); err != nil {
		return ast.Int{}, err
	}

	i, err := ast.NewInt(
		p.mover.Token().Value().String(),
		p.mover.Token().Location(),
	)
	if err != nil {
		return ast.Int{}, errors.Wrap(err, "parse int")
	}

	return i, nil
}

func (p *Parser) parseFloat() (ast.Float, error) {
	if err := p.require(token.Float); err != nil {
		return ast.Float{}, err
	}

	f, err := ast.NewFloat(
		p.mover.Token().Value().String(),
		p.mover.Token().Location(),
	)
	if err != nil {
		return ast.Float{}, errors.Wrap(err, "parse float")
	}

	return f, nil
}

func (p *Parser) parseEnv() (ast.Env, error) {
	if err := p.require(token.Dollar); err != nil {
		return ast.Env{}, err
	}

	dollarToken := p.mover.Token()

	p.mover.Next()

	if err := p.require(token.Ident); err != nil {
		return ast.Env{}, err
	}

	return ast.NewEnv(
		ast.NewIdent(
			p.mover.Token().Value().String(),
			p.mover.Token().Location(),
		),
		types.NewLocation(
			dollarToken.Location().Start(),
			p.mover.Token().Location().End(),
		),
	), nil
}

func (p *Parser) parseArray() (ast.Array, error) {
	if err := p.require(token.LBracket); err != nil {
		return ast.Array{}, err
	}

	start := p.mover.Token().Location().Start()

	p.mover.Next()

	elements := make([]ast.Expression, 0)

	for !p.match(token.RBracket) {
		expr, err := p.parseExpression()
		switch {
		case err == nil:
		case errors.Is(err, ErrTokenMismatch):
			return ast.Array{}, NewErrExpectedNode("expression")
		default:
			return ast.Array{}, errors.Wrap(err, "parse expression")
		}

		elements = append(elements, expr)
	}

	return ast.NewArray(
		elements,
		types.NewLocation(
			start,
			p.mover.Token().Location().End(),
		),
	), nil
}
