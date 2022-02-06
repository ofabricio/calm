package calm

// Or tests each matcher and returns
// true if one of them return true.
func Or(ms ...Matcher) MatcherFunc {
	return func(c *Code) bool {
		for _, m := range ms {
			if m.Run(c) {
				return true
			}
		}
		return false
	}
}

// And tests each matcher and returns
// true if all of them return true.
func And(ms ...Matcher) MatcherFunc {
	return func(c *Code) bool {
		for _, m := range ms {
			if !m.Run(c) {
				return false
			}
		}
		return true
	}
}

// Not negates the current matcher. True
// becomes false and false becomes true.
func (m MatcherFunc) Not() MatcherFunc {
	return func(c *Code) bool {
		return !m(c)
	}
}

// True forces the current matcher to return true.
func (m MatcherFunc) True() MatcherFunc {
	return func(c *Code) bool {
		return m(c) || true
	}
}

// False forces the current matcher to return false.
func (m MatcherFunc) False() MatcherFunc {
	return func(c *Code) bool {
		return m(c) && false
	}
}

// True forces the current matcher to return true.
func True() MatcherFunc {
	return func(c *Code) bool {
		return true
	}
}

// False forces the current matcher to return false.
func False() MatcherFunc {
	return func(c *Code) bool {
		return false
	}
}
