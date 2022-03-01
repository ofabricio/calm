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

		i := c.Mark()
		var tk string
		ok := tc.mf.On(Grab(&tk)).Run(c)
		e := c.Mark()

		assert.Equal(t, tc.ok, ok, tc.in)
		assert.Equal(t, tc.ex, tk, tc.in)
		assert.Equal(t, tc.ex, c.Token(i, e).Text, tc.in)
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

func TestMatcherFunc_Number(t *testing.T) {

	tt := []struct {
		in string
		ok bool
		mf MatcherFunc
		ex string
	}{
		{"", false, Number(), ""},
		{"1", true, Number(), "1"},
		{"1.", false, Number(), ""},
		{"0.0", true, Number(), "0.0"},
		{"1.234", true, Number(), "1.234"},
		{"1.2e", false, Number(), ""},
		{"1.2E10", true, Number(), "1.2E10"},
		{"1.2e2", true, Number(), "1.2e2"},
		{"1.5e-3", true, Number(), "1.5e-3"},
		{"1.77e+4", true, Number(), "1.77e+4"},
		{"-20.45", true, Number(), "-20.45"},
	}

	for _, tc := range tt {

		c := New(tc.in)

		var tk string
		ok := tc.mf.On(Grab(&tk)).Run(c)

		assert.Equal(t, tc.ok, ok, tc.in)
		assert.Equal(t, tc.ex, tk, tc.in)
	}
}

func TestTag(t *testing.T) {

	tt := []struct {
		in string
		ok bool
		mf MatcherFunc
		ex string
	}{
		{"", false, Tag("<", ">"), ""},
		{"<", false, Tag("<", ">"), ""},
		{">", false, Tag("<", ">"), ""},
		{"[]", false, Tag("<", ">"), ""},
		{"<a", false, Tag("<", ">"), ""},
		{"<>", true, Tag("<", ">"), "<>"},
		{"<a>", true, Tag("<", ">"), "<a>"},
		{"<a b>", true, Tag("<", ">"), "<a b>"},
		{"{}", true, Tag("{", "}"), "{}"},
		{"{ab}", true, Tag("{", "}"), "{ab}"},
		{"{ab}x", true, Tag("{", "}"), "{ab}"},
		{"{a{b}c}", true, Tag("{", "}"), "{a{b}c}"},
		{"{a}{b}", true, Tag("{", "}"), "{a}"},
		{"{1}{2}", true, Tag("{", "}").OneToMany(), "{1}{2}"},
		{"{1}{2}{3}", true, Tag("{", "}").OneToMany(), "{1}{2}{3}"},
		{"{1}{{23}}", true, Tag("{", "}").OneToMany(), "{1}{{23}}"},
		{"{{{}}}", true, Tag("{", "}").OneToMany(), "{{{}}}"},
		{"{{a}}", true, Tag("{{", "}}").OneToMany(), "{{a}}"},
		{"{{{a{{{b}}}cc}}}", true, Tag("{{{", "}}}").OneToMany(), "{{{a{{{b}}}cc}}}"},
		{"{{{a{{{b}}}cc}}}", true, Tag("{{", "}}").OneToMany(), "{{{a{{{b}}}cc}}"},
		{"{ab\ncd}", true, Tag("{", "}").OneToMany(), "{ab\ncd}"},
		{"{a\nb{\nc}d}", true, Tag("{", "}").OneToMany(), "{a\nb{\nc}d}"},
	}

	for _, tc := range tt {

		c := New(tc.in)

		var tk string
		ok := tc.mf.On(Grab(&tk)).Run(c)

		assert.Equal(t, tc.ok, ok, tc.in)
		assert.Equal(t, tc.ex, tk, tc.in)
	}
}

func TestTag_Scan(t *testing.T) {

	tt := []struct {
		in string
		ok bool
		mf MatcherFunc
		ex []string
	}{
		{"Hello, {name}! You got {sign}{value}.", true, Tag("{", "}"), []string{"{name}", "{sign}", "{value}"}},
		{"Empty [] Not [1] Nested [ 1, [2, 3] ].", true, Tag("[", "]"), []string{"[]", "[1]", "[ 1, [2, 3] ]"}},
		{"print(1, 2); print(1, call()); println('hello', 3)", true, Tag("(", ")"), []string{"(1, 2)", "(1, call())", "('hello', 3)"}},
		{`Doc { "a": "b" } and subdoc { "a": { "b": 3 } }`, true, Tag("{", "}"), []string{`{ "a": "b" }`, `{ "a": { "b": 3 } }`}},
		{`<p>Click <a href="#">here</a></p>`, true, Tag("<", ">"), []string{`<p>`, `<a href="#">`, `</a>`, `</p>`}},
		{"This /* is */ a /* comment */", true, Tag("/*", "*/"), []string{"/* is */", "/* comment */"}},
		{"This /* is /* a nice */ nested */ comment", true, Tag("/*", "*/"), []string{"/* is /* a nice */ nested */"}},
	}

	for _, tc := range tt {

		c := New(tc.in)

		var tk []string

		oks := tc.mf.On(Grabs(&tk)).Scan(c)

		assert.Equal(t, tc.ok, oks, tc.in)
		assert.Equal(t, tc.ex, tk, tc.in)
	}
}

func TestScan(t *testing.T) {

	c := New("aaa")

	var tk []string
	ok1 := S("a").On(Grabs(&tk)).Scan(c)
	ok2 := S("a").Scan(c)

	assert.True(t, ok1)
	assert.Equal(t, []string{"a", "a", "a"}, tk)
	assert.False(t, ok2)
}
