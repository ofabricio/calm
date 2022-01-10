package calm

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOr(t *testing.T) {

	tt := []struct {
		in string
		ok bool
	}{
		{"0", true},
		{"1", true},
		{"2", false},
	}

	for _, tc := range tt {

		c := New(tc.in)

		ok := c.Run(Or(S("0"), S("1")))

		assert.Equal(t, tc.ok, ok, tc.in)
	}
}

func TestAnd(t *testing.T) {

	tt := []struct {
		in string
		ok bool
	}{
		{"01", true},
		{"00", false},
		{"10", false},
		{"11", false},
	}

	for _, tc := range tt {

		c := New(tc.in)

		ok := c.Run(And(S("0"), S("1")))

		assert.Equal(t, tc.ok, ok, tc.in)
	}
}

func TestTrueFalseNot(t *testing.T) {

	tt := []struct {
		in string
		ok bool
		mf MatcherFunc
	}{
		{"a", true, S("a").True()},
		{"b", true, S("a").True()},
		{"a", false, S("a").False()},
		{"b", false, S("a").False()},
		{"a", false, S("a").Not()},
		{"b", true, S("a").Not()},
	}

	for _, tc := range tt {

		c := New(tc.in)

		ok := c.Run(tc.mf)

		assert.Equal(t, tc.ok, ok, tc.in)
	}
}
