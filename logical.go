package calm

// Or tests each matcher and returns
// true if one of them return true.
func Or(ms ...MatcherFunc) MatcherFunc {
	return func(c *Code) bool {
		for _, m := range ms {
			if m(c) {
				return true
			}
		}
		return false
	}
}

// And tests each matcher and returns
// true if all of them return true.
func And(ms ...MatcherFunc) MatcherFunc {
	return func(c *Code) bool {
		for _, m := range ms {
			if !m(c) {
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

// If runs 'then' if 'cond' is true or 'elze' if 'cond' is false.
func If(cond MatcherFunc, then MatcherFunc, elze MatcherFunc) MatcherFunc {
	return func(c *Code) bool {
		if cond(c) {
			return then(c)
		}
		return elze(c)
	}
}
