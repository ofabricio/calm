# Calm

This is a library to scan, parse and tokenize text.

[![Go](https://github.com/ofabricio/calm/actions/workflows/go.yml/badge.svg)](https://github.com/ofabricio/calm/actions/workflows/go.yml)

### Note

Use [issues](https://github.com/ofabricio/calm/issues) only to report bugs.
Use [discussions](https://github.com/ofabricio/calm/discussions) for everything else.
Don't open PRs.

## Examples

### Tokenizer

Example of a very simple tokenizer for Go code.

```go
package main

import (
    "fmt"
    "unicode"
    . "github.com/ofabricio/calm"
)

func main() {

    code := New(`
    package main

    import "fmt"

    func main() {
        fmt.Println("Hello, 世界")
    }`)

    print := func(t Token) {
        fmt.Printf("Pos: %-3d Line: %-2d Column: %-3d Token: %s\n", t.Pos, t.Row, t.Col, t.Text)
    }

    spac := F(unicode.IsSpace)
    strg := String(`"`).On(print)
    word := F(unicode.IsLetter).OneToMany().On(print)
    rest := Next().On(print)

    root := Or(spac, strg, word, rest).ZeroToMany()

    ok := root.Run(code)

    fmt.Println(ok)

    // Output:
    // Pos: 2   Line: 2  Column: 2   Token: package
    // Pos: 10  Line: 2  Column: 10  Token: main
    // Pos: 17  Line: 4  Column: 2   Token: import
    // Pos: 24  Line: 4  Column: 9   Token: "fmt"
    // Pos: 32  Line: 6  Column: 2   Token: func
    // Pos: 37  Line: 6  Column: 7   Token: main
    // Pos: 41  Line: 6  Column: 11  Token: (
    // Pos: 42  Line: 6  Column: 12  Token: )
    // Pos: 44  Line: 6  Column: 14  Token: {
    // Pos: 48  Line: 7  Column: 3   Token: fmt
    // Pos: 51  Line: 7  Column: 6   Token: .
    // Pos: 52  Line: 7  Column: 7   Token: Println
    // Pos: 59  Line: 7  Column: 14  Token: (
    // Pos: 60  Line: 7  Column: 15  Token: "Hello, 世界"
    // Pos: 75  Line: 7  Column: 26  Token: )
    // Pos: 78  Line: 8  Column: 2   Token: }
    // true
}
```

### Parser

Example of a very simple parser for Go code.

```go
package main

import (
    "fmt"
    "unicode"
    . "github.com/ofabricio/calm"
)

func main() {

    src := New(`
    // You can edit this code!
    // Click here and start typing.
    package main

    import "fmt"

    func main() {
        fmt.Println("Hello, 世界")
    }`)

    ws := SOr(" \t").OneToMany()
    wz := F(unicode.IsSpace).ZeroToMany()
    name := F(unicode.IsLetter).OneToMany()
    strg := String(`"`).Leaf("Str")

    comment := And(S("//"), Until(Eq("\n"))).Leaf("Comment")
    pkgDef := S("package").Leaf("Pkg").Child( ws, name.Leaf("Name"))
    impDef := S("import").Leaf("Imp").Child( ws, strg)

    fnCall := And(name.Leaf("Pkg"), S("."), name.Leaf("Name"), S("("), strg, S(")")).Group("Call")
    fnBody := Or(wz.False(), fnCall, wz.False()).ZeroToMany().Group("Body")
    fnDef := S("func").Leaf("Fun").Child( ws, name.Leaf("Name"), wz, S("()"), wz,
        S("{"), fnBody, S("}"))

    root := Or(
        wz.False(),
        comment,
        pkgDef,
        impDef,
        fnDef,
    ).ZeroToMany()

    var ast AST
    ok := root.Tree(&ast).Run(src)

    fmt.Println("Ok:", ok)
    fmt.Println(ast.Print("short"))

    // Output:
    // Ok: true
    // Root [
    //     Comment // You can edit this code!
    //     Comment // Click here and start typing.
    //     Pkg package [
    //         Name main
    //     ]
    //     Imp import [
    //         Str "fmt"
    //     ]
    //     Fun func [
    //         Name main
    //         Body [
    //             Call [
    //                 Pkg fmt
    //                 Name Println
    //                 Str "Hello, 世界"
    //             ]
    //         ]
    //     ]
    // ]
}
```

See more examples in the [example](/example) folder.

## How it works

When a match happens the cursor moves to the next character.

If there is no match the cursor doesn't move.

There are many kinds of operators: conditional execution, repetition, recursion, etc.

Not all operators move the cursor.

Some operators can move the cursor back and forth.

### Deeper look

Suppose the code below and the input `HelloWorld`:

```go
And(S("Hello"), S("World"))
```

The starting state is this:

```
HelloWorld
^
```

`And` operator will test its arguments applying an AND logic to them.
The first argument is `S("Hello")`.
`S` will test `Hello` against the first characters of the input,
if they match it moves the cursor and returns true:

```
HelloWorld
     ^
```

`And` saw that the first argument returned true so it moves to the next argument `S("World")`.
`S` will test `World` against the characters in the current position of the input,
if they match it moves the cursor and returns true:

```
HelloWorld
          ^
```

If there was no match the cursor would stay on the `W` character and `S` would return false.

`And` sees that both arguments returned true, so it returns true and we are at the end of the input.

## Operators

#### Matcher

- [x] [S](#S)
- [x] [SR](#SR)
- [x] [SOr](#SOr)
- [x] [F](#F)
- [x] [R](#R)

#### Tester

- [x] [Eq](#Eq)
- [x] [EqF](#EqF)
- [x] [More](#More)

#### Logical

- [x] [Or](#Or)
- [x] [And](#And)
- [x] [Not](#Not)
- [x] [True](#True)
- [x] [False](#False)
- [x] [If](#If)

#### Repetition

- [x] [ZeroToMany](#ZeroToMany)
- [x] [OneToMany](#OneToMany)
- [x] [ZeroToOne](#ZeroToOne)
- [x] [Min](#Min)
- [x] [Until](#Until)
- [x] [While](#While)

#### Movement

- [x] [Next](#Next)
- [x] [Undo](#Undo)

#### Event

- [x] [On](#On)

#### Grabber

- [x] [Grab](#Grab)
- [x] [Grabs](#Grabs)
- [x] [Emit](#Emit)
- [x] [Emits](#Emits)
- [x] [Index](#Index)
- [x] [Indexes](#Indexes)
- [x] [ToInt](#ToInt)
- [x] [ToFloat](#ToFloat)
- [x] [Back Reference](#Back-Reference)

#### Util

- [x] [Scan](#Scan)
- [x] [Debug](#Debug)
- [x] [String](#String)
- [x] [Number](#Number)
- [x] [Json](#Json)
- [x] [Tag](#Tag)

#### Recursion

- [x] [Recursive](#Recursive)

#### Tree ([AST](#AST))

- [x] [Tree](#Tree)
- [x] [Leaf](#Leaf)
- [x] [Root](#Root)
- [x] [Enter](#Enter)
- [x] [Leave](#Leave)
- [x] [Child](#Child)
- [x] [Group](#Group)

### S

S tests if the current token matches a string and moves the position if true.

```go
m := New("hello world")

a := S("hello").Run(m)
b := S(" ").Run(m)
c := S("hello").Run(m)
d := S("world").Run(m)

fmt.Println(a, b, c, d) // true true false true
```

### SR

SR behaves exactly the same as [S](#S) but receives a string reference as argument.
This allows [Back Reference](#Back-Reference).

### SOr

SOr tests if the current token matches any character of the string
and moves the position if true.

```go
a := SOr("1a").Run(New("1"))
b := SOr("1a").Run(New("a"))
c := SOr("1a").Run(New("@"))

fmt.Println(a, b, c) // true true false
```

Note that `SOr("ab")` is exactly the same as `Or(S("a"), S("b"))`.

### F

F tests if the current character matches a rune function and moves the position if true.

```go
m := New("hi5")

a := F(unicode.IsLetter).Run(m) // h
b := F(unicode.IsLetter).Run(m) // i
c := F(unicode.IsLetter).Run(m) // 5

fmt.Println(a, b, c) // true true false
```

### R

R tests if the current token matches a regular expression and moves the position if true.

```go
c := New("hi5")

a := R("hi\\d").Run(c)

fmt.Println(a) // true
```

### Eq

Eq tests if the current token matches a string. It does not move the cursor.

```go
c := New("hello")

a := Eq("hello").Run(c)
b := Eq("hello").Run(c)

fmt.Println(a, b) // true true
```

### EqF

EqF tests the current character against a rune function. It does not move the position.

```go
c := New("H1")

a := EqF(unicode.IsLetter).Run(c)
b := S("H1").Run(c)

fmt.Println(a, b) // true true
```

### More

More runs the current operator only if there are more characters to match (so the cursor can move).

It is used to prevent overflow in some operations.

### Or

Or returns true if one of its arguments return true.

```go
c := New("apple")

ok := Or(S("grape"), S("apple")).Run(c)

fmt.Println(ok) // true
```

### And

And returns true if all of its arguments return true.

```go
c := New("hello world")

ok := And(S("hello"), S(" "), S("world")).Run(c)

fmt.Println(ok) // true
```

### Not

Not negates the current operator. True becomes false and vice-versa.

```go
c := New("hello")

ok := S("hello").Not().Run(c)

fmt.Println(ok) // false
```

Note that the cursor advances even though it returns false.

### True

True forces the current operator to return true.

```go
c := New("hello")

ok := S("world").True().Run(c)

fmt.Println(ok) // true
```

### False

False forces the current operator to return false.

```go
c := New("hello")

ok := S("hello").False().Run(c)

fmt.Println(ok) // false
```

### If

If runs the second argument if the first argument returns true
or runs the third argument if the first argument returns false.

```go
m := If(S("."), S("1"), S("0"))

a := m.Run(New(".1"))
b := m.Run(New("0"))

fmt.Println(a, b) // true true
```

### ZeroToMany

ZeroToMany matches zero to many tokens. It is equivalent to the regex symbol `*`.

```go
m := And(S("hello"), S(" ").ZeroToMany(), S("world"))

a := m.Run(New("helloworld"))
b := m.Run(New("hello world"))
c := m.Run(New("hello  world"))

fmt.Println(a, b, c) // true true true
```

### OneToMany

OneToMany matches one to many tokens. It is equivalent to the regex symbol `+`.

```go
m := And(S("hello"), S(" ").OneToMany(), S("world"))

a := m.Run(New("helloworld"))
b := m.Run(New("hello world"))
c := m.Run(New("hello  world"))

fmt.Println(a, b, c) // false true true
```

### ZeroToOne

ZeroToOne matches an optional token. It is equivalent to the regex symbol `?`.

```go
m := And(S("colo"), S("u").ZeroToOne(), S("r"))

a := m.Run(New("color"))
b := m.Run(New("colour"))

fmt.Println(a, b) // true true
```

### Min

Min matches a minimum number of tokens.

```go
m := S("a").Min(2)

a := m.Run(New("a"))
b := m.Run(New("aa"))
c := m.Run(New("aaa"))

fmt.Println(a, b, c) // false true true
```

### Until

Until matches until some matcher return true.

```go
m := Until(Eq(","), Eq("."))

a := m.Run(New(","))
b := m.Run(New("ab,"))
c := m.Run(New("abcd."))

fmt.Println(a, b, c) // false true true
```

### While

While matches while any matcher returns true.

```go
m := While(Eq("0"), Eq("1"))

a := m.Run(New("01100"))
b := m.Run(New("hello"))

fmt.Println(a, b) // true false
```

### Next

Next moves to the next character when the current matcher returns true.

```go
c := New("hello world")

ok := And(S("hello").Next(), S("world")).Run(c)

fmt.Println(ok) // true
```

There is also a static version of `Next`.

### Undo

Undo sends the cursor back to the
beginning of the current matcher if it
returns false.

```go
c := New("1+2")

ok := Or(
    And(S("1"), S("*"), S("2")).Undo(),
    And(S("1"), S("+"), S("2")),
).Run(c)

fmt.Println(ok) // true
```

### Grab

Grab captures the current token string.

```go
c := New("abc123")

var t string
F(unicode.IsLetter).OneToMany().On(Grab(&t)).Run(c)

fmt.Println(t) // abc
```

### Grabs

Grabs captures the current token string and adds it to a slice.

```go
c := New("abc123")

var ts []string
F(unicode.IsLetter).On(Grabs(&ts)).OneToMany().Run(c)

fmt.Println(ts) // [a b c]
```

### Emit

Emit captures the current token.

```go
c := New("abc123")

var t Token
F(unicode.IsLetter).OneToMany().On(Emit(&t)).Run(c)

fmt.Println(t) // {abc 0 1 1}
```

### Emits

Emits captures the current token and adds it to a slice.

```go
c := New("abc123")

var ts []Token
F(unicode.IsLetter).On(Emits(&ts)).OneToMany().Run(c)

fmt.Println(ts) // [{a 0 1 1} {b 1 1 2} {c 2 1 3}]
```

### Index

Index captures the current token position.

```go
c := New("abc")

var pos int
And(S("a"), S("b").On(Index(&pos)), S("c")).Run(c)

fmt.Println(pos) // 1
```

### Indexes

Index captures the current token position and adds it to a slice.

```go
c := New("abc")

var pos []int
And(S("a"), S("b").On(Indexes(&pos)), S("c").On(Indexes(&pos))).Run(c)

fmt.Println(pos) // [1 2]
```

### ToInt

ToInt captures the current token and converts it to integer.

```go
c := New("123")

var v int
Next().OneToMany().On(ToInt(&v)).Run(c)

fmt.Println(v) // 123
```

### ToFloat

ToFloat captures the current token and converts it to float.

```go
c := New("1.2")

var v float64
Next().OneToMany().On(ToFloat(&v)).Run(c)

fmt.Println(v) // 1.2
```

### On

On calls a function with the current token when the current operator returns true.

```go
c := New("hello")

f := func(t Token) {
    fmt.Println(t) // {hello 0 1 1}
}

S("hello").On(f).Run(c)
```

### Back Reference

It is possible to create a back reference with [Grab](#Grab) + [SR](#SR).

```go
var quote string

m := And(SOr(`"'`).On(Grab(&quote)), F(unicode.IsLetter).OneToMany(), SR(&quote))

a := m.Run(New(`"hello"`))
b := m.Run(New(`'hello'`))
c := m.Run(New(`"hello'`))

fmt.Println(a, b, c) // true true false
```

### Scan

Scan scans the input from start to end.

Unlike `Run`, that stops scanning when a matcher returns false,
`Scan` keeps going up to the end of the input no matter if a
matcher doesn't match.

```go
c := New("Hello, {name}! You have {count} messages!")

var n []string

ok := Tag("{", "}").On(Grabs(&n)).Scan(c)

fmt.Println(ok, n) // true [{name} {count}]
```

### Debug

Debug prints debug info to the stdout.

```go
c := New("Hi")

F(unicode.IsLetter).Debug().OneToMany().Run(c)

// [debug] Match: true  Token: 'H' Pos: 0 Row: 1 Col: 1
// [debug] Match: true  Token: 'i' Pos: 1 Row: 1 Col: 2
// [debug] Match: false Token: ''  Pos: 2 Row: 1 Col: 3
```

### String

String parses a common string definition. It allows quote escaping.

```go
c := New(`They said "Wow!" and "This is cool!" when they saw this.`)

var quotes []string

ok := String(`"`).On(Grabs(&quotes)).Scan(c)

fmt.Println(ok, quotes) // true ["Wow!" "This is cool!"]
```

### Number

Number matches numbers.

```go
c := New("Heard of 3.1415? What about 0, 1, 1, 2, 3 sequence? Isn't 2e3 a cool notation?")

var n []string

ok := Number().On(Grabs(&n)).Scan(c)

fmt.Println(ok, n) // true [3.14159 0 1 1 2 3 2e3]
```

### Json

Json matches a json.

```go
c := New(`Use either { "hello": "world" } or { "foo": "bar" }.`)

var jsons []string

ok := Json().On(Grabs(&jsons)).Scan(c)

fmt.Println(ok, jsons) // true [{ "hello": "world" } { "foo": "bar" }]
```

### Tag

Tag matches a tag.

```go
c := New("Hello, {name}! You have {count} messages!")

var n []string

ok := Tag("{", "}").On(Grabs(&n)).Scan(c)

fmt.Println(ok, n) // true [{name} {count}]
```

### Recursive

Recursive allows a recursive call of a matcher.

```go
code := New("0+1*(2+3)*4")

term, setTerm := Recursive()
expr, setExpr := Recursive()

value := F(unicode.IsNumber)
factor := Or(And(S("("), expr, S(")")), value)
setTerm(Or(And(factor, S("*"), term).Undo(), factor))
setExpr(Or(And(term, S("+"), expr).Undo(), term))

ok := expr.Run(code)

fmt.Println(ok) // true
```

But it's possible to avoid using this operator by using this approach:

```go
code := New("0+1*(2+3)*4")

var term, expr, factor MatcherFunc

value := F(unicode.IsNumber)

factor = func(c *Code) bool {
    return Or(And(S("("), expr, S(")")), value).Run(c)
}

term = func(c *Code) bool {
    return Or(And(factor, S("*"), term).Undo(), factor).Run(c)
}

expr = func(c *Code) bool {
    return Or(And(term, S("+"), expr).Undo(), term).Run(c)
}

ok := expr.Run(code)

fmt.Println(ok) // true
```

This operator is handy, but it can be a pain when capturing tokens.
See [here](example/expression_test.go) an example on how to recursively
parse an expression without using this operator.

But for an easier and more advanced way to capture tokens see the [AST](#AST) section.

## AST

You can parse a text into an AST (Abstract Syntax Tree).

An AST node has three fields.

- The `Type` field is a string to categorize nodes. You provide this information when building a tree.
- The `Name` field is of type `Token` and holds information about a captured token (Text, Line, etc).
- The `Args` field is a slice of children nodes.

A tree always starts with a default root node of type `"Root"`.

A [Leaf](#Leaf) node is the building block of a tree. This operator builds leaf nodes. Without it a tree would have empty nodes.

A tree is only valid if the scanner returns `true`.

#### How it works

Let's take the example below. It has the input text `2+4`.

```go
code := New("2+4")
root := And(S("2"), S("+"), S("4"))

var ast AST
oks := root.Tree(&ast).Run(code)

fmt.Println(oks, ast)
// true { "type": "Root" }
```

If you run the code above the only thing you see is a `Root` node.
That's because we didn't build any [Leaf](#Leaf) node yet.
Let's fix this.

> **Note** that we print the `ast` in JSON format.
But from now on we will print using a shorter, nicer form with the `ast.Print()` function.

```go
root := And(S("2"), S("+").Leaf("Op"), S("4"))
fmt.Println(oks, ast.Print("short-inline"))
// Root [ Op + ]
```

There we go, now we have a Leaf node.
The `Leaf` operator always go with a matcher (`S` in this case). Let's capture the other matchers.

```go
root := And(S("2").Leaf("Val"), S("+").Leaf("Op"), S("4").Leaf("Val"))
// Root [ Val 2, Op +, Val 4 ]
```

As you can see `Leaf` adds nodes in the parent node. Remember there is always a default `Root` node.

This is not a fancy tree, it's just an introduction.
There are more operators to compose an advanced tree.
They are described in their own section below.

### Tree

Tree is where you get the AST result.

```go
code := New("2+4")
root := And(S("2").Leaf("Val"), S("+").Leaf("Op"), S("4").Leaf("Val"))

var ast AST
oks := root.Tree(&ast).Run(code)

fmt.Println(oks, ast.Print("short-inline"))
// Root [ Val 2, Op +, Val 4 ]
```

Note that `Tree` captured the `Root` node along with the `Leaf` nodes.

### Leaf

Leaf is the basic operator to build a node.
It creates a leaf node in the AST.
Without it the tree would have empty nodes.

```go
src := New("2+4")

root := And(S("2").Leaf("Val"), S("+").Leaf("Op"), S("4").Leaf("Val"))

var ast AST
oks := root.Tree(&ast).Run(src)

fmt.Println(oks, ast.Print("short-inline"))
// Root [ Val 2, Op +, Val 4 ]
```

### Root

Root converts three [Leaf](#Leaf) nodes into a binary branch.
In other words, it will set the middle node as a root
node and add the left and right nodes as its children.

```go
src := New("2+4")

root := Root(S("2").Leaf("Val"), S("+").Leaf("Op"), S("4").Leaf("Val"))

var ast AST
oks := root.Tree(&ast).Run(src)

fmt.Println(oks, ast.Print("short-inline"))
// Root [ Op + [ Val 2, Val 4 ] ]
```

> Note to self: maybe this operator needs a better name.

### Enter

Enter makes the selected [Leaf](#Leaf) node a root node,
so the next nodes are added as its children.
In other words it increases the depth of the tree in that node.
It is usually used on a [Leaf](#Leaf) node, for example `.Leaf("Func").Enter()`.

Remember to always [Leave](#Leave) after an `Enter`.

```go
src := New("print(1)")

root := And(S("print").Leaf("FnCall").Enter(), S("("), S("1").Leaf("Val"), S(")")).Leave()

var ast AST
oks := root.Tree(&ast).Run(src)

fmt.Println(oks, ast.Print("short-inline"))
// Root [ FnCall print [ Val 1 ] ]
```

> Note to self: maybe this operator needs a better name.

### Leave

Leave is the opposite of [Enter](#Enter).
Useful to restore an AST depth.
Leave should be used in a parent node.

```go
src := New("abcd")

cod := And(
    And(S("a").Leaf("Char").Enter(), S("b").Leaf("Char")).Leave(),
    And(S("c").Leaf("Char").Enter(), S("d").Leaf("Char")).Leave(),
)

var ast AST
oks := cod.Tree(&ast).Run(src)

fmt.Println(oks, ast.Print("short-inline"))
// Root [ Char a [ Char b ], Char c [ Char d ] ]
```

> Note to self: maybe this operator needs a better name.

### Child

Child makes nodes children of a node.
It is a shorter from of `Enter` + `Leave`.

```go
src := New("abcd")

cod := And(
    S("a").Leaf("Char").Child(
        S("b").Leaf("Char"),
        S("c").Leaf("Char"),
    ),
    S("d").Leaf("Char"),
)

var ast AST
oks := cod.Tree(&ast).Run(src)

fmt.Println(oks, ast.Print("short-inline"))
// Root [ Char a [ Char b, Char c ], Char d ]
```

### Group

Group groups nodes inside a node.

```go
src := New("ab12")

cod := And(
    And(S("a").Leaf("L"), S("b").Leaf("L")).Group("Letters"),
    And(S("1").Leaf("N"), S("2").Leaf("N")).Group("Numbers"),
)

var ast AST
oks := cod.Tree(&ast).Run(src)

fmt.Println(oks, ast.Print("short-inline"))
// Root [ Letters [ L a, L b ], Numbers [ N 1, N 2 ] ]
```

### AST Examples

Parsing an expression.

```go
src := New("2+3*4+(5+6)")

term, setTerm := Recursive()
expr, setExpr := Recursive()

value := F(unicode.IsDigit).Leaf("Val")
factor := Or(And(S("("), expr, S(")")), value)
setTerm(Or(Root(factor, S("*").Leaf("Expr"), term).Undo(), factor))
setExpr(Or(Root(term, S("+").Leaf("Expr"), expr).Undo(), term))

var ast AST
oks := expr.Tree(&ast).Run(src)

fmt.Println(oks, ast.Print("short-inline"))
// Root [ Expr + [ Val 2, Expr + [ Expr * [ Val 3, Val 4 ], Expr + [ Val 5, Val 6 ] ] ] ]
```

More examples [here](/example/expression_ast_test.go).
