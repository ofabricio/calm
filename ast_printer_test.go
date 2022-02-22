package calm

import (
	"fmt"
	"os"
	"unicode"
)

func ExamplePrintTree_Short() {

	var ast AST
	oks := parseExpr(&ast)

	fmt.Println(oks)
	PrintTree(os.Stdout, "short", &ast)

	// Output:
	// true
	// Root [
	//     Expr + [
	//         Val 2
	//         Expr + [
	//             Expr * [
	//                 Val 3
	//                 Val 4
	//             ]
	//             Expr + [
	//                 Val 5
	//                 Val 6
	//             ]
	//         ]
	//     ]
	// ]
}

func ExamplePrintTree_ShortInline() {

	var ast AST
	oks := parseExpr(&ast)

	fmt.Println(oks)
	PrintTree(os.Stdout, "short-inline", &ast)

	// Output:
	// true
	// Root [ Expr + [ Val 2, Expr + [ Expr * [ Val 3, Val 4 ], Expr + [ Val 5, Val 6 ] ] ] ]
}

func ExamplePrintTree_Nice() {

	var ast AST
	oks := parseExpr(&ast)

	fmt.Println(oks)
	PrintTree(os.Stdout, "nice", &ast)

	// Output:
	// true
	// Root [ + [ 2, + [ * [ 3, 4 ], + [ 5, 6 ] ] ] ]
}

func ExamplePrintTree_Json() {

	var ast AST
	oks := parseExpr(&ast)

	fmt.Println(oks)
	PrintTree(os.Stdout, "json", &ast)

	// Output:
	// true
	// {
	//     "type": "Root",
	//     "args": [
	//         {
	//             "type": "Expr",
	//             "name": "+",
	//             "args": [
	//                 {
	//                     "type": "Val",
	//                     "name": "2"
	//                 },
	//                 {
	//                     "type": "Expr",
	//                     "name": "+",
	//                     "args": [
	//                         {
	//                             "type": "Expr",
	//                             "name": "*",
	//                             "args": [
	//                                 {
	//                                     "type": "Val",
	//                                     "name": "3"
	//                                 },
	//                                 {
	//                                     "type": "Val",
	//                                     "name": "4"
	//                                 }
	//                             ]
	//                         },
	//                         {
	//                             "type": "Expr",
	//                             "name": "+",
	//                             "args": [
	//                                 {
	//                                     "type": "Val",
	//                                     "name": "5"
	//                                 },
	//                                 {
	//                                     "type": "Val",
	//                                     "name": "6"
	//                                 }
	//                             ]
	//                         }
	//                     ]
	//                 }
	//             ]
	//         }
	//     ]
	// }
}

func ExamplePrintTree_JsonInline() {

	var ast AST
	oks := parseExpr(&ast)

	fmt.Println(oks)
	PrintTree(os.Stdout, "json-inline", &ast)

	// Output:
	// true
	// { "type": "Root", "args": [{ "type": "Expr", "name": "+", "args": [{ "type": "Val", "name": "2" }, { "type": "Expr", "name": "+", "args": [{ "type": "Expr", "name": "*", "args": [{ "type": "Val", "name": "3" }, { "type": "Val", "name": "4" }] }, { "type": "Expr", "name": "+", "args": [{ "type": "Val", "name": "5" }, { "type": "Val", "name": "6" }] }] }] }] }
}

func ExamplePrintTree_Json_String_Encoding() {

	var ast AST
	oks := String(`"`).Leaf("Str").Tree(&ast).Run(New(`"Hello"`))

	fmt.Println(oks)
	PrintTree(os.Stdout, "json-inline", &ast)

	// Output:
	// true
	// { "type": "Root", "args": [{ "type": "Str", "name": "\"Hello\"" }] }
}

func parseExpr(ast *AST) bool {
	term, setTerm := Recursive()
	expr, setExpr := Recursive()

	value := F(unicode.IsDigit).Leaf("Val")
	factor := Or(And(S("("), expr, S(")")), value)
	setTerm(Or(And(factor, S("*").Leaf("Expr"), term).Root(), factor))
	setExpr(Or(And(term, S("+").Leaf("Expr"), expr).Root(), term))
	return expr.Tree(ast).Run(New("2+3*4+(5+6)"))
}
