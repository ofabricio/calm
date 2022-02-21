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
	pkgDef := And(S("package").Leaf("PkgDef").Enter(), ws, name.Leaf("Ident"))
	impDef := And(S("import").Leaf("ImpDef").Enter(), ws, strg)

	fnCall := And(name.Leaf("Pkg"), S("."), name.Leaf("Ident"), S("("), strg, S(")")).Group("FnCall")
	fnBody := Or(wz.False(), fnCall).ZeroToMany().Group("Body")
	fnDef := And(S("func").Leaf("FnDef").Enter(), ws, name.Leaf("Ident"), wz, S("()"), wz,
		S("{"), fnBody, wz, S("}"))

	root := Or(
		wz.False(),
		comment,
		pkgDef,
		impDef,
		fnDef,
	).ZeroToMany()

	var ast Ast
	ok := root.Tree(&ast).Run(src)

	fmt.Println("Ok:", ok)
	fmt.Println(ast.Print("short"))

	// Output:
	// Ok: true
	// Root [
	//     Comment // You can edit this code!
	//     Comment // Click here and start typing.
	//     PkgDef package [
	//         Ident main
	//     ]
	//     ImpDef import [
	//         Str "fmt"
	//     ]
	//     FnDef func [
	//         Ident main
	//         Body [
	//             FnCall [
	//                 Pkg fmt
	//                 Ident Println
	//                 Str "Hello, 世界"
	//             ]
	//         ]
	//     ]
	// ]
}
