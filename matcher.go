package calm

// S tests if the current token matches a
// string and moves the position if true.
func S(s string) MatcherFunc {
	return func(c *Code) bool {
		return c.Match(s)
	}
}

// S tests if the current token matches a string
// reference and moves the position if true.
func SR(s *string) MatcherFunc {
	return func(c *Code) bool {
		return c.Match(*s)
	}
}

// F tests the current character against a rune
// function and moves the position if true.
func F(fn func(rune) bool) MatcherFunc {
	return MatcherFunc(func(c *Code) bool {
		return c.MatchF(fn)
	}).More()
}

// Eq tests if the current token equals a
// string, but does not move the position.
func Eq(s string) MatcherFunc {
	return func(c *Code) bool {
		return c.Equal(s)
	}
}

// More runs the current matcher only if
// there are more characters to match.
func (m MatcherFunc) More() MatcherFunc {
	return func(c *Code) bool {
		if c.More() {
			return m(c)
		}
		return false
	}
}

// Run implements the Matcher interface.
func (m MatcherFunc) Run(c *Code) bool {
	return m(c)
}

type MatcherFunc func(*Code) bool

type Matcher interface {
	Run(*Code) bool
}
