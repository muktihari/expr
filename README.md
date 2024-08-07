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

## Usage

### Bind

For binding variables into expr string, see [Bind](./bind/README.md).

```go
s := "{price} - ({price} * {discount-percentage})"
v, _ := bind.Bind(s,
    "price", 100,
    "discount-percentage", 0.1,
)
fmt.Println(v) // "100 - (100 * 0.1)"
```

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
    v, err := expr.Any(str)
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
  - "(1 \* 10) > -2" -> true
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
    v, err := expr.Complex128(str)
    if err != nil {
        panic(err)
    }
    fmt.Printf("%f", v) // (4+3i)
```

### Float64

- Float64 parses the given expr string into float64 as a result. e.g:
  - "2 + 2" -> 4
  - "2.2 + 2" -> 4.2
  - "10 \* -5 + (-5.5)" -> -55.5
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
  - "10 + ((-5 \* -10) / -10) - 2" -> 3
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

## Benchmark

Benchmark results for evaluating simple math expression in comparison to [github.com/expr-lang/expr](github.com/expr-lang/expr). Please note that this library only offers simple expression evaluation, while expr-lang may offer richer features. The purpose of this benchmark is to demonstrate how effective this library is at handling simple use case scenarios.

```js
goos: darwin; goarch: amd64; pkg: benchmark
cpu: Intel(R) Core(TM) i5-5257U CPU @ 2.70GHz
BenchmarkExprLangExpr-4    68866  17455 ns/op  12835 B/op  70 allocs/op
BenchmarkMuktihariExpr-4  417950   2812 ns/op    872 B/op  24 allocs/op
```

Code:

```go
package benchmark_test

import (
	"testing"

	exprlang "github.com/expr-lang/expr"
	"github.com/muktihari/expr"
	"github.com/muktihari/expr/bind"
)

func BenchmarkExprLangExpr(b *testing.B) {
	for i := 0; i < b.N; i++ {
		const code = `price - (price * discountPercentage)`
		env := map[string]interface{}{
			"price":              10.0,
			"discountPercentage": 0.15,
		}
		program, _ := exprlang.Compile(code, exprlang.Env(env))
		val, _ := exprlang.Run(program, env)
		if expected := float64(8.5); expected != val {
			b.Fatalf("expected: %v, got: %v", expected, val)
		}
	}
}

func BenchmarkMuktihariExpr(b *testing.B) {
	for i := 0; i < b.N; i++ {
		const code = `{price} - ({price} * {discountPercentage})`
		s, _ := bind.Bind(code,
			"price", 10.0,
			"discountPercentage", 0.15,
		)
		val, _ := expr.Any(s)
		if expected := float64(8.5); expected != val {
			b.Fatalf("expected: %v, got: %v", expected, val)
		}
	}
}
```
