package calm

// Here tells the current position.
func (c *core) Here() int {
	return c.pos
}

// Tail returns the content from the
// current position to the end.
func (c *core) Tail() string {
	return c.src[c.pos:]
}

// Take returns the content from position a to b.
func (c *core) Take(a, b int) string {
	return c.src[a:b]
}

// More tells if there are more characters to scan.
func (c *core) More() bool {
	return c.pos < len(c.src)
}

// move moves the position to a new place.
func (c *core) move(pos int) {
	if pos > len(c.src) {
		pos = len(c.src)
	}
	c.pos = pos
}

type core struct {
	src string // Source code.
	pos int    // Position.
}
