package calm

import (
	"testing"
	"unicode"

	"github.com/stretchr/testify/assert"
)

func TestOn(t *testing.T) {

	tt := []struct {
		in       string
		ok       bool
		expToken string
		expIdx   int
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

		ok := And(S("a"), S("b").On(on), S("c")).Run(c)

		assert.Equal(t, tc.ok, ok, tc.in)
		assert.Equal(t, tc.expToken, tk.Text, tc.in)
		assert.Equal(t, tc.expIdx, tk.Pos, tc.in)
	}
}

func TestEmit(t *testing.T) {

	c := New("ab\ncdef\ngh")

	var tk1, tk2 Token
	ok1 := S("ab\ncd").On(Emit(&tk1)).Run(c)
	ok2 := S("ef\ngh").On(Emit(&tk2)).Run(c)

	assert.True(t, ok1)
	assert.Equal(t, Token{Text: "ab\ncd", Pos: 0, Row: 1, Col: 1}, tk1)
	assert.True(t, ok2)
	assert.Equal(t, Token{Text: "ef\ngh", Pos: 5, Row: 2, Col: 3}, tk2)
}

func TestEmits(t *testing.T) {

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
	ok := Or(F(unicode.IsLetter), S("\n")).On(Emits(&tks)).OneToMany().Run(c)

	assert.True(t, ok)
	assert.Equal(t, exp, tks)
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
		ok := And(S("a"), S("b").On(Grab(&token)), S("c")).Run(c)

		assert.Equal(t, tc.ok, ok, tc.in)
		assert.Equal(t, tc.ex, token, tc.in)
	}
}

func TestGrabs(t *testing.T) {

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

		ok := F(unicode.IsLetter).On(Grabs(&tokens)).OneToMany().Run(c)

		assert.Equal(t, tc.ok, ok, tc.in)
		assert.Equal(t, tc.ex, tokens, tc.in)
	}
}

func TestIndex(t *testing.T) {

	tt := []struct {
		in     string
		ok     bool
		expIdx int
	}{
		{"abc", true, 1},
		{"axc", false, 0},
	}

	for _, tc := range tt {

		c := New(tc.in)

		var idx int
		ok := And(S("a"), S("b").On(Index(&idx)), S("c")).Run(c)

		assert.Equal(t, tc.ok, ok, tc.in)
		assert.Equal(t, tc.expIdx, idx, tc.in)
	}
}

func TestIndexes(t *testing.T) {

	tt := []struct {
		in string
		ok bool
		ex []int
	}{
		{"", false, nil},
		{"abc", true, []int{0, 1, 2}},
	}

	for _, tc := range tt {

		c := New(tc.in)

		var idx []int

		ok := F(unicode.IsLetter).On(Indexes(&idx)).OneToMany().Run(c)

		assert.Equal(t, tc.ok, ok, tc.in)
		assert.Equal(t, tc.ex, idx, tc.in)
	}
}

func TestToInt(t *testing.T) {

	tt := []struct {
		in string
		ok bool
		ex int
	}{
		{"0", true, 0},
		{"1", true, 1},
		{"33", true, 33},
	}

	for _, tc := range tt {

		c := New(tc.in)

		var v int
		ok := Next().OneToMany().On(ToInt(&v)).Run(c)

		assert.Equal(t, tc.ok, ok, tc.in)
		assert.Equal(t, tc.ex, v, tc.in)
	}
}

func TestToFloat(t *testing.T) {

	tt := []struct {
		in string
		ok bool
		ex float64
	}{
		{"0", true, 0},
		{"1", true, 1},
		{"0.0", true, 0},
		{"1.0", true, 1},
		{"3.3", true, 3.3},
	}

	for _, tc := range tt {

		c := New(tc.in)

		var v float64
		ok := Next().OneToMany().On(ToFloat(&v)).Run(c)

		assert.Equal(t, tc.ok, ok, tc.in)
		assert.Equal(t, tc.ex, v, tc.in)
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
		ok := And(SOr(`"'`).On(Grab(&quote)), S("a"), SR(&quote)).Run(c)

		assert.Equal(t, tc.ok, ok, tc.in)
	}
}
