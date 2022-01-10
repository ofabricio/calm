package calm

import "fmt"

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
		fmt.Printf("[debug] Match: %v Token: '%s' Pos: %d End: %d\n", okz, tkn.Text, tkn.Pos, end.pos)
		return okz
	}
}
