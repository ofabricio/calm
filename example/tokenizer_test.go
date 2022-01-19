package example

import (
	"fmt"
	"unicode"

	. "github.com/ofabricio/calm"
)

func Example() {

	code := New(`
	package main

	import "fmt"

	func main() {
		fmt.Println("Hello, 世界")
	}`)

	print := func(t Token) {
		fmt.Printf("Pos: %-3d Line: %-2d Column: %-3d Token: %s\n", t.Pos, t.Row, t.Col, t.Text)
	}

	spac := F(unicode.IsSpace)
	strg := String(`"`).On(print)
	word := F(unicode.IsLetter).OneToMany().On(print)
	rest := Next().On(print)

	root := Or(spac, strg, word, rest).ZeroToMany()

	ok := root.Run(code)

	fmt.Println(ok)

	// Output:
	// Pos: 2   Line: 2  Column: 2   Token: package
	// Pos: 10  Line: 2  Column: 10  Token: main
	// Pos: 17  Line: 4  Column: 2   Token: import
	// Pos: 24  Line: 4  Column: 9   Token: "fmt"
	// Pos: 32  Line: 6  Column: 2   Token: func
	// Pos: 37  Line: 6  Column: 7   Token: main
	// Pos: 41  Line: 6  Column: 11  Token: (
	// Pos: 42  Line: 6  Column: 12  Token: )
	// Pos: 44  Line: 6  Column: 14  Token: {
	// Pos: 48  Line: 7  Column: 3   Token: fmt
	// Pos: 51  Line: 7  Column: 6   Token: .
	// Pos: 52  Line: 7  Column: 7   Token: Println
	// Pos: 59  Line: 7  Column: 14  Token: (
	// Pos: 60  Line: 7  Column: 15  Token: "Hello, 世界"
	// Pos: 75  Line: 7  Column: 26  Token: )
	// Pos: 78  Line: 8  Column: 2   Token: }
	// true
}
