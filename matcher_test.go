package calm

import (
	"testing"
	"unicode"

	"github.com/stretchr/testify/assert"
)

func TestSR(t *testing.T) {

	tt := []struct {
		in string
		ok bool
	}{
		{"a", true},
		{"b", false},
	}

	for _, tc := range tt {

		c := New(tc.in)

		a := "a"
		ok := c.Run(SR(&a))

		assert.Equal(t, tc.ok, ok, tc.in)
	}
}

func TestS_F(t *testing.T) {

	tt := []struct {
		in string
		ok bool
		mf MatcherFunc
	}{
		{"a", true, S("a")},
		{"b", false, S("a")},
		{"abc", true, S("abc")},
		{"cba", false, S("abc")},
		{"b", true, F(unicode.IsLetter)},
		{"1", false, F(unicode.IsLetter)},
	}

	for _, tc := range tt {

		c := New(tc.in)

		ok := c.Run(tc.mf)

		assert.Equal(t, tc.ok, ok, tc.in)
	}
}

func TestEq(t *testing.T) {

	tt := []struct {
		in string
		ok bool
		mf MatcherFunc
	}{
		{"", false, Eq("")},
		{"", false, Eq("a")},
		{"a", false, Eq("")},
		{"a", true, Eq("a")},
		{"b", false, Eq("a")},
		{"abc", true, Eq("abc")},
		{"cba", false, Eq("abc")},
	}

	for _, tc := range tt {

		c := New(tc.in)

		ok := c.Run(tc.mf)

		assert.Equal(t, tc.ok, ok, tc.in)
	}
}
