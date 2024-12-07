package types

type Location struct {
	start Position
	end   Position
}

func (l Location) SetStart(start Position) Location {
	l.start = start
	return l
}

func (l Location) SetEnd(end Position) Location {
	l.end = end
	return l
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

func NewInitialLocation() Location {
	return Location{start: NewInitialPosition(), end: NewInitialPosition()}
}

type Position struct {
	line   uint
	column uint
	pos    uint
}

func NewPosition(line uint, column uint, pos uint) Position {
	return Position{line: line, column: column, pos: pos}
}

func NewInitialPosition() Position {
	return Position{
		line:   1,
		column: 0,
		pos:    0,
	}
}

func (p Position) Pos() uint {
	return p.pos
}

func (p Position) IncrPos() Position {
	p.pos++
	return p
}

func (p Position) AddPos(amount uint) Position {
	p.pos += amount
	return p
}

func (p Position) Line() uint {
	return p.line
}

func (p Position) IncrLine() Position {
	p.line++
	return p
}

func (p Position) Column() uint {
	return p.column
}

func (p Position) IncrColumn() Position {
	p.column++
	return p
}

func (p Position) AddColumn(amount uint) Position {
	p.column += amount
	return p
}

func (p Position) ResetColumn() Position {
	p.column = 0
	return p
}
