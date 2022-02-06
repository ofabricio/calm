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

		ok := Or(S("0"), S("1")).Run(c)

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

		ok := And(S("0"), S("1")).Run(c)

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
		{"", true, True()},
		{"a", false, S("a").False()},
		{"b", false, S("a").False()},
		{"", false, False()},
		{"a", false, S("a").Not()},
		{"b", true, S("a").Not()},
	}

	for _, tc := range tt {

		c := New(tc.in)

		ok := tc.mf.Run(c)

		assert.Equal(t, tc.ok, ok, tc.in)
	}
}

func TestIf(t *testing.T) {

	tt := []struct {
		in string
		ok bool
		mf MatcherFunc
	}{
		{".", true, If(S("."), True(), False())},
		{"@", false, If(S("."), True(), False())},
	}

	for _, tc := range tt {

		c := New(tc.in)

		ok := tc.mf.Run(c)

		assert.Equal(t, tc.ok, ok, tc.in)
	}
}
