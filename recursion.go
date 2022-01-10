package calm

// Recursive allows recursive calls of the current matcher.
func (m MatcherFunc) Recursive(r *Wrap) MatcherFunc {
	r.mf = m
	return m
}

func (w Wrap) Run(c *Code) bool {
	return w.mf(c)
}

type Wrap struct {
	mf MatcherFunc
}
