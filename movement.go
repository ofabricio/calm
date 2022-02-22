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
	return MatcherFunc(func(c *Code) bool {
		if m(c) {
			c.Next()
			return true
		}
		return false
	}).More()
}

// Undo sends the cursor back to the
// beginning of the current matcher
// if it returns false.
func (m MatcherFunc) Undo() MatcherFunc {
	return MatcherFunc(func(c *Code) bool {
		if ini := c.Mark(); !m(c) {
			c.Back(ini)
			return false
		}
		return true
	}).undoAST()
}
