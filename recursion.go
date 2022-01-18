package calm

// Recursive allows recursive call of a matcher.
func Recursive() (ref MatcherFunc, set func(MatcherFunc) MatcherFunc) {
	var m MatcherFunc
	ref = func(c *Code) bool {
		return m(c)
	}
	set = func(mf MatcherFunc) MatcherFunc {
		m = mf
		return mf
	}
	return
}
