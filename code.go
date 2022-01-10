package calm

import (
	"strings"
	"unicode/utf8"
)

func New(src string) *Code {
	return &Code{core: core{src: src}, rowcol: rowcol{1, 1}}
}

// Run implements the Matcher interface.
func (c *Code) Run(m Matcher) bool {
	return m.Run(c)
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
		c.Next()
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
	return Token{Text: c.Take(ini.pos, end.pos), Pos: ini.pos, Row: ini.row, Col: ini.col}
}

// Next moves the position to the next character.
func (c *Code) Next() {
	r, _ := c.decodeRune()
	c.advanceC(r)
}

func (c *Code) advance(s string) {
	for _, r := range s {
		c.advanceC(r)
	}
}

func (c *Code) advanceC(r rune) {
	c.incCol()
	if r == '\n' {
		c.incRow()
	}
	c.move(c.Here() + utf8.RuneLen(r))
}

func (c *Code) decodeRune() (rune, int) {
	return utf8.DecodeRuneInString(c.Tail())
}

type Code struct {
	core
	rowcol
}
