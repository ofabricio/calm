package calm

// ZeroToMany matches zero or many tokens.
// It is equivalent to the regex '*' symbol.
func (m MatcherFunc) ZeroToMany() MatcherFunc {
	return m.Min(0)
}

// OneToMany matches one or many tokens.
// It is equivalent to the regex '+' symbol.
func (m MatcherFunc) OneToMany() MatcherFunc {
	return m.Min(1)
}

// ZeroToOne matches an optional token.
// It is equivalent to the regex '?' symbol.
func (m MatcherFunc) ZeroToOne() MatcherFunc {
	return m.True().More()
}

// Min matches a minimum number of tokens.
func (t MatcherFunc) Min(n int) MatcherFunc {
	return MatcherFunc(func(c *Code) bool {
		i := 0
		for t(c) {
			i++
		}
		return i >= n
	}).More()
}

// Until matches until some matcher return true.
// Note that it only advances the position by one character,
// so be careful when using matchers with more than one
// character like S("abc").
func Until(or ...Matcher) MatcherFunc {
	return Or(or...).Not().Next().OneToMany()
}

// While matches while any matcher returns true.
// Note that it only advances the position by one character,
// so be careful when using matchers with more than one
// character like S("abc").
func While(or ...Matcher) MatcherFunc {
	return Or(or...).Next().OneToMany()
}
