package calm

import (
	"testing"
	"unicode"

	"github.com/stretchr/testify/assert"
)

func TestAst(t *testing.T) {

	// Given.

	src := New(`
	func One(a int) {
		A()
		func Two(b int) {
			B()
			C()
		}
		D()
	}`)

	exp := `{ "type": "Root", "args": [{ "type": "FnDef", "name": "func", "args": [{ "type": "Name", "name": "One" }, { "type": "Args", "args": [{ "type": "Var", "name": "a" }, { "type": "Type", "name": "int" }] }, { "type": "Body", "args": [{ "type": "FnCall", "name": "A" }, { "type": "FnDef", "name": "func", "args": [{ "type": "Name", "name": "Two" }, { "type": "Args", "args": [{ "type": "Var", "name": "b" }, { "type": "Type", "name": "int" }] }, { "type": "Body", "args": [{ "type": "FnCall", "name": "B" }, { "type": "FnCall", "name": "C" }] }] }, { "type": "FnCall", "name": "D" }] }] }] }`

	// When.

	ws := SOr(" \t").OneToMany()
	wz := F(unicode.IsSpace).ZeroToMany()
	name := F(unicode.IsLetter).OneToMany()

	fnArgs := And(name.Leaf("Var"), ws, name.Leaf("Type")).ZeroToOne()
	fnCall := And(name.Leaf("FnCall"), S("()"))

	var fnDefn, fnBody MatcherFunc

	fnBody = func(c *Code) bool {
		return Or(wz.False(), fnCall.Undo(), fnDefn.Undo()).ZeroToMany().Run(c)
	}

	fnDefn = func(c *Code) bool {
		return And(wz, S("func").Leaf("FnDef").Enter(), ws, name.Leaf("Name"), wz, S("("), fnArgs.Group("Args"), S(")"), wz, S("{"), wz, fnBody.Group("Body"), wz, S("}"), wz).Run(c)
	}

	var ast Ast
	ok := fnDefn.Tree(&ast).Run(src)

	// Then.

	assert.True(t, ok)
	assert.Equal(t, exp, ast.String())
}

func TestAst_Expression(t *testing.T) {

	// Given.

	tt := []struct {
		inp string
		exp string
	}{
		{"2+3", `{ "type": "Root", "args": [{ "type": "BinExpr", "name": "+", "args": [{ "type": "Value", "name": "2" }, { "type": "Value", "name": "3" }] }] }`},
		{"2+3*4", `{ "type": "Root", "args": [{ "type": "BinExpr", "name": "+", "args": [{ "type": "Value", "name": "2" }, { "type": "BinExpr", "name": "*", "args": [{ "type": "Value", "name": "3" }, { "type": "Value", "name": "4" }] }] }] }`},
		{"2*3+4", `{ "type": "Root", "args": [{ "type": "BinExpr", "name": "+", "args": [{ "type": "BinExpr", "name": "*", "args": [{ "type": "Value", "name": "2" }, { "type": "Value", "name": "3" }] }, { "type": "Value", "name": "4" }] }] }`},
		{"2*(3+4)*5", `{ "type": "Root", "args": [{ "type": "BinExpr", "name": "*", "args": [{ "type": "Value", "name": "2" }, { "type": "BinExpr", "name": "*", "args": [{ "type": "BinExpr", "name": "+", "args": [{ "type": "Value", "name": "3" }, { "type": "Value", "name": "4" }] }, { "type": "Value", "name": "5" }] }] }] }`},
	}

	for _, tc := range tt {

		src := New(tc.inp)

		// When.

		var term, expr, factor func(c *Code) bool

		value := F(unicode.IsNumber).Leaf("Value")

		factor = func(c *Code) bool {
			return Or(And(S("("), expr, S(")")), value).Run(c)
		}

		term = func(c *Code) bool {
			return Or(And(factor, S("*").Leaf("BinExpr"), term).Root(), factor).Run(c)
		}

		expr = func(c *Code) bool {
			return Or(And(term, S("+").Leaf("BinExpr"), expr).Root(), term).Run(c)
		}

		var ast Ast
		ok := MatcherFunc(expr).Tree(&ast).Run(src)

		// Then.

		assert.True(t, ok)
		assert.Equal(t, tc.exp, ast.String())
	}
}

func TestLeaf(t *testing.T) {

	// Given.

	src := New("abc")

	exp := `{ "type": "Root", "args": [{ "type": "L", "name": "a" }, { "type": "L", "name": "b" }, { "type": "L", "name": "c" }] }`

	// When.

	var ast Ast
	ok := F(unicode.IsLetter).Leaf("L").OneToMany().Tree(&ast).Run(src)

	// Then.

	assert.True(t, ok)
	assert.Equal(t, exp, ast.String())
}

func TestLeaf_When_False(t *testing.T) {

	// Given.

	src := New("a23")

	exp := `{ "type": "Root", "args": [{ "type": "L", "name": "a" }] }`

	// When.

	var ast Ast
	ok := F(unicode.IsLetter).Leaf("L").OneToMany().Tree(&ast).Run(src)

	// Then.

	assert.True(t, ok)
	assert.Equal(t, exp, ast.String())
}

func TestChild(t *testing.T) {

	// Given.

	src := New("abcde")

	exp := `{ "type": "Root", "args": [{ "type": "L", "name": "a", "args": [{ "type": "L", "name": "b" }, { "type": "L", "name": "c" }] }, { "type": "L", "name": "d" }, { "type": "L", "name": "e" }] }`

	// When.

	root := And(
		S("a").Leaf("L").Child(
			S("b").Leaf("L"),
			S("c").Leaf("L"),
		),
		S("d").Leaf("L"),
		S("e").Leaf("L"),
	)

	var ast Ast
	ok := root.Tree(&ast).Run(src)

	// Then.

	assert.True(t, ok)
	assert.Equal(t, exp, ast.String())
}

func TestEnter_And_Leave_With_And(t *testing.T) {

	// Given.

	src := New("abcde")

	exp := `{ "type": "Root", "args": [{ "type": "L", "name": "a", "args": [{ "type": "L", "name": "b", "args": [{ "type": "L", "name": "c" }] }] }, { "type": "L", "name": "d" }, { "type": "L", "name": "e" }] }`

	// When.

	root := And(
		And(
			S("a").Leaf("L").Enter(),
			S("b").Leaf("L").Enter(),
			S("c").Leaf("L"),
		),
		S("d").Leaf("L"),
		S("e").Leaf("L"),
	)

	var ast Ast
	ok := root.Tree(&ast).Run(src)

	// Then.

	assert.True(t, ok)
	assert.Equal(t, exp, ast.String())
}

func TestLeave_When_False(t *testing.T) {

	// Given.

	src := New("ab")

	exp := `{ "type": "Root", "args": [{ "type": "L", "name": "a" }, { "type": "L", "name": "b" }] }`

	// When.

	root := Or(
		And(
			S("a").Leaf("L").Enter(),
			S("x").Leaf("L"),
		),
		And(
			S("a").Leaf("L"),
			S("b").Leaf("L"),
		),
	)

	var ast Ast
	ok := root.Tree(&ast).Run(src)

	// Then.

	assert.True(t, ok)
	assert.Equal(t, exp, ast.String())
}

func TestRoot(t *testing.T) {

	// Given.

	src := New("2+3")

	exp := `{ "type": "Root", "args": [{ "type": "Op", "name": "+", "args": [{ "type": "V", "name": "2" }, { "type": "V", "name": "3" }] }] }`

	// When.

	root := And(S("2").Leaf("V"), S("+").Leaf("Op"), S("3").Leaf("V")).Root()

	var ast Ast
	ok := root.Tree(&ast).Run(src)

	// Then.

	assert.True(t, ok)
	assert.Equal(t, exp, ast.String())
}

func TestRoot_When_False(t *testing.T) {

	// Given.

	src := New("2+3")

	exp := `{ "type": "Root", "args": [{ "type": "Op", "name": "+", "args": [{ "type": "V", "name": "2" }, { "type": "V", "name": "3" }] }] }`

	// When.

	root := Or(
		And(S("2").Leaf("V"), S("*").Leaf("Op"), S("3").Leaf("V")).Root(),
		And(S("2").Leaf("V"), S("+").Leaf("Op"), S("3").Leaf("V")).Root(),
	)

	var ast Ast
	ok := root.Tree(&ast).Run(src)

	// Then.

	assert.True(t, ok)
	assert.Equal(t, exp, ast.String())
}

func TestGroup(t *testing.T) {

	// Given.

	src := New("2+3")

	exp := `{ "type": "Root", "args": [{ "type": "Group", "args": [{ "type": "V", "name": "2" }, { "type": "Op", "name": "+" }, { "type": "V", "name": "3" }] }] }`

	// When.

	root := And(S("2").Leaf("V"), S("+").Leaf("Op"), S("3").Leaf("V"))

	var ast Ast
	ok := root.Group("Group").Tree(&ast).Run(src)

	// Then.

	assert.True(t, ok)
	assert.Equal(t, exp, ast.String())
}

func TestGroup_When_False(t *testing.T) {

	// Given.

	src := New("2+3")

	exp := `{ "type": "Root", "args": [{ "type": "Group", "args": [{ "type": "V", "name": "2" }, { "type": "Op", "name": "+" }, { "type": "V", "name": "3" }] }] }`

	// When.

	root := Or(
		And(S("2").Leaf("V"), S("*").Leaf("Op"), S("3").Leaf("V")).Group("Group"),
		And(S("2").Leaf("V"), S("+").Leaf("Op"), S("3").Leaf("V")).Group("Group"),
	)

	var ast Ast
	ok := root.Tree(&ast).Run(src)

	// Then.

	assert.True(t, ok)
	assert.Equal(t, exp, ast.String())
}

func TestGroup_Inside_Group(t *testing.T) {

	// Given.

	src := New("abc")

	exp := `{ "type": "Root", "args": [{ "type": "Group", "args": [{ "type": "Group", "args": [{ "type": "L", "name": "b" }] }] }] }`

	// When.

	root := And(S("a"), S("b").Leaf("L").Group("Group"), S("c")).Group("Group")

	var ast Ast
	ok := root.Tree(&ast).Run(src)

	// Then.

	assert.True(t, ok)
	assert.Equal(t, exp, ast.String())
}
