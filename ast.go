package calm

import (
	"fmt"
	"strings"
)

// Tree grabs a node by its parent.
func (m MatcherFunc) Tree(a *Ast) MatcherFunc {
	return func(c *Code) bool {
		if m(c) {
			*a = *c.ast
			return true
		}
		return false
	}
}

// Leaf builds a leaf node.
func (m MatcherFunc) Leaf(Type string) MatcherFunc {
	return func(c *Code) bool {
		if ini := c.Mark(); m(c) {
			leaf := &Ast{Type: Type, Name: c.Token(ini, c.Mark())}
			c.ast.Args = append(c.ast.Args, leaf)
			return true
		}
		return false
	}
}

// Root converts three Leaf nodes into a binary branch.
// In other words, it will set the middle node as a root
// node and add the left and right nodes as its children.
// Example: [ 2, +, 4 ] becomes [ + [ 2, 4 ] ].
func (m MatcherFunc) Root() MatcherFunc {
	return func(c *Code) bool {
		parent := c.ast
		if m(c) {
			a := parent.Args[len(parent.Args)-3]
			o := parent.Args[len(parent.Args)-2]
			b := parent.Args[len(parent.Args)-1]
			o.Args = append(o.Args, a, b)
			parent.Args = append(parent.Args[:len(parent.Args)-3], o)
			return true
		}
		return false
	}
}

// Enter makes the selected node a root node,
// so the next nodes are added as its children.
// In other words it increases the depth of the
// tree in that node. Enter is usually used on
// a Leaf node, for example .Leaf("Func").Enter()
func (m MatcherFunc) Enter() MatcherFunc {
	return func(c *Code) bool {
		if m(c) {
			parent := c.ast
			c.ast = parent.Right()
			return true
		}
		return false
	}
}

// Child makes nodes children of a node.
func (m MatcherFunc) Child(ms ...MatcherFunc) MatcherFunc {
	return And(m.Enter(), And(ms...))
}

// Leave is the opposite of Enter. Useful
// to restore an AST depth. Make sure to
// Leave to a parent node.
func (m MatcherFunc) Leave() MatcherFunc {
	return func(c *Code) bool {
		parent := c.ast
		ok := m(c)
		c.ast = parent
		return ok
	}
}

// Group groups AST nodes inside an AST node.
func (m MatcherFunc) Group(Type string) MatcherFunc {
	return MatcherFunc(func(c *Code) bool {
		parent := c.ast
		group := &Ast{Type: Type}
		c.ast = group
		if m(c) {
			parent.Args = append(parent.Args, group)
			return true
		}
		return false
	}).Leave()
}

// undoAst sends the AST back to the
// beginning of the current matcher
// if it returns false.
func (m MatcherFunc) undoAst() MatcherFunc {
	return MatcherFunc(func(c *Code) bool {
		// We create a new "Undo" node that is
		// used to add the next nodes.
		// On both conditions (true or false) we
		// discard "Undo" node with Leave(), but
		// the children nodes of "Undo" must be
		// copied to the parent node when true.
		parent := c.ast
		undo := &Ast{Type: "Undo"}
		c.ast = undo
		if m(c) {
			parent.Args = append(parent.Args, undo.Args...)
			return true
		}
		return false
	}).Leave()
}

type Ast struct {
	Type string
	Name Token
	Args []*Ast
}

// Left returns the leftmost node.
func (a *Ast) Left() *Ast {
	return a.Args[0]
}

// Right returns the rightmost node.
func (a *Ast) Right() *Ast {
	return a.Args[len(a.Args)-1]
}

// String returns a JSON string representation of the AST.
func (a Ast) String() string {
	var name string
	var args string
	if a.Name.Text != "" {
		name = fmt.Sprintf(`, "name": "%s"`, a.Name.Text)
	}
	if len(a.Args) != 0 {
		var argz []string
		for _, n := range a.Args {
			argz = append(argz, n.String())
		}
		args = fmt.Sprintf(`, "args": [%s]`, strings.Join(argz, ", "))
	}
	return fmt.Sprintf(`{ "type": "%s"%s%s }`, a.Type, name, args)
}

// Print returns a short string representation of the AST.
func (a *Ast) Print() string {
	var name string
	var args string
	if a.Name.Text != "" {
		name = " " + a.Name.Text
	}
	if len(a.Args) != 0 {
		var argz []string
		for _, n := range a.Args {
			argz = append(argz, n.Print())
		}
		args = fmt.Sprintf(` [ %s ]`, strings.Join(argz, ", "))
	}
	return fmt.Sprintf(`%s%s%s`, a.Type, name, args)
}

// Walk traverses an AST.
func Walk(v Visitor, node *Ast) {
	if v = v.Visit(node); v == nil {
		return
	}
	for _, n := range node.Args {
		Walk(v, n)
	}
	v.Visit(nil)
}

type Visitor interface {
	Visit(*Ast) Visitor
}
