package example

import (
	"fmt"
	"strconv"
	"unicode"

	. "github.com/ofabricio/calm"
)

func ExampleAST() {

	src := New("6+5*(4+3)*2")

	var term, expr, factor MatcherFunc

	value := F(unicode.IsNumber).Leaf("Value")

	factor = func(c *Code) bool {
		return Or(And(S("("), expr, S(")")), value).Run(c)
	}

	term = func(c *Code) bool {
		return Or(And(factor, S("*").Leaf("BinExpr"), term).Root(), factor).Run(c)
	}

	expr = func(c *Code) bool {
		return Or(And(term, S("+").Leaf("BinExpr"), expr).Root(), term).Run(c)
	}

	// When.

	var ast AST
	ok := MatcherFunc(expr).Tree(&ast).Run(src)

	fmt.Println(ast.Print("short"))
	fmt.Println("Result:", calcResult(&ast))
	fmt.Println("Ok:", ok)

	// Output:
	// Root [
	//     BinExpr + [
	//         Value 6
	//         BinExpr * [
	//             Value 5
	//             BinExpr * [
	//                 BinExpr + [
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
