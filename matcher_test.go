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
		ok := SR(&a).Run(c)

		assert.Equal(t, tc.ok, ok, tc.in)
	}
}

func TestS_F_R(t *testing.T) {

	tt := []struct {
		in string
		ok bool
		mf MatcherFunc
	}{
		{"a", true, S("a")},
		{"b", false, S("a")},
		{"abc", true, S("abc")},
		{"cba", false, S("abc")},
		{"a1c", true, R("\\w\\d\\w")},
		{"a1c", false, R("\\d\\w\\d")},
		{"b", true, F(unicode.IsLetter)},
		{"1", false, F(unicode.IsLetter)},
	}

	for _, tc := range tt {

		c := New(tc.in)

		ok := tc.mf.Run(c)

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

		ok := tc.mf.Run(c)

		assert.Equal(t, tc.ok, ok, tc.in)
	}
}
