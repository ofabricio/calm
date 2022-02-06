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

#### Recursion

- [x] [Recursive](#Recursive)

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

### F

F tests if the current character matches a rune function and moves the position if true.

```go
m := New("hi5")

a := F(unicode.IsLetter).Run(m) // h
b := F(unicode.IsLetter).Run(m) // i
c := F(unicode.IsLetter).Run(m) // 5

fmt.Println(a, b, c) // true true false
```

### SR

SR behaves exactly the same as [S](#S) but receives a string reference as argument.
This allows [Back Reference](#Back-Reference).

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

Note that the cursor advances even though it returns false.

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

### Recursive

Recursive allows recursive call of a matcher.

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

This operator is handy for small validations,
but it can be a pain when capturing tokens.
See [here](example/expression_test.go) an example on how to recursively
parse an expression without using this operator.

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

m := And(Or(S(`"`), S(`'`)).On(Grab(&quote)), F(unicode.IsLetter).OneToMany(), SR(&quote))

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
