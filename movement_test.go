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

		ok := c.Run(tc.mf)

		assert.Equal(t, tc.ok, ok, tc.in)
		assert.Equal(t, tc.ex, c.Token(a, c.Mark()).Text, tc.in)
	}
}

func TestRewind_and_Undo(t *testing.T) {

	tt := []struct {
		in string
		ok bool
		mf MatcherFunc
	}{
		{"a", true, S("a").Not().Rewind().Not()},
		{"a", true, S("a").Undo()},
	}

	for _, tc := range tt {

		c := New(tc.in)

		ok := c.Run(And(tc.mf, S("a")))

		assert.Equal(t, tc.ok, ok, tc.in)
	}
}
