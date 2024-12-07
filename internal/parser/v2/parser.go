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

	start := types.NewPosition(0, 0, 0)

	if len(imports) > 0 {
		start = imports[0].Location().Start()
	}

	object, err := p.parseObject()
	if err != nil {
		return ast.File{}, errors.Wrap(err, "parse object")
	}

	end := object.Location().End()

	return ast.NewFile(imports, object, types.NewLocation(start, end)), nil
}

func (p *Parser) parseImports() ([]ast.Import, error) {
	imports := make([]ast.Import, 0)

	imp, err := p.parseImport()
	for {
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
		// TODO: тут нужно решить проблему локации токенов. У них тоже должна быть локация.
		ast.NewIdent(importName.Value().String(), importName.Location()),
		ast.NewPath(importPath.Value().String()),
		types.NewLocation(importName.Location().Start(), importPath.Location().End()),
	), nil
}

func (p *Parser) parseObject() (ast.Object, error) {
	return ast.Object{}, nil
}
