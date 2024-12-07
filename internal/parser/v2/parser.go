package v2

import (
	"github.com/atmxlab/atmcfg/internal/parser/ast"
	"github.com/atmxlab/atmcfg/internal/types"
	"github.com/atmxlab/atmcfg/internal/types/token"
	"github.com/atmxlab/atmcfg/pkg/errors"
)

type TokenMover interface {
	Token() token.Token
	Next() TokenMover
	Prev() TokenMover
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

	spreads := make([]ast.Spread, 0)
	entries := make([]ast.Entry, 0)

	for {
		entry, err := p.parseEntry()
		switch {
		case err == nil:
			entries = append(entries, entry)
			continue
		case errors.Is(err, ErrTokenMismatch):
		default:
			return ast.Object{}, errors.Wrap(err, "parse entry")
		}

		spread, err := p.parseSpread()
		switch {
		case err == nil:
			spreads = append(spreads, spread)
			continue
		case errors.Is(err, ErrTokenMismatch):
		default:
			return ast.Object{}, errors.Wrap(err, "parse spread")
		}

		break
	}

	p.mover.Next()

	if err := p.require(token.RBrace); err != nil {
		return ast.Object{}, errors.Wrap(err, "object expected rbrace in the object end position")
	}

	end := p.mover.Token().Location().End()

	return ast.NewObject(
		spreads,
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
	}

	return ast.NewVar(
		idents,
		types.NewLocation(
			idents[0].Location().Start(),
			idents[len(idents)-1].Location().End(),
		),
	), nil
}
