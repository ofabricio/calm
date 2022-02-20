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

- For a very simple expression parser example see [example/expression_test.go](example/expression_test.go).
- For a very simple Go code parser example see [example/parser_test.go](example/parser_test.go).

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

- [x] [Debug](#Debug)
- [x] [String](#String)
- [x] [Json](#Json)
- [x] [Number](#Number)

#### Recursion

- [x] [Recursive](#Recursive)

#### Tree ([AST](#AST))

- [x] [Leaf](#Leaf)
- [x] [Root](#Root)
- [x] [Enter](#Enter)
- [x] [Leave](#Leave)
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

The code below illustrates `Undo` behavior,
but there is no need to use `Undo` here
since `And` already calls `Undo`.

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

strg := String(`"`).On(Grabs(&quotes))
root := Or(strg, Next()).OneToMany()

ok := root.Run(c)

fmt.Println(ok, quotes) // true ["Wow!" "This is cool!"]
```

### Json

Json matches a json.

```go
c := New(`Use either { "hello": "world" } or { "foo": "bar" }.`)

var jsons []string

jsns := Json().On(Grabs(&jsons))
root := Or(jsns, Next()).OneToMany()

ok := root.Run(c)

fmt.Println(ok, jsons) // true [{ "hello": "world" } { "foo": "bar" }]
```

### Number

Number matches numbers.

```go
c := New("Heard of 3.1415? What about 0, 1, 1, 2, 3 sequence? Isn't 2e3 a cool notation?")

var n []string

numb := Number().On(Grabs(&n))
root := Or(numb, Next()).OneToMany()

ok := root.Run(c)

fmt.Println(ok, n) // true [3.14159 0 1 1 2 3 2e3]
```

### Recursive

Recursive allows a recursive call of a matcher.

```go
code := New("0+1*(2+3)*4")

term, setTerm := Recursive()
expr, setExpr := Recursive()

value := F(unicode.IsNumber)
factor := Or(And(S("("), expr, S(")")), value)
setTerm(Or(And(factor, S("*"), term), factor))
setExpr(Or(And(term, S("+"), expr), term))

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
    return Or(And(factor, S("*"), term), factor).Run(c)
}

expr = func(c *Code) bool {
    return Or(And(term, S("+"), expr), term).Run(c)
}

ok := expr.Run(code)

fmt.Println(ok) // true
```

This operator is handy for small validations,
but it can be a pain when capturing tokens.
See [here](example/expression_test.go) an example on how to recursively
parse an expression without using this operator.

For an easier and more advanced way to capture tokens see the [AST](#AST) section.

> Note to self: maybe I should remove this operator.

## AST

You can parse an input code into an AST (Abstract Syntax Tree).

For that, start the scanner with a root node:

```go
src := New("2+4")              // 1)
ast := Root("Program")         // 2)
oks := ast.Run(src, S("2+4"))  // 3)
fmt.Println(oks, ast)
```

- In line 1) we create the scanner, as we have doing since the beginning.
- In line 2) we create a root node of type "Program". Can be any string.
  It is used to categorize a node.
- In line 3) we run the scanner as an AST.

If you run the code above nothing will happen. That's because we need to
set up a few operators to capture the values we are interested in.

When dealing with a tree there are many operators to use.

### Leaf

Leaf is the basic operator to capture a node.
It creates a leaf node in the AST.
Without it the other operators won't work.

```go
src := New("2+4")

cod := And(S("2").Leaf("N"), S("+").Leaf("Op"), S("4").Leaf("N"))

ast := Root("Program")
oks := ast.Run(src, cod)

fmt.Println(oks, ast)
// { "type": "Program", "args": [
//     { "type": "N", "name": "2" }, { "type": "Op", "name": "+" }, { "type": "N", "name": "4" }] } 
```

### Root

Root converts three [Leaf](#Leaf) nodes into a binary branch.
In other words, it will set the middle node as a root
node and add the left and right nodes as its children.

```go
src := New("2+4")

cod := And(S("2").Leaf("N"), S("+").Leaf("Op"), S("4").Leaf("N")).Root()

ast := Root("Program")
oks := ast.Run(src, cod)

fmt.Println(oks, ast)
// { "type": "Program", "args": [
//     { "type": "Op", "name": "+", "args": [
//         { "type": "N", "name": "2" }, { "type": "N", "name": "4" }] }] }
```

> Note to self: maybe this operator needs a better name.

### Enter

Enter makes the selected [Leaf](#Leaf) node a root node,
so the next nodes are added as its children.
In other words it increases the depth of the tree in that node.
It is usually used on a [Leaf](#Leaf) node, for example `.Leaf("Func").Enter()`.

```go
src := New("print { 1 }")

cod := And(S("print").Leaf("FnCall").Enter(), S(" { "), S("1").Leaf("N"), S(" }"))

ast := Root("Program")
oks := ast.Run(src, cod)

fmt.Println(oks, ast)
// { "type": "Program", "args": [
//     { "type": "FnCall", "name": "print", "args": [{ "type": "N", "name": "1" }] }] }
```

> Note to self: maybe this operator needs a better name.

### Leave

Leave is the opposite of [Enter](#Enter).
Useful to restore an AST depth.
Make sure to Leave to a parent node.

```go
src := New("print { 1 } print { 1 } ")

fun := And(S("print").Leaf("FnCall").Enter(), S(" { "), S("1").Leaf("N"), S(" } "))
cod := And(fun, fun)

ast := Root("Program")
oks := ast.Run(src, cod)

fmt.Println(oks, ast)
// { "type": "Program", "args": [
//     { "type": "FnCall", "name": "print", "args": [{ "type": "N", "name": "1" }] },
//     { "type": "FnCall", "name": "print", "args": [{ "type": "N", "name": "1" }] }] }
```

> Note to self: maybe this operator needs a better name.

### Group

Group groups nodes inside a node.

```go
src := New("print(1, 2)")

num := F(unicode.IsDigit).Leaf("N")
arg := And(num, S(", ").ZeroToOne()).ZeroToMany().Group("Args")
cod := And(S("print").Leaf("FnCall").Enter(), S("("), arg, S(")"))

ast := Root("Program")
oks := ast.Run(src, cod)

fmt.Println(oks, ast)
// { "type": "Program", "args": [
//     { "type": "FnCall", "name": "print", "args": [
//         { "type": "Args", "args": [{ "type": "N", "name": "1" }, { "type": "N", "name": "2" }] }] }] }
```

### Examples

Parsing an expression.

```go
src := New("2+3*4+(5+6)")

term, setTerm := Recursive()
expr, setExpr := Recursive()

value := F(unicode.IsNumber).Leaf("N")
factor := Or(And(S("("), expr, S(")")), value)
setTerm(Or(And(factor, S("*").Leaf("BinExpr"), term).Root(), factor))
setExpr(Or(And(term, S("+").Leaf("BinExpr"), expr).Root(), term))

ast := Root("Program")
ok := ast.Run(src, expr)

fmt.Println(ok, ast)
```

More examples [here](/example/expression_ast_test.go).
