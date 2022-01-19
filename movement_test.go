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
		{"de", false, Next().Next().Next(), "de"},
	}

	for _, tc := range tt {

		c := New(tc.in)
		a := c.Mark()

		ok := tc.mf.Run(c)

		assert.Equal(t, tc.ok, ok, tc.in)
		assert.Equal(t, tc.ex, c.Token(a, c.Mark()).Text, tc.in)
	}
}

func TestUndo(t *testing.T) {

	tt := []struct {
		in string
		ok bool
		mf MatcherFunc
	}{
		{"a", true, S("a").Not().Undo().Not()},
	}

	for _, tc := range tt {

		c := New(tc.in)

		ok := And(tc.mf, S("a")).Run(c)

		assert.Equal(t, tc.ok, ok, tc.in)
	}
}
