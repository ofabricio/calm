package example

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"

	. "github.com/ofabricio/calm"
)

func ExampleExpression() {

	src := "6+5*(4+3)*2"

	c := New(src)

	exp := &Expression{}

	ok := c.Run(expr(exp))

	exp.print(0)
	fmt.Println("Result:", exp.result())
	fmt.Println("Ok:", ok)

	// Output:
	// +
	// -- 6
	// -- *
	// ---- 5
	// ---- *
	// ------ +
	// -------- 4
	// -------- 3
	// ------ 2
	// Result: 76
	// Ok: true
}

type Expression struct {
	V string
	L *Expression
	R *Expression
}

func expr(e *Expression) MatcherFunc {
	return func(c *Code) bool {
		l := &Expression{}
		r := &Expression{}
		return Or(And(term(l), S("+"), expr(r)).Rewind().On(func(Token) {
			e.L, e.R, e.V = l, r, "+"
		}), term(e)).Run(c)
	}
}

func term(e *Expression) MatcherFunc {
	l := &Expression{}
	r := &Expression{}
	return Or(And(factor(l), S("*"), expr(r)).Rewind().On(func(Token) {
		e.L, e.R, e.V = l, r, "*"
	}), factor(e))
}

func factor(e *Expression) MatcherFunc {
	return Or(And(S("("), expr(e), S(")")), value(e))
}

func value(e *Expression) MatcherFunc {
	return F(unicode.IsNumber).Grab(&e.V)
}

func (e *Expression) print(pad int) {
	if e == nil {
		return
	}
	fmt.Println(strings.Repeat("-", pad), e.V)
	e.L.print(pad + 2)
	e.R.print(pad + 2)
}

func (e *Expression) result() int {
	if e == nil {
		return 0
	}
	switch e.V {
	case "+":
		return e.L.result() + e.R.result()
	case "*":
		return e.L.result() * e.R.result()
	default:
		i, _ := strconv.Atoi(e.V)
		return i
	}
}
