package expr

import (
	"go/ast"
	"go/parser"

	"github.com/muktihari/expr/boolean"
	"github.com/muktihari/expr/float"
	"github.com/muktihari/expr/integer"
)

// Int parses the given expr string into int as a result. e.g: "2 + 2" -> 4, "2.2 + 2" -> 4.
//
// Supported operators:
//
// Arithmetic: [+, -, *, /, %]
//
// Bitwise: [&, |, ^, &^, <<, >>]
func Int(str string) (int, error) {
	expr, err := parser.ParseExpr(str)
	if err != nil {
		return 0, err
	}

	visitor := integer.NewVisitor()
	ast.Walk(visitor, expr)
	return visitor.Result()
}

// Float64 parses the given expr string into float64 as a result . e.g: "2 + 2" -> 4, "2.2 + 2" -> 4.2
//
// Supported operators:
//
// Arithmetic: [+, -, *, /]
func Float64(str string) (float64, error) {
	expr, err := parser.ParseExpr(str)
	if err != nil {
		return 0, err
	}

	visitor := float.NewVisitor()
	ast.Walk(visitor, expr)
	return visitor.Result()
}

// Bool parses the given expr string into boolean as a result. e.g: "1 < 2" -> true, "1 > 2" -> false, "true || false" -> true.
//
// However, arithmetic is not supported at the moment, it WON'T DO "1 + 2 > 1" -> true [X]
//
// Supported operators:
//
// Comparison: [==, !=, <, <=, >, >=]
//
// Logical: [&&, ||]
func Bool(str string) (bool, error) {
	expr, err := parser.ParseExpr(str)
	if err != nil {
		return false, err
	}

	visitor := boolean.NewVisitor()
	ast.Walk(visitor, expr)
	return visitor.Result()
}
