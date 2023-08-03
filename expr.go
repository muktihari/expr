package expr

import (
	"go/ast"
	"go/parser"
	"strconv"
)

// Any parses the given expr string into any type it returns as a result. e.g:
//   - "1 < 2" -> true
//   - "true || false" -> true
//   - "2 + 2" -> 4
//   - "4 << 10" -> 4906
//   - "2.2 + 2" -> 4.2
//   - "(2+1i) + (2+2i)" -> (4+3i)
//   - ""abc" == "abc"" -> true
//   - ""abc"" -> "abc"
//
// - Supported operators:
//   - Comparison: [==, !=, <, <=, >, >=]
//   - Logical: [&&, ||, !]
//   - Arithmetic: [+, -, *, /, %] (% operator does not work for complex number)
//   - Bitwise: [&, |, ^, &^, <<, >>] (only work for integer values)
func Any(str string) (interface{}, error) {
	expr, err := parser.ParseExpr(str)
	if err != nil {
		return nil, err
	}

	v := NewVisitor()
	ast.Walk(v, expr)

	if err := v.Err(); err != nil {
		return nil, err
	}

	switch v.Kind() {
	case KindInt:
		v, _ := strconv.ParseInt(v.Value(), 0, 64)
		return v, nil
	case KindFloat:
		v, _ := strconv.ParseFloat(v.Value(), 64)
		vInt := int64(v)
		if v == float64(vInt) {
			return vInt, nil
		}
		return v, nil
	case KindImag:
		v, _ := strconv.ParseComplex(v.Value(), 128)
		return v, nil
	case KindBoolean:
		v, _ := strconv.ParseBool(v.Value())
		return v, nil
	default:
		return v.Value(), nil
	}
}

// Bool parses the given expr string into boolean as a result. e.g:
//   - "1 < 2" -> true
//   - "1 > 2" -> false
//   - "true || false" -> true
//   - "true && !false" -> true
//
// - Arithmetic operation are supported. e.g:
//   - "1 + 2 > 1" -> true
//   - "(1 * 10) > -2" -> true
//
// - Supported operators:
//   - Comparison: [==, !=, <, <=, >, >=]
//   - Logical: [&&, ||, !]
//   - Arithmetic: [+, -, *, /, %] (% operator does not work for complex number)
//   - Bitwise: [&, |, ^, &^, <<, >>] (only work for integer values)
func Bool(str string) (bool, error) {
	expr, err := parser.ParseExpr(str)
	if err != nil {
		return false, err
	}

	v := NewVisitor()
	ast.Walk(v, expr)

	if err := v.Err(); err != nil {
		return false, err
	}

	switch v.kind {
	case KindBoolean:
		val, _ := strconv.ParseBool(v.Value())
		return val, nil
	}

	return false, ErrValueTypeMismatch
}

//	Complex128 parses the given expr string into complex128 as a result. e.g:
//	   - "(2+1i) + (2+2i)" -> (4+3i)
//	   - "(2.2+1i) + 2" -> (4.2+1i)
//	   - "2 + 2" -> (4+0i)
//
// - Supported operators:
//   - Arithmetic: [+, -, *, /]
func Complex128(str string) (complex128, error) {
	expr, err := parser.ParseExpr(str)
	if err != nil {
		return 0, err
	}

	v := NewVisitor(WithNumericType(NumericTypeComplex))
	ast.Walk(v, expr)

	if err := v.Err(); err != nil {
		return 0, err
	}

	switch v.Kind() {
	case KindImag:
		v, _ := strconv.ParseComplex(v.Value(), 128)
		return v, nil
	case KindInt:
		v, _ := strconv.ParseInt(v.Value(), 0, 64)
		return complex(float64(v), 0), nil
	case KindFloat:
		v, _ := strconv.ParseFloat(v.Value(), 64)
		return complex(v, 0), nil
	}

	return 0, ErrValueTypeMismatch
}

// Float64 parses the given expr string into float64 as a result. e.g:
//   - "2 + 2" -> 4
//   - "2.2 + 2" -> 4.2
//   - "10 * -5 + (-5.5)" -> -55.5
//   - "10.0 % 2.6" -> 2.2
//
// - Supported operators:
//   - Arithmetic: [+, -, *, /, %]
func Float64(str string) (float64, error) {
	expr, err := parser.ParseExpr(str)
	if err != nil {
		return 0, err
	}

	v := NewVisitor(WithNumericType(NumericTypeFloat))
	ast.Walk(v, expr)

	if err := v.Err(); err != nil {
		return 0, err
	}

	switch v.Kind() {
	case KindInt:
		v, _ := strconv.ParseInt(v.Value(), 0, 64)
		return float64(v), nil
	case KindFloat:
		v, _ := strconv.ParseFloat(v.Value(), 64)
		return v, nil
	}

	return 0, ErrValueTypeMismatch
}

// - Int64 parses the given expr string into int64 as a result. e.g:
//   - "2 + 2" -> 4
//   - "2.2 + 2" -> 4
//   - "10 + ((-5 * -10) / -10) - 2" -> 3
//
// - Supported operators:
//   - Arithmetic: [+, -, *, /, %]
//   - Bitwise: [&, |, ^, &^, <<, >>]
func Int64(str string) (int64, error) {
	return parseStringExprIntoInt64(str, true)
}

// Int64Strict is shorthand for Int64(str) but when x / y and y == 0, it will return ErrIntegerDividedByZero
func Int64Strict(str string) (int64, error) {
	return parseStringExprIntoInt64(str, false)
}

// Int is shorthand for Int64(str) with its result will be converted into int.
func Int(str string) (int, error) {
	v, err := Int64(str)
	if err != nil {
		return 0, err
	}
	return int(v), nil
}

func parseStringExprIntoInt64(str string, allowIntegerDividedByZero bool) (int64, error) {
	expr, err := parser.ParseExpr(str)
	if err != nil {
		return 0, err
	}

	v := NewVisitor(
		WithNumericType(NumericTypeInt),
		WithAllowIntegerDividedByZero(allowIntegerDividedByZero),
	)
	ast.Walk(v, expr)

	if err := v.Err(); err != nil {
		return 0, err
	}

	switch v.Kind() {
	case KindInt:
		v, _ := strconv.ParseInt(v.Value(), 0, 64)
		return v, nil
	case KindFloat:
		v, _ := strconv.ParseFloat(v.Value(), 64)
		return int64(v), nil
	}

	return 0, ErrValueTypeMismatch
}
