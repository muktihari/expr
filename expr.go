// Copyright 2020-2023 The Expr Authors
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package expr

import (
	"go/parser"
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

	var v Visitor
	v.options = defaultOptions()
	v.options.allowIntegerDividedByZero = true
	v.options.numericType = NumericTypeAuto

	v.Visit(expr)
	if err := v.Err(); err != nil {
		return nil, err
	}

	switch v.value.Kind() {
	case KindBoolean:
		return v.value.Bool(), nil
	case KindInt:
		return v.value.Int64(), nil
	case KindFloat:
		val := v.value.Float64()
		if val == float64(int64(val)) {
			return int64(val), nil
		}
		return val, nil
	case KindImag:
		return v.value.Complex128(), nil
	case KindString:
		return v.value.String(), nil
	default:
		return nil, nil
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

	var v Visitor
	v.options = defaultOptions()

	v.Visit(expr)
	if err := v.Err(); err != nil {
		return false, err
	}

	if v.value.Kind() == KindBoolean {
		return v.value.Bool(), nil
	}

	return false, ErrValueTypeMismatch
}

// Complex128 parses the given expr string into complex128 as a result. e.g:
//   - "(2+1i) + (2+2i)" -> (4+3i)
//   - "(2.2+1i) + 2" -> (4.2+1i)
//   - "2 + 2" -> (4+0i)
//
// - Supported operators:
//   - Arithmetic: [+, -, *, /]
func Complex128(str string) (complex128, error) {
	expr, err := parser.ParseExpr(str)
	if err != nil {
		return 0, err
	}

	var v Visitor
	v.options = defaultOptions()
	v.options.numericType = NumericTypeComplex

	v.Visit(expr)
	if err := v.Err(); err != nil {
		return 0, err
	}

	switch v.value.Kind() {
	case KindImag:
		return v.value.Complex128(), nil
	case KindFloat:
		return complex(v.value.Float64(), 0), nil
	case KindInt:
		return complex(float64(v.value.Int64()), 0), nil
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

	var v Visitor
	v.options = defaultOptions()
	v.options.numericType = NumericTypeFloat

	v.Visit(expr)
	if err := v.Err(); err != nil {
		return 0, err
	}

	switch v.value.Kind() {
	case KindImag:
		return real(v.value.Complex128()), nil
	case KindFloat:
		return v.value.Float64(), nil
	case KindInt:
		return float64(v.value.Int64()), nil
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

	var v Visitor
	v.options = defaultOptions()
	v.options.allowIntegerDividedByZero = allowIntegerDividedByZero
	v.options.numericType = NumericTypeInt

	v.Visit(expr)
	if err := v.Err(); err != nil {
		return 0, err
	}

	switch v.value.Kind() {
	case KindImag:
		return int64(real(v.value.Complex128())), nil
	case KindFloat:
		return int64(v.value.Float64()), nil
	case KindInt:
		return v.value.Int64(), nil
	}

	return 0, ErrValueTypeMismatch
}
