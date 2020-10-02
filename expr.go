package expr

import (
	"go/ast"
	"go/parser"

	"github.com/muktihari/expr/boolean"
	"github.com/muktihari/expr/float"
	"github.com/muktihari/expr/integer"
)

// Int parses the given expr string into int as a result.
// - e.g:
// 	- "2 + 2" -> 4
// 	- "2.2 + 2" -> 4
// 	- "10 + ((-5 * -10) / -10) - 2" -> 3
//
// - Supported operators:
// 	- Arithmetic: [+, -, *, /, %]
// 	- Bitwise: [&, |, ^, &^, <<, >>]
// - Notes:
//  - << and >> operators are not permitted to be used in signed integer for go version less than 1.13.x.
//  - Reference: golang.org/doc/go1.13#language
//  - Even if bitwise is supported, the priority operation is not granted, any bit operation is advised to be put in parentheses.
func Int(str string) (int, error) {
	expr, err := parser.ParseExpr(str)
	if err != nil {
		return 0, err
	}

	visitor := integer.NewVisitor()
	ast.Walk(visitor, expr)
	return visitor.Result()
}

// Float64 parses the given expr string into float64 as a result.
// - e.g:
// 	- "2 + 2" -> 4
// 	- "2.2 + 2" -> 4.2
//	- "10 * -5 + (-5.5)" -> -55.5
//
// - Supported operators:
// 	- Arithmetic: [+, -, *, /]
func Float64(str string) (float64, error) {
	expr, err := parser.ParseExpr(str)
	if err != nil {
		return 0, err
	}

	visitor := float.NewVisitor()
	ast.Walk(visitor, expr)
	return visitor.Result()
}

// Bool parses the given expr string into boolean as a result.
// - e.g:
// 	- "1 < 2" -> true
// 	- "1 > 2" -> false
// 	- "true || false" -> true
//
// - Arithmetic operation are supported. e.g:
// 	- "1 + 2 > 1" -> true
//	- "(1 * 10) > -2" -> true
// - Supported operators:
// 	- Comparison: [==, !=, <, <=, >, >=]
// 	- Logical: [&&, ||, !]
// 	- Arithmetic: [+, -, *, /, %] (the % operator is only work for interger operation)
func Bool(str string) (bool, error) {
	expr, err := parser.ParseExpr(str)
	if err != nil {
		return false, err
	}

	visitor := boolean.NewVisitor()
	ast.Walk(visitor, expr)
	return visitor.Result()
}
