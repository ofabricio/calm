package calm

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCoreMoveHereTailMore(t *testing.T) {

	s := New("abc")

	s.move(0)
	assert.Equal(t, 0, s.Here())
	assert.Equal(t, "abc", s.Tail())
	assert.True(t, s.More())

	s.move(1)
	assert.Equal(t, 1, s.Here())
	assert.Equal(t, "bc", s.Tail())
	assert.True(t, s.More())

	s.move(2)
	assert.Equal(t, 2, s.Here())
	assert.Equal(t, "c", s.Tail())
	assert.True(t, s.More())

	s.move(3)
	assert.Equal(t, 3, s.Here())
	assert.Equal(t, "", s.Tail())
	assert.False(t, s.More())

	// Test overflow.
	s.move(4)
	assert.Equal(t, 3, s.Here())
	assert.Equal(t, "", s.Tail())
	assert.False(t, s.More())
}

func TestCoreTake(t *testing.T) {

	s := New("abc")

	assert.Equal(t, "a", s.Take(0, 1))
	assert.Equal(t, "b", s.Take(1, 2))
	assert.Equal(t, "c", s.Take(2, 3))
	assert.Equal(t, "abc", s.Take(0, 3))
	assert.Equal(t, "ab", s.Take(0, 2))
	assert.Equal(t, "bc", s.Take(1, 3))
}
