package calm

import (
	"testing"
	"unicode"

	"github.com/stretchr/testify/assert"
)

func TestRepetition(t *testing.T) {

	tt := []struct {
		in string
		ok bool
		mf MatcherFunc
		ex string
	}{
		// ZeroToMany
		{"", false, S("a").ZeroToMany(), ""},
		{"b", true, S("a").ZeroToMany(), ""},
		{"a", true, S("a").ZeroToMany(), "a"},
		{"aa", true, S("a").ZeroToMany(), "aa"},
		{"aaa", true, S("a").ZeroToMany(), "aaa"},
		{"aaab", true, S("a").ZeroToMany(), "aaa"},
		// OneToMany
		{"", false, S("a").OneToMany(), ""},
		{"b", false, S("a").OneToMany(), ""},
		{"a", true, S("a").OneToMany(), "a"},
		{"aa", true, S("a").OneToMany(), "aa"},
		{"aaa", true, S("a").OneToMany(), "aaa"},
		{"aaab", true, S("a").OneToMany(), "aaa"},
		// ZeroToOne
		{"", false, S("a").ZeroToOne(), ""},
		{"b", true, S("a").ZeroToOne(), ""},
		{"a", true, S("a").ZeroToOne(), "a"},
		{"aa", true, S("a").ZeroToOne(), "a"},
		// Min
		{"", false, S("a").Min(0), ""},
		{"a", true, S("a").Min(0), "a"},
		{"aa", true, S("a").Min(0), "aa"},
		{"", false, S("a").Min(1), ""},
		{"a", true, S("a").Min(1), "a"},
		{"aa", true, S("a").Min(1), "aa"},
		{"", false, S("a").Min(2), ""},
		{"a", false, S("a").Min(2), "a"}, // This should fail if we put Rewind in Min.
		{"aa", true, S("a").Min(2), "aa"},
		{"aaa", true, S("a").Min(2), "aaa"},
		// Until
		{"", false, Until(Eq("a")), ""},
		{"a", false, Until(Eq("a")), ""},
		{"x", true, Until(Eq("a")), "x"},
		{"xa", true, Until(Eq("a")), "x"},
		{"xxa", true, Until(Eq("a")), "xx"},
		// While
		{"", false, While(Eq("0")), ""},
		{"a", false, While(Eq("0")), ""},
		{"0", true, While(Eq("0")), "0"},
		{"01", true, While(Eq("0"), Eq("1")), "01"},
		{"0110201", true, While(Eq("0"), Eq("1")), "0110"},
	}

	for _, tc := range tt {

		c := New(tc.in)

		ok := c.Run(tc.mf)

		assert.Equal(t, tc.ok, ok, tc.in)
		assert.Equal(t, tc.ex, c.Take(0, c.Here()), tc.in)
	}
}

func Test_End_Of_Source_Code(t *testing.T) {

	tt := []struct {
		in string
		mf MatcherFunc
		a  bool
		b  bool
	}{
		{"a", S("a"), true, false},
		{"b", F(unicode.IsPrint), true, false},
		//
		{"", S("c").ZeroToMany(), false, false},
		{"c", S("c").ZeroToMany(), true, false},
		{"cc", S("c").ZeroToMany(), true, false},
		//
		{"", S("d").OneToMany(), false, false},
		{"d", S("d").OneToMany(), true, false},
		{"dd", S("d").OneToMany(), true, false},
		//
		{"", S("e").ZeroToOne(), false, false},
		{"e", S("e").ZeroToOne(), true, false},
		{"f", S("e").ZeroToOne(), true, true},
		//
		{"", Eq("g"), false, false},
		{"g", Eq("g"), true, true},
		// Special case equivalent to "Until".
		{"", Eq("x").Not().Next().OneToMany(), false, false},
		{"h", S("x").Not().Next().OneToMany(), true, false},
		//
		{"", Until(Eq("x")), false, false},
		{"h", Until(Eq("x")), true, false},
		//
		{"", While(Eq("x")), false, false},
		{"x", While(Eq("x")), true, false},
		{"xx", While(Eq("x")), true, false},
		//
		{"", Eq("i").Next(), false, false},
		{"i", Eq("i").Next(), true, false},
	}

	for _, tc := range tt {
		c := New(tc.in)

		a := c.Run(tc.mf)
		b := c.Run(tc.mf)

		assert.Equal(t, tc.a, a, tc.in)
		assert.Equal(t, tc.b, b, tc.in)
	}
}
