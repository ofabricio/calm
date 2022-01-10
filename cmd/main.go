package main

import (
	"fmt"
	"unicode"

	. "github.com/ofabricio/calm"
)

func main() {

	src := New("Hello, World!")

	var words []string

	word := F(unicode.IsLetter).OneToMany().GrabMany(&words)
	only := Or(word, Next()).OneToMany()

	ok := src.Run(only)

	fmt.Println(ok, words)
	// true [Hello World]
}
