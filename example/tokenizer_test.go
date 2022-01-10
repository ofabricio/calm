package example

import (
	"fmt"
	"unicode"

	. "github.com/ofabricio/calm"
)

func Example() {

	src := `
	package main

	import "fmt"

	func main() {
		fmt.Println("Hello, 世界")
	}`

	c := New(src)

	print := func(t Token) {
		fmt.Printf("'%s' Row: %d Col: %d\n", t.Text, t.Row, t.Col)
	}

	newl := S("\n")
	spac := F(unicode.IsSpace)
	strs := String(`"`).On(print)
	word := F(unicode.IsLetter).OneToMany().On(print)
	rest := F(unicode.IsLetter).Not().Next().On(print)

	code := Or(newl, spac, strs, word, rest).ZeroToMany()

	ok := c.Run(code)

	fmt.Println(ok)

	// Output:
	// 'package' Row: 2 Col: 2
	// 'main' Row: 2 Col: 10
	// 'import' Row: 4 Col: 2
	// '"fmt"' Row: 4 Col: 9
	// 'func' Row: 6 Col: 2
	// 'main' Row: 6 Col: 7
	// '(' Row: 6 Col: 11
	// ')' Row: 6 Col: 12
	// '{' Row: 6 Col: 14
	// 'fmt' Row: 7 Col: 3
	// '.' Row: 7 Col: 6
	// 'Println' Row: 7 Col: 7
	// '(' Row: 7 Col: 14
	// '"Hello, 世界"' Row: 7 Col: 15
	// ')' Row: 7 Col: 26
	// '}' Row: 8 Col: 2
	// true
}
