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
		ok := tc.mf.On(Grab(&tk)).Run(c)

		assert.Equal(t, tc.ok, ok, tc.in)
		assert.Equal(t, tc.ex, tk, tc.in)
	}
}

func ExampleMatcherFunc_Debug() {

	c := New("a")

	S("a").Debug().Run(c)

	// Output:
	// [debug] Match: true  Token: 'a' Pos: 0 Row: 1 Col: 1
}

func TestJson(t *testing.T) {

	tt := []struct {
		in string
		ok bool
	}{
		// Valid.
		{`{}`, true},
		{`{ }`, true},
		{`[]`, true},
		{`[ ]`, true},
		{`[1]`, true},
		{`["1"]`, true},
		{`[1.0]`, true},
		{`[1.0e+2]`, true},
		{`[1.2e-3]`, true},
		{`[1e-3]`, true},
		{`[1E-3]`, true},
		{`[ 1, 2]`, true},
		{`[ 1, "a", {}]`, true},
		{`[ 1, "a", { "b": 2 }]`, true},
		{`{ "a": "b" }`, true},
		{`{ "a": 1 }`, true},
		{`{ "a": true }`, true},
		{`{ "a": false }`, true},
		{`{ "a": null }`, true},
		{`{ "a": {} }`, true},
		{`{ "a": [] }`, true},
		{`{ "a": [{}] }`, true},
		{`{ "a": [2] }`, true},
		{`{ "a": [2], "b": { "c": 3 } }`, true},
		{`{ "a": { "b": { "c": "d" } } }`, true},
		{`{ "a": { "b": { "c": [] } } }`, true},
		// Invalid.
		{`[`, false},
		{`{`, false},
		{`{ "a": 1, "b": { "c": 2 }`, false},
		{`[1, 2`, false},
	}

	for _, tc := range tt {

		c := New(tc.in)

		ok := Json().Run(c)

		assert.Equal(t, tc.ok, ok, tc.in)
	}
}
