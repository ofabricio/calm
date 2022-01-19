package calm

import (
	"testing"
	"unicode"

	"github.com/stretchr/testify/assert"
)

func TestEqual(t *testing.T) {

	c := New("abc")

	assert.False(t, c.Equal(""))
	assert.True(t, c.Equal("a"))
	assert.True(t, c.Equal("ab"))
	assert.True(t, c.Equal("abc"))
	assert.False(t, c.Equal("abcd"))
	assert.False(t, c.Equal("x"))
}

func TestMatch(t *testing.T) {

	c := New("abcd")

	assert.False(t, c.Match(""))
	assert.True(t, c.Match("a"))
	assert.True(t, c.Match("b"))
	assert.False(t, c.Match("x"))
	assert.True(t, c.Match("cd"))
	assert.False(t, c.Match(""))
}

func TestMatchF(t *testing.T) {

	c := New("ab")

	assert.False(t, c.MatchF(unicode.IsNumber))
	assert.True(t, c.MatchF(unicode.IsLetter))
	assert.False(t, c.MatchF(unicode.IsNumber))
	assert.True(t, c.MatchF(unicode.IsLetter))
	assert.False(t, c.MatchF(unicode.IsLetter))
}
func TestNextTailMore(t *testing.T) {

	s := New("a世c")

	assert.Equal(t, "a世c", s.Tail())
	assert.True(t, s.More())

	s.Next()
	assert.Equal(t, "世c", s.Tail())
	assert.True(t, s.More())

	s.Next()
	assert.Equal(t, "c", s.Tail())
	assert.True(t, s.More())

	s.Next()
	assert.Equal(t, "", s.Tail())
	assert.False(t, s.More())

	// Test overflow.
	s.Next()
	assert.Equal(t, "", s.Tail())
	assert.False(t, s.More())
}

func TestMarkBackToken(t *testing.T) {

	c := New("a")

	ini := c.Mark()

	c.Next()

	end := c.Mark()

	assert.Equal(t, "a", c.Token(ini, end).Text)
}
