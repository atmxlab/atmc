package types

type Location struct {
	start Position
	end   Position
}

func (l Location) Start() Position {
	return l.start
}

func (l Location) End() Position {
	return l.end
}

func NewLocation(start Position, end Position) Location {
	return Location{start: start, end: end}
}

type Position struct {
	line   uint
	column uint
	pos    uint
}

func NewPosition(line uint, column uint, pos uint) Position {
	return Position{line: line, column: column, pos: pos}
}

func NewInitialPosition() *Position {
	return &Position{
		line:   1,
		column: 0,
		pos:    0,
	}
}

func (p *Position) Clone() *Position {
	return NewPosition(p.line, p.column, p.pos)
}

func (p *Position) Pos() uint {
	return p.pos
}

func (p *Position) IncrPos() {
	p.pos++
}

func (p *Position) AddPos(amount uint) {
	p.pos += amount
}

func (p *Position) Line() uint {
	return p.line
}

func (p *Position) IncrLine() {
	p.line++
}

func (p *Position) Column() uint {
	return p.column
}

func (p *Position) IncrColumn() {
	p.column++
}

func (p *Position) AddColumn(amount uint) {
	p.column += amount
}

func (p *Position) ResetColumn() {
	p.column = 0
}
