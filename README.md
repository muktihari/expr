![expr](./expr.png)

[![expr](https://img.shields.io/badge/go-reference-blue.svg?style=flat)](https://pkg.go.dev/github.com/muktihari/expr)
[![Go Report Card](https://goreportcard.com/badge/github.com/muktihari/expr?)](https://goreportcard.com/report/github.com/muktihari/expr)
[![Travis Widget](https://travis-ci.org/muktihari/expr.svg?branch=master)](https://travis-ci.org/github/muktihari/expr)

Expr is a string expression parser in go. Not a fancy eval, just a simple and lightweight expr parser.

```
"1 + 1" -> 2
"1 < 2" -> true
"true && false" -> false
```

## Usage
#### Boolean
```go
    str := "((1 < 2 && 3 > 4) || 1 == 1) && 4 < 5"
    v, err := expr.Bool(str)
    if err != nil {
        panic(err) // err is error invalid expression
    }
    fmt.Printf("%t", v) // true
```

#### Float64
```go
    str := "((2 * 2) * (8 + 2) * 2) + 2.56789"
    v, err := expr.Float64(str)
    if err != nil {
        panic(err)
    }
    fmt.Printf("%f", v) // 82.56789
```

#### Integer
```go
    str := "((2 * 2) * (8 + 2) * 2) + 2.56789"
    v, err := expr.Int(str)
    if err != nil {
        panic(err)
    }
    fmt.Printf("%d", v) // 82
```

*simple, isn't it?*

## License
Expr is released under [Apache Licence 2.0](https://www.apache.org/licenses/LICENSE-2.0)
