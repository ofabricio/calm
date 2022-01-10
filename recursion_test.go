package calm

import (
	"testing"
	"unicode"

	"github.com/stretchr/testify/assert"
)

func TestRecursive(t *testing.T) {

	tt := []struct {
		in string
		ok bool
	}{
		{"0", true},
		{"0+", false},
		{"0+1", true},
		{"0+1*", false},
		{"0+1*2", true},
		{"0+1*(2+3)*4", true},
	}

	for _, tc := range tt {

		c := New(tc.in)

		var term Wrap
		var expr Wrap

		value := F(unicode.IsNumber)
		factor := Or(And(S("("), &expr, S(")")), value)

		Or(And(factor, S("*"), &term).Not().Rewind().Not(), factor).Recursive(&term)
		Or(And(&term, S("+"), &expr).Not().Rewind().Not(), &term).Recursive(&expr)

		ok := c.Run(And(&expr, F(unicode.IsPrint).Not()))

		assert.Equal(t, tc.ok, ok, tc.in)
		assert.Equal(t, tc.in, c.Take(0, c.Here()), tc.in)
	}
}
