package calm

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMatcherFunc_String(t *testing.T) {

	tt := []struct {
		in string
		ok bool
		mf MatcherFunc
		ex string
	}{
		// Valid.
		{`""`, true, String(`"`), `""`},
		{`"\""`, true, String(`"`), `"\""`},
		{`''`, true, String(`'`), `''`},
		{`"hello world"`, true, String(`"`), `"hello world"`},
		{`'hello world'`, true, String(`'`), `'hello world'`},
		{`"hello\"world"`, true, String(`"`), `"hello\"world"`},
		{`"hello\world"`, true, String(`"`), `"hello\world"`},
		{`"hello\"world\""`, true, String(`"`), `"hello\"world\""`},
		{`"\"hello world\""`, true, String(`"`), `"\"hello world\""`},
		{`"hello"world"`, true, String(`"`), `"hello"`},
		// Invalid.
		{`'"`, false, String(`'`), ``},
		{`"\"`, false, String(`"`), ``},
		{`\"hello\"`, false, String(`"`), ``},
		{`"hello\"`, false, String(`"`), ``},
		{"\"hello\nworld\"", false, String(`"`), ``},
		{`"a"`, false, String(`'`), ``},
		{`"b"`, false, String(`'`), ``},
		{`"c'`, false, String(`'`), ``},
		{`"d'`, false, String(`"`), ``},
	}

	for _, tc := range tt {

		c := New(tc.in)

		var tk string
		ok := c.Run(tc.mf.On(Grab(&tk)))

		assert.Equal(t, tc.ok, ok, tc.in)
		assert.Equal(t, tc.ex, tk, tc.in)
	}
}

func ExampleMatcherFunc_Debug() {

	c := New("a")

	c.Run(S("a").Debug())

	// Output:
	// [debug] Match: true Token: 'a' Pos: 0 End: 1
}
