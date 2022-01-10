package calm

import (
	"testing"
	"unicode"

	"github.com/stretchr/testify/assert"
)

func TestCodeEqual(t *testing.T) {

	c := New("abc")

	assert.False(t, c.Equal(""))
	assert.True(t, c.Equal("a"))
	assert.True(t, c.Equal("ab"))
	assert.True(t, c.Equal("abc"))
	assert.False(t, c.Equal("abcd"))
	assert.False(t, c.Equal("x"))
}

func TestCodeMatch(t *testing.T) {

	c := New("abcd")

	assert.False(t, c.Match(""))
	assert.True(t, c.Match("a"))
	assert.True(t, c.Match("b"))
	assert.False(t, c.Match("x"))
	assert.True(t, c.Match("cd"))
	assert.False(t, c.Match(""))
}

func TestCodeMatchF(t *testing.T) {

	c := New("ab")

	assert.False(t, c.MatchF(unicode.IsNumber))
	assert.True(t, c.MatchF(unicode.IsLetter))
	assert.False(t, c.MatchF(unicode.IsNumber))
	assert.True(t, c.MatchF(unicode.IsLetter))
	assert.False(t, c.MatchF(unicode.IsLetter))
}

func TestCodeMark(t *testing.T) {

	c := New("a")

	a := c.Mark()
	c.Match("a")
	b := c.Mark()

	assert.Equal(t, Mark{pos: 0, row: 1, col: 1}, a)
	assert.Equal(t, Mark{pos: 1, row: 1, col: 2}, b)
}

func TestCodeBack(t *testing.T) {

	c := New("a")

	a := c.Mark()

	assert.True(t, c.Match("a"))

	c.Back(a)

	assert.True(t, c.Match("a"))
}

func TestCodeToken(t *testing.T) {

	c := New("a")

	ini := c.Mark()

	c.Match("a")

	end := c.Mark()

	assert.Equal(t, "a", c.Token(ini, end).Text)
}

func TestCodeNext(t *testing.T) {

	c := New("a世c")

	aa := c.Mark()
	c.Next()
	bb := c.Mark()
	c.Next()
	cc := c.Mark()
	c.Next()
	dd := c.Mark()

	assert.Equal(t, "a", c.Token(aa, bb).Text)
	assert.Equal(t, "世", c.Token(bb, cc).Text)
	assert.Equal(t, "c", c.Token(cc, dd).Text)
}
