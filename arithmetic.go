// Copyright 2023 The Expr Authors
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
	"fmt"
	"go/ast"
	"go/token"
	"math"

	"github.com/muktihari/expr/internal/conv"
)

func arithmetic(v, vx, vy *Visitor, binaryExpr *ast.BinaryExpr) {
	// numeric guards:
	if vx.kind <= numeric_beg || vx.kind >= numeric_end {
		v.err = newArithmeticNonNumericError(vx, binaryExpr.X)
		return
	}
	if vy.kind <= numeric_beg || vy.kind >= numeric_end {
		v.err = newArithmeticNonNumericError(vy, binaryExpr.Y)
		return
	}

	switch v.options.numericType {
	case NumericTypeComplex:
		calculateComplex(v, vx, vy, binaryExpr)
		return
	case NumericTypeFloat:
		calculateFloat(v, vx, vy, binaryExpr)
		return
	case NumericTypeInt:
		calculateInt(v, vx, vy, binaryExpr)
		return
	}

	// NumericTypeAuto: Auto figure out value type
	if vx.kind == KindImag || vy.kind == KindImag {
		calculateComplex(v, vx, vy, binaryExpr)
		return
	}

	// calculate other types as float64
	calculateFloat(v, vx, vy, binaryExpr)
}

func newArithmeticNonNumericError(v *Visitor, e ast.Expr) error {
	s := conv.FormatExpr(e)
	return &SyntaxError{
		Msg: "result of \"" + s + "\" is \"" + fmt.Sprintf("%v", v.value) + "\" which is not a number",
		Pos: v.pos,
		Err: ErrArithmeticOperation,
	}
}

func calculateComplex(v, vx, vy *Visitor, binaryExpr *ast.BinaryExpr) {
	v.kind = KindImag
	x := parseComplex(vx.value, vx.kind)
	y := parseComplex(vy.value, vy.kind)

	switch binaryExpr.Op {
	case token.ADD:
		v.value = x + y
	case token.SUB:
		v.value = x - y
	case token.MUL:
		v.value = x * y
	case token.QUO:
		v.value = x / y
	case token.REM:
		v.kind = KindIllegal
		v.err = &SyntaxError{
			Msg: "operator \"" + binaryExpr.Op.String() + "\" is not supported to do arithmetic on complex number",
			Pos: int(binaryExpr.OpPos),
			Err: ErrArithmeticOperation,
		}
		return
	}
}

func calculateFloat(v, vx, vy *Visitor, binaryExpr *ast.BinaryExpr) {
	v.kind = KindFloat
	x := parseFloat(vx.value, vx.kind)
	y := parseFloat(vy.value, vy.kind)

	switch binaryExpr.Op {
	case token.ADD:
		v.value = x + y
	case token.SUB:
		v.value = x - y
	case token.MUL:
		v.value = x * y
	case token.QUO:
		v.value = x / y
	case token.REM:
		v.value = math.Mod(x, y)
	}
}

func calculateInt(v, vx, vy *Visitor, binaryExpr *ast.BinaryExpr) {
	v.kind = KindInt
	x := parseInt(vx.value, vx.kind)
	y := parseInt(vy.value, vy.kind)

	switch binaryExpr.Op {
	case token.ADD:
		v.value = x + y
	case token.SUB:
		v.value = x - y
	case token.MUL:
		v.value = x * y
	case token.QUO:
		if y == 0 {
			if v.options.allowIntegerDividedByZero {
				v.value = int64(0)
				return
			}
			v.kind = KindIllegal
			v.err = &SyntaxError{
				Msg: "could not divide x with zero y, allowIntegerDividedByZero == false",
				Pos: vy.pos,
				Err: ErrIntegerDividedByZero,
			}
			return
		}
		v.value = x / y
	case token.REM:
		v.value = x % y
	}
}

func parseComplex(v interface{}, kind Kind) complex128 { // kind must be numeric
	switch kind {
	case KindImag:
		v := v.(complex128)
		return v
	case KindFloat:
		v := v.(float64)
		return complex(v, 0)
	default: // INT: 0xFF, 0b1010, 0o77, 071, 90
		v := v.(int64)
		return complex(float64(v), 0)
	}
}

func parseFloat(v interface{}, kind Kind) float64 { // kind must be numeric
	switch kind {
	case KindImag:
		v := v.(complex128)
		return real(v)
	case KindFloat:
		v := v.(float64)
		return v
	default: // INT: 0xFF, 0b1010, 0o77, 071, 90
		v := v.(int64)
		return float64(v)
	}
}

func parseInt(v interface{}, kind Kind) int64 { // kind must be numeric
	switch kind {
	case KindImag:
		v := v.(complex128)
		return int64(real(v))
	case KindFloat:
		v := v.(float64)
		return int64(v)
	default: // INT: 0xFF, 0b1010, 0o77, 071, 90
		v := v.(int64)
		return v
	}
}
