package example

import (
	"fmt"
	"unicode"

	. "github.com/ofabricio/calm"
)

func ExampleGoProgram() {

	src := `
	// You can edit this code!
	// Click here and start typing.
	package main

	import "fmt"

	func main() {
		fmt.Println("Hello, 世界")
	}`

	c := New(src)

	p := &GoProgram{}

	ok := c.Run(p.program())

	fmt.Println("Package:", p.Package)
	fmt.Println("Imports:", p.Imports)
	fmt.Println("Comments:", p.Comments)
	fmt.Println("Functions:", p.Functions)
	fmt.Println("Ok:", ok)

	// Output:
	// Package: main
	// Imports: ["fmt"]
	// Comments: [// You can edit this code! // Click here and start typing.]
	// Functions: [{main [fmt.Println("Hello, 世界")]}]
	// Ok: true
}

func (t *GoProgram) program() MatcherFunc {
	return And(
		t.comment().True(),
		t.packageName(&t.Package),
		t.imports(&t.Imports).ZeroToMany(),
		t.globals().ZeroToMany(),
	)
}

func (t *GoProgram) globals() MatcherFunc {
	return Or(
		t.comment().False(),
		t.function(),
		t.stop(),
	)
}

func (t *GoProgram) function() MatcherFunc {
	// TODO: note the t.sw() pattern below.
	// How could it be handled? Should it?
	var f GoFunction
	return And(
		t.ws(),
		S("func "),
		t.ws(),
		t.name().Grab(&f.Name),
		t.ws(),
		S("("),
		t.ws(),
		S(")"),
		t.ws(),
		S("{"),
		t.statement(&f.Body).ZeroToMany(),
		t.ws(),
		S("}"),
	).On(func(Token) {
		t.Functions = append(t.Functions, f)
	})
}

func (t *GoProgram) statement(body *[]string) MatcherFunc {
	return And(t.ws(), Until(Eq("\n"), Eq("}")).GrabMany(body))
}

func (t *GoProgram) imports(name *[]string) MatcherFunc {
	return And(t.ws(), S("import "), String(`"`).GrabMany(name))
}

func (t *GoProgram) packageName(name *string) MatcherFunc {
	return And(t.ws(), S("package "), t.name().Grab(name))
}

func (t *GoProgram) comment() MatcherFunc {
	return And(t.ws(), And(S("//"), Until(Eq("\n")).True()).GrabMany(&t.Comments)).ZeroToMany()
}

func (t *GoProgram) name() MatcherFunc {
	return F(unicode.IsLetter).OneToMany()
}

func (t *GoProgram) ws() MatcherFunc {
	return F(unicode.IsSpace).ZeroToMany()
}

func (t *GoProgram) stop() MatcherFunc {
	return F(unicode.IsPrint).False()
}

type GoProgram struct {
	Package   string
	Imports   []string
	Comments  []string
	Functions []GoFunction
}

type GoFunction struct {
	Name string
	Body []string
}
