# The ᚱune Programming Language

Congratulations. You have just discovered the silly programming language called "rune". It is named after the fact that all lexical tokens consist of single characters (called "runes" in go). This decision makes it very easy to write a lexer for the language, and result in some fun restrictions for the design of the language. You need to think about how to

- add comments
- use whitespace and line breaks for formatting
- represent strings (e.g. `"foo"`)
- represent multi-digit numbers (e.g. `42`)
- represent floating point numbers (e.g. `3.14`)
- use non-ascii operators (e.g. `<=`)

## Installing

You can just do a `go build` to produce an executable.

## The language

In Rune everything is an expression. It supports `+`, `-`, `*` and `/`.

    >>> 1+2*3
    --> 7

assignment

    >>> a:1
    --> 1

sequences of expressions

    >>> a:1;b:2;a+b
    --> 3
