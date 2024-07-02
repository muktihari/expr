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
	if vx.value.Kind() <= numeric_beg || vx.value.Kind() >= numeric_end {
		v.err = newArithmeticNonNumericError(vx, binaryExpr.X)
		return
	}
	if vy.value.Kind() <= numeric_beg || vy.value.Kind() >= numeric_end {
		v.err = newArithmeticNonNumericError(vy, binaryExpr.Y)
		return
	}

	switch v.options.numericType {
	case NumericTypeAuto:
		if vx.value.Kind() == KindImag || vy.value.Kind() == KindImag {
			calculateComplex(v, parseComplex(vx.value), parseComplex(vy.value), binaryExpr.Op, binaryExpr.OpPos)
			return
		}
		calculateFloat(v, parseFloat(vx.value), parseFloat(vy.value), binaryExpr.Op)
		return
	case NumericTypeComplex:
		calculateComplex(v, parseComplex(vx.value), parseComplex(vy.value), binaryExpr.Op, binaryExpr.OpPos)
		return
	case NumericTypeFloat:
		calculateFloat(v, parseFloat(vx.value), parseFloat(vy.value), binaryExpr.Op)
		return
	case NumericTypeInt:
		calculateInt(v, parseInt(vx.value), parseInt(vy.value), vy.pos, binaryExpr.Op)
		return
	}
}

func newArithmeticNonNumericError(v *Visitor, e ast.Expr) error {
	s := conv.FormatExpr(e)
	return &SyntaxError{
		Msg: "result of \"" + s + "\" is \"" + fmt.Sprintf("%v", v.value.Any()) + "\" which is not a number",
		Pos: v.pos,
		Err: ErrArithmeticOperation,
	}
}

func calculateComplex(v *Visitor, x, y complex128, op token.Token, opPos token.Pos) {
	v.value.SetKind(KindImag)
	switch op {
	case token.ADD:
		v.value = complex128Value(x + y)
	case token.SUB:
		v.value = complex128Value(x - y)
	case token.MUL:
		v.value = complex128Value(x * y)
	case token.QUO:
		v.value = complex128Value(x / y)
	case token.REM:
		v.value = value{}
		v.err = &SyntaxError{
			Msg: "operator \"" + op.String() + "\" is not supported to do arithmetic on complex number",
			Pos: int(opPos),
			Err: ErrArithmeticOperation,
		}
	}
}

func calculateFloat(v *Visitor, x, y float64, op token.Token) {
	v.value.SetKind(KindFloat)
	switch op {
	case token.ADD:
		v.value = float64Value(x + y)
	case token.SUB:
		v.value = float64Value(x - y)
	case token.MUL:
		v.value = float64Value(x * y)
	case token.QUO:
		v.value = float64Value(x / y)
	case token.REM:
		v.value = float64Value(math.Mod(x, y))
	}
}

func calculateInt(v *Visitor, x, y int64, yPos int, op token.Token) {
	v.value.SetKind(KindInt)
	switch op {
	case token.ADD:
		v.value = int64Value(x + y)
	case token.SUB:
		v.value = int64Value(x - y)
	case token.MUL:
		v.value = int64Value(x * y)
	case token.QUO:
		if y == 0 {
			if v.options.allowIntegerDividedByZero {
				v.value = int64Value(0)
				return
			}
			v.value = value{}
			v.err = &SyntaxError{
				Msg: "could not divide x with zero y, allowIntegerDividedByZero == false",
				Pos: yPos,
				Err: ErrIntegerDividedByZero,
			}
			return
		}
		v.value = int64Value(x / y)
	case token.REM:
		v.value = int64Value(x % y)
	}
}

func parseComplex(val value) complex128 { // kind must be numeric
	switch val.Kind() {
	case KindImag:
		return val.Complex128()
	case KindFloat:
		return complex(val.Float64(), 0)
	case KindInt:
		return complex(float64(val.Int64()), 0)
	}
	return 0
}

func parseFloat(val value) float64 {
	switch val.Kind() {
	case KindImag:
		return real(val.Complex128())
	case KindFloat:
		return val.Float64()
	case KindInt:
		return float64(val.Int64())
	}
	return 0
}

func parseInt(val value) int64 {
	switch val.Kind() {
	case KindImag:
		return int64(real(val.Complex128()))
	case KindFloat:
		return int64(val.Float64())
	case KindInt:
		return val.Int64()
	}
	return 0
}
