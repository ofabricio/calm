package calm

import (
	"strings"
	"unicode/utf8"
)

func New(src string) *Code {
	return &Code{src: src, row: 1, col: 1}
}

// Equal tests if the string matches
// with the current position.
// It does not advances the position.
func (c *Code) Equal(s string) bool {
	if s == "" {
		return false
	}
	return strings.HasPrefix(c.Tail(), s)
}

// Match tests if the string matches
// with the current position and
// advances the position if true.
func (c *Code) Match(s string) bool {
	if c.Equal(s) {
		c.advance(s)
		return true
	}
	return false
}

// MatchF tests if f function matches
// with the current position and
// advances the position if true.
func (c *Code) MatchF(f func(rune) bool) bool {
	if r, _ := c.decodeRune(); f(r) {
		c.advanceC(r)
		return true
	}
	return false
}

// Mark marks the current position.
func (c *Code) Mark() Mark {
	return Mark{pos: c.pos, row: c.row, col: c.col}
}

// Back sends the position back to a mark.
func (c *Code) Back(m Mark) {
	c.pos = m.pos
	c.row = m.row
	c.col = m.col
}

// Token returns the token between ini and end.
func (c *Code) Token(ini, end Mark) Token {
	return Token{Text: c.src[ini.pos:end.pos], Pos: ini.pos, Row: ini.row, Col: ini.col}
}

// Next moves the position to the next character.
func (c *Code) Next() {
	r, _ := c.decodeRune()
	c.advanceC(r)
}

// Tail returns the content from the
// current position to the end.
func (c *Code) Tail() string {
	return c.src[c.pos:]
}

// More tells if there are more characters to scan.
func (c *Code) More() bool {
	return c.pos < len(c.src)
}

func (c *Code) advance(s string) {
	for _, r := range s {
		c.advanceC(r)
	}
}

func (c *Code) advanceC(r rune) {
	if c.pos < len(c.src) {
		c.rowcol(r)
		c.pos += utf8.RuneLen(r)
	}
}

func (c *Code) decodeRune() (rune, int) {
	return utf8.DecodeRuneInString(c.Tail())
}

func (c *Code) rowcol(r rune) {
	c.col++
	if r == '\n' {
		c.row++
		c.col = 1
	}
}

type Code struct {
	src string // Source code.
	pos int    // Position/Index/Offset/Cursor.
	row int    // Current line.
	col int    // Current column.
}

// Mark represents a mark in the code.
type Mark struct {
	pos int
	row int
	col int
}

// Token represents a token of the code.
type Token struct {
	Text string
	Pos  int
	Row  int
	Col  int
}
