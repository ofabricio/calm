# Calm

This is a library to scan, match, parse and tokenize text.

[![Go](https://github.com/ofabricio/calm/actions/workflows/go.yml/badge.svg)](https://github.com/ofabricio/calm/actions/workflows/go.yml)

### Note

Use [issues](https://github.com/ofabricio/calm/issues) only to report bugs.
Use [discussions](https://github.com/ofabricio/calm/discussions) for everything else.
Don't open PRs.

### Example

```go
package main

import (
    "fmt"
    "unicode"
    . "github.com/ofabricio/calm"
)

func main() {

    src := New("Hello, World!")

    var words []string

    word := F(unicode.IsLetter).OneToMany().GrabMany(&words)
    only := Or(word, Next()).OneToMany()

    ok := src.Run(only)

    fmt.Println(ok, words)
    // true [Hello World]
}
```

See more examples in the [example](/example) folder.

## How it works

Each time a match happens it moves the cursor to the next character.

If there is no match the cursor doesn't move.

Only `S` and `F` operators move the cursor on a match.

There are many kinds of operators: conditional execution, repetition, recursion, etc.

Some operators can move the cursor back and forth.

## Operators

#### Matcher

- [x] [S](#S)
- [x] [SR](#SR)
- [x] [F](#F)

#### Tester

- [x] [Eq](#Eq)
- [x] [More](#More)

#### Logical

- [x] [Or](#Or)
- [x] [And](#And)
- [x] [Not](#Not)
- [x] [True](#True)
- [x] [False](#False)

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
- [x] [Rewind](#Rewind)

#### Grabber

- [x] [Grab](#Grab)
- [x] [GrabMany](#GrabMany)
- [x] [GrabUndo](#GrabUndo)
- [x] [GrabPos](#GrabPos)
- [x] [Emit](#Emit)
- [x] [EmitMany](#EmitMany)
- [x] [EmitUndo](#EmitUndo)
- [x] [On](#On)
- [x] [Back Reference](#Back-Reference)

#### Util

- [x] [Debug](#Debug)
- [x] [String](#String)

## S

S tests if the current token matches a string and moves the position if true.

```go
m := New("hello world")

a := m.Run(S("hello"))
b := m.Run(S(" "))
c := m.Run(S("hello"))
d := m.Run(S("world"))

fmt.Println(a, b, c, d) // true true false true
```

## F

F tests if the current character matches a rune function and moves the position if true.

```go
m := New("hi5")

a := m.Run(F(unicode.IsLetter)) // h
b := m.Run(F(unicode.IsLetter)) // i
c := m.Run(F(unicode.IsLetter)) // 5

fmt.Println(a, b, c) // true true false
```

## SR

SR behaves exactly the same as `S` but receives a string reference as argument.

## Eq

Eq tests if the current token matches a string. It does not move the cursor.

```go
c := New("hello")

a := c.Run(Eq("hello"))
b := c.Run(Eq("hello"))

fmt.Println(a, b) // true true
```

## More

More runs the current operator only if there are more characters to match (so the cursor can move).

It is used to prevent overflow in some operations.

## Or

Or returns true if one of its arguments return true.

```go
c := New("apple")

ok := c.Run(Or(S("grape"), S("apple")))

fmt.Println(ok) // true
```

## And

And returns true if all of its arguments return true.

```go
c := New("hello world")

ok := c.Run(And(S("hello"), S(" "), S("world")))

fmt.Println(ok) // true
```

## Not

Not negates the current operator. True becomes false and vice-versa.

```go
c := New("hello")

ok := c.Run(S("hello").Not())

fmt.Println(ok) // false
```

Note that the cursor advances even though it returns false.

## True

True forces the current operator to return true.

```go
c := New("hello")

ok := c.Run(S("world").True())

fmt.Println(ok) // true
```

## False

False forces the current operator to return false.

```go
c := New("hello")

ok := c.Run(S("hello").False())

fmt.Println(ok) // false
```

Note that the cursor advances even though it returns false.

## ZeroToMany

ZeroToMany matches zero to many tokens. It is equivalent to the regex symbol `*`.

```go
m := And(S("hello"), S(" ").ZeroToMany(), S("world"))

a := New("helloworld").Run(m)
b := New("hello world").Run(m)
c := New("hello  world").Run(m)

fmt.Println(a, b, c) // true true true
```

## OneToMany

OneToMany matches one to many tokens. It is equivalent to the regex symbol `+`.

```go
m := And(S("hello"), S(" ").OneToMany(), S("world"))

a := New("helloworld").Run(m)
b := New("hello world").Run(m)
c := New("hello  world").Run(m)

fmt.Println(a, b, c) // false true true
```

## ZeroToOne

ZeroToOne matches an optional token. It is equivalent to the regex symbol `?`.

```go
m := And(S("colo"), S("u").ZeroToOne(), S("r"))

a := New("color").Run(m)
b := New("colour").Run(m)

fmt.Println(a, b) // true true
```

## Min

Min matches a minimum number of tokens.

```go
m := S("a").Min(2)

a := New("a").Run(m)
b := New("aa").Run(m)
c := New("aaa").Run(m)

fmt.Println(a, b, c) // false true true
```

## Until

Until matches until some matcher return true.

```go
m := Until(Eq(","), Eq("."))

a := New(",").Run(m)
b := New("ab,").Run(m)
c := New("abcd.").Run(m)

fmt.Println(a, b, c) // false true true
```

Note that Until advances the position by one character.
Be careful when using as argument matchers with more
than one character like `S("abc")`.

## While

While matches while any matcher returns true.

```go
m := While(Eq("0"), Eq("1"))

a := New("01100").Run(m)
b := New("hello").Run(m)

fmt.Println(a, b) // true false
```

Note that While advances the position by one character.
Be careful when using as argument matchers with more
than one character like `S("abc")`.

## Recursive

Recursive allows recursive calls of the current matcher.

```go
c := New("0+1*(2+3)*4")

var term Wrap
var expr Wrap

// <expr> = <term> '+' <expr> | <term>
// <term> = <factor> '*' <term> | <factor>
// <factor> = '(' <expr> ')' | <value>

value := F(unicode.IsNumber)
factor := Or(And(S("("), &expr, S(")")), value)
Or(And(factor, S("*"), &term).Not().Rewind().Not(), factor).Recursive(&term)
Or(And(&term, S("+"), &expr).Not().Rewind().Not(), &term).Recursive(&expr)

ok := c.Run(&expr)

fmt.Println(ok) // true
```

## Next

Next moves to the next character when the current matcher returns true.

```go
c := New("hello world")

ok := c.Run(And(S("hello").Next(), S("world")))

fmt.Println(ok) // true
```

There is also a static version of `Next`.

## Rewind

Rewind rewinds the cursor back to the begining of the matched token.

```go
c := New("hello world")

ok := c.Run(And(S("hello").Rewind(), S("hello world")))

fmt.Println(ok) // true
```

## Grab

Grab captures the current token string.

```go
c := New("abc123")

var t string
c.Run(F(unicode.IsLetter).OneToMany().Grab(&t))

fmt.Println(t) // abc
```

## GrabMany

GrabMany captures the current token string and adds it to a slice.

```go
c := New("abc123")

var ts []string
c.Run(F(unicode.IsLetter).GrabMany(&ts).OneToMany())

fmt.Println(ts) // [a b c]
```

Note that GrabMany might repeat tokens depending on the logic used.

## GrabUndo

GrabUndo can be used to undo (remove) tokens grabbed by GrabMany.
This is usually used along with [Rewind](#Rewind),
since when a match rewinds you might want to discard grabbed tokens.
See an example [here](grabber_test.go).

## GrabPos

GrabPos captures the current token position.

```go
c := New("abc")

var pos int
c.Run(And(S("a"), S("b").GrabPos(&pos), S("c")))

fmt.Println(pos) // 1
```

## Emit

Emit captures the current token.

```go
c := New("abc123")

var t Token
c.Run(F(unicode.IsLetter).OneToMany().Emit(&t))

fmt.Println(t) // {abc 0 1 1}
```

## EmitMany

EmitMany captures the current token and adds it to a slice.

```go
c := New("abc123")

var ts []Token
c.Run(F(unicode.IsLetter).EmitMany(&ts).OneToMany())

fmt.Println(ts) // [{a 0 1 1} {b 1 1 2} {c 2 1 3}]
```

Note that EmitMany might repeat tokens depending on the logic used.

## EmitUndo

EmitUndo can be used to undo (remove) tokens emitted by EmitMany.
This is usually used along with [Rewind](#Rewind),
since when a match rewinds you might want to discard emitted tokens.
See an example [here](grabber_test.go).

## On

On calls a function with the token when the current operator matches.

```go
c := New("hello")

f := func(t Token) {
    fmt.Println(t) // {hello 0 1 1}
}

c.Run(S("hello").On(f))
```

## Back Reference

It is possible to create a back reference with [Grab](#Grab) + [SR](#SR).

```go
var quote string

code := And(Or(S(`"`), S(`'`)).Grab(&quote), F(unicode.IsLetter).OneToMany(), SR(&quote))

a := New(`"hello"`).Run(code)
b := New(`'hello'`).Run(code)
c := New(`"hello'`).Run(code)

fmt.Println(a, b, c) // true true false
```

## Debug

Debug is an util that prints debug info to the stdout.

```go
c := New("Hi")

c.Run(F(unicode.IsLetter).Debug().OneToMany())

// [debug] Match: true  Token: 'H' Pos: 0 End: 1
// [debug] Match: true  Token: 'i' Pos: 1 End: 2
// [debug] Match: false Token: ''  Pos: 2 End: 2
```

## String

String is an util that parses a common string definition.
It allows quote escaping.

```go
c := New(`He said "Wow!" and "This is cool!" when he saw this.`)

var ts []string
code := Or(String(`"`).GrabMany(&ts), Next()).OneToMany()

c.Run(code)

fmt.Println(ts)
// ["Wow!" "This is cool!"]
```
