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
	case NumericTypeAuto:
		if vx.kind == KindImag || vy.kind == KindImag {
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
		Msg: "result of \"" + s + "\" is \"" + fmt.Sprintf("%v", v.value) + "\" which is not a number",
		Pos: v.pos,
		Err: ErrArithmeticOperation,
	}
}

func calculateComplex(v *Visitor, x, y complex128, op token.Token, opPos token.Pos) {
	v.kind = KindImag
	switch op {
	case token.ADD:
		v.value = x + y
	case token.SUB:
		v.value = x - y
	case token.MUL:
		v.value = x * y
	case token.QUO:
		v.value = x / y
	case token.REM:
		v.err = &SyntaxError{
			Msg: "operator \"" + op.String() + "\" is not supported to do arithmetic on complex number",
			Pos: int(opPos),
			Err: ErrArithmeticOperation,
		}
	}
}

func calculateFloat(v *Visitor, x, y float64, op token.Token) {
	v.kind = KindFloat
	switch op {
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

func calculateInt(v *Visitor, x, y int64, yPos int, op token.Token) {
	v.kind = KindInt
	switch op {
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
			v.err = &SyntaxError{
				Msg: "could not divide x with zero y, allowIntegerDividedByZero == false",
				Pos: yPos,
				Err: ErrIntegerDividedByZero,
			}
			return
		}
		v.value = x / y
	case token.REM:
		v.value = x % y
	}
}

func parseComplex(value interface{}) complex128 { // kind must be numeric
	switch val := value.(type) {
	case complex128:
		return val
	case float64:
		return complex(val, 0)
	case int64:
		return complex(float64(val), 0)
	}
	return 0
}

func parseFloat(value interface{}) float64 {
	switch val := value.(type) {
	case complex128:
		return real(val)
	case float64:
		return val
	case int64:
		return float64(val)
	}
	return 0
}

func parseInt(value interface{}) int64 {
	switch val := value.(type) {
	case complex128:
		return int64(real(val))
	case float64:
		return int64(val)
	case int64:
		return val
	}
	return 0
}
