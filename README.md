![expr](./icon.svg)

![GitHub Workflow Status](https://github.com/muktihari/expr/workflows/CI/badge.svg)
[![Go Dev Reference](https://img.shields.io/badge/go.dev-reference-007d9c?logo=go&logoColor=white&style=flat-square)](https://pkg.go.dev/github.com/muktihari/expr)
[![CodeCov](https://codecov.io/gh/muktihari/expr/branch/master/graph/badge.svg)](https://codecov.io/gh/muktihari/expr)
[![Go Report Card](https://goreportcard.com/badge/github.com/muktihari/expr)](https://goreportcard.com/report/github.com/muktihari/expr)


Expr is a simple, lightweight and performant programming toolkit for evaluating basic mathematical expression and boolean expression. The resulting value is one of these following primitive types: string, boolean, numerical (complex, float, integer).

## Supported Numerical Notations
```js
- Binary (base-2)       : 0b1011
- Octal (base-8)        : 0o13 or 013
- Decimal (base-10)     : 11
- Hexadecimal (base-16) : 0xB 
- Scientific            : 11e0
```

## Expression Examples
```js
"1 + 1"             -> 2
"1.0 / 2"           -> 0.5
"2 < 2 + 2"         -> true
"true && !false"    -> true
"4 << 10"           -> 4096
"0b0100 << 0b1010"  -> 4096 (0b1000000000000)
"0x4 << 0xA"        -> 4096 (0x1000)
"0o4 << 0o12"       -> 4096 (0o10000)
"0b1000 | 0b1001"   -> 9 (0b1001)
"(2+1i) + (2+2i)"   -> (4+3i)
"0x4 << 0xA > 1024" -> true
```

## Usage

### Any
- Any parses the given expr string into any type it returns as a result. e.g:
    - "1 < 2" -> true
    - "true || false" -> true
    - "2 + 2" -> 4
    - "4 << 10" -> 4906
    - "2.2 + 2" -> 4.2
    - "(2+1i) + (2+2i)" -> (4+3i)
    - ""abc" == "abc"" -> true
    - ""abc"" -> "abc"
- Supported operators:
    - Comparison: [==, !=, <, <=, >, >=]
    - Logical: [&&, ||, !]
    - Arithmetic: [+, -, *, /, %] (% operator does not work for complex number)
    - Bitwise: [&, |, ^, &^, <<, >>] (only work for integer values)

```go
    str := "(2+1i) + (2+2i)"
    v, err := expr.Bool(str)
    if err != nil {
        panic(err)
    }
    fmt.Printf("%v", v) // (4+3i)
```

### Boolean
- Bool parses the given expr string into boolean as a result. e.g:
    - "1 < 2" -> true
    - "1 > 2" -> false
    - "true || false" -> true
    - "true && !false" -> true
- Arithmetic operation are supported. e.g:
    - "1 + 2 > 1" -> true
    - "(1 * 10) > -2" -> true
- Supported operators:
    - Comparison: [==, !=, <, <=, >, >=]
    - Logical: [&&, ||, !]
    - Arithmetic: [+, -, *, /, %] (% operator does not work for complex number)
    - Bitwise: [&, |, ^, &^, <<, >>] (only work for integer values)

```go
    str := "((1 < 2 && 3 > 4) || 1 == 1) && 4 < 5"
    v, err := expr.Bool(str)
    if err != nil {
        panic(err)
    }
    fmt.Printf("%t", v) // true
```

### Complex128
- Complex128 parses the given expr string into complex128 as a result. e.g:
    - "(2+1i) + (2+2i)" -> (4+3i)
    - "(2.2+1i) + 2" -> (4.2+1i)
    - "2 + 2" -> (4+0i)
- Supported operators:
    - Arithmetic: [+, -, *, /]

```go
    str := "(2+1i) + (2+2i)"
    v, err := expr.Float64(str)
    if err != nil {
        panic(err)
    }
    fmt.Printf("%f", v) // (4+3i)
```

### Float64
- Float64 parses the given expr string into float64 as a result. e.g:
    - "2 + 2" -> 4
    - "2.2 + 2" -> 4.2
    - "10 * -5 + (-5.5)" -> -55.5
    - "10.0 % 2.6" -> 2.2
- Supported operators:
    - Arithmetic: [+, -, *, /, %]

```go
    str := "((2 * 2) * (8 + 2) * 2) + 2.56789"
    v, err := expr.Float64(str)
    if err != nil {
        panic(err)
    }
    fmt.Printf("%f", v) // 82.56789
```

### Int64
- Int64 parses the given expr string into int64 as a result. e.g:
    - "2 + 2" -> 4
    - "2.2 + 2" -> 4
    - "10 + ((-5 * -10) / -10) - 2" -> 3
- Supported operators:
    - Arithmetic: [+, -, *, /, %]
    - Bitwise: [&, |, ^, &^, <<, >>]

```go
    str := "((2 * 2) * (8 + 2) * 2) + 2.56789"
    v, err := expr.Int64(str)
    if err != nil {
        panic(err)
    }
    fmt.Printf("%d", v) // 82
```

### Int64Strict
- Int64Strict is shorthand for Int64(str) but when x / y and y == 0, it will return ErrIntegerDividedByZero

```go
    str := "12 + 24 - 10/0"
    v, err := expr.Int64Strict(str)
    if err != nil {
        // err == ErrIntegerDividedByZero
    }
    fmt.Printf("%d", v) // 0
```

### Int
- Int is shorthand for Int64(str) with its result will be converted into int.

```go
    str := "1 + 10"
    v, err := expr.Int(str)
    if err != nil {
        panic(err)
    }
    fmt.Printf("%d", v) // 11
```

## License
Expr is released under [Apache Licence 2.0](https://www.apache.org/licenses/LICENSE-2.0)
