package example

import (
	"fmt"
	"unicode"

	. "github.com/ofabricio/calm"
)

func ExampleGoCodeParsing() {

	src := New(`
	// You can edit this code!
	// Click here and start typing.
	package main

	import "fmt"

	func main() {
		fmt.Println("Hello, 世界")
	}`)

	ws := SOr(" \t").OneToMany()
	wz := F(unicode.IsSpace).ZeroToMany()
	name := F(unicode.IsLetter).OneToMany()
	strg := String(`"`).Leaf("Str")

	comment := And(S("//"), Until(Eq("\n"))).Leaf("Comment")
	pkgDef := And(S("package").Leaf("Pkg").Enter(), ws, name.Leaf("Name"))
	impDef := And(S("import").Leaf("Imp").Enter(), ws, strg)

	fnCall := And(name.Leaf("Pkg"), S("."), name.Leaf("Name"), S("("), strg, S(")")).Group("Call")
	fnBody := Or(wz.False(), fnCall).ZeroToMany().Group("Body")
	fnDef := And(S("func").Leaf("Fun").Enter(), ws, name.Leaf("Name"), wz, S("()"), wz,
		S("{"), fnBody, wz, S("}"))

	root := Or(
		wz.False(),
		comment,
		pkgDef,
		impDef,
		fnDef,
	).ZeroToMany()

	var ast AST
	ok := root.Tree(&ast).Run(src)

	fmt.Println("Ok:", ok)
	fmt.Println(ast.Print("short"))

	// Output:
	// Ok: true
	// Root [
	//     Comment // You can edit this code!
	//     Comment // Click here and start typing.
	//     Pkg package [
	//         Name main
	//     ]
	//     Imp import [
	//         Str "fmt"
	//     ]
	//     Fun func [
	//         Name main
	//         Body [
	//             Call [
	//                 Pkg fmt
	//                 Name Println
	//                 Str "Hello, 世界"
	//             ]
	//         ]
	//     ]
	// ]
}
