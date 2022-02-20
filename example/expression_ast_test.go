package example

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"

	. "github.com/ofabricio/calm"
)

func ExampleAst() {

	src := New("6+5*(4+3)*2")

	var term, expr, factor func(c *Code) bool

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

	var ast Ast
	ok := MatcherFunc(expr).Tree(&ast).Run(src)

	fmt.Println(printTree(&ast, 0))
	fmt.Println("Result:", calcResult(&ast))
	fmt.Println("Ok:", ok)

	// Output:
	// -- +
	// ---- 6
	// ---- *
	// ------ 5
	// ------ *
	// -------- +
	// ---------- 4
	// ---------- 3
	// -------- 2
	// Result: 76
	// Ok: true
}

func printTree(a *Ast, pad int) string {
	if a == nil {
		return ""
	}
	var args string
	for _, v := range a.Args {
		args += "\n" + printTree(v, pad+2)
	}
	return strings.Repeat("-", pad) + " " + a.Name.Text + args
}

func calcResult(a *Ast) int {
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
