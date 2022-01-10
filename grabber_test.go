package calm

import (
	"testing"
	"unicode"

	"github.com/stretchr/testify/assert"
)

func TestEmit(t *testing.T) {

	c := New("ab\ncdef\ngh")

	var tk1, tk2 Token
	ok1 := c.Run(S("ab\ncd").Emit(&tk1))
	ok2 := c.Run(S("ef\ngh").Emit(&tk2))

	assert.True(t, ok1)
	assert.Equal(t, Token{Text: "ab\ncd", Pos: 0, Row: 1, Col: 1}, tk1)
	assert.True(t, ok2)
	assert.Equal(t, Token{Text: "ef\ngh", Pos: 5, Row: 2, Col: 3}, tk2)
}

func TestEmitMany(t *testing.T) {

	exp := []Token{
		{Text: "a", Pos: 0, Row: 1, Col: 1},
		{Text: "b", Pos: 1, Row: 1, Col: 2},
		{Text: "\n", Pos: 2, Row: 1, Col: 3},
		{Text: "c", Pos: 3, Row: 2, Col: 1},
		{Text: "d", Pos: 4, Row: 2, Col: 2},
		{Text: "e", Pos: 5, Row: 2, Col: 3},
		{Text: "f", Pos: 6, Row: 2, Col: 4},
		{Text: "\n", Pos: 7, Row: 2, Col: 5},
		{Text: "g", Pos: 8, Row: 3, Col: 1},
		{Text: "h", Pos: 9, Row: 3, Col: 2},
	}

	c := New("ab\ncdef\ngh")

	var tks []Token
	ok := c.Run(Or(F(unicode.IsLetter), S("\n")).EmitMany(&tks).OneToMany())

	assert.True(t, ok)
	assert.Equal(t, exp, tks)
}

func TestEmitUndo(t *testing.T) {

	exp := []Token{
		{Text: "0", Pos: 0, Row: 1, Col: 1},
		{Text: "+", Pos: 1, Row: 1, Col: 2},
		{Text: "1", Pos: 2, Row: 1, Col: 3},
		{Text: "*", Pos: 3, Row: 1, Col: 4},
		{Text: "(", Pos: 4, Row: 1, Col: 5},
		{Text: "2", Pos: 5, Row: 1, Col: 6},
		{Text: "+", Pos: 6, Row: 1, Col: 7},
		{Text: "3", Pos: 7, Row: 1, Col: 8},
		{Text: ")", Pos: 8, Row: 1, Col: 9},
		{Text: "*", Pos: 9, Row: 1, Col: 10},
		{Text: "4", Pos: 10, Row: 1, Col: 11},
	}

	c := New("0+1*(2+3)*4")

	var tokens []Token

	var term Wrap
	var expr Wrap

	value := F(unicode.IsNumber).EmitMany(&tokens)
	factor := Or(And(S("(").EmitMany(&tokens), &expr, S(")").EmitMany(&tokens)), value)

	Or(And(factor, S("*").EmitMany(&tokens), &term).Not().Rewind().EmitUndo(&tokens).Not(), factor).Recursive(&term)
	Or(And(&term, S("+").EmitMany(&tokens), &expr).Not().Rewind().EmitUndo(&tokens).Not(), &term).Recursive(&expr)

	ok := c.Run(&expr)

	assert.True(t, ok)
	assert.Equal(t, exp, tokens)
}

func TestGrab(t *testing.T) {

	tt := []struct {
		in string
		ok bool
		ex string
	}{
		{"abc", true, "b"},
		{"axc", false, ""},
	}

	for _, tc := range tt {

		c := New(tc.in)

		var token string
		ok := c.Run(And(S("a"), S("b").Grab(&token), S("c")))

		assert.Equal(t, tc.ok, ok, tc.in)
		assert.Equal(t, tc.ex, token, tc.in)
	}
}

func TestGrabMany(t *testing.T) {

	tt := []struct {
		in string
		ok bool
		ex []string
	}{
		{"", false, nil},
		{"abc", true, []string{"a", "b", "c"}},
	}

	for _, tc := range tt {

		c := New(tc.in)

		var tokens []string

		ok := c.Run(F(unicode.IsLetter).GrabMany(&tokens).OneToMany())

		assert.Equal(t, tc.ok, ok, tc.in)
		assert.Equal(t, tc.ex, tokens, tc.in)
	}
}

func TestGrabUndo(t *testing.T) {

	exp := []string{"0", "+", "1", "*", "(", "2", "+", "3", ")", "*", "4"}

	c := New("0+1*(2+3)*4")

	var tokens []string

	var term Wrap
	var expr Wrap

	value := F(unicode.IsNumber).GrabMany(&tokens)
	factor := Or(And(S("(").GrabMany(&tokens), &expr, S(")").GrabMany(&tokens)), value)

	Or(And(factor, S("*").GrabMany(&tokens), &term).Not().Rewind().GrabUndo(&tokens).Not(), factor).Recursive(&term)
	Or(And(&term, S("+").GrabMany(&tokens), &expr).Not().Rewind().GrabUndo(&tokens).Not(), &term).Recursive(&expr)

	ok := c.Run(&expr)

	assert.True(t, ok)
	assert.Equal(t, exp, tokens)
}

func TestGrabPos(t *testing.T) {

	tt := []struct {
		in     string
		ok     bool
		expPos int
	}{
		{"abc", true, 1},
		{"axc", false, 0},
	}

	for _, tc := range tt {

		c := New(tc.in)

		var pos int
		ok := c.Run(And(S("a"), S("b").GrabPos(&pos), S("c")))

		assert.Equal(t, tc.ok, ok, tc.in)
		assert.Equal(t, tc.expPos, pos, tc.in)
	}
}

func TestBackReference_With_Grab_And_SR(t *testing.T) {

	tt := []struct {
		in string
		ok bool
	}{
		{`"a"`, true},
		{`'a'`, true},
		{`'a"`, false},
	}

	for _, tc := range tt {

		c := New(tc.in)

		var quote string
		ok := c.Run(And(Or(S("\""), S("'")).Grab(&quote), S("a"), SR(&quote)))

		assert.Equal(t, tc.ok, ok, tc.in)
	}
}

func TestOn(t *testing.T) {

	tt := []struct {
		in       string
		ok       bool
		expToken string
		expIni   int
	}{
		{"abc", true, "b", 1},
		{"axc", false, "", 0},
	}

	for _, tc := range tt {

		c := New(tc.in)

		var tk Token

		on := func(tkn Token) {
			tk = tkn
		}

		ok := c.Run(And(S("a"), S("b").On(on), S("c")))

		assert.Equal(t, tc.ok, ok, tc.in)
		assert.Equal(t, tc.expToken, tk.Text, tc.in)
		assert.Equal(t, tc.expIni, tk.Pos, tc.in)
	}
}
