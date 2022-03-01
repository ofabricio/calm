package calm

import (
	"testing"
	"unicode"

	"github.com/stretchr/testify/assert"
)

func TestAST(t *testing.T) {

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

	exp := "Root [ Func func [ Name One, Args [ Var a, Type int ], Body [ Call A, Func func [ Name Two, Args [ Var b, Type int ], Body [ Call B, Call C ] ], Call D ] ] ]"

	// When.

	ws := SOr(" \t").OneToMany()
	wz := F(unicode.IsSpace).ZeroToMany()
	name := F(unicode.IsLetter).OneToMany()

	fnDefn, setFnDef := Recursive()
	fnArgs := And(name.Leaf("Var"), ws, name.Leaf("Type")).ZeroToOne()
	fnCall := And(name.Leaf("Call"), S("()"))
	fnBody := Or(wz.False(), fnDefn, fnCall).ZeroToMany()
	setFnDef(S("func").Leaf("Func").Child(ws, name.Leaf("Name"), wz, S("("), fnArgs.Group("Args"), S(")"), wz, S("{"), fnBody.Group("Body"), S("}"), wz))

	var ast AST
	ok := And(wz, fnDefn).Tree(&ast).Run(src)

	// Then.

	assert.True(t, ok)
	assert.Equal(t, exp, ast.Print("short-inline"))
}

func TestAST_Expression(t *testing.T) {

	// Given.

	tt := []struct {
		inp string
		exp string
	}{
		{"2+3", "Root [ BinExpr + [ Value 2, Value 3 ] ]"},
		{"2+3*4", "Root [ BinExpr + [ Value 2, BinExpr * [ Value 3, Value 4 ] ] ]"},
		{"2*3+4", "Root [ BinExpr + [ BinExpr * [ Value 2, Value 3 ], Value 4 ] ]"},
		{"2*(3+4)*5", "Root [ BinExpr * [ Value 2, BinExpr * [ BinExpr + [ Value 3, Value 4 ], Value 5 ] ] ]"},
	}

	for _, tc := range tt {

		src := New(tc.inp)

		// When.

		term, setTerm := Recursive()
		expr, setExpr := Recursive()

		value := F(unicode.IsNumber).Leaf("Value")

		factor := Or(And(S("("), expr, S(")")), value)

		setTerm(Or(Root(factor, S("*").Leaf("BinExpr"), term).Undo(), factor))
		setExpr(Or(Root(term, S("+").Leaf("BinExpr"), expr).Undo(), term))

		var ast AST
		ok := expr.Tree(&ast).Run(src)

		// Then.

		assert.True(t, ok)
		assert.Equal(t, tc.exp, ast.Print("short-inline"))
	}
}

func TestLeaf(t *testing.T) {

	// Given.

	src := New("abc")

	exp := "Root [ L a, L b, L c ]"

	// When.

	var ast AST
	ok := F(unicode.IsLetter).Leaf("L").OneToMany().Tree(&ast).Run(src)

	// Then.

	assert.True(t, ok)
	assert.Equal(t, exp, ast.Print("short-inline"))
}

func TestLeaf_When_False(t *testing.T) {

	// Given.

	src := New("a23")

	exp := "Root [ L a ]"

	// When.

	var ast AST
	ok := F(unicode.IsLetter).Leaf("L").OneToMany().Tree(&ast).Run(src)

	// Then.

	assert.True(t, ok)
	assert.Equal(t, exp, ast.Print("short-inline"))
}

func TestChild(t *testing.T) {

	// Given.

	src := New("abcde")

	exp := "Root [ L a [ L b, L c ], L d, L e ]"

	// When.

	root := And(
		S("a").Leaf("L").Child(
			S("b").Leaf("L"),
			S("c").Leaf("L"),
		),
		S("d").Leaf("L"),
		S("e").Leaf("L"),
	)

	var ast AST
	ok := root.Tree(&ast).Run(src)

	// Then.

	assert.True(t, ok)
	assert.Equal(t, exp, ast.Print("short-inline"))
}

func TestEnter_And_Leave_With_And(t *testing.T) {

	// Given.

	src := New("abcde")

	exp := "Root [ L a [ L b [ L c ] ], L d, L e ]"

	// When.

	root := And(
		And(
			S("a").Leaf("L").Enter(),
			S("b").Leaf("L").Enter(),
			S("c").Leaf("L"),
		).Leave(),
		S("d").Leaf("L"),
		S("e").Leaf("L"),
	)

	var ast AST
	ok := root.Tree(&ast).Run(src)

	// Then.

	assert.True(t, ok)
	assert.Equal(t, exp, ast.Print("short-inline"))
}

func TestLeave_When_False(t *testing.T) {

	// Given.

	src := New("ab")

	exp := "Root [ L a, L b ]"

	// When.

	root := Or(
		And(
			S("a").Leaf("L").Enter(),
			S("x").Leaf("L"),
		).Leave(),
		S("b").Leaf("L"),
	)

	var ast AST
	ok := root.Tree(&ast).Run(src)

	// Then.

	assert.True(t, ok)
	assert.Equal(t, exp, ast.Print("short-inline"))
}

func TestRoot(t *testing.T) {

	// Given.

	src := New("2+3")

	exp := "Root [ Op + [ V 2, V 3 ] ]"

	// When.

	root := Root(S("2").Leaf("V"), S("+").Leaf("Op"), S("3").Leaf("V"))

	var ast AST
	ok := root.Tree(&ast).Run(src)

	// Then.

	assert.True(t, ok)
	assert.Equal(t, exp, ast.Print("short-inline"))
}

func TestRoot_When_False(t *testing.T) {

	// Given.

	src := New("2+3")

	exp := "Root [ Op + [ V 2, V 3 ] ]"

	// When.

	root := Or(
		Root(S("2").Leaf("V"), S("*").Leaf("Op"), S("3").Leaf("V")).Undo(),
		Root(S("2").Leaf("V"), S("+").Leaf("Op"), S("3").Leaf("V")).Undo(),
	)

	var ast AST
	ok := root.Tree(&ast).Run(src)

	// Then.

	assert.True(t, ok)
	assert.Equal(t, exp, ast.Print("short-inline"))
}

func TestGroup(t *testing.T) {

	// Given.

	src := New("2+3")

	exp := "Root [ Group [ V 2, Op +, V 3 ] ]"

	// When.

	root := And(S("2").Leaf("V"), S("+").Leaf("Op"), S("3").Leaf("V"))

	var ast AST
	ok := root.Group("Group").Tree(&ast).Run(src)

	// Then.

	assert.True(t, ok)
	assert.Equal(t, exp, ast.Print("short-inline"))
}

func TestGroup_When_False(t *testing.T) {

	// Given.

	src := New("2+3")

	exp := "Root [ Group [ V 2, Op +, V 3 ] ]"

	// When.

	root := Or(
		And(S("2").Leaf("V"), S("*").Leaf("Op"), S("3").Leaf("V")).Group("Group").Undo(),
		And(S("2").Leaf("V"), S("+").Leaf("Op"), S("3").Leaf("V")).Group("Group"),
	)

	var ast AST
	ok := root.Tree(&ast).Run(src)

	// Then.

	assert.True(t, ok)
	assert.Equal(t, exp, ast.Print("short-inline"))
}

func TestGroup_Inside_Group(t *testing.T) {

	// Given.

	src := New("abc")

	exp := "Root [ Group [ Group [ L b ] ] ]"

	// When.

	root := And(S("a"), S("b").Leaf("L").Group("Group"), S("c")).Group("Group")

	var ast AST
	ok := root.Tree(&ast).Run(src)

	// Then.

	assert.True(t, ok)
	assert.Equal(t, exp, ast.Print("short-inline"))
}
