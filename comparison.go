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
	v.kind = KindBoolean
	switch x := vx.value.(type) {
	case complex128:
		// Treat y as complex number since x is a complex number.
		switch y := vy.value.(type) {
		case complex128:
			compareComplex(v, x, y, binaryExpr.Op, binaryExpr.OpPos)
			return
		case float64:
			compareComplex(v, x, complex(y, 0), binaryExpr.Op, binaryExpr.OpPos)
			return
		case int64:
			compareComplex(v, x, complex(float64(y), 0), binaryExpr.Op, binaryExpr.OpPos)
			return
		}
	case float64:
		switch y := vy.value.(type) {
		case complex128: // Treat x as complex number since y is a complex number.
			compareComplex(v, complex(x, 0), y, binaryExpr.Op, binaryExpr.OpPos)
			return
		case float64:
			compareFloat(v, x, y, binaryExpr.Op)
			return
		case int64: // Treat y as float64 since x is a float64.
			compareFloat(v, x, float64(y), binaryExpr.Op)
			return
		}
	case int64:
		// Treat x as y's type, since int64 hierarchy is at the bottom.
		switch y := vy.value.(type) {
		case complex128:
			compareComplex(v, complex(float64(x), 0), y, binaryExpr.Op, binaryExpr.OpPos)
			return
		case float64:
			compareFloat(v, float64(x), y, binaryExpr.Op)
			return
		case int64:
			compareInt(v, x, y, binaryExpr)
			return
		}
	case bool:
		y, ok := vy.value.(bool)
		if ok {
			compareBoolean(v, x, y, binaryExpr.Op, binaryExpr.OpPos)
			return
		}
	case string:
		y, ok := vy.value.(string)
		if ok {
			compareString(v, x, y, binaryExpr.Op)
			return
		}
	}
	v.kind = KindIllegal
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
		v.value = x == y
	case token.NEQ:
		v.value = x != y
	default:
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
func compareComplex(v *Visitor, x, y complex128, op token.Token, opPos token.Pos) {
	switch op {
	case token.EQL:
		v.value = x == y
	case token.NEQ:
		v.value = x != y
	default:
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

func compareInt(v *Visitor, x, y int64, binaryExpr *ast.BinaryExpr) {
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
