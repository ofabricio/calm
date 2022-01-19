package main

import (
	"fmt"
	"unicode"

	. "github.com/ofabricio/calm"
)

func main() {

	code := New("Hello, World!")

	var words []string

	word := F(unicode.IsLetter).OneToMany().On(Grabs(&words))
	root := Or(word, Next()).OneToMany()

	ok := root.Run(code)

	fmt.Println(ok, words)
	// true [Hello World]
}
