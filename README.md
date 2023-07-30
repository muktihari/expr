![expr](./icon.svg)

![GitHub Workflow Status](https://github.com/muktihari/expr/workflows/CI/badge.svg)
[![Go Dev Reference](https://img.shields.io/badge/go.dev-reference-007d9c?logo=go&logoColor=white&style=flat-square)](https://pkg.go.dev/github.com/muktihari/expr)
[![CodeCov](https://codecov.io/gh/muktihari/expr/branch/master/graph/badge.svg)](https://codecov.io/gh/muktihari/expr)
[![Go Report Card](https://goreportcard.com/badge/github.com/muktihari/expr)](https://goreportcard.com/report/github.com/muktihari/expr)


Expr is a string expression parser in go. Not a fancy eval, just a simple and lightweight expr parser. Bool, Float64 and Int (with bitwise opperators) are available.

```go
"1 + 1" -> 2
"1 < 2 + 2" -> true
"true && !false" -> true
```

## Usage
#### Boolean
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
    - Arithmetic: [+, -, *, /, %] *(the % operator is only work for interger operation)*

```go
    str := "((1 < 2 && 3 > 4) || 1 == 1) && 4 < 5"
    v, err := expr.Bool(str)
    if err != nil {
        panic(err)
    }
    fmt.Printf("%t", v) // true
```

#### Float64
- Float64 parses the given expr string into float64 as a result. e.g:
    - "2 + 2" -> 4
    - "2.2 + 2" -> 4.2
    - "10 * -5 + (-5.5)" -> -55.5
- Supported operators:
    - Arithmetic: [+, -, *, /]

```go
    str := "((2 * 2) * (8 + 2) * 2) + 2.56789"
    v, err := expr.Float64(str)
    if err != nil {
        panic(err)
    }
    fmt.Printf("%f", v) // 82.56789
```

#### Integer
- Int parses the given expr string into int as a result. e.g:
    - "2 + 2" -> 4
    - "2.2 + 2" -> 4
    - "10 + ((-5 * -10) / -10) - 2" -> 3
- Supported operators:
    - Arithmetic: [+, -, *, /, %]
    - Bitwise: [&, |, ^, &^, <<, >>] (signed integer)
- Notes: 
    - << and >> operators are not permitted to be used in signed integer for go version less than 1.13.x.
    - Reference: [https://golang.org/doc/go1.13#language](https://golang.org/doc/go1.13#language)
    - Even if bitwise is supported, the priority operation is not granted, any bit operation is advised to be put in parentheses.

```go
    str := "((2 * 2) * (8 + 2) * 2) + 2.56789"
    v, err := expr.Int(str)
    if err != nil {
        panic(err)
    }
    fmt.Printf("%d", v) // 82
```


## License
Expr is released under [Apache Licence 2.0](https://www.apache.org/licenses/LICENSE-2.0)
