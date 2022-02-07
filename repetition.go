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
	return m.True()
}

// Min matches a minimum number of tokens.
func (t MatcherFunc) Min(n int) MatcherFunc {
	return func(c *Code) bool {
		i := 0
		for t(c) {
			i++
		}
		return i >= n
	}
}

// Until matches until some matcher return true.
func Until(or ...MatcherFunc) MatcherFunc {
	return Or(or...).Not().Next().OneToMany()
}

// While matches while any matcher returns true.
func While(or ...MatcherFunc) MatcherFunc {
	return Or(or...).Next().OneToMany()
}
