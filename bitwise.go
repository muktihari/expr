package expr

import (
	"fmt"
	"go/ast"
	"go/token"

	"github.com/muktihari/expr/conv"
)

func bitwise(v, vx, vy *Visitor, binaryExpr *ast.BinaryExpr) {
	if v.options.numericType != NumericTypeAuto && v.options.numericType != NumericTypeInt {
		v.err = &SyntaxError{
			Msg: "could not do bitwise operation: numeric type is treated as non-integer",
			Pos: int(binaryExpr.OpPos),
			Err: ErrBitwiseOperation,
		}
		return
	}

	if vx.kind != KindInt {
		v.err = newBitwiseNonIntegerError(vx, binaryExpr.X)
		return
	}

	if vy.kind != KindInt {
		v.err = newBitwiseNonIntegerError(vy, binaryExpr.Y)
		return
	}

	x := vx.value.(int64)
	y := vy.value.(int64)

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
