package expr

import (
	"fmt"
	"go/ast"
	"go/token"

	"github.com/muktihari/expr/conv"
)

func bitwise(v, vx, vy *Visitor, binaryExpr *ast.BinaryExpr) {
	// No matters what options, having bolean here is invalid.
	if vx.kind == KindBoolean {
		v.err = newBitwiseNonIntegerError(vx, binaryExpr.X)
		return
	}
	if vy.kind == KindBoolean {
		v.err = newBitwiseNonIntegerError(vx, binaryExpr.Y)
		return
	}

	switch v.options.numericType {
	case NumericTypeAuto:
		// NumericTypeAuto: check whether both values are represent integers
		x := parseFloat(vx.value, vx.kind)
		y := parseFloat(vy.value, vy.kind)

		if x != float64(int64(x)) {
			v.err = newBitwiseNonIntegerError(vx, binaryExpr.X)
			return
		}

		if y != float64(int64(y)) {
			v.err = newBitwiseNonIntegerError(vy, binaryExpr.Y)
			return
		}
	case NumericTypeFloat, NumericTypeComplex:
		v.err = newBitwiseNonIntegerError(vy, binaryExpr.Y)
		return
	}

	x := parseInt(vx.value, vx.kind)
	y := parseInt(vy.value, vy.kind)

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
