package calm

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNext(t *testing.T) {

	tt := []struct {
		in string
		ok bool
		mf MatcherFunc
		ex string
	}{
		{"", false, Eq("a").Next(), ""},
		{"b", true, Eq("b").Next(), "b"},
		{"c", false, Eq("b").Next(), ""},
		{"bc", true, Eq("b").Next(), "b"},
		{"de", true, Next(), "d"},
		{"de", true, And(Next(), Next()), "de"},
		{"de", true, Next().Next(), "de"},
	}

	for _, tc := range tt {

		c := New(tc.in)

		ok := c.Run(tc.mf)

		assert.Equal(t, tc.ok, ok, tc.in)
		assert.Equal(t, tc.ex, c.Take(0, c.Here()), tc.in)
	}
}

func TestRewind(t *testing.T) {

	tt := []struct {
		in string
		ok bool
		mf MatcherFunc
	}{
		{"a", true, S("a").Rewind()},
		{"bb", false, And(S("b"), S("a")).Not().Rewind().Not()},
	}

	for _, tc := range tt {

		c := New(tc.in)

		ok := c.Run(tc.mf)

		assert.Equal(t, tc.ok, ok, tc.in)
		assert.Equal(t, 0, c.Here(), tc.in)
	}
}
