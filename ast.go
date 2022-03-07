package calm

import (
	"bytes"
)

// Tree grabs a node by its parent.
func (m MatcherFunc) Tree(a *AST) MatcherFunc {
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
			leaf := &AST{Type: Type, Name: c.Token(ini, c.Mark())}
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
func Root(left, root, right MatcherFunc) MatcherFunc {
	m := AND(left, root, right)
	return func(c *Code) bool {
		parent := c.ast
		if m(c) {
			if len(parent.Args) >= 3 {
				args := parent.Args[len(parent.Args)-3:]
				a, o, b := args[0], args[1], args[2]
				o.Args = append(o.Args, a, b)
				parent.Args = append(parent.Args[:len(parent.Args)-3], o)
			}
			return true
		}
		return false
	}
}

// enter makes the selected node a root node,
// so the next nodes are added as its children.
// In other words it increases the depth of the
// tree in that node. enter is usually used on
// a Leaf node, for example .Leaf("Func").enter()
func (m MatcherFunc) enter() MatcherFunc {
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
	return And(m.enter(), And(ms...)).leave()
}

// leave is the opposite of Enter. Useful
// to restore an AST depth. Make sure to
// leave to a parent node.
func (m MatcherFunc) leave() MatcherFunc {
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
		group := &AST{Type: Type}
		c.ast = group
		if m(c) {
			parent.Args = append(parent.Args, group)
			return true
		}
		return false
	}).leave()
}

// undoAST sends the AST back to the
// beginning of the current matcher
// if it returns false.
func (m MatcherFunc) undoAST() MatcherFunc {
	return MatcherFunc(func(c *Code) bool {
		// We create a new "Undo" node that is
		// used to add the next nodes.
		// On both conditions (true or false) we
		// discard "Undo" node with Leave(), but
		// the children nodes of "Undo" must be
		// copied to the parent node when true.
		parent := c.ast
		undo := &AST{Type: "Undo"}
		c.ast = undo
		if m(c) {
			parent.Args = append(parent.Args, undo.Args...)
			return true
		}
		return false
	}).leave()
}

type AST struct {
	Type string
	Name Token
	Args []*AST
}

// Left returns the leftmost node.
func (a *AST) Left() *AST {
	return a.Args[0]
}

// Right returns the rightmost node.
func (a *AST) Right() *AST {
	return a.Args[len(a.Args)-1]
}

// String returns a JSON string representation
// of the AST.
func (a AST) String() string {
	return a.Print("json-inline")
}

// Print returns a string representation
// of the AST given a format.
func (a *AST) Print(format string) string {
	var buf bytes.Buffer
	PrintTree(&buf, format, a)
	return buf.String()
}

// Walk traverses an AST.
func Walk(v Visitor, node *AST) {
	if v = v.Visit(node); v == nil {
		return
	}
	for _, n := range node.Args {
		Walk(v, n)
	}
	v.Visit(nil)
}

type Visitor interface {
	Visit(*AST) Visitor
}
