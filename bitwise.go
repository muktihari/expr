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

	"github.com/muktihari/expr/internal/conv"
)

func bitwise(v, vx, vy *Visitor, binaryExpr *ast.BinaryExpr) {
	// numeric guards:
	if vx.kind <= numeric_beg || vx.kind >= numeric_end {
		v.err = newBitwiseNonIntegerError(vx, binaryExpr.X)
		return
	}
	if vy.kind <= numeric_beg || vy.kind >= numeric_end {
		v.err = newBitwiseNonIntegerError(vy, binaryExpr.Y)
		return
	}

	var x, y int64
	switch v.options.numericType {
	case NumericTypeFloat:
		v.err = newBitwiseNonIntegerError(v, binaryExpr)
		return
	case NumericTypeComplex:
		v.err = newBitwiseNonIntegerError(v, binaryExpr)
		return
	case NumericTypeInt:
		x = parseInt(vx.value)
		y = parseInt(vy.value)
	case NumericTypeAuto:
		var ok bool
		x, ok = convertToInt64(vx.value)
		if !ok {
			v.err = newBitwiseNonIntegerError(vx, binaryExpr.X)
			return
		}
		y, ok = convertToInt64(vy.value)
		if !ok {
			v.err = newBitwiseNonIntegerError(vy, binaryExpr.Y)
			return
		}
	}

	v.kind = KindInt
	switch binaryExpr.Op {
	case token.AND:
		v.value = x & y
	case token.OR:
		v.value = x | y
	case token.XOR:
		v.value = x ^ y
	case token.AND_NOT:
		v.value = x &^ y
	case token.SHL:
		v.value = x << y
	case token.SHR:
		v.value = x >> y
	}
}

func newBitwiseNonIntegerError(v *Visitor, e ast.Expr) error {
	s := conv.FormatExpr(e)
	return &SyntaxError{
		Msg: "result value of \"" + s + "\" is \"" + fmt.Sprintf("%v", v.value) + "\" which is not an integer",
		Pos: int(e.Pos()),
		Err: ErrBitwiseOperation,
	}
}

func convertToInt64(value interface{}) (int64, bool) {
	switch val := value.(type) {
	case float64:
		if float64(int64(val)) == val { // only if it doesn't have decimal.
			return int64(val), true
		}
	case int64:
		return val, true
	}
	return 0, false
}
