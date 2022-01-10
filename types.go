package calm

// Run implements the Matcher interface.
func (m MatcherFunc) Run(c *Code) bool {
	return m(c)
}

type MatcherFunc func(*Code) bool

type Matcher interface {
	Run(*Code) bool
}

// Mark represents a mark in the code.
type Mark struct {
	pos int
	row int
	col int
}

// Token represents a token of the code.
type Token struct {
	Text string
	Pos  int
	Row  int
	Col  int
}

func (t *rowcol) incRow() {
	t.row++
	t.col = 1
}

func (t *rowcol) incCol() {
	t.col++
}

type rowcol struct {
	row int
	col int
}
