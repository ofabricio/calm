package calm

import (
	"fmt"
	"unicode"
)

// String matches a string text given a quote.
func String(quote string) MatcherFunc {
	return AND(
		S(quote),
		Until(S(`\`+quote).False(), Eq(quote), Eq("\n")).True(),
		S(quote),
	)
}

// Tag matches a tag.
func Tag(open, close string) MatcherFunc {
	tag, setTag := Recursive()
	body := Or(Until(Eq(open), Eq(close)), tag)
	return setTag(AND(S(open), body.ZeroToMany(), S(close)))
}

// Json matches a Json.
func Json() MatcherFunc {
	// BNF from https://www.json.org
	jsn, setJsn := Recursive()
	wz := F(unicode.IsSpace).ZeroToMany()
	value := And(wz, Or(S("true"), S("false"), S("null"), Number(), String("\""), jsn), wz)
	objField := And(wz, String("\""), wz, S(":"), value)
	emptyObj := AND(S("{"), wz, S("}"))
	emptyArr := AND(S("["), wz, S("]"))
	obj := AND(S("{"), objField, AND(S(","), objField).ZeroToMany(), S("}"))
	arr := AND(S("["), value, AND(S(","), value).ZeroToMany(), S("]"))
	return setJsn(Or(emptyObj, obj, emptyArr, arr))
}

// Number matches a number.
func Number() MatcherFunc {
	digits := F(unicode.IsDigit).OneToMany()
	integer := And(S("-").ZeroToOne(), digits)
	sign := SOr("+-").ZeroToOne()
	exponent := If(SOr("Ee"), And(sign, digits), True())
	fraction := If(S("."), digits, True())
	return AND(integer, fraction, exponent)
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

// Scan scans the input from start to end.
func (m MatcherFunc) Scan(c *Code) bool {
	return Or(m, Next()).OneToMany().Run(c)
}
