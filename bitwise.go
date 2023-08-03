package expr

import (
	"go/ast"
	"go/token"
	"strconv"

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
	}

	x, _ := strconv.ParseInt(vx.value, 0, 64)
	y, _ := strconv.ParseInt(vy.value, 0, 64)

	v.kind = KindInt
	switch binaryExpr.Op {
	case token.AND:
		v.value = strconv.FormatInt(x&y, 10)
	case token.OR:
		v.value = strconv.FormatInt(x|y, 10)
	case token.XOR:
		v.value = strconv.FormatInt(x^y, 10)
	case token.AND_NOT:
		v.value = strconv.FormatInt(x&^y, 10)
	case token.SHL:
		v.value = strconv.FormatInt(x<<y, 10)
	case token.SHR:
		v.value = strconv.FormatInt(x>>y, 10)
	}
}

func newBitwiseNonIntegerError(v *Visitor, e ast.Expr) error {
	s := conv.FormatExpr(e)
	return &SyntaxError{
		Msg: "result value of \"" + s + "\" is \"" + v.value + "\" which is not an integer",
		Pos: int(e.Pos()),
		Err: ErrBitwiseOperation,
	}
}
