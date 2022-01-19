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

		term, setTerm := Recursive()
		expr, setExpr := Recursive()

		value := F(unicode.IsNumber)
		factor := Or(And(S("("), expr, S(")")), value)
		setTerm(Or(And(factor, S("*"), term).Rewind(), factor))
		setExpr(Or(And(term, S("+"), expr).Rewind(), term))

		ok := c.Run(And(expr, Next().Not()))

		assert.Equal(t, tc.ok, ok, tc.in)
	}
}
