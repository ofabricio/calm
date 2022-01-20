package calm

import "strconv"

// On calls f with the current token
// when the current matcher returns true.
func (m MatcherFunc) On(f func(Token)) MatcherFunc {
	return func(c *Code) bool {
		if ini := c.Mark(); m(c) {
			token := c.Token(ini, c.Mark())
			f(token)
			return true
		}
		return false
	}
}

// Emit captures the current token.
func Emit(t *Token) func(Token) {
	return func(tk Token) {
		*t = tk
	}
}

// Emits captures the current token and adds it to a slice.
func Emits(ts *[]Token) func(Token) {
	return func(t Token) {
		*ts = append(*ts, t)
	}
}

// Grab captures the current token.
func Grab(s *string) func(Token) {
	return func(t Token) {
		*s = t.Text
	}
}

// Grabs captures the current token and adds it to a slice.
func Grabs(s *[]string) func(Token) {
	return func(t Token) {
		*s = append(*s, t.Text)
	}
}

// Index captures the current token position.
func Index(i *int) func(Token) {
	return func(t Token) {
		*i = t.Pos
	}
}

// Indexes captures the current token position and adds it to a slice.
func Indexes(is *[]int) func(Token) {
	return func(t Token) {
		*is = append(*is, t.Pos)
	}
}

// ToInt captures the current token and converts it to integer.
func ToInt(v *int) func(Token) {
	return func(t Token) {
		*v, _ = strconv.Atoi(t.Text)
	}
}

// ToFloat captures the current token and converts it to float.
func ToFloat(v *float64) func(Token) {
	return func(t Token) {
		*v, _ = strconv.ParseFloat(t.Text, 64)
	}
}
