# The ᚱune Programming Language

Congratulations. You have just discovered the silly programming language called "rune". It is named after the fact that all lexical tokens consist of single characters (called "runes" in go, hence the name). This decision makes it very easy to write a lexer for the language, and result in some fun restrictions for the design of the language. You need to think about how to

- add comments
- use whitespace for formatting
- represent strings (e.g. `"foo"`)
- represent multi-digit numbers (e.g. `42`)
- represent floating point numbers (e.g. `3.14`)
- use non-ascii operators (e.g. `<=`)

## Installing

You can just do a `go build` to produce an executable that you can run.

## The language

Here is a short introduction to the language.

### Basics

Single digit integers can be represented directly in rune.

    1
    --> 1

Multidigit integers can be obtained using the standard arithmic operators.

    5*5*4
    --> 100

Floating point numbers can not be expressed directly, but can be obtained by applying the `.` operator on two integers. Note that this is an operator, so it applies to any expressions evaluating to an integer.

    9.2
    --> 9.2
    3.(2*9*9*9-6*6-7)
    --> 3.1415

Boolean values can be obtained by applying comparison operators, but there is no way to define them explicitely. There are no multi-character operators, but `a not equal b` can be written as `!(a=b)` and `a less than or equal b` can be written as `!(a>b)`.

    !(1=1)
    --> false
    2<3
    --> true

Conditional operators are lazily evaluated (short-circuit), meaning that the second expression is only avaluated if necessary.

    1=0&1=1/0
    --> false
    1=1|1=1/0
    --> true

Variables can be defined using the `$` statement. Note that `$` should be read as `set`, and is not part of the resulting variable name.

    $a=3
    --> 3
    a+1
    --> 4
