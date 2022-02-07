package calm

import (
	"fmt"
	"unicode"
)

// String matches a string text given a quote.
func String(quote string) MatcherFunc {
	return And(
		S(quote),
		Until(S(`\`+quote).False(), Eq(quote), Eq("\n")).True(),
		S(quote),
	)
}

// Debug prints debug info to the stdout.
func (m MatcherFunc) Debug() MatcherFunc {
	return func(c *Code) bool {
		ini := c.Mark()
		okz := m(c)
		end := c.Mark()
		tkn := c.Token(ini, end)
		fmt.Printf("[debug] Match: %-5t Token: %-3s Pos: %d Row: %d Col: %d\n", okz, "'"+tkn.Text+"'", tkn.Pos, tkn.Row, tkn.Col)
		return okz
	}
}

// Json matches a Json.
func Json() MatcherFunc {
	// BNF from https://www.json.org
	ws := F(unicode.IsSpace).ZeroToMany()
	value := MatcherFunc(func(c *Code) bool {
		return And(ws, Or(S("true"), S("false"), S("null"), Number(), String("\""), Json()), ws).Run(c)
	})
	objField := And(ws, String("\""), ws, S(":"), value)
	emptyObj := And(S("{"), ws, S("}")).Undo()
	emptyArr := And(S("["), ws, S("]")).Undo()
	obj := And(S("{"), objField, And(S(","), objField).ZeroToMany(), S("}"))
	arr := And(S("["), value, And(S(","), value).ZeroToMany(), S("]"))
	return Or(emptyObj, obj, emptyArr, arr)
}

// Number matches numbers.
func Number() MatcherFunc {
	integer := And(S("-").ZeroToOne(), F(unicode.IsDigit).OneToMany())
	sign := SOr("+-").ZeroToOne()
	digits := F(unicode.IsDigit).OneToMany()
	exponent := If(SOr("Ee"), And(sign, digits), True())
	fraction := If(S("."), digits, True())
	return And(integer, fraction, exponent)
}
