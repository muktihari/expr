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

// comparison compares visitor X and visitor Y values. Numeric hierarchy will apply: complex128 > float64 > int64.
func comparison(v, vx, vy *Visitor, binaryExpr *ast.BinaryExpr) {
	v.value.SetKind(KindBoolean)
	switch vx.value.Kind() {
	case KindImag:
		// Treat y as complex number since x is a complex number.
		switch vy.value.Kind() {
		case KindImag:
			compareComplex(v, vx.value.Complex128(), vy.value.Complex128(), binaryExpr.Op, binaryExpr.OpPos)
			return
		case KindFloat:
			compareComplex(v, vx.value.Complex128(), complex(vy.value.Float64(), 0), binaryExpr.Op, binaryExpr.OpPos)
			return
		case KindInt:
			compareComplex(v, vx.value.Complex128(), complex(float64(vy.value.Int64()), 0), binaryExpr.Op, binaryExpr.OpPos)
			return
		}
	case KindFloat:
		switch vy.value.Kind() {
		case KindImag: // Treat x as complex number since y is a complex number.
			compareComplex(v, complex(vx.value.Float64(), 0), vy.value.Complex128(), binaryExpr.Op, binaryExpr.OpPos)
			return
		case KindFloat:
			compareFloat(v, vx.value.Float64(), vy.value.Float64(), binaryExpr.Op)
			return
		case KindInt: // Treat y as float64 since x is a float64.
			compareFloat(v, vx.value.Float64(), float64(vy.value.Int64()), binaryExpr.Op)
			return
		}
	case KindInt:
		// Treat x as y's type, since int64 hierarchy is at the bottom.
		switch vy.value.Kind() {
		case KindImag:
			compareComplex(v, complex(float64(vx.value.Int64()), 0), vy.value.Complex128(), binaryExpr.Op, binaryExpr.OpPos)
			return
		case KindFloat:
			compareFloat(v, float64(vx.value.Int64()), vy.value.Float64(), binaryExpr.Op)
			return
		case KindInt:
			compareInt(v, vx.value.Int64(), vy.value.Int64(), binaryExpr)
			return
		}
	case KindBoolean:
		if vy.value.Kind() == KindBoolean {
			compareBoolean(v, vx.value.Bool(), vy.value.Bool(), binaryExpr.Op, binaryExpr.OpPos)
			return
		}
	case KindString:
		if vy.value.Kind() == KindString {
			compareString(v, vx.value.String(), vy.value.String(), binaryExpr.Op)
			return
		}
	}
	v.value.SetKind(KindIllegal)
	v.err = newComparisonNonComparableError(v, binaryExpr) // Catch non-comparable values.
}

func newComparisonNonComparableError(v *Visitor, e ast.Expr) error {
	return &SyntaxError{
		Msg: fmt.Sprintf("expression %q is not comparable", conv.FormatExpr(e)),
		Pos: v.pos,
		Err: ErrComparisonOperation,
	}
}

func compareBoolean(v *Visitor, x, y bool, op token.Token, opPos token.Pos) {
	switch op {
	case token.EQL:
		v.value = boolValue(x == y)
	case token.NEQ:
		v.value = boolValue(x != y)
	default:
		v.value = value{}
		v.err = &SyntaxError{
			Msg: "operator \"" + op.String() + "\" is not supported for comparing boolean values",
			Pos: int(opPos),
			Err: ErrUnsupportedOperator,
		}
		return
	}
}

func compareString(v *Visitor, x, y string, op token.Token) {
	switch op {
	case token.EQL:
		v.value = boolValue(x == y)
	case token.NEQ:
		v.value = boolValue(x != y)
	case token.GTR:
		v.value = boolValue(x > y)
	case token.GEQ:
		v.value = boolValue(x >= y)
	case token.LSS:
		v.value = boolValue(x < y)
	case token.LEQ:
		v.value = boolValue(x <= y)
	}
}

// IEEE 754 says that only NaNs satisfy f != f.
func compareComplex(v *Visitor, x, y complex128, op token.Token, opPos token.Pos) {
	switch op {
	case token.EQL:
		v.value = boolValue(x == y)
	case token.NEQ:
		v.value = boolValue(x != y)
	default:
		v.value = value{}
		v.err = &SyntaxError{
			Msg: "operator \"" + op.String() + "\" is not supported for comparing complex numbers",
			Pos: int(opPos),
			Err: ErrUnsupportedOperator,
		}
		return
	}
}

func compareFloat(v *Visitor, x, y float64, op token.Token) {
	switch op {
	case token.EQL:
		v.value = boolValue(x == y)
	case token.NEQ:
		v.value = boolValue(x != y)
	case token.GTR:
		v.value = boolValue(x > y)
	case token.GEQ:
		v.value = boolValue(x >= y)
	case token.LSS:
		v.value = boolValue(x < y)
	case token.LEQ:
		v.value = boolValue(x <= y)
	}
}

func compareInt(v *Visitor, x, y int64, binaryExpr *ast.BinaryExpr) {
	switch binaryExpr.Op {
	case token.EQL:
		v.value = boolValue(x == y)
	case token.NEQ:
		v.value = boolValue(x != y)
	case token.GTR:
		v.value = boolValue(x > y)
	case token.GEQ:
		v.value = boolValue(x >= y)
	case token.LSS:
		v.value = boolValue(x < y)
	case token.LEQ:
		v.value = boolValue(x <= y)
	}
}
