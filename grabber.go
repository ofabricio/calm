package calm

// Emit captures the current token and fills t with it.
func (m MatcherFunc) Emit(t *Token) MatcherFunc {
	return m.On(func(tk Token) {
		*t = tk
	})
}

// EmitMany captures the current token and adds it to a slice.
// Note that tokens might repeat depending on the logic.
func (m MatcherFunc) EmitMany(ts *[]Token) MatcherFunc {
	return m.On(func(t Token) {
		*ts = append(*ts, t)
	})
}

func (m MatcherFunc) EmitUndo(ts *[]Token) MatcherFunc {
	return func(c *Code) bool {
		if ini := len(*ts); m(c) {
			*ts = (*ts)[0:ini]
			return true
		}
		return false
	}
}

// Grab captures the current token and fills s with it.
func (m MatcherFunc) Grab(s *string) MatcherFunc {
	return m.On(func(t Token) {
		*s = t.Text
	})
}

// GrabMany captures the current token and adds it to a slice.
// Note that tokens might repeat depending on the logic.
func (m MatcherFunc) GrabMany(s *[]string) MatcherFunc {
	return m.On(func(t Token) {
		*s = append(*s, t.Text)
	})
}

func (m MatcherFunc) GrabUndo(ts *[]string) MatcherFunc {
	return func(c *Code) bool {
		if ini := len(*ts); m(c) {
			*ts = (*ts)[0:ini]
			return true
		}
		return false
	}
}

// Grab captures the current token position and fills p with it.
func (m MatcherFunc) GrabPos(p *int) MatcherFunc {
	return m.On(func(t Token) {
		*p = t.Pos
	})
}

// On calls f with the matched token
// when the current matcher matches.
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
