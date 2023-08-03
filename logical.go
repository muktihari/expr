package expr

import (
	"go/ast"
	"go/token"
	"strconv"

	"github.com/muktihari/expr/conv"
)

func logical(v, vx, vy *Visitor, binaryExpr *ast.BinaryExpr) {
	if vx.kind != KindBoolean {
		v.err = newLogicalNonBooleanError(vx, binaryExpr.X)
		return
	}

	if vy.kind != KindBoolean {
		v.err = newLogicalNonBooleanError(vy, binaryExpr.Y)
		return
	}

	x, _ := strconv.ParseBool(vx.value)
	y, _ := strconv.ParseBool(vy.value)

	v.kind = KindBoolean
	if binaryExpr.Op == token.LAND {
		v.value = strconv.FormatBool(x && y)
		return
	}

	v.value = strconv.FormatBool(x || y) // token.LOR
}

func newLogicalNonBooleanError(v *Visitor, e ast.Expr) error {
	s := conv.FormatExpr(e)
	return &SyntaxError{
		Msg: "result of \"" + s + "\" is \"" + v.value + "\" which is not a boolean",
		Pos: v.pos,
		Err: ErrLogicalOperation,
	}
}
