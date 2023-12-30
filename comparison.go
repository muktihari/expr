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

func comparison(v, vx, vy *Visitor, binaryExpr *ast.BinaryExpr) {
	v.kind = KindBoolean
	if vx.kind == KindBoolean && vy.kind == KindBoolean {
		compareMustBoolean(v, vx, vy, binaryExpr)
		return
	}
	if vx.kind == KindString && vy.kind == KindString {
		compareMustString(v, vx, vy, binaryExpr)
		return
	}

	// numeric can be compare one another e.g. 0.4 < 1 -> true
	// numeric guards:
	if vx.kind <= numeric_beg || vx.kind >= numeric_end {
		v.err = newComparisonNonNumericError(vx, binaryExpr.X)
		return
	}
	if vy.kind <= numeric_beg || vy.kind >= numeric_end {
		v.err = newComparisonNonNumericError(vy, binaryExpr.Y)
		return
	}

	if vx.kind == KindImag || vy.kind == KindImag {
		compareComplex(v, vx, vy, binaryExpr)
		return
	}
	if vx.kind == KindFloat || vy.kind == KindFloat {
		compareFloat(v, vx, vy, binaryExpr)
		return
	}
	if vx.kind == KindInt || vy.kind == KindInt {
		compareInt(v, vx, vy, binaryExpr)
		return
	}
}

func newComparisonNonNumericError(v *Visitor, e ast.Expr) error {
	s := conv.FormatExpr(e)
	return &SyntaxError{
		Msg: "result of \"" + s + "\" is \"" + fmt.Sprintf("%v", v.value) + "\" which is not a number",
		Pos: v.pos,
		Err: ErrComparisonOperation,
	}
}

func compareMustBoolean(v, vx, vy *Visitor, binaryExpr *ast.BinaryExpr) {
	x := vx.value.(bool)
	y := vy.value.(bool)

	switch binaryExpr.Op {
	case token.EQL:
		v.value = x == y
	case token.NEQ:
		v.value = x != y
	default:
		v.kind = KindIllegal
		v.err = &SyntaxError{
			Msg: "operator \"" + binaryExpr.Op.String() + "\" is not supported for comparing boolean values",
			Pos: int(binaryExpr.OpPos),
			Err: ErrUnsupportedOperator,
		}
		return
	}
}

func compareMustString(v, vx, vy *Visitor, binaryExpr *ast.BinaryExpr) {
	x := vx.value.(string)
	y := vy.value.(string)

	switch binaryExpr.Op {
	case token.EQL:
		v.value = x == y
	case token.NEQ:
		v.value = x != y
	case token.GTR:
		v.value = x > y
	case token.GEQ:
		v.value = x >= y
	case token.LSS:
		v.value = x < y
	case token.LEQ:
		v.value = x <= y
	}
}

// IEEE 754 says that only NaNs satisfy f != f.
func compareComplex(v, vx, vy *Visitor, binaryExpr *ast.BinaryExpr) {
	x := parseComplex(vx.value, vx.kind)
	y := parseComplex(vy.value, vy.kind)

	switch binaryExpr.Op {
	case token.EQL:
		v.value = x == y
	case token.NEQ:
		v.value = x != y
	default:
		v.kind = KindIllegal
		v.err = &SyntaxError{
			Msg: "operator \"" + binaryExpr.Op.String() + "\" is not supported for comparing complex numbers",
			Pos: int(binaryExpr.OpPos),
			Err: ErrUnsupportedOperator,
		}
		return
	}
}

func compareFloat(v, vx, vy *Visitor, binaryExpr *ast.BinaryExpr) {
	x := parseFloat(vx.value, vx.kind)
	y := parseFloat(vy.value, vy.kind)

	switch binaryExpr.Op {
	case token.EQL:
		v.value = x == y
	case token.NEQ:
		v.value = x != y
	case token.GTR:
		v.value = x > y
	case token.GEQ:
		v.value = x >= y
	case token.LSS:
		v.value = x < y
	case token.LEQ:
		v.value = x <= y
	}
}

func compareInt(v, vx, vy *Visitor, binaryExpr *ast.BinaryExpr) {
	x := parseInt(vx.value, vx.kind)
	y := parseInt(vy.value, vy.kind)

	switch binaryExpr.Op {
	case token.EQL:
		v.value = x == y
	case token.NEQ:
		v.value = x != y
	case token.GTR:
		v.value = x > y
	case token.GEQ:
		v.value = x >= y
	case token.LSS:
		v.value = x < y
	case token.LEQ:
		v.value = x <= y
	}
}
