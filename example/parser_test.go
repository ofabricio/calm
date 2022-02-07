package example

import (
	"fmt"
	"unicode"

	. "github.com/ofabricio/calm"
)

func ExampleGoProgram() {

	src := New(`
	// You can edit this code!
	// Click here and start typing.
	package main

	import "fmt"

	func main() {
		fmt.Println("Hello, 世界")
	}`)

	p := &GoProgram{}

	ok := p.program().Run(src)

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
		t.packageName(&t.Package),
		t.imports(&t.Imports).ZeroToMany(),
		t.globals().ZeroToMany(),
	)
}

func (t *GoProgram) globals() MatcherFunc {
	return Or(
		t.function(),
		t.stop(),
	)
}

func (t *GoProgram) function() MatcherFunc {
	// TODO: note the t.sw() pattern below.
	// How could it be handled? Should it?
	var f GoFunction
	return And(
		t.comment().True(),
		t.ws(),
		S("func "),
		t.ws(),
		t.name().On(Grab(&f.Name)),
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
		f = GoFunction{}
	})
}

func (t *GoProgram) statement(body *[]string) MatcherFunc {
	return And(t.ws(), Until(Eq("\n"), Eq("}")).On(Grabs(body)))
}

func (t *GoProgram) imports(name *[]string) MatcherFunc {
	return And(t.comment().True(), t.ws(), S("import "), String(`"`).On(Grabs(name)))
}

func (t *GoProgram) packageName(name *string) MatcherFunc {
	return And(t.comment().True(), t.ws(), S("package "), t.name().On(Grab(name)))
}

func (t *GoProgram) comment() MatcherFunc {
	return And(t.ws(), And(S("//"), Until(Eq("\n")).True()).On(Grabs(&t.Comments))).ZeroToMany()
}

func (t *GoProgram) name() MatcherFunc {
	return F(unicode.IsLetter).OneToMany()
}

func (t *GoProgram) ws() MatcherFunc {
	return F(unicode.IsSpace).ZeroToMany()
}

func (t *GoProgram) stop() MatcherFunc {
	return Next().False()
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
