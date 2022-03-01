package example

import (
	"fmt"
	"strconv"
	"unicode"

	. "github.com/ofabricio/calm"
)

func ExampleAST() {

	src := New("6+5*(4+3)*2")

	term, setTerm := Recursive()
	expr, setExpr := Recursive()

	value := F(unicode.IsNumber).Leaf("Value")
	factor := Or(And(S("("), expr, S(")")), value)
	setTerm(Or(And(factor, S("*").Leaf("Expr"), term).Undo().Root(), factor))
	setExpr(Or(And(term, S("+").Leaf("Expr"), expr).Undo().Root(), term))

	// When.

	var ast AST
	ok := expr.Tree(&ast).Run(src)

	fmt.Println(ast.Print("short"))
	fmt.Println("Result:", calcResult(&ast))
	fmt.Println("Ok:", ok)

	// Output:
	// Root [
	//     Expr + [
	//         Value 6
	//         Expr * [
	//             Value 5
	//             Expr * [
	//                 Expr + [
	//                     Value 4
	//                     Value 3
	//                 ]
	//                 Value 2
	//             ]
	//         ]
	//     ]
	// ]
	// Result: 76
	// Ok: true
}

func calcResult(a *AST) int {
	if a.Type == "Root" {
		res := 0
		for _, v := range a.Args {
			res += calcResult(v)
		}
		return res
	}
	switch a.Name.Text {
	case "+":
		return calcResult(a.Left()) + calcResult(a.Right())
	case "*":
		return calcResult(a.Left()) * calcResult(a.Right())
	default:
		i, _ := strconv.Atoi(a.Name.Text)
		return i
	}
}
