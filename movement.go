package calm

// Next moves to the next character.
func Next() MatcherFunc {
	return MatcherFunc(func(c *Code) bool {
		c.Next()
		return true
	}).More()
}

// Next moves to the next character when
// the current matcher returns true.
func (m MatcherFunc) Next() MatcherFunc {
	return And(m, Next())
}

// Rewind sends the cursor back to the
// begining of the current matcher if
// it returns false.
func (m MatcherFunc) Rewind() MatcherFunc {
	return m.Not().Undo().Not()
}

// Undo sends the cursor back to the
// begining of the current matcher
// if it returns true.
func (m MatcherFunc) Undo() MatcherFunc {
	return func(c *Code) bool {
		if ini := c.Mark(); m(c) {
			c.Back(ini)
			return true
		}
		return false
	}
}
