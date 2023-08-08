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

	// NumericTypeAuto: check whether both values are represent integers
	if v.options.numericType == NumericTypeAuto {
		x := parseFloat(vx.value, vx.kind)
		f := parseFloat(vy.value, vy.kind)

		if x != float64(int64(x)) {
			v.err = newBitwiseNonIntegerError(vx, binaryExpr.X)
			return
		}

		if f != float64(int64(f)) {
			v.err = newBitwiseNonIntegerError(vy, binaryExpr.Y)
			return
		}
	}

	x := parseInt(vx.value, vx.kind)
	y := parseInt(vy.value, vy.kind)

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
